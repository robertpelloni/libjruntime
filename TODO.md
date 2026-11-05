# TODO

## Immediate Next Steps
- [ ] Implement `java.lang.StringBuilder` in `libjruntime` for complex string concatenations.
- [ ] Add constructor chaining (`this(...)`, `super(...)`) support in `VisitConstructorDeclaration`.
- [ ] Implement static initialization blocks (`static { ... }`) mapping to C++ static anonymous execution helpers.
- [ ] Add continuous integration GitHub Actions workflow running `go test` and `tests/e2e_test.sh`.

## Maintenance & Optimization
- [x] Eliminate compiler warnings on `JAVA_FINALLY` guard variable using `[[maybe_unused]]`.
- [x] Fix duplicate method body emission and duplicate logging in `visitor.go`.
- [x] Automatically track and emit inter-class `#include` statements in generated C++ files.
- [x] Implement native console output via `libjruntime/System.h` and transpiler mapping for `System.out.println`.
- [x] Implement control flow statement visitor (`if-else`, `while`, `for`).
- [x] Implement class constructor translation in `VisitConstructorDeclaration`.
- [x] Implement class inheritance (`extends SuperClass`) in `VisitClassDeclaration`.
- [x] Implement interface declaration and `implements` clause multiple inheritance support.
- [x] Implement array wrapper `java::lang::Array<T>` and array transpilation syntax.
- [x] Implement static fields and static method declarations and member access without null checks.
