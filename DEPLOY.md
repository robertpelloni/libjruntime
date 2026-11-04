# DEPLOYMENT & ENVIRONMENT SETUP MANUAL

## Requirements
- **Operating System**: Linux (Ubuntu 20.04+, Debian 11+, or compatible POSIX environment)
- **C++ Compiler**: GCC (g++ 10+) or Clang (clang++ 11+) supporting `-std=c++20`
- **Go Toolchain**: Go 1.20+
- **Build Tools**: `make`, `autoconf`, `automake`, `libtool`

## Step-by-Step Build & Setup

### 1. Initialize Git Submodules
```bash
git submodule update --init --recursive
```

### 2. Compile Boehm Garbage Collector (`bdwgc`)
```bash
cd bdwgc
./autogen.sh
./configure
make -j$(nproc)
cd ..
```
*Note: This creates shared libraries in `bdwgc/.libs/libgc.so`.*

### 3. Build the Go Transpiler Engine (`j2cpp-engine`)
```bash
cd compiler-engine
go build -o j2cpp-engine .
cd ..
```

### 4. Running the Transpilation Pipeline
Transpile all Java source files from a source directory into a C++ output directory:
```bash
./compiler-engine/j2cpp-engine -src=/path/to/java/sources -out=./output
```

### 5. Compiling Native C++ Applications
Compile the transpiled C++ files together with `libjruntime` and Boehm GC:
```bash
g++ -std=c++20 output/*.cpp \
    -I./libjruntime \
    -I./bdwgc/include \
    -L./bdwgc/.libs \
    -lgc \
    -Wl,-rpath,./bdwgc/.libs \
    -o native_app
```

### 6. Executing Tests
- Run Go unit tests:
  ```bash
  cd compiler-engine && go test -v ./... && cd ..
  ```
- Run End-to-End integration test suite:
  ```bash
  ./tests/e2e_test.sh
  ```
