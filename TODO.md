# TODO

## Immediate Next Steps
- [ ] Add constructor chaining (`this(...)`, `super(...)`) support in `VisitConstructorDeclaration`.
- [ ] Implement `java.util.HashMap` associative container in `libjruntime`.
- [ ] Support explicit type casting and `instanceof` expressions.
- [ ] Add continuous integration GitHub Actions workflow running `go test` and `tests/e2e_test.sh`.

## Maintenance & Optimization
- [x] Implement enhanced `for-each` loops in statement visitor (`for (Type x : collection)`) with C++ range iterators.
- [x] Implement dynamic collections in `libjruntime/ArrayList.h` (`ArrayList<T>`) with dynamic GC allocation.
- [x] Implement primitive wrapper classes in `libjruntime/Boxes.h` (`Integer`, `Long`, `Boolean`) and visitor mappings.
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
- [x] Implement `java::lang::StringBuilder` with chained appends and `toString()`.
- [x] Implement static initialization blocks (`static { ... }`) mapped to self-executing lambda runners.
- [x] Implement `java::lang::Math` in `libjruntime/Math.h` and visitor mapping for Math methods/constants.
