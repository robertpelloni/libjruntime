# ROADMAP

## Milestone 1: Core Foundation & Toolchain (v1.0.x) - [Completed]
- [x] Initial submodules and runtime headers setup (`grammars-v4`, `bdwgc`).
- [x] AST transpiler engine in Go (`j2cpp-engine`) mapping Java syntax to C++20.
- [x] Base native runtime types: `Object`, `String`, `Exceptions`, `FinallyGuard`, and `GC`.
- [x] Automated unit test suite (`visitor_test.go`) and end-to-end integration test harness (`tests/e2e_test.sh`).
- [x] Control flow statement transpilation: `if`, `else`, `while`, and `for` loops.
- [x] Basic console I/O emulation: `System.out.println` and `System.out.print` (`System.h`).
- [x] Class constructor transpilation (default, overloaded, parameterized constructors).
- [x] Class inheritance (`extends SuperClass`), dynamic dispatch, and polymorphism.
- [x] Interface declarations (`interface Foo`) and implementation (`implements Bar`) with pure virtual C++ bases.
- [x] Array types (`T[]`), creation (`new T[n]`), indexing (`arr[i]`), and bounds safety (`Array.h`).
- [x] Static fields, static methods, and class member access (`ClassName::member`).
- [x] Native string manipulation: `java.lang.StringBuilder` with chained appends.
- [x] Static initialization blocks (`static { ... }`) mapped to self-executing lambda runners.

## Milestone 2: Language Completeness (v1.1.0)
- [ ] Constructor chaining (`this(...)`, `super(...)`) and field default initializers.
- [ ] Enhanced `for-each` loops and switch / pattern matching statements.
- [ ] Explicit type casting and `instanceof` expressions.

## Milestone 3: Runtime Library Expansion (v1.2.0)
- [ ] Core collections: `ArrayList`, `HashMap`.
- [ ] File I/O (`FileInputStream`, `FileOutputStream`, `File`).
- [ ] `java.lang.Math` and primitive wrappers (`Integer`, `Boolean`, etc.).

## Milestone 4: Multi-threading & Advanced Tooling (v2.0.0)
- [ ] `synchronized` block and method syntax translation backed by `std::mutex*` monitor.
- [ ] Real-time transpilation dashboard and interactive developer CLI tools.
