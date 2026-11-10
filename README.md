# libjruntime

A Go-based Transpiler CLI and a modern C++20 Native Runtime Library designed to compile and run unmodified Java source code natively in C++ using conservative garbage collection (Boehm GC).

## Architecture & Overview

Goal: Run unmodified Java source code natively by transpiling it to C++20.

### Submodules & Dependencies
For deep architectural insights, build configuration flags, memory management hooks, and grammar details, see the dedicated [SUBMODULES.md](SUBMODULES.md) guide. For internal transpiler design, ANTLR4 AST visitor lowerings, and symbol resolution mechanics, see [compiler-engine/ARCHITECTURE.md](compiler-engine/ARCHITECTURE.md).

### SECTION 1: ANTLR4 Grammar Integration & Transpiler CLI Setup

The Go-based transpiler front-end relies on the official Java ANTLR4 grammar specifications maintained in the `grammars-v4` git submodule (`grammars-v4/java/java`):
- `JavaLexer.g4`: Complete lexical specifications defining token patterns, unicode escapes, identifiers, operators, keywords (`class`, `interface`, `extends`, `implements`, `synchronized`, `default`, `static`, etc.), and numeric/string/character literals.
- `JavaParser.g4`: Context-free grammar rules defining Java compilation units (`compilationUnit`), class and interface declarations, member fields, constructors, methods, control flow statements (`if`, `while`, `for`, `switch`), lambdas, and expression trees.

#### 1.1 Grammar Integration Architecture & Workflow
The transpilation pipeline connects the ANTLR4 grammar sources, generated Go parser package, and C++ emitter:

```
┌─────────────────────────────────┐       ┌────────────────────────────────┐
│   grammars-v4/java/java/*.g4    │ ────> │ scripts/generate_parser.sh     │
└─────────────────────────────────┘       └────────────────────────────────┘
                                                           │
                                                           ▼ (ANTLR 4.13.1 Go target & sed normalizer)
┌─────────────────────────────────┐       ┌────────────────────────────────┐
│ compiler-engine/visitor.go      │ <──── │ compiler-engine/parser/*.go    │
│ (CppTranspilerVisitor)          │       │ (Lexer, Parser, Listener, Base)│
└─────────────────────────────────┘       └────────────────────────────────┘
                 │
                 ▼ (AST Traversal, Type Mapping & C++20 Lowering)
┌─────────────────────────────────┐
│ C++20 Header (.h) & Body (.cpp) │
└─────────────────────────────────┘
```

#### 1.2 Initial Transpiler CLI Stream Setup & Bootstrap Code
The Go CLI initializes the ANTLR4 input stream, token stream, and parser using standard Go runtime bindings (`github.com/antlr4-go/antlr/v4`). Here is the concrete setup used to parse arbitrary Java source strings or disk files:

```go
package main

import (
    "fmt"
    "log"
    "github.com/antlr4-go/antlr/v4"
    "github.com/user/j2cpp-engine/parser"
)

// TranspileJava parses Java source and emits C++20 header and body code
func TranspileJava(javaSource, outDir string) (headerCode, cppCode string, err error) {
    // 1. Initialize character input stream (or antlr.NewFileStream(filepath) for disk files)
    input := antlr.NewInputStream(javaSource)

    // 2. Tokenize input with JavaLexer and common token stream on TokenDefaultChannel
    lexer := parser.NewJavaLexer(input)
    tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

    // 3. Construct JavaParser instance and enable AST construction
    p := parser.NewJavaParser(tokens)
    p.BuildParseTrees = true

    // 4. Attach error listeners for syntax validation
    p.RemoveErrorListeners()
    p.AddErrorListener(antlr.NewConsoleErrorListener())

    // 5. Parse starting at compilationUnit root grammar rule
    tree := p.CompilationUnit()

    // 6. Execute AST reduction with CppTranspilerVisitor configured with target output directory
    visitor := NewCppTranspilerVisitor(outDir)
    visitor.Visit(tree)

    return visitor.GetHeaderCode(), visitor.GetCppCode(), nil
}
```

#### 1.3 Step-by-Step Transpiler CLI Stream Setup
1. **Character Stream Creation**: `antlr.NewInputStream(javaSource)` converts raw in-memory strings into unicode codepoints, while `antlr.NewFileStream(path)` reads directly from disk.
2. **Lexical Scanning**: `parser.NewJavaLexer` translates the character stream into discrete tokens according to `JavaLexer.g4`.
3. **Common Token Stream**: `antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)` buffers lexed tokens, filtering out comments and whitespace to channel 0.
4. **Parser Instantiation**: `parser.NewJavaParser(tokens)` creates an LL(*) predictive parser instance with `BuildParseTrees = true`.
5. **Parse Tree Reduction**: Calling `p.CompilationUnit()` begins top-down parse tree derivation, passing the root context to `visitor.Visit(tree)`.

