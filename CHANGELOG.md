# Changelog

## 1.0.40
- Implemented `libjruntime/LinkedList.h` providing `java::util::LinkedList<T>`:
  - Backed by doubly linked nodes allocated with `GC_MALLOC`.
  - Sequential, queue, and deque methods: `add`, `addFirst`, `addLast`, `get`, `getFirst`, `getLast`, `remove`, `removeFirst`, `removeLast`, `size`, `isEmpty`, `clear`, `contains`.
  - Forward range-based iterator support (`begin()`, `end()`).
- Updated `compiler-engine/visitor.go` type mappings and header includes for `LinkedList<T>`.
- Added unit test `TestTranspileLinkedList` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration test coverage for `LinkedList` in `tests/e2e_test.sh`.
- Bumped version to `1.0.40` in `VERSION`.

## 1.0.39
- Implemented `libjruntime/Collections.h` providing standard `java::util::Collections` algorithms:
  - `sort(ArrayList<T>*)` using `std::sort`.
  - `sort(ArrayList<T>*, Compare)` with custom comparator lambda support.
  - `reverse(ArrayList<T>*)` using `std::reverse`.
  - `shuffle(ArrayList<T>*)` using pseudo-random swap shuffle.
  - `min(ArrayList<T>*)` and `max(ArrayList<T>*)` elements with empty check safety.
  - `fill(ArrayList<T>*, T)` element population.
  - `frequency(ArrayList<T>*, T)` count calculation.
- Integrated `Collections` into `compiler-engine/visitor.go` type mappings and header includes.
- Added unit test `TestTranspileCollections` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration test coverage for `Collections` algorithms in `tests/e2e_test.sh`.
- Bumped version to `1.0.39` in `VERSION`.

## 1.0.38
- Implemented embedded interactive web dashboard studio in `compiler-engine`:
  - Added `-dashboard` and `-port` flags in `compiler-engine/main.go`.
  - Added single-page responsive Web UI with live side-by-side Java-to-C++ playground editor in `compiler-engine/dashboard.go`.
  - Implemented `/api/transpile` and `/api/status` HTTP endpoints for on-the-fly code generation and runtime health inspection.
  - Added unit tests in `compiler-engine/dashboard_test.go` (`TestDashboardStatus`, `TestDashboardTranspile`).
- Bumped version to `1.0.38` in `VERSION`.

## 1.0.37
- Created `libjruntime/HashSet.h` implementing `java::util::HashSet<E>`:
  - Backed by an internal `java::util::HashMap<E, Object*>` for O(1) membership operations.
  - Dynamically maintained flat array for `begin()` and `end()` range-based iteration support.
  - Methods: `add`, `contains`, `remove`, `size`, `isEmpty`, `clear`.
- Updated `compiler-engine/visitor.go` type mappings and header auto-inclusion for `HashSet<T>` and `Set<T>`.
- Added unit test `TestTranspileHashSet` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration and native execution test coverage for `HashSet<T>` in `tests/e2e_test.sh`.
- Bumped version to `1.0.37` in `VERSION`.

## 1.0.36
- Created `libjruntime/Optional.h` implementing `java::util::Optional<T>`:
  - Static factory methods `of`, `ofNullable`, `empty`.
  - Value accessors and checks `isPresent`, `isEmpty`, `get`, `orElse`, `ifPresent`.
  - Transformation and filtering methods `filter`, `map`.
- Updated `compiler-engine/visitor.go` type mappings and header auto-inclusion for `Optional<T>`.
- Added unit test `TestTranspileOptional` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration and native execution test coverage for `Optional<T>` in `tests/e2e_test.sh`.
- Bumped version to `1.0.36` in `VERSION`.

## 1.0.35
- Created `libjruntime/Class.h` implementing reflection basics (`java::lang::Class<T>`):
  - `getName()`, `getSimpleName()`, `isInstance(Object*)`, and static factory `Class<T>::forType(...)`.
  - Added `getClass()` virtual member method on `java::lang::Object`.
