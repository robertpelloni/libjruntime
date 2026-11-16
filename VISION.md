# VISION & ARCHITECTURAL FOUNDATION

## Ultimate Goal
`libjruntime` enables running unmodified Java source code natively as compiled modern C++ (C++20), achieving bare-metal performance, predictable execution latencies, and minimal memory footprints while completely bypassing standard JVM bytecode interpretation or heavy JIT runtimes.

## Core Architectural Pillars

### 1. Direct Source-to-Source Transpilation (`compiler-engine`)
- Parses pristine Java language syntax using ANTLR4 Go bindings and official grammar specifications (`grammars-v4`).
- Converts classes to deterministic `.h` (interface/layout) and `.cpp` (method implementation) source pairs.
- Preserves object semantics: dynamic dispatch mapped to C++ `virtual` methods, primitive types mapped to fixed-width standards (`int32_t`, `int64_t`, `bool`).
- Automatically resolves inter-class dependencies and generates header inclusion graphs.

### 2. High-Performance Conservative Garbage Collection
- Transparent memory management using Boehm-Demers-Weiser conservative garbage collector (`bdwgc`).
- `java_new<T>(...)` translates allocations directly to `GC_MALLOC`, providing seamless automatic memory reclamation without reference-cycle memory leaks or reference counting overhead.

### 3. Native Java Semantics Emulation (`libjruntime`)
- **`java::lang::Object`**: Root of class hierarchy with support for polymorphic comparisons (`equals`), identity/hash generation (`hashCode`), string formatting (`toString`), and lazily initialized atomic monitors backing `synchronized`.
- **`java::lang::String`**: First-class immutable string object interoperable with C++ string facilities.
- **Safety Injection**: Member access automatically guarded via `JAVA_NULL_CHECK(obj)` throwing `java::lang::NullPointerException`.
- **Deterministic Exception Handling & Scope Cleanup**: Exception mapping with RAII cleanup wrappers (`JAVA_FINALLY`) providing exact `try-catch-finally` Java semantics in native C++.

### 4. Zero-Friction Developer Experience
- Single-command transpilation via `j2cpp-engine`.
- Standard C++ compiler toolchains (GCC, Clang) link generated code with zero runtime dependencies beyond Boehm GC.
