# libjruntime

**This repository contains no code.** Its only file is the document that
used to be this README: a set of step-by-step prompt blueprints for having
an LLM build a Java-to-C++ transpiler. Neither the transpiler nor the
runtime library was ever committed (one session, two commits, 2026-06-21).

## The design in the document

Goal: run *unmodified* Java source natively by transpiling it to C++.

1. **Go transpiler CLI** — ANTLR4 with the Java grammar (from the
   `grammars-v4` repo, to be added as a submodule) generating Go parser
   bindings; a `CppTranspilerVisitor` walks the AST and emits one `.h` /
   `.cpp` pair per Java class, mapping primitives to fixed-width C types
   (`int` → `int32_t`, `long` → `int64_t`, …), making instance methods
   `virtual`, and rewriting `new Foo()` as `java_new<Foo>()`.
2. **libjruntime (C++ runtime)** — header-only base types in namespace
   `java::lang`: an `Object` with `hashCode()`/`equals()`/`toString()` and
   a lazily-initialized mutex backing `synchronized`; a `String` wrapping
   `std::string`; allocation routed through the Boehm GC (`bdwgc`, also to
   be added as a submodule) via `GC_MALLOC`.
3. **Safety injection** — the emitter wraps member access in a
   `JAVA_NULL_CHECK` macro that throws a native `NullPointerException`,
   and synthesizes Java `finally` semantics with an RAII cleanup guard,
   since C++ has no `finally`.
4. **Build plan** — `./j2cpp-engine -src=<java project> -out=./output`,
   then `clang++ -std=c++20 ./output/*.cpp ./libjruntime/runtime.cpp -lgc`.

## Status

Blueprint only. The two submodule links named in the plan
(`grammars-v4`, `bdwgc`) exist as separate mirror repos under this
account, but they were never added here, and no component of the pipeline
exists. Sibling experiment: `jvm-cpp-runtime`, which does contain a small
implemented Java runtime — for the opposite translation direction.
