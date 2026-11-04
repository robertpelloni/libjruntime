# Changelog

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
