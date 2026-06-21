# libjruntime
Go-based Transpiler CLI and a tiny C++ Native Runtime Library to run unmodified Java source code natively as/in C++

Add https://github.com/robertpelloni/grammars-v4 and https://github.com/robertpelloni/bdwgc as submodules.

Step-by-Step Prompt Blueprints for the LLM
You will split the architecture into a series of highly granular prompts. Feed these to your LLM sequentially, verifying execution at each stage.

Prompt 1: Setting up the Go ANTLR4 Front-End Engine
Plaintext
Context: We are building a custom source-to-source compiler from unmodified Java to native modern C++ using Go as the engine.
Task: Write a Go application setup script and entry point.
Requirements:
1. Target the ANTLR4 Java parser grammar files (JavaLexer.g4 and JavaParser.g4).
2. Generate Go-based bindings using the ANTLR tool chain: `antlr4 -Dlanguage=Go -visitor -o parser JavaLexer.g4 JavaParser.g4`.
3. Provide a standard main.go file that reads a specified target directory full of standard Java source files, invokes the lexer/parser, and steps into a custom AST structural tree walker struct named `CppTranspilerVisitor`.
4. Output placeholder console logs indicating when classes, methods, and field declarations are discovered during AST traversal.
Prompt 2: Crafting the libjruntime Core C++ Header Files
Plaintext
Context: This is the native target runtime for our Java-to-C++ compiler. We need to implement the core object structures.
Task: Write standard C++20 header files representing the base Java framework types.
Requirements:
1. Define a namespace named `java::lang`.
2. Implement an `Object` class. It must have virtual methods for `hashCode()`, `equals(Object* obj)`, and `toString()`.
3. The `Object` class must hold a raw pointer to a lightweight mutex `std::mutex*` that is lazily initialized only when requested, to back Java's native `synchronized` keyword footprint.
4. Define a `String` class inheriting from `Object` that wraps a `std::string_view` or `std::string` internally and overloads `operator+` for Java-style concatenation.
5. Create a macro or inline template function named `java_new<T>(...)` that maps allocation internal mechanics directly to the Boehm GC library allocation engine: `GC_MALLOC(sizeof(T))` or `new (GC_MALLOC(sizeof(T))) T(...)`.
Prompt 3: Writing the Go AST-to-C++ Code Emitter
Plaintext
Context: We have a working ANTLR4 Java AST visitor setup in Go. Now we need to emit actual C++ code strings.
Task: Implement the Visitor methods in Go to map Java syntax to structural C++.
Requirements:
1. Class declarations: Map `public class Foo extends Bar` into separate C++ files: `Foo.h` and `Foo.cpp`. Add appropriate header inclusion guard statements.
2. Ensure fields and variable instantiations map from Java primitives to deterministic C types (e.g., `int` -> `int32_t`, `boolean` -> `bool`, `long` -> `int64_t`).
3. Methods: Translate typical method declarations. Because Java uses dynamic dispatch by default, prefix all instance member functions in the generated C++ class declarations with the `virtual` keyword.
4. Allocations: Intercept allocations like `new Foo(args)` and rewrite them explicitly as `java_new<Foo>(args)` to wire up our conservative garbage collector.
Prompt 4: Handling Scoped Controls (Null Checks & Exceptions)
Plaintext
Context: To preserve unmodified Java syntax guarantees, we need to handle runtime safety properties in C++.
Task: Enhance the Go code emitter logic to inject safety layers into statements.
Requirements:
1. Member access interception: When parsing an expression like `obj.method()`, evaluate the type. Inject a macro check inline: `JAVA_NULL_CHECK(obj)->method()`. Define this macro in C++ to throw a native `java::lang::NullPointerException` if the pointer value evaluates to null.
2. Exceptions: Map Java `try {} catch (Exception e) {} finally {}` blocks cleanly. Because modern C++ does not natively contain a `finally` construct, synthesize a clean solution by wrapping the body execution using an RAII structural cleanup guard pattern (e.g., executing a local lambda expression in a custom destructor wrapper when leaving scope) to guarantee completion.
Step 3: Deployment & Build System Execution Plan
Once your LLM generates the components, tie them together into a unified build tool.

The Compilation Directory Structure:

/compiler-engine (Go CLI)
/libjruntime (C++ Source + Boehm GC headers)
/output (Target output directory for generated .h/.cpp)
Transpilation Phase: Run your Go CLI tool against the source directory.

Bash
./j2cpp-engine -src=/path/to/java/project -out=./output
Native Compilation Phase: Compile the whole thing using your host toolchain (Clang or GCC), linking against the local runtime code and the conservative GC package:

Bash
clang++ -std=c++20 ./output/*.cpp ./libjruntime/runtime.cpp -lgc -o native_app