- Updated `compiler-engine/visitor.go` type mappings and header auto-inclusion for `Class<T>`.
- Added support in `VisitPrimary` for Java class literal expressions (`TypeName.class`).
- Added unit test `TestTranspileClassLiteral` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration and native execution test coverage for class literals and reflection in `tests/e2e_test.sh`.
- Bumped version to `1.0.35` in `VERSION`.

## 1.0.34
- Implemented default interface methods (`default ReturnType method(...) { ... }`) and static interface methods (`static ReturnType method(...) { ... }`) in `compiler-engine/visitor.go`.
- Emitted inlined implementations in interface headers while preserving pure virtual declarations for abstract interface methods.
- Added unit test `TestTranspileDefaultInterfaceMethod` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration and native execution test coverage for default and static interface methods in `tests/e2e_test.sh`.
- Bumped version to `1.0.34` in `VERSION`.

## 1.0.33
- Created `libjruntime/Stream.h` implementing `java::util::stream::Stream<T>` with intermediate and terminal operations:
  - `filter`, `map`, `limit`, `skip`, `count`, `forEach`, and `toList`.
- Added `stream()` member function to `java::util::ArrayList<T>`.
- Updated `compiler-engine/visitor.go` type mappings and header inclusions for `Stream`.
- Added unit test `TestTranspileStreamAPI` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration and native execution test coverage for Stream pipelines in `tests/e2e_test.sh`.
- Bumped version to `1.0.33` in `VERSION`.

## 1.0.32
- Created `libjruntime/Functions.h` implementing core functional interfaces under `java::util::function`:
  - `Function<T, R>` with `R apply(T t)`
  - `Consumer<T>` with `void accept(T t)`
  - `Supplier<T>` with `T get()`
  - `Predicate<T>` with `bool test(T t)`
  - Adapter closures and factory helpers (`make_function`, `make_consumer`, `make_supplier`, `make_predicate`).
- Updated `compiler-engine/visitor.go` type mappings to route `Function`, `Consumer`, `Supplier`, and `Predicate` to their corresponding `java::util::function::*` templates and automatically include `Functions.h`.
- Added unit test `TestTranspileFunctionalInterfaces` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration and native execution test coverage for functional interfaces in `tests/e2e_test.sh`.
- Bumped version to `1.0.32` in `VERSION`.

## 1.0.30
- Implemented method reference syntax transpilation (`Type::method`, `expr::method`, `Type::new`) in `compiler-engine/visitor.go`.
- Lowered method references into universal forwarding C++ lambda closures (`[=](auto&&... args) { ... }`).
- Added unit test `TestTranspileMethodReference` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration and native execution test coverage for method references in `tests/e2e_test.sh`.
- Bumped version to `1.0.30` in `VERSION`.

## 1.0.29
- Added support for transpiling Java lambda expressions into C++20 lambda expressions (`[=](...) { ... }`).
- Supported single and multi-parameter lambdas with explicit types and inferred (`auto`) types.
- Supported expression-bodied and block-bodied lambdas.
- Added unit test `TestTranspileLambdaExpression` in `compiler-engine/visitor_test.go`.
- Added end-to-end integration test coverage for lambda expressions in `tests/e2e_test.sh`.
- Bumped version to `1.0.29` in `VERSION`.

## 1.0.28
- Added continuous integration workflow `.github/workflows/ci.yml` running Go unit tests and end-to-end native execution tests across pull requests and pushes.
- Bumped version to `1.0.28` in `VERSION`.

## 1.0.27
- Enhanced `README.md` with comprehensive step-by-step documentation for building `j2cpp-engine`, submodule preparation scripts, CLI flags, native compilation, and end-to-end testing.
- Bumped version to `1.0.27` in `VERSION`.

