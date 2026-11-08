# Compiler Engine Architecture & Visitor Blueprint

This document details the internal design, ANTLR4 grammar pipeline, and code emission architecture of the Go-based transpiler CLI (`compiler-engine`).

---

## 1. Overview & System Pipeline

The `compiler-engine` is a standalone Go program (`j2cpp-engine`) that converts Java source code into modern C++20 header (`.h`) and implementation (`.cpp`) files.

```
                  ┌──────────────────────────────────────────────┐
                  │              Input Java Files                │
                  └──────────────────────┬───────────────────────┘
                                         │
                                         ▼
                  ┌──────────────────────────────────────────────┐
                  │    ANTLR4 Lexer & Parser (parser/*.go)       │
                  │       (Generated from grammars-v4)           │
                  └──────────────────────┬───────────────────────┘
                                         │
                                         ▼
                  ┌──────────────────────────────────────────────┐
                  │ Parse Tree (compilationUnit -> classBody...) │
                  └──────────────────────┬───────────────────────┘
                                         │
                                         ▼
                  ┌──────────────────────────────────────────────┐
                  │    CppTranspilerVisitor (visitor.go)         │
                  │   - Type Mapping & Symbol Resolution         │
                  │   - Null Safety & Memory Management (GC)     │
                  │   - RAII Guard & Exception Synthesis         │
                  └──────────────────────┬───────────────────────┘
                                         │
                     ┌───────────────────┴───────────────────┐
                     ▼                                       ▼
          ┌─────────────────────┐                 ┌─────────────────────┐
          │  Header File (.h)   │                 │ Implementation (.cpp)│
          │  - Class definition │                 │  - Method bodies    │
          │  - Virtual vtables  │                 │  - Expressions      │
          │  - Default fields   │                 │  - RAII guards      │
          └─────────────────────┘                 └─────────────────────┘
```

---

## 2. ANTLR4 Integration Mechanics

The transpiler leverages official Java grammars from `grammars-v4/java/java`:
- **`JavaLexer.g4`**: Tokenizes input character streams into tokens (`IF`, `WHILE`, `IDENTIFIER`, `STRING_LITERAL`, etc.).
- **`JavaParser.g4`**: Builds an AST according to the Java language specification.

### Code Generation (`scripts/generate_parser.sh`)
```bash
java -jar /usr/local/lib/antlr-4.13.1-complete.jar \
     -Dlanguage=Go -visitor -package parser -o compiler-engine/parser \
     grammars-v4/java/java/JavaLexer.g4 grammars-v4/java/java/JavaParser.g4
```
The output Go package in `compiler-engine/parser/` exposes:
- `JavaLexer`: The scanner.
- `JavaParser`: The parser returning `ICompilationUnitContext`.
- `JavaParserVisitor`: Interface defining visitor callbacks for every grammar production.
- `BaseJavaParserVisitor`: Default no-op implementation embedded by `CppTranspilerVisitor`.

### Initializing the ANTLR4 Parser Stream
The parsing entrypoint in `compiler-engine` sets up the lexer and parser pipeline as follows:
```go
// 1. Create character stream from input source code
inputStream := antlr.NewInputStream(javaSource)

// 2. Lexical analysis: generate token stream
lexer := parser.NewJavaLexer(inputStream)
tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

// 3. Syntactic analysis: construct parser instance
javaParser := parser.NewJavaParser(tokenStream)

// 4. Parse starting from root rule compilationUnit
tree := javaParser.CompilationUnit()

// 5. Walk tree with custom transpiler visitor
visitor := NewCppTranspilerVisitor(className, superName)
visitor.Visit(tree)
```

---

## 3. Core Visitor: `CppTranspilerVisitor`

Defined in `compiler-engine/visitor.go`, `CppTranspilerVisitor` walks the AST and accumulates C++ code in string buffers (`headerCode` and `cppCode`).

### State & Buffers
```go
type CppTranspilerVisitor struct {
    parser.BaseJavaParserVisitor
    headerCode       strings.Builder
    cppCode          strings.Builder
    currentClass     string
    currentSuper     string
    currentIsStatic  bool
    staticFields     map[string]bool
    included         map[string]bool
    currentLoopDepth int
}
```

### Key AST Lowerings

1. **Class Declarations (`VisitClassDeclaration`)**:
   - Superclasses (`extends SuperClass`) map to `public SuperClass`. Default base is `public java::lang::Object`.
   - Interfaces (`implements InterfaceA, InterfaceB`) map to multiple virtual base inheritance.
   - Outputs `#pragma once` and standard `#include` dependencies for runtime headers (`Object.h`, `String.h`, `GC.h`, etc.).

