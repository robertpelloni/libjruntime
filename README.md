# libjruntime

A Go-based Transpiler CLI and a modern C++20 Native Runtime Library designed to compile and run unmodified Java source code natively in C++ using conservative garbage collection (Boehm GC).

## Architecture & Overview

Goal: Run unmodified Java source code natively by transpiling it to C++20.

1. **Go Transpiler Engine (`compiler-engine/`)**
   - Uses ANTLR4 with Java grammar (`grammars-v4`).
   - Walks AST with `CppTranspilerVisitor`.
   - Generates `.h` and `.cpp` files for each Java class.
   - Converts primitive types (`int` -> `int32_t`, `boolean` -> `bool`, `long` -> `int64_t`).
   - Generates `virtual` methods for dynamic dispatch.
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

### 1. Build Boehm GC Submodule
```bash
cd bdwgc
./autogen.sh
./configure
make -j$(nproc)
cd ..
```

### 2. Build the Go Transpiler CLI
```bash
cd compiler-engine
go build -o j2cpp-engine .
cd ..
```

### 3. Transpile Java Sources
```bash
./compiler-engine/j2cpp-engine -src=/path/to/java/source -out=./output
```

### 4. Compile Native App
```bash
g++ -std=c++20 output/*.cpp -I./libjruntime -I./bdwgc/include -L./bdwgc/.libs -lgc -Wl,-rpath,./bdwgc/.libs -o native_app
```