## 1.0.26
- Created `libjruntime/Thread.h` implementing `java::lang::Runnable` and `java::lang::Thread` with multi-threading support.
- Registered worker threads with Boehm GC using `GC_allow_register_threads()`, `GC_register_my_thread()`, and `GC_unregister_my_thread()` for safe GC tracking across threads.
- Implemented `start()`, `join()`, `run()`, `isAlive()`, and `Thread::sleep(millis)` methods.
- Updated `compiler-engine/visitor.go` to map `Thread` and `Runnable` types, auto-include `Thread.h`, and support `implements Runnable`.
- Added unit test `TestTranspileThreadAndRunnable` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime asynchronous worker thread execution with Boehm GC.

## 1.0.25
- Implemented `synchronized` statement block transpilation (`synchronized (expr) { ... }`) lowering to `std::lock_guard<std::mutex> _java_sync_lock(*JAVA_NULL_CHECK(expr)->getMonitor());` using `Object` atomic monitor mutexes.
- Implemented `synchronized` method modifier detection in `compiler-engine/visitor.go` injecting `this->getMonitor()` lock guards at the entry of synchronized methods.
- Added unit test `TestTranspileSynchronized` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime synchronized block monitor acquisition and safe execution.

## 1.0.24
- Implemented primitive wrapper classes `Double` and `Float` in `libjruntime/Boxes.h` with `valueOf`, unboxing (`doubleValue()`, `floatValue()`, `intValue()`, `longValue()`), `hashCode()`, `equals()`, and `toString()`.
- Updated `compiler-engine/visitor.go` with type mappings for `double`, `float`, `char`, `short`, `Double`, and `Float`, object creation, and static `Double.valueOf(...)` / `Float.valueOf(...)` call routing.
- Added unit test `TestTranspileDoubleAndFloat` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime `Double` and `Float` boxing, unboxing, casting, and hashing.

## 1.0.23
- Created `libjruntime/FileIO.h` implementing `java::io::File`, `java::io::FileInputStream`, and `java::io::FileOutputStream` backed by standard C/POSIX file streams and Boehm GC.
- Updated `compiler-engine/visitor.go` to support `byte` (`int8_t`), `byte[]` (`Array<int8_t>*`), and File I/O types (`File`, `FileInputStream`, `FileOutputStream`) with automatic inclusion of `FileIO.h`.
- Added unit test `TestTranspileFileIO` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime file creation, byte writing, byte reading, and file deletion.

## 1.0.22
- Implemented `throw` statement transpilation (`throw expr;`) in `compiler-engine/visitor.go`.
- Added constructor overloads in `libjruntime/Exceptions.h` accepting `java::lang::String*` message arguments across `Exception`, `NullPointerException`, and `ArrayIndexOutOfBoundsException`.
- Added exception type mappings in `VisitObjectCreationExpression`.
- Added unit test `TestTranspileThrowStatement` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime custom exception throwing and catching.

## 1.0.21
- Implemented constructor chaining (`this(...)`, `super(...)`) in `compiler-engine/visitor.go` (`VisitConstructorDeclaration`).
- Detected leading delegating and super constructor calls in constructor bodies and translated them to C++ member initializer lists (`: Base(args)` / `: Sibling(args)`).
- Prevented duplicate body execution of chaining constructor calls.
- Added unit test `TestTranspileConstructorChaining` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime constructor delegation and base constructor execution.

## 1.0.20
- Created `libjruntime/HashMap.h` implementing `java::util::HashMap<K, V>` key-value associative dictionary with dynamic bucket array resizing and Boehm GC allocation.
- Implemented `put(k, v)`, `get(k)`, `containsKey(k)`, `remove(k)`, `size()`, `isEmpty()`, and `clear()`, supporting both primitive keys and `Object`-derived reference keys through `equals(Object*)` and `hashCode()`.
- Updated `compiler-engine/visitor.go` with type mappings for `HashMap` and parameterized `HashMap<K, V>`, object allocation (`java_new<java::util::HashMap<...>>`), and `#include "HashMap.h"` header generation.
- Added unit test `TestTranspileHashMap` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime native `HashMap` entry insertion, key lookup, containment, and removal.

