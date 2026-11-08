# libjruntime

A Go-based Transpiler CLI and a modern C++20 Native Runtime Library designed to compile and run unmodified Java source code natively in C++ using conservative garbage collection (Boehm GC).

## Architecture & Overview

Goal: Run unmodified Java source code natively by transpiling it to C++20.

### Submodules & Dependencies
For deep architectural insights, build configuration flags, memory management hooks, and grammar details, see the dedicated [SUBMODULES.md](SUBMODULES.md) guide. For internal transpiler design, ANTLR4 AST visitor lowerings, and symbol resolution mechanics, see [compiler-engine/ARCHITECTURE.md](compiler-engine/ARCHITECTURE.md).

### SECTION 1: ANTLR4 Grammar Integration & Transpiler CLI Setup

The Go-based transpiler front-end relies on the official Java ANTLR4 grammar specifications maintained in the `grammars-v4` git submodule (`grammars-v4/java/java`):
- `JavaLexer.g4`: Complete lexical specifications defining token patterns, unicode escapes, identifiers, operators, keywords, and numeric/string/character literals.
- `JavaParser.g4`: Context-free grammar rules defining Java compilation units, class and interface declarations, member fields, constructors, methods, control flow statements, and expression trees.

#### 1.1 Grammar Integration Architecture
```
┌─────────────────────────────────┐       ┌────────────────────────────────┐
│   grammars-v4/java/java/*.g4    │ ────> │ scripts/generate_parser.sh     │
└─────────────────────────────────┘       └────────────────────────────────┘
                                                           │
                                                           ▼ (ANTLR 4.13.1 Go target)
┌─────────────────────────────────┐       ┌────────────────────────────────┐
│ compiler-engine/visitor.go      │ <──── │ compiler-engine/parser/*.go    │
│ (CppTranspilerVisitor)          │       │ (Lexer, Parser, Listener, Base)│
└─────────────────────────────────┘       └────────────────────────────────┘
                 │
                 ▼ (AST Traversal & Semantic Mapping)
┌─────────────────────────────────┐
│ C++20 Header (.h) & Body (.cpp) │
└─────────────────────────────────┘
```

#### 1.2 Initial Transpiler CLI Stream Setup & Bootstrap Code
The Go CLI initializes the ANTLR4 input stream, token stream, and parser using standard Go runtime bindings. Here is the concrete setup used to parse arbitrary Java source strings:

```go
package main

import (
    "fmt"
    "github.com/antlr4-go/antlr/v4"
    "github.com/user/j2cpp-engine/parser"
)

// TranspileJava parses Java source and emits C++20 header and body code
func TranspileJava(javaSource, className, superName string) (headerCode, cppCode string, err error) {
    // 1. Initialize character input stream
    input := antlr.NewInputStream(javaSource)

    // 2. Tokenize input with JavaLexer
    lexer := parser.NewJavaLexer(input)
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

    // 3. Construct JavaParser instance
    p := parser.NewJavaParser(tokens)
    p.BuildParseTrees = true

    // 4. Attach error listeners if strict syntax checking is desired
    // p.RemoveErrorListeners()
    // p.AddErrorListener(antlr.NewConsoleErrorListener())

    // 5. Parse starting at compilationUnit rule
    tree := p.CompilationUnit()

    // 6. Execute AST reduction with CppTranspilerVisitor
    visitor := NewCppTranspilerVisitor(className, superName)
    visitor.Visit(tree)

    return visitor.GetHeaderCode(), visitor.GetCppCode(), nil
}
```

#### 1.3 Visitor vs. Listener Design Pattern
ANTLR4 provides two primary tree traversal paradigms:
- **Parse Tree Listeners (`JavaParserListener`)**: Passive, event-driven callbacks triggered as a walker enters and exits tree nodes. Listeners do not control traversal flow and cannot easily synthesize values up the tree.
- **Parse Tree Visitors (`CppTranspilerVisitor`)**: Explicit, active reduction where `visitor.Visit(node)` returns intermediate string expressions and AST values. This enables localized code generation (e.g. evaluating an `ExpressionList` or constructor arguments before emitting the enclosing method call).

