# Changelog

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