#### 1.4 Initial Structure of the Go Transpiler CLI & `CppTranspilerVisitor`
The transpiler CLI is organized into three primary units:
1. `compiler-engine/main.go`: CLI parsing (`-src`, `-out`, `-dashboard`, `-port`), recursive file discovery, and orchestrating the ANTLR4 pipeline.
2. `compiler-engine/visitor.go`: Implements `CppTranspilerVisitor`, embedding `parser.BaseJavaParserVisitor`.
3. `compiler-engine/parser/`: Generated Go lexer and parser bindings from `grammars-v4/java/java`.

#### CppTranspilerVisitor Struct & Lifecycle
The core AST visitor struct tracks code emission buffers and inter-class dependency includes:

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

// NewCppTranspilerVisitor creates a new visitor instance initialized with empty buffers
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

// GetHeaderCode returns emitted header (.h) contents
func (v *CppTranspilerVisitor) GetHeaderCode() string { return v.headerCode.String() }

// GetCppCode returns emitted implementation (.cpp) contents
func (v *CppTranspilerVisitor) GetCppCode() string { return v.cppCode.String() }
```

#### Core AST Traversal Logic Example
```go
// VisitClassDeclaration extracts identifiers, base class, and interfaces
func (v *CppTranspilerVisitor) VisitClassDeclaration(ctx *parser.ClassDeclarationContext) interface{} {
    className := ctx.IDENTIFIER().GetText()
    v.currentClass = className
    superClass := "java::lang::Object"
    if ctx.EXTENDS() != nil && ctx.TypeType() != nil {
        superClass = mapType(ctx.TypeType().GetText())
    }
    v.headerCode.WriteString(fmt.Sprintf("class %s : public %s {\npublic:\n", className, superClass))
    v.Visit(ctx.ClassBody())
    v.headerCode.WriteString("};\n")
    return nil
}