#### 1.4 ANTLR4 Parsing & Code Generation Pipeline
1. **Lexical Analysis (`JavaLexer`)**: Stream of raw Java code is tokenized into typed tokens (`IDENTIFIER`, `INT`, `RETURN`, `IF`, etc.).
2. **Syntactic Analysis (`JavaParser`)**: Tokens form a parse tree starting from root rule `compilationUnit`:
   - `typeDeclaration` ➔ `classDeclaration` | `interfaceDeclaration`
   - `classBody` ➔ `classBodyDeclaration` ➔ `memberDeclaration` (`methodDeclaration`, `constructorDeclaration`, `fieldDeclaration`)
   - `statement` ➔ `block`, `statementExpression`, `ifThenElse`, `whileStatement`, `forStatement`, `synchronizedStatement`
   - `expression` ➔ `primary`, `expression '.' IDENTIFIER`, `methodCall`, `lambdaExpression`, `methodReference`
3. **AST Visitor Traversal (`CppTranspilerVisitor`)**:
   `CppTranspilerVisitor` embeds `parser.BaseJavaParserVisitor` to implement selective override hooks:
   - `VisitCompilationUnit`: Scans package/import metadata and initiates class lowering.
   - `VisitClassDeclaration`: Extracts class names, superclasses (`extends`), implemented interfaces (`implements`), emitting modern C++ class declarations.
   - `VisitInterfaceDeclaration`: Emits pure virtual abstract C++ base classes or inlined default/static methods.
   - `VisitConstructorDeclaration`: Translates constructors and constructor chaining (`this(...)`, `super(...)`) into C++ member initializer lists.
   - `VisitMethodDeclaration`: Emits virtual and static method signatures into `.h` and implementations into `.cpp`.
   - `VisitExpression`: Translates arithmetic, comparisons, ternary expressions (`cond ? a : b`), method invocations, and object allocations (`new T(...)` ➔ `java_new<T>(...)`).
   - `VisitLambdaExpression` & `VisitMethodReferenceExpression`: Lowers Java lambdas and method references into C++20 universal generic forwarding closures `[=](auto&&... args) { ... }`.

#### 1.5 Type Conversion & Lowering Matrix
The transpiler maps Java types to modern C++20 types and garbage-collected pointers:
| Java Source Type | Transpiled C++20 Type | Header Dependency |
|:---|:---|:---|
| `int`, `long`, `boolean` | `int32_t`, `int64_t`, `bool` | `<cstdint>` |
| `byte`, `short`, `char` | `int8_t`, `int16_t`, `char16_t` | `<cstdint>` |
| `float`, `double` | `float`, `double` | Standard scalar |
| `String` | `java::lang::String*` | `String.h` |
| `StringBuilder` | `java::lang::StringBuilder*` | `StringBuilder.h` |
| `T[]` (arrays) | `java::lang::Array<T>*` | `Array.h` |
| `ArrayList<T>` | `java::util::ArrayList<T>*` | `ArrayList.h` |
| `LinkedList<T>` | `java::util::LinkedList<T>*` | `LinkedList.h` |
| `HashMap<K, V>` | `java::util::HashMap<K, V>*` | `HashMap.h` |
| `HashSet<E>` | `java::util::HashSet<E>*` | `HashSet.h` |
| Boxed Primitives (`Integer`, `Long`, etc.) | `java::lang::*` pointers | `Boxes.h` |
| `AtomicInteger`, `AtomicLong`, `AtomicBoolean` | `java::util::concurrent::atomic::*` | `Atomic.h` |
| `ReentrantLock` | `java::util::concurrent::locks::ReentrantLock*` | `Lock.h` |

#### Automated Parser Regeneration
Parser Go bindings can be re-generated at any time using `scripts/generate_parser.sh`:
```bash
./scripts/generate_parser.sh
```
This runs `antlr-4.13.1-complete.jar` targeting Go with visitor generation enabled (`-Dlanguage=Go -visitor -package parser -o compiler-engine/parser`) and normalizes receiver predicates.