## 1.0.19
- Implemented ternary operator expression transpilation (`cond ? thenBranch : elseBranch`) in `compiler-engine/visitor.go` (`VisitTernaryExpression`).
- Added parenthesized expression preservation in `VisitPrimary` (`(expr)`).
- Added unit test `TestTranspileTernaryExpression` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime ternary operator evaluations with Boehm GC.

## 1.0.18
- Implemented `switch` statement transpilation (`switch (expr) { case val: ... default: ... }`) with `case`, `default`, `break`, and `continue` statement support in `compiler-engine/visitor.go`.
- Implemented `VisitUnaryOperatorExpression` in `compiler-engine/visitor.go` to support unary prefix operators (e.g., negative literals `-1`, `!`, `~`, `+`).
- Added unit test `TestTranspileSwitchStatement` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime `switch` statement evaluations and branch selection.

## 1.0.17
- Expanded `README.md` with a comprehensive section detailing ANTLR4 grammar integration with `grammars-v4` and Go transpiler CLI setup/usage instructions.
- Sourced version string `1.0.17` in `VERSION`.

## 1.0.16
- Created `scripts/build_submodules.sh` automating submodule recursive initialization and native Boehm GC (`bdwgc`) compilation.
- Created `scripts/generate_parser.sh` automating ANTLR4 Go target generation from `grammars-v4` Java grammar sources (`JavaLexer.g4`, `JavaParser.g4`) into `compiler-engine/parser/`.
- Updated `DEPLOY.md` and `README.md` to document automated scripts for developers.

## 1.0.15
- Implemented `instanceof` expression transpilation (`expr instanceof TargetType`) into `(dynamic_cast<TargetType*>(expr) != nullptr)`.
- Implemented explicit type casting expressions (`(TargetType) expr`) into `dynamic_cast<TargetType*>(expr)` for reference types and `static_cast<PrimitiveType>(expr)` for primitive types.
- Added Java `null` literal mapping to C++ `nullptr` in `VisitLiteral`.
- Added delayed header generation in `VisitClassDeclaration` to ensure all referenced and casted types are included in `<ClassName>.h`.
- Added unit test `TestTranspileCastAndInstanceOf` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime `instanceof` checks and polymorphic downcasting with Boehm GC.

## 1.0.14
- Implemented enhanced for-each loop transpilation (`for (Type var : iterable)`) in `compiler-engine/visitor.go` mapping to C++ range-based loops (`for (Type var : *JAVA_NULL_CHECK(iterable))`).
- Added standard iterator accessors `begin()`, `end()`, `cbegin()`, and `cend()` to `libjruntime/Array.h` and `libjruntime/ArrayList.h` for seamless C++ standard iteration.
- Added unit test `TestTranspileEnhancedForLoop` to `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime native array and `ArrayList` element iteration and summation.

## 1.0.13
- Created `libjruntime/ArrayList.h` implementing `java::util::ArrayList<T>` dynamic collection with Boehm GC-backed resizing.
- Implemented `add(element)`, `add(index, element)`, `get(index)`, `set(index, element)`, `remove(index)`, `size()`, `isEmpty()`, `clear()`, `contains(element)`, and bounds checking throwing `ArrayIndexOutOfBoundsException`.
- Updated `compiler-engine/visitor.go` to support `ArrayList` and parameterized `ArrayList<T>` type mappings, header auto-inclusion (`ArrayList.h`), and instantiation (`java_new<java::util::ArrayList<...>>`).
- Added unit test `TestTranspileArrayList` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime `ArrayList` creation, element addition, indexed retrieval, removal, and bounds exception handling in native C++20.

## 1.0.12
- Created `libjruntime/Boxes.h` introducing primitive boxed classes `java::lang::Integer`, `java::lang::Boolean`, and `java::lang::Long`.
- Implemented static factory methods `Integer::valueOf(int32_t)`, `Boolean::valueOf(bool)`, `Long::valueOf(int64_t)`.
- Implemented primitive extraction methods `intValue()`, `booleanValue()`, `longValue()`, and numeric conversion `doubleValue()`.
- Implemented `hashCode()`, `equals(Object*)`, and `toString()` on all boxed types.
- Enhanced `compiler-engine/visitor.go` to recognize primitive wrapper types, map type references to `java::lang::Integer*`, `java::lang::Boolean*`, `java::lang::Long*`, auto-include `#include "Boxes.h"`, and route `Integer.valueOf(...)`, `Boolean.valueOf(...)`, `Long.valueOf(...)`.
- Added unit test `TestTranspilePrimitiveWrappers` in `compiler-engine/visitor_test.go` and verified with Go test suite.
- Updated `tests/e2e_test.sh` with `BoxHelper` and assertions validating native allocation, boxing, unboxing, equals, and string conversion with Boehm GC.

