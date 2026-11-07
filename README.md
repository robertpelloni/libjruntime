# libjruntime

A Go-based Transpiler CLI and a modern C++20 Native Runtime Library designed to compile and run unmodified Java source code natively in C++ using conservative garbage collection (Boehm GC).

## Architecture & Overview

Goal: Run unmodified Java source code natively by transpiling it to C++20.

### ANTLR4 Grammar Integration & Transpiler CLI Setup

The Go-based transpiler front-end relies on the official Java ANTLR4 grammar definitions located in the `grammars-v4` git submodule (`grammars-v4/java/java`):
- `JavaLexer.g4`: Lexical specifications defining tokens, operators, keywords, and literals.
- `JavaParser.g4`: Production rules defining Java compilation units, declarations, statements, and expressions.

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

Execute transpilation across any directory of `.java` source files:
```bash
./compiler-engine/j2cpp-engine -src=/path/to/src -out=/path/to/out
```
- `-src`: Directory containing input `.java` files (default: `./src`).
- `-out`: Directory where modern C++20 header (`.h`) and implementation (`.cpp`) files will be emitted (default: `./out`).

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
