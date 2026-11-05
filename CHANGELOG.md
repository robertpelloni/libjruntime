# Changelog

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
