# INTERNAL CODEBASE MEMORY & ARCHITECTURAL OBSERVATIONS

## Repository Structure & Components
- **`compiler-engine/`**: Go module translating Java source to modern C++20.
  - `main.go`: CLI parsing `-src` and `-out` directories and orchestrating ANTLR file traversal.
  - `visitor.go`: `CppTranspilerVisitor` extending ANTLR's `BaseJavaParserVisitor`. Emits `.h` declarations and `.cpp` definitions, manages header inclusion dependencies, and translates Java constructors, methods, fields, statements/expressions (`if`, `while`, `for`, `try-catch-finally`, arithmetic/relational operators).
  - `parser/`: Generated Go parser and lexer bindings from ANTLR4 Java grammar.
- **`libjruntime/`**: Header-only native C++ runtime library in `java::lang` namespace.
  - `Object.h`: Base `Object` containing lazy atomic monitor initialization (`std::mutex*`).
  - `String.h`: Standard string wrapper with concatenation overloads and `toString()`.
  - `System.h`: Standard I/O output stream wrapper (`PrintStream`) backing `java::lang::System::out`.
  - `GC.h`: Memory allocation via `GC_MALLOC` wrapped in `java_new<T>()`.
  - `Exceptions.h`: Native exception types `Exception` and `NullPointerException`.
  - `FinallyGuard.h`: RAII guard wrapper `JAVA_FINALLY` annotated with `[[maybe_unused]]`.
  - `Runtime.h`: Runtime safety macros including `JAVA_NULL_CHECK`.
- **`bdwgc/`**: Git submodule pointing to Boehm GC repository.
- **`grammars-v4/`**: Git submodule containing standard ANTLR4 grammar definitions.
- **`tests/`**: Integration and regression tests.
  - `e2e_test.sh`: End-to-end pipeline transpiling Java sample code, compiling with `g++ -std=c++20`, linking with `libgc`, and verifying assertions, constructors, and console output at runtime.

## Key Design Patterns & Invariants
- Class transpilation produces two files per Java class: `<ClassName>.h` and `<ClassName>.cpp`.
- Constructors map to C++ constructors with matching parameter types, instantiated through `java_new<T>(...)`.
- Class methods are marked `virtual` by default to preserve Java's late binding dispatch.
- When generating `.cpp` files, referenced non-primitive classes are tracked and emitted as `#include "<Class>.h"`.
- All native allocations must route through `java_new<T>` instead of raw `new` to remain visible to `bdwgc`.