// VisitMethodDeclaration lowers instance virtual and static member methods
func (v *CppTranspilerVisitor) VisitMethodDeclaration(ctx *parser.MethodDeclarationContext) interface{} {
    methodName := ctx.IDENTIFIER().GetText()
    returnType := mapType(ctx.TypeTypeOrVoid().GetText())
    params := v.extractParameters(ctx.FormalParameters())

    // Virtual dispatch for instance methods in headers
    qualifier := "virtual "
    if v.currentIsStatic {
        qualifier = "static "
    }
    v.headerCode.WriteString(fmt.Sprintf("    %s%s %s(%s);\n", qualifier, returnType, methodName, params))

    // Definition emission into .cpp
    v.cppCode.WriteString(fmt.Sprintf("%s %s::%s(%s) {\n", returnType, v.currentClass, methodName, params))
    v.Visit(ctx.MethodBody())
    v.cppCode.WriteString("}\n\n")
    return nil
}
```

#### 1.4 Visitor vs. Listener Design Pattern
ANTLR4 provides two primary tree traversal paradigms:
- **Parse Tree Listeners (`JavaParserListener`)**: Passive, event-driven callbacks triggered as a walker enters and exits tree nodes. Listeners do not control traversal flow and cannot easily synthesize values up the tree.
- **Parse Tree Visitors (`CppTranspilerVisitor`)**: Explicit, active reduction where `visitor.Visit(node)` returns intermediate string expressions and AST values. This enables localized code generation (e.g. evaluating an `ExpressionList` or constructor arguments before emitting the enclosing method call).

#### 1.5 ANTLR4 Parsing & Code Generation Pipeline
1. **Lexical Analysis (`JavaLexer`)**:
   Stream of raw Java code is tokenized into typed tokens (`IDENTIFIER`, `INT`, `RETURN`, `IF`, etc.).
   - Literals: `DECIMAL_LITERAL`, `FLOAT_LITERAL`, `BOOL_LITERAL`, `STRING_LITERAL`, `TEXT_BLOCK`.
   - Keywords: `CLASS`, `INTERFACE`, `EXTENDS`, `IMPLEMENTS`, `SYNCHRONIZED`, `DEFAULT`, `STATIC`.
2. **Syntactic Rule Hierarchy (`JavaParser`)**:
   Tokens form a typed parse tree starting from root rule `compilationUnit`:
   - `compilationUnit` ➔ `packageDeclaration?` `importDeclaration*` `typeDeclaration*` `EOF`
   - `typeDeclaration` ➔ `classDeclaration` | `interfaceDeclaration` | `recordDeclaration`
   - `classDeclaration` ➔ `CLASS IDENTIFIER (EXTENDS typeType)? (IMPLEMENTS typeList)? classBody`
   - `interfaceDeclaration` ➔ `INTERFACE IDENTIFIER (EXTENDS typeList)? interfaceBody`
   - `classBody` ➔ `classBodyDeclaration*` ➔ `memberDeclaration`
   - `memberDeclaration` ➔ `methodDeclaration` | `constructorDeclaration` | `fieldDeclaration` | `genericMethodDeclaration`
   - `statement` ➔ `block`, `statementExpression`, `ifThenElse`, `whileStatement`, `forStatement`, `synchronizedStatement`, `switchStatement`, `tryStatement`
   - `expression` ➔ `primary`, `expression '.' IDENTIFIER`, `methodCall`, `lambdaExpression`, `methodReference`, `ternaryExpression`
3. **AST Visitor Traversal (`CppTranspilerVisitor`)**:
   `CppTranspilerVisitor` embeds `parser.BaseJavaParserVisitor` to implement selective override hooks:
   - `VisitCompilationUnit`: Scans package/import metadata and initiates class lowering.
   - `VisitClassDeclaration`: Extracts class names, superclasses (`extends`), implemented interfaces (`implements`), emitting modern C++ class declarations.
   - `VisitInterfaceDeclaration`: Emits pure virtual abstract C++ base classes or inlined default/static methods.
   - `VisitConstructorDeclaration`: Translates constructors and constructor chaining (`this(...)`, `super(...)`) into C++ member initializer lists.
   - `VisitFieldDeclaration`: Detects static vs instance modifiers and applies Java default initializers (`0`, `false`, `0.0`, `nullptr`).
   - `VisitMethodDeclaration`: Emits virtual and static method signatures into `.h` and definitions into `.cpp`.
   - `VisitStatement`: Lowers `if-else`, `while`, standard `for`, enhanced `for-each`, `switch`, synchronized blocks, and `try-catch-finally` with RAII guards.
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

#### Automated Parser Regeneration & Go Parser Bindings
Parser Go bindings can be re-generated at any time using `scripts/generate_parser.sh`:
```bash
./scripts/generate_parser.sh
```
This runs `antlr-4.13.1-complete.jar` targeting Go with visitor generation enabled (`-Dlanguage=Go -visitor -package parser -o compiler-engine/parser`) from grammar definitions in `grammars-v4/java/java/` and automatically normalizes receiver predicates (`sed -i 's/\bthis\./p\./g' compiler-engine/parser/*.go`).

Generated bindings in `compiler-engine/parser/`:
- `java_lexer.go`: Lexer scanning characters into tokens.
- `java_parser.go`: Parser building AST trees from tokens.
- `java_parser_base_visitor.go`: Base no-op visitor embedded by `CppTranspilerVisitor`.
- `java_parser_visitor.go`: Visitor interface with typed methods for every Java production.

#### 1.6 CLI Entrypoint & Flag Specifications
The CLI binary (`compiler-engine/main.go`) handles arguments via Go's `flag` standard library:
- `-src`: Root directory containing input `.java` source files to recursively discover and transpile.
- `-out`: Target directory where `.h` and `.cpp` files will be generated.
- `-dashboard`: Boolean flag to launch the embedded Web Dashboard Studio.
- `-port`: HTTP port to bind the dashboard server (default: `8080`).

Build the CLI binary:
```bash
cd compiler-engine && go build -o j2cpp-engine . && cd ..
```

Execute batch transpilation across any directory of `.java` source files:
```bash
./compiler-engine/j2cpp-engine -src=/path/to/src -out=/path/to/out
```
In batch mode, the CLI:
1. Validates that `-src` and `-out` are non-empty directory paths.
2. Ensures the target output directory exists via `os.MkdirAll(*outDir, os.ModePerm)`.
3. Traverses the source directory recursively with `filepath.WalkDir`.
4. Parses and transpiles every discovered `.java` file into `<ClassName>.h` and `<ClassName>.cpp`.

Or run the interactive Web Dashboard Studio:
```bash
./compiler-engine/j2cpp-engine -dashboard -port=8080
```

#### 1.7 Initial Architecture of the Go Transpiler CLI & ANTLR4 Integration
The Go transpiler CLI operates as an end-to-end multi-pass compiler driver connecting ANTLR4 grammar analysis with modern C++ emission:

1. **ANTLR4 Go Runtime Binding**:
   - `github.com/antlr4-go/antlr/v4` provides the underlying runtime engine for lexing, token streaming, and LL(*) predictive parsing.
   - The Go lexer (`parser.NewJavaLexer(input)`) converts raw Unicode codepoints into typed tokens against `grammars-v4/java/java/JavaLexer.g4`.
2. **Token Buffering & Channel Isolation**:
   - `antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)` buffers lexed tokens while cleanly filtering hidden tokens (whitespace, line comments, block docstrings) away from the parser channel.
3. **Parser Initialization & Rule Invocation**:
   - `parser.NewJavaParser(tokens)` instantiates the syntactic rule engine.
   - Calling `p.CompilationUnit()` parses the stream starting at the top-level compilation unit rule (`compilationUnit`), building an AST parse tree in memory.
4. **Syntax Error Handling & Listener Attachments**:
   - Custom diagnostic listeners implementing `antlr.ErrorListener` (or `antlr.NewConsoleErrorListener()`) intercept syntax errors, lexical anomalies, and mismatched tokens, providing exact line/column diagnostics.
5. **AST Visitor Traversal & Code Synthesis**:
   - `CppTranspilerVisitor` walks the AST nodes (`VisitClassDeclaration`, `VisitMethodDeclaration`, etc.).
   - Header definitions (`.h`) and implementation bodies (`.cpp`) are emitted in tandem into in-memory string buffers and written to disk or returned to the interactive dashboard.

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

### Prerequisites & Initial Environment Setup
Before building and running `libjruntime` or the Go transpiler CLI, ensure your host environment satisfies the following dependencies:

| Toolchain / Dependency | Minimum Version | Required For | Verification Command |
|:---|:---|:---|:---|
| **Go Toolchain** | 1.20+ | Compiling the `j2cpp-engine` CLI & dashboard | `go version` |
| **C++ Compiler** | GCC 10+ / Clang 11+ | Compiling C++20 transpiled code & runtime | `g++ --version` or `clang++ --version` |
| **Autotools Suite** | Autoconf 2.69+, Automake 1.15+, Libtool 2.4+ | Configuring & building `bdwgc` shared library | `autoconf --version && libtool --version` |
| **GNU Make** | 4.0+ | Building submodules & orchestrating compilation | `make --version` |
| **Java Runtime (JRE/JDK)** | 11+ | Optional ANTLR4 parser regeneration | `java -version` |
| **Git** | 2.25+ | Submodule cloning and version control management | `git --version` |

#### Installing Dependencies

##### Ubuntu / Debian
```bash
sudo apt-get update
sudo apt-get install -y build-essential autoconf automake libtool make \
    openjdk-17-jre-headless golang git
```

##### Fedora / RHEL
```bash
sudo dnf groupinstall "Development Tools"
sudo dnf install autoconf automake libtool golang java-17-openjdk-headless git
```

##### macOS (Homebrew)
```bash
brew install go autoconf automake libtool make openjdk@17 git
```

### Step 1: Initialize & Build Submodules
The project relies on two core Git submodules configured in `.gitmodules`:
- `grammars-v4`: Official ANTLR4 Java lexer and parser grammars (`https://github.com/robertpelloni/grammars-v4`).
- `bdwgc`: Boehm-Demers-Weiser Conservative Garbage Collector (`https://github.com/robertpelloni/bdwgc`).

#### Submodule Incorporation & Architecture Overview
1. **Verify Submodule Configuration**: Inspect `.gitmodules` to confirm valid tracking paths and repository URLs.
2. **Recursive Checkout**: Fetch and checkout all submodule references:
   ```bash
   git submodule update --init --recursive
   ```
3. **Verify Commit Status**: Check initialized commit hashes (`git submodule status`):
   - `bdwgc`: `1d38c6d905a83b006df843440b380e51a2119871`
   - `grammars-v4`: `5a9a89fa7274878bdb93877c4be5d56859c08e5e`
4. **Compile Boehm GC Shared Library**:
   Run the automated preparation script:
   ```bash
   ./scripts/build_submodules.sh
   ```
   Or execute manual Autotools configuration:
   ```bash
   cd bdwgc
   ./autogen.sh
   ./configure --enable-threads=posix --enable-shared --disable-static
   make -j$(nproc)
   cd ..
   ```
5. **ANTLR4 Code Generation**:
   Generate Go parser bindings from `grammars-v4/java/java/`:
   ```bash
   ./scripts/generate_parser.sh
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