## 1.0.11
- Created `libjruntime/Math.h` implementing `java::lang::Math` with mathematical constants (`PI`, `E`) and utility functions (`abs`, `min`, `max`, `sqrt`, `pow`, `floor`, `ceil`, `round`).
- Enhanced `compiler-engine/visitor.go` to route `Math` member references and method invocations directly to `java::lang::Math::member`.
- Added unit test `TestTranspileMathOperations` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime `java.lang.Math` operations in native C++20.
- Synchronized workspace documentation and bumped version to `1.0.11`.

## 1.0.10
- Implemented static initialization blocks (`static { ... }`) in `compiler-engine/visitor.go` transpiled to self-executing static lambdas run at module initialization time.
- Added tracking and automatic qualification of static class fields inside static blocks and class methods (`ClassName::field`).
- Added unit test `TestTranspileStaticInitBlock` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting static initialization block execution prior to application entry.
- Synchronized workspace documentation and bumped version to `1.0.10`.

## 1.0.9
- Created `libjruntime/StringBuilder.h` implementing `java::lang::StringBuilder` with chained `append` overloads for strings, characters, numbers, booleans, and objects, and `toString()` backed by Boehm GC.
- Updated `compiler-engine/visitor.go` with `StringBuilder` type mapping, object creation translation, and header inclusion.
- Added unit test `TestTranspileStringBuilder` to `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime `StringBuilder` append chaining and string creation.
- Synchronized workspace documentation and bumped version to `1.0.9`.

## 1.0.8
- Implemented static method transpilation (`static returnType methodName(...)`) in `compiler-engine/visitor.go`.
- Implemented static field transpilation (`inline static fieldType fieldName = init;`) in C++ headers.
- Enhanced static member invocation and field access in `VisitMethodCallExpression` and `VisitMemberReferenceExpression` (`ClassName::member`) bypassing instance null checks.
- Added unit test `TestTranspileStaticMembers` to `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime static field reads and static method calculations.
- Synchronized workspace documentation and bumped version to `1.0.8`.

## 1.0.7
- Created `libjruntime/Array.h` defining `java::lang::Array<T>` with bounds checking, length property, and Boehm GC allocation.
- Added `ArrayIndexOutOfBoundsException` in `libjruntime/Exceptions.h`.
- Supported array type mapping (`T[]` -> `java::lang::Array<T>*`) and array creation (`new T[n]`) in `compiler-engine/visitor.go`.
- Implemented square bracket indexing (`VisitSquareBracketExpression`) for array access and mutation.
- Updated `JAVA_FINALLY` in `libjruntime/FinallyGuard.h` to accept variadic macro arguments with line-unique guard identifiers.
- Added unit test `TestTranspileArrayOperations` in `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test suite `tests/e2e_test.sh` asserting runtime array creation, assignment, length, and bounds checking exceptions.
- Synchronized workspace documentation and bumped version to `1.0.7`.

## 1.0.6
- Implemented Java interface declaration transpilation (`VisitInterfaceDeclaration`) emitting C++ abstract base classes with virtual destructors.
- Implemented pure virtual method signature generation (`VisitInterfaceMethodDeclaration`) for interface methods.
- Supported `implements` clause in class declarations with multiple inheritance and automatic interface header inclusion.
- Added unit test `TestTranspileInterfaceAndImplements` to `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test harness `tests/e2e_test.sh` asserting runtime polymorphic interface dispatch.
- Synchronized workspace documentation and bumped version to `1.0.6`.

