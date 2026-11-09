# Submodule Integration & Incorporation Guide: bdwgc & grammars-v4

This document outlines the architecture, configuration, git submodule incorporation, integration mechanics, and future evolution of the two core Git submodules used by `libjruntime`:
1. **Boehm-Demers-Weiser Conservative Garbage Collector (`bdwgc`)**
2. **Official ANTLR4 Grammars Repository (`grammars-v4`)**

---

## 0. Git Submodule Incorporation & Management Workflow

Both submodules are tracked in `.gitmodules`:
```ini
[submodule "grammars-v4"]
	path = grammars-v4
	url = https://github.com/robertpelloni/grammars-v4
[submodule "bdwgc"]
	path = bdwgc
	url = https://github.com/robertpelloni/bdwgc
```

### Initial Cloning & Recursive Initialization
When cloning `libjruntime` fresh from source control:
```bash
# Clone repository and recursively initialize all submodules
git clone --recursive https://github.com/robertpelloni/libjruntime.git
```

If the repository has already been cloned without submodules:
```bash
# Initialize and fetch submodules recursively
git submodule update --init --recursive
```

### Keeping Submodules Updated
To pull the latest commits for submodules tracked from remote forks:
```bash
# Fetch and merge updates from the remote tracking branches
git submodule update --remote --merge
```

### Adding New Submodules
To incorporate additional libraries as git submodules:
```bash
git submodule add <repository-url> <path>
git commit -m "chore: add <submodule-name> submodule"
```

---

## 1. Boehm Garbage Collector (`bdwgc`)

### Architecture & Role
The Boehm-Demers-Weiser Garbage Collector provides automatic conservative memory management for all transpiled Java objects running natively in C++. Java requires non-manual memory deallocation with automatic cyclic garbage collection; `bdwgc` fulfills this requirement without needing explicit reference counting cycles or complex compiler instrumentation.

### Directory Layout
- `bdwgc/include/`: Public headers, primarily `gc.h` and platform-specific definitions.
- `bdwgc/include/gc_pthread_redirects.h`: Interception headers for pthread calls.
- `bdwgc/.libs/`: Output directory containing compiled shared libraries (`libgc.so`).
- `bdwgc/tests/`: Diagnostic tests and multi-threading verification routines.

### Build Configuration & Flags
The automated setup script (`scripts/build_submodules.sh`) compiles `bdwgc` using GNU Autotools:
```bash
cd bdwgc
./autogen.sh
./configure --enable-threads=posix --enable-shared --disable-static
make -j$(nproc)
```

**Key Configuration Options:**
- `--enable-threads=posix`: Essential for multi-threaded Java applications (`java.lang.Thread`). Configures GC stop-the-world signal handlers (`SIGPWR`/`SIGXCPU` on Linux) to safely suspend and inspect threads.
- `--enable-shared`: Compiles `libgc.so` for dynamic linking and small executable footprint.
- `--disable-static`: Prevents redundant static archive generation, accelerating builds.

### Runtime Integration in `libjruntime`
- **Object Allocation (`libjruntime/GC.h`)**:
  All Java object allocations are wrapped via `java_new<T>(...)`, which invokes `GC_MALLOC(sizeof(T))` and executes C++ in-place placement new `new (ptr) T(std::forward<Args>(args)...);`.
- **Multi-threading Stack Registration (`libjruntime/Thread.h`)**:
  When a thread spawns via `java.lang.Thread`, the runtime calls:
  - `GC_allow_register_threads()`: Permitting dynamically created threads to attach to the collector.
  - `GC_register_my_thread(&sb)`: Registering the thread's stack bottom and context with the GC before executing user code.
  - `GC_unregister_my_thread()`: Detaching the thread upon termination.

### Future Integration Possibilities
1. **Incremental Collection (`GC_enable_incremental()`)**: Minimize GC pause times for interactive applications and game engines.
2. **Custom Finalizers (`GC_register_finalizer_no_order`)**: Emulate Java `Object.finalize()` or modern `java.lang.ref.Cleaner` mechanisms.
3. **Generational Garbage Collection**: Enable GC write-barrier support for high-throughput batch workloads.

---

## 2. ANTLR4 Grammars (`grammars-v4`)

### Architecture & Role
`grammars-v4` is the official repository of community-maintained ANTLR4 grammar files. `libjruntime` integrates the Java grammar located in `grammars-v4/java/java`.

### Grammar Specifications
- **`JavaLexer.g4`**:
  - Defines lexical modes, whitespace rules, identifiers, and token definitions.
  - Recognizes keywords (`class`, `interface`, `extends`, `implements`, `synchronized`, `default`, `static`).
  - Recognizes literals (hex, binary, float, double, string, character, and text blocks).
- **`JavaParser.g4`**:
  - Context-free production rules for compilation units, types, statements, and expressions.
  - Explicit production rules for modern Java features: `lambdaExpression`, `methodReference`, and default interface methods.

### Go Parser Generation Pipeline
ANTLR4 parser Go bindings are generated via `scripts/generate_parser.sh`:
```bash
java -jar /usr/local/lib/antlr-4.13.1-complete.jar \
    -Dlanguage=Go \
    -visitor \
    -package parser \
    -o compiler-engine/parser \
    grammars-v4/java/java/JavaLexer.g4 \
    grammars-v4/java/java/JavaParser.g4
```

**Normalization Step:**
ANTLR4's Go target generates parser receiver variables like `p` while some predicate helper methods reference `this`. `scripts/generate_parser.sh` automatically normalizes `this` predicates to `p` using `sed` to ensure clean Go compilation:
```bash
sed -i 's/\bthis\./p\./g' compiler-engine/parser/*.go
```

### Initial Transpiler CLI Stream Setup
The Go CLI imports the generated `parser` package and initializes the ANTLR stream:
```go
import (
    "github.com/antlr4-go/antlr/v4"
    "github.com/user/j2cpp-engine/parser"
)

input := antlr.NewInputStream(javaSource)
lexer := parser.NewJavaLexer(input)
tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
p := parser.NewJavaParser(tokens)
p.BuildParseTrees = true
tree := p.CompilationUnit()
```

### Runtime Integration in `compiler-engine`
- `compiler-engine/parser/`: Houses the generated lexer and parser.
- `compiler-engine/visitor.go`: Implements `CppTranspilerVisitor`, subclassing `parser.BaseJavaParserVisitor`.
  - Walks the ANTLR4 parse tree starting from `VisitCompilationUnit`.
  - Lowers AST nodes into clean, idiomatic C++20 headers (`.h`) and implementation files (`.cpp`).

### Future Integration Possibilities
1. **Java 17+ Grammar Features**: Pattern matching for `switch`, record classes (`record Point(int x, int y)`), and sealed classes (`sealed interface Shape permits Circle, Square`).
2. **Incremental AST Caching**: Cache parse tree hashes to enable sub-millisecond re-transpilation of large Java codebases during live dashboard editing.
3. **Syntax Tree Visualizer**: Wire the dashboard to render AST parse tree graphs alongside the C++ output.