2. **Constructors & Initializers (`VisitConstructorDeclaration`)**:
   - Detects leading `this(...)` or `super(...)` calls in the constructor body.
   - Lowers them into C++ constructor member initializer lists (`: Sibling(args)` or `: Super(args)`).
   - Generates constructor signature in `.h` and definition in `.cpp`.

3. **Field Declarations & Defaults (`VisitFieldDeclaration`)**:
   - Evaluates explicit initializers if present.
   - Synthesizes Java-compliant default values when uninitialized:
     - `bool` ➔ `false`
     - numeric/char ➔ `0`
     - float/double ➔ `0.0f` / `0.0`
     - pointer/reference ➔ `nullptr`
   - Static fields are emitted with `inline static`.

4. **Methods & Polymorphism (`VisitMethodDeclaration`)**:
   - Methods default to `virtual` in C++ for dynamic dispatch matching Java semantics.
   - Static methods are emitted as `static` member functions.

5. **Memory Allocation & Null Safety**:
   - Object instantiation `new T(...)` is lowered to `java_new<T>(...)`, which routes allocations through Boehm GC (`GC_MALLOC`).
   - Field and method dereferences are wrapped with `JAVA_NULL_CHECK(obj)->method(...)` to throw `NullPointerException` on null dereference without evaluation side effects.

6. **Exceptions & RAII Finally (`VisitTryStatement`)**:
   - `try-catch-finally` constructs emit native C++ `try` and `catch (java::lang::Exception* e)` blocks.
   - Java `finally { ... }` blocks are converted to local RAII cleanup guards using `JAVA_FINALLY([&] { ... });`, guaranteeing execution regardless of early returns or thrown exceptions.

7. **Lambdas & Method References**:
   - Lambdas `(args) -> expr` map to C++20 lambda closures `[=](params) { return expr; }`.
   - Method references `Type::method` or `expr::method` map to generic forwarding lambdas `[=](auto&&... args) { return ...; }`.

---

## 3.1 Type System & Symbol Resolution Matrix

`CppTranspilerVisitor` uses deterministic type mapping functions (`mapType`) to maintain strict C++20 equivalence:

```go
func mapType(javaType string) string {
    // Array types: T[] -> java::lang::Array<T>*
    if strings.HasSuffix(javaType, "[]") {
        elem := mapType(strings.TrimSuffix(javaType, "[]"))
        return fmt.Sprintf("java::lang::Array<%s>*", elem)
    }

    // Generic collection types: List<T>, Map<K, V>, Set<E>
    if strings.HasPrefix(javaType, "ArrayList<") { ... }
    if strings.HasPrefix(javaType, "LinkedList<") { ... }
    if strings.HasPrefix(javaType, "HashMap<") { ... }
    if strings.HasPrefix(javaType, "HashSet<") { ... }

    // Primitive & standard mappings
    switch javaType {
    case "int": return "int32_t"
    case "boolean": return "bool"
    case "long": return "int64_t"
    case "byte": return "int8_t"
    case "short": return "int16_t"
    case "char": return "char16_t"
    case "float": return "float"
    case "double": return "double"
    case "String": return "java::lang::String*"
    case "StringBuilder": return "java::lang::StringBuilder*"
    case "AtomicInteger": return "java::util::concurrent::atomic::AtomicInteger*"
    case "AtomicBoolean": return "java::util::concurrent::atomic::AtomicBoolean*"
    case "AtomicLong": return "java::util::concurrent::atomic::AtomicLong*"
    case "ReentrantLock": return "java::util::concurrent::locks::ReentrantLock*"
    default:
        // User-defined class types map to pointer references
        return javaType + "*"
    }
}
```

---

## 4. Execution Modes & CLI

The CLI (`compiler-engine/main.go`) operates in two modes:

### Batch Compilation
```bash
./compiler-engine/j2cpp-engine -src=/path/to/java/src -out=/path/to/output
```
Scans the source directory recursively for `.java` files, parses each file using `JavaParser`, runs `CppTranspilerVisitor`, and writes `<ClassName>.h` and `<ClassName>.cpp` pairs to the output directory.

### Interactive Web Dashboard Studio
```bash
./compiler-engine/j2cpp-engine -dashboard -port=8080
```
Launches an embedded HTTP server hosting a responsive single-page IDE:
- Interactive code editor with live syntax transpilation.
- Side-by-side header (`.h`) and implementation (`.cpp`) tabbed outputs.
- Tooltips explaining transpiler behavior and lowering rules.
- REST endpoints:
  - `POST /api/transpile`: Dynamically transpiles code snippets.
  - `GET /api/status`: Returns compiler engine version, features, and runtime capabilities.