#### CLI Transpiler Usage
Build the CLI binary:
```bash
cd compiler-engine && go build -o j2cpp-engine . && cd ..
```

Execute batch transpilation across any directory of `.java` source files:
```bash
./compiler-engine/j2cpp-engine -src=/path/to/src -out=/path/to/out
```
- `-src`: Directory containing input `.java` files (default: `./src`).
- `-out`: Directory where modern C++20 header (`.h`) and implementation (`.cpp`) files will be emitted (default: `./out`).

Or run the interactive Web Dashboard Studio:
```bash
./compiler-engine/j2cpp-engine -dashboard -port=8080
```

1. **Go Transpiler Engine (`compiler-engine/`)**
   - Uses ANTLR4 with Java grammar (`grammars-v4`).
   - Walks AST with `CppTranspilerVisitor`.
   - Generates `.h` and `.cpp` files for each Java class.
   - Converts primitive types (`int` -> `int32_t`, `boolean` -> `bool`, `long` -> `int64_t`).
   - Generates `virtual` methods for dynamic dispatch.
   - Maps collections (`ArrayList<T>`), primitive wrapper classes (`Integer`, `Boolean`, `Long`), and math functions (`Math`).
   - Injects `JAVA_NULL_CHECK(obj)` macros for null safety.
   - Translates `new T(...)` into `java_new<T>(...)`.
   - Synthesizes `finally` blocks via RAII `JAVA_FINALLY` cleanup guards.

2. **libjruntime C++ Runtime (`libjruntime/`)**
   - Header-only base types in namespace `java::lang`.
   - `Object`: Base class with `hashCode()`, `equals()`, `toString()`, and lazy monitor mutex initialization for `synchronized`.
   - `String`: String wrapper implementing Java object semantics and operator overloading.
   - `Exceptions`: C++ exception hierarchy mapping Java exceptions (`NullPointerException`, `Exception`).
   - `GC.h`: Boehm GC memory allocation abstraction `java_new<T>`.
   - `FinallyGuard.h`: RAII guard executing scope exit cleanup.

## Building and Running

### Prerequisites
- **Go Toolchain**: 1.20+
- **C++ Compiler**: GCC (g++ 10+) or Clang (clang++ 11+) with C++20 support
- **Build Tools**: Autoconf, Automake, Libtool, Make
- **Java Runtime**: JRE/JDK 11+ (for optional ANTLR4 parser regeneration)

### Step 1: Initialize & Build Submodules
Run the automated submodule preparation script to initialize Git submodules (`grammars-v4`, `bdwgc`) and build the Boehm Garbage Collector shared library:
```bash
./scripts/build_submodules.sh
```
Or build manually:
```bash
git submodule update --init --recursive
cd bdwgc && ./autogen.sh && ./configure --enable-threads=posix --enable-shared --disable-static && make -j$(nproc) && cd ..
```

### Step 2: Build the `j2cpp-engine` Transpiler CLI
Compile the Go transpiler CLI into a standalone binary:
```bash
cd compiler-engine
go build -o j2cpp-engine .
cd ..
```

To run transpiler unit tests:
```bash
cd compiler-engine && go test -v ./... && cd ..
```

### Step 3: Transpile Java Source Code
Run `j2cpp-engine` to parse Java source files and generate C++20 header (`.h`) and implementation (`.cpp`) files:
```bash
./compiler-engine/j2cpp-engine -src=/path/to/java/source -out=./output
```
Command options:
- `-src`: Directory containing input `.java` source files (required).
- `-out`: Target directory where transpiled `.h` and `.cpp` files are emitted (required).

### Step 4: Compile Native C++ Application
Compile the emitted C++ files linking against `libjruntime` and Boehm GC:
```bash
g++ -std=c++20 output/*.cpp \
    -I./libjruntime \
    -I./bdwgc/include \
    -L./bdwgc/.libs \
    -lgc \
    -Wl,-rpath,./bdwgc/.libs \
    -o native_app
```

Run the compiled native executable:
```bash
./native_app
```

### Step 5: End-to-End Validation
Execute the automated end-to-end integration and native execution test harness:
```bash
./tests/e2e_test.sh
```
