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

## 2. ANTLR4 Grammar & Submodule Integration Mechanics

The transpiler frontend leverages official Java grammars maintained inside the `grammars-v4` git submodule at `grammars-v4/java/java`:
- **`JavaLexer.g4`**: Complete lexical specification defining unicode identifiers, keywords (`class`, `interface`, `extends`, `implements`, `synchronized`, `default`, `static`), literals, and operators.
- **`JavaParser.g4`**: Comprehensive context-free production rules for compilation units, class and interface declarations, member fields, constructors, methods, control flow statements, lambdas, and expressions.

### 2.1 Submodule Checkout & Workspace Initialization
To initialize the `grammars-v4` submodule alongside Boehm GC (`bdwgc`), the CLI setup workflow requires recursive git submodule synchronization:
```bash
# Initialize and fetch grammars-v4 and bdwgc submodules
git submodule update --init --recursive
```
The grammar files are verified to exist at:
- `grammars-v4/java/java/JavaLexer.g4`
- `grammars-v4/java/java/JavaParser.g4`

### 2.2 Go Parser Code Generation Pipeline (`scripts/generate_parser.sh`)
ANTLR4 parser Go bindings are generated via the automated script `scripts/generate_parser.sh`. The generator invokes the ANTLR 4.13.1 tool targeting the Go runtime:
```bash
java -jar /usr/local/lib/antlr-4.13.1-complete.jar \
     -Dlanguage=Go \
     -visitor \
     -package parser \
     -o compiler-engine/parser \
     grammars-v4/java/java/JavaLexer.g4 \
     grammars-v4/java/java/JavaParser.g4
```

**Predicate Receiver Normalization:**
When ANTLR4 generates Go code for grammars containing semantic predicates, some predicate methods may reference `this` instead of the Go receiver identifier `p`. `scripts/generate_parser.sh` automatically applies an automated normalization step using `sed`:
```bash
sed -i 's/\bthis\./p\./g' compiler-engine/parser/*.go
```

The output Go package in `compiler-engine/parser/` exposes:
- `JavaLexer`: Tokenizer scanning input character streams into categorized tokens (`IDENTIFIER`, `DECIMAL_LITERAL`, `STRING_LITERAL`, `CLASS`, `INTERFACE`, etc.).
- `JavaParser`: Predictive LL(*) recursive-descent parser returning parse tree contexts starting at root rule `compilationUnit`:
  - `ICompilationUnitContext`: Root unit parsing package declarations, imports, and top-level types.
  - `IClassDeclarationContext`: Class header, superclass (`extends`), interfaces (`implements`), and body.
  - `IInterfaceDeclarationContext`: Interface definitions and signatures.
  - `IClassBodyDeclarationContext`: Methods, constructors, fields, and static initializers.
  - `IStatementContext`: Control flow, blocks, synchronized locks, and exception blocks.
  - `IExpressionContext`: Expressions, lambda syntax, and method references.
- `JavaParserVisitor`: Interface defining visitor callback signatures for all grammar rules.
- `BaseJavaParserVisitor`: Default no-op implementation embedded by `CppTranspilerVisitor`.

### 2.3 Initializing the ANTLR4 Parser Stream in Go
The CLI bootstrap sets up the ANTLR4 input stream, token stream, and parser pipeline:
```go
// 1. Create character stream from input source code (or antlr.NewFileStream for disk files)
inputStream := antlr.NewInputStream(javaSource)

// 2. Lexical analysis: instantiate lexer and wrap in common token stream
lexer := parser.NewJavaLexer(inputStream)
tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

// 3. Syntactic analysis: construct parser instance with parse tree construction enabled
javaParser := parser.NewJavaParser(tokenStream)
javaParser.BuildParseTrees = true

// Optional error listener hook for diagnostics
// javaParser.RemoveErrorListeners()
// javaParser.AddErrorListener(antlr.NewConsoleErrorListener())

// 4. Parse starting from root rule compilationUnit
tree := javaParser.CompilationUnit()

// 5. Walk tree with custom transpiler visitor configured with output directory
visitor := NewCppTranspilerVisitor(outDir)
visitor.Visit(tree)
```

