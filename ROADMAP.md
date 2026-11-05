# ROADMAP

## Milestone 1: Core Foundation & Toolchain (v1.0.x) - [Completed]
- [x] Initial submodules and runtime headers setup (`grammars-v4`, `bdwgc`).
- [x] AST transpiler engine in Go (`j2cpp-engine`) mapping Java syntax to C++20.
- [x] Base native runtime types: `Object`, `String`, `Exceptions`, `FinallyGuard`, and `GC`.
- [x] Automated unit test suite (`visitor_test.go`) and end-to-end integration test harness (`tests/e2e_test.sh`).
- [x] Control flow statement transpilation: `if`, `else`, `while`, and `for` loops.
- [x] Basic console I/O emulation: `System.out.println` and `System.out.print` (`System.h`).
- [x] Class constructor transpilation (default, overloaded, parameterized constructors).

## Milestone 2: Language Completeness (v1.1.0)
- [ ] Support class inheritance (`extends`), interface implementation (`implements`), and method overrides.
- [ ] Static fields, static methods, and class initialization blocks (`static { ... }`).
- [ ] Constructor chaining (`this(...)`, `super(...)`) and field default initializers.
- [ ] Enhanced `for-each` loops and switch / pattern matching statements.

## Milestone 3: Runtime Library Expansion (v1.2.0)
- [ ] Core collections: `ArrayList`, `HashMap`, `Arrays`.
- [ ] File I/O (`FileInputStream`, `FileOutputStream`, `File`).
- [ ] `java.lang.Math`, `java.lang.StringBuilder`, and primitive wrappers (`Integer`, `Boolean`, etc.).

## Milestone 4: Multi-threading & Advanced Tooling (v2.0.0)
- [ ] `synchronized` block and method syntax translation backed by `std::mutex*` monitor.
- [ ] Real-time transpilation dashboard and interactive developer CLI tools.
