# HANDOFF & ARCHITECTURAL SUMMARY

## Milestone Reached: v1.0.22 - Explicit `throw` Statement & Parameterized Exceptions

### What Was Accomplished in v1.0.22
1. **Throw Statement Transpilation**:
   - Implemented `ctx.THROW()` in `VisitStatement` (`compiler-engine/visitor.go`) lowering `throw expr;`.
   - Added exception class mappings (`Exception`, `NullPointerException`, `ArrayIndexOutOfBoundsException`) in `VisitObjectCreationExpression`.
2. **Exception Runtime Constructor Support**:
   - Extended `libjruntime/Exceptions.h` with constructor overloads taking `java::lang::String*` message arguments, converting them safely to `std::string message`.
3. **Verification**:
   - Added unit test `TestTranspileThrowStatement` in `compiler-engine/visitor_test.go`.
   - Expanded end-to-end integration test runner `tests/e2e_test.sh` asserting runtime custom exception throwing and catching with Boehm GC.
   - All tests passing cleanly.

### Context for Successor Agents
- Project compiles clean with Go and modern C++20.
- Runtime headers reside in `libjruntime/`.
- Native executable tests link with `bdwgc/.libs/libgc.so` (`-lgc`).
- Scripts for submodules and parser regeneration reside in `scripts/`.
- Next roadmap milestones: concurrency primitives (`Thread`, `Runnable`), file I/O runtime support (`FileInputStream`, `FileOutputStream`).
