# TODO

## Immediate Next Steps
- [ ] Implement `System.out.println(...)` mapping in visitor to native C++ standard output.
- [ ] Add constructor transpilation support to `CppTranspilerVisitor`.
- [ ] Support `for`, `while`, and `if-else` branching control structures in statement visitor.
- [ ] Add continuous integration GitHub Actions workflow running `go test` and `tests/e2e_test.sh`.

## Maintenance & Optimization
- [x] Eliminate compiler warnings on `JAVA_FINALLY` guard variable using `[[maybe_unused]]`.
- [x] Fix duplicate method body emission and duplicate logging in `visitor.go`.
- [x] Automatically track and emit inter-class `#include` statements in generated C++ files.