---

## 3. Core Visitor: `CppTranspilerVisitor`

Defined in `compiler-engine/visitor.go`, `CppTranspilerVisitor` walks the AST and accumulates C++ code in string buffers (`headerCode` and `cppCode`).

### State & Buffers
```go
type CppTranspilerVisitor struct {
    *parser.BaseJavaParserVisitor
    outDir                string
    currentClass          string
    currentSuper          string
    headerCode            *strings.Builder
    cppCode               *strings.Builder
    included              map[string]bool
    currentIsStatic       bool
    currentIsSynchronized bool
    staticInitCount       int
    staticFields          map[string]bool
}

func NewCppTranspilerVisitor(outDir string) *CppTranspilerVisitor {
    return &CppTranspilerVisitor{
        BaseJavaParserVisitor: &parser.BaseJavaParserVisitor{},
        outDir:                outDir,
        headerCode:           &strings.Builder{},
        cppCode:              &strings.Builder{},
        included:             make(map[string]bool),
        staticFields:         make(map[string]bool),
    }
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

## 4. Execution Modes & CLI Entrypoint

The CLI entrypoint is defined in `compiler-engine/main.go`. It handles command-line arguments using Go's `flag` standard library, manages input validation, and dispatches to either the interactive web dashboard or the recursive batch compiler.

### Command-Line Arguments & Parsing Logic
```go
func main() {
    srcDir := flag.String("src", "", "Source directory containing Java files")
    outDir := flag.String("out", "", "Output directory for C++ files")
    dashboard := flag.Bool("dashboard", false, "Start interactive web dashboard studio")
    port := flag.Int("port", 8080, "Port for web dashboard studio")
    flag.Parse()

    // 1. Interactive Studio Dispatch
    if *dashboard {
        if err := startDashboardServer(*port); err != nil {
            log.Fatalf("Dashboard server failed: %v", err)
        }
        return
    }

    // 2. Argument Validation for Batch Mode
    if *srcDir == "" || *outDir == "" {
        fmt.Println("Usage: j2cpp-engine -src=<source_directory> -out=<output_directory> [-dashboard -port=8080]")
        os.Exit(1)
    }

    // 3. Ensure Output Directory Exists
    err := os.MkdirAll(*outDir, os.ModePerm)
    if err != nil {
        log.Fatalf("Failed to create output directory: %v", err)
    }

    // 4. Recursive Source Directory Traversal
    err = filepath.WalkDir(*srcDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() && strings.HasSuffix(d.Name(), ".java") {
            processJavaFile(path, *outDir)
        }
        return nil
    })

    if err != nil {
        log.Fatalf("Error walking source directory: %v", err)
    }
}
```

### File Processing Pipeline (`processJavaFile`)
For each `.java` file discovered during recursive directory walking:
```go
func processJavaFile(javaFile string, outDir string) {
    fmt.Printf("Processing %s\n", javaFile)

    input, err := antlr.NewFileStream(javaFile)
    if err != nil {
        log.Printf("Failed to read file %s: %v", javaFile, err)
        return
    }

    lexer := parser.NewJavaLexer(input)
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
    p := parser.NewJavaParser(stream)

    tree := p.CompilationUnit()

    visitor := NewCppTranspilerVisitor(outDir)
    visitor.Visit(tree)
}
```

### CLI Execution Examples

#### 1. Recursive Batch Compilation
```bash
./compiler-engine/j2cpp-engine -src=/path/to/java/src -out=/path/to/output
```
Scans the source directory recursively for `.java` files, parses each file using `JavaParser`, runs `CppTranspilerVisitor`, and writes `<ClassName>.h` and `<ClassName>.cpp` pairs to the output directory.

#### 2. Interactive Web Dashboard Studio
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
