# ROADMAP

## Milestone 1: Core Foundation & Toolchain (v1.0.x) - [Completed]
- [x] Initial submodules and runtime headers setup (`grammars-v4`, `bdwgc`).
- [x] Automated submodule build and parser generation toolchain scripts (`scripts/build_submodules.sh`, `scripts/generate_parser.sh`).
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
- [x] Mathematical utilities: `java.lang.Math` constants (`PI`, `E`) and functions (`abs`, `min`, `max`, `sqrt`, `pow`, etc.).
- [x] Primitive wrapper classes (`Integer`, `Long`, `Boolean`) and boxed methods (`libjruntime/Boxes.h`).
- [x] Dynamic collections runtime (`java.util.ArrayList<T>` in `libjruntime/ArrayList.h`).
- [x] Enhanced for-each loops (`for (T x : collection)`) with range-based loop lowering.

- [x] Explicit type casting (`(Type)expr`) and runtime `instanceof` expressions.
- [x] `switch` statements (`switch (expr) { case ...: default: ... }`), `break`, `continue`, and unary expressions.
- [x] Ternary expressions (`cond ? a : b`).

- [x] Constructor chaining (`this(...)`, `super(...)`) and C++ member initializer lowering.
- [x] Explicit `throw` statements and parameterized exception constructor overloads.

## Milestone 2: Language Completeness (v1.1.0)
- [ ] Field default initializers in generated constructors.

- [x] Key-value associative container: `java.util.HashMap<K, V>` in `libjruntime/HashMap.h`.

- [x] File I/O primitives (`File`, `FileInputStream`, `FileOutputStream`) in `libjruntime/FileIO.h`.

- [x] Extended wrapper classes (`Double`, `Float`) in `libjruntime/Boxes.h`.

- [x] `synchronized` block and method syntax translation backed by `Object` atomic monitor mutexes.

- [x] Concurrency primitives (`Thread`, `Runnable`) in `libjruntime/Thread.h`.

## Milestone 3: Runtime Library Expansion (v1.2.0)
- [x] Functional interfaces and lambda expressions (`(args) -> expr / { ... }`).
- [x] Method references (`Class::method`, `expr::method`, `Class::new`).
- [x] Nested / inner classes and anonymous class support.
- [x] Functional interfaces runtime headers (`libjruntime/Functions.h` - `Function`, `Consumer`, `Supplier`, `Predicate`).
- [x] Stream API runtime headers (`libjruntime/Stream.h` - `java.util.stream.Stream<T>`).
- [x] Default and static interface methods transpilation.
- [x] Reflection basics (`java.lang.Class<T>` runtime header and class literal transpilation).
- [x] Container abstractions (`java.util.Optional<T>` runtime header and transpiler mappings).
- [x] Set collections (`java.util.HashSet<E>` runtime header and transpiler mappings).
- [x] Continuous integration automated testing workflow (`.github/workflows/ci.yml`).

## Milestone 4: Multi-threading & Advanced Tooling (v2.0.0)
- [x] Real-time transpilation dashboard and interactive developer CLI tools (`compiler-engine -dashboard`).
