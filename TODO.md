# TODO

## Immediate Next Steps
- [ ] Add explicit constructor declaration and instantiation mappings in `CppTranspilerVisitor`.
- [ ] Implement support for class inheritance (`extends`) and virtual method overrides in C++ headers.
- [ ] Implement `java.lang.StringBuilder` in `libjruntime` for complex string concatenations.
- [ ] Add continuous integration GitHub Actions workflow running `go test` and `tests/e2e_test.sh`.

## Maintenance & Optimization
- [x] Eliminate compiler warnings on `JAVA_FINALLY` guard variable using `[[maybe_unused]]`.
- [x] Fix duplicate method body emission and duplicate logging in `visitor.go`.
- [x] Automatically track and emit inter-class `#include` statements in generated C++ files.
- [x] Implement native console output via `libjruntime/System.h` and transpiler mapping for `System.out.println`.
- [x] Implement control flow statement visitor (`if-else`, `while`, `for`).