## 1.0.5
- Implemented class inheritance (`extends SuperClass`) in `compiler-engine/visitor.go`.
- Added automatic header inclusion of parent classes in generated `<ClassName>.h`.
- Wrapped Java string literals in `java_new<java::lang::String>(...)` via `VisitLiteral` for strict C++ type safety when passing or returning strings.
- Refactored `JAVA_NULL_CHECK` into an inline C++ template function `check_not_null<T>` in `libjruntime/Runtime.h` to eliminate duplicate expression evaluation.
- Added unit test `TestTranspileInheritance` in `compiler-engine/visitor_test.go`.
- Updated end-to-end integration test runner `tests/e2e_test.sh` asserting runtime inheritance, polymorphism, and virtual method overrides.
- Synchronized workspace documentation and bumped version to `1.0.5`.

## 1.0.4
- Implemented constructor transpilation (`VisitConstructorDeclaration`) in `compiler-engine/visitor.go`.
- Added support for overloaded and parameterized constructors mapped to C++ constructors.
- Added unit test `TestTranspileConstructor` to `compiler-engine/visitor_test.go`.
- Updated end-to-end integration test harness `tests/e2e_test.sh` to verify default and parameterized constructor runtime initialization with `java_new<T>(...)`.
- Synchronized workspace documentation and bumped version to `1.0.4`.

## 1.0.3
- Created `libjruntime/System.h` providing `java::lang::PrintStream` and static `System::out` supporting printing and line printing for strings, primitives, and objects.
- Enhanced `compiler-engine/visitor.go` to handle `System.out.println` / `System.out.print` without injecting redundant null checks on static `System`.
- Implemented control flow transpilation in `VisitStatement`: `if-else` branching, `while` loops, and standard `for` loops.
- Added support for post-increment and post-decrement expressions (`VisitPostIncrementDecrementOperatorExpression`).
- Added unit tests for `System.out.println` and control flow structures to `compiler-engine/visitor_test.go`.
- Expanded end-to-end integration test suite (`tests/e2e_test.sh`) verifying control flow and native standard output execution.
- Updated documentation across `ROADMAP.md`, `TODO.md`, `MEMORY.md`, `HANDOFF.md`, and `VERSION`.

## 1.0.2
- Fixed duplicate method body statement emission and duplicate logging in `compiler-engine/visitor.go`.
- Implemented parameter list extraction and binary operator expression (`VisitBinaryOperatorExpression`) AST parsing.
- Added automatic dependent header tracking and `#include` generation for referenced classes.
- Updated `libjruntime/FinallyGuard.h` with `[[maybe_unused]]` attribute to eliminate unused-variable compiler warnings.
- Included `<cstdint>` in `libjruntime/Object.h` for `std::intptr_t`.
- Updated `libjruntime/String.h` and `libjruntime/Object.h` to use `java_new<String>` for consistent Boehm GC memory tracking.
- Added automated unit test suite (`compiler-engine/visitor_test.go`).
- Added end-to-end integration and native execution test suite (`tests/e2e_test.sh`).
- Updated core project governance documentation (`VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `ROADMAP.md`, `TODO.md`, `HANDOFF.md`).

## 1.0.1
- Integrated Go ANTLR4 transpiler engine (`j2cpp-engine`).
- Integrated C++ runtime framework (`libjruntime`) mapping `java::lang::Object` and Boehm GC.
- Synchronized `bdwgc` and `grammars-v4` submodules.

## 1.0.0
- Initial Project Setup.
