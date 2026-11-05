#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMP_DIR="$(mktemp -d /tmp/j2cpp_e2e_XXXXXX)"
trap 'rm -rf "${TEMP_DIR}"' EXIT

echo "=== Building Boehm GC if not already built ==="
if [ ! -f "${REPO_ROOT}/bdwgc/.libs/libgc.so" ]; then
    (cd "${REPO_ROOT}/bdwgc" && ./autogen.sh && ./configure && make -j"$(nproc)")
fi

echo "=== Building compiler-engine CLI ==="
(cd "${REPO_ROOT}/compiler-engine" && go build -o "${TEMP_DIR}/j2cpp-engine" .)

echo "=== Preparing test Java source ==="
mkdir -p "${TEMP_DIR}/src"

cat << 'EOF' > "${TEMP_DIR}/src/MathHelper.java"
public class MathHelper {
    public static int FACTOR = 7;

    public static int multiplyByFactor(int n) {
        return n * FACTOR;
    }
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/Describable.java"
public interface Describable {
    String describe();
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/Animal.java"
public class Animal implements Describable {
    int legs;

    public Animal() {
        this.legs = 4;
    }

    public int getLegCount() {
        return this.legs;
    }

    public String speak() {
        return "Generic sound";
    }

    public String describe() {
        return "An animal with legs";
    }
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/Dog.java"
public class Dog extends Animal {
    String breed;

    public Dog(String breed) {
        this.legs = 4;
        this.breed = breed;
    }

    public String speak() {
        return "Woof";
    }

    public String describe() {
        return "A loyal dog";
    }
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/SampleApp.java"
public class SampleApp {
    int counter;
    int baseValue;

    public SampleApp() {
        this.baseValue = 100;
    }

    public SampleApp(int base) {
        this.baseValue = base;
    }

    public int getBaseValue() {
        return this.baseValue;
    }

    public int compute(int a, int b) {
        return a + b;
    }

    public int calculateSum(int limit) {
        int total = 0;
        for (int i = 0; i < limit; i++) {
            total = total + i;
        }
        return total;
    }

    public int checkBranch(int val) {
        if (val > 10) {
            return 1;
        } else {
            return 0;
        }
    }

    public int[] createArray(int size) {
        int[] arr = new int[size];
        for (int i = 0; i < size; i++) {
            arr[i] = i * 10;
        }
        return arr;
    }

    public void testException() {
        try {
            int val = 42;
        } catch (Exception e) {
            int err = 1;
        } finally {
            int done = 2;
        }
    }
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/Main.java"
public class Main {
    public void start() {
        SampleApp app = new SampleApp(250);
        app.compute(10, 20);
        app.testException();
        System.out.println("Main pipeline running!");
    }
}
EOF

echo "=== Transpiling Java to C++ ==="
mkdir -p "${TEMP_DIR}/out"
"${TEMP_DIR}/j2cpp-engine" -src="${TEMP_DIR}/src" -out="${TEMP_DIR}/out"

echo "=== Creating C++ Test Entrypoint ==="
cat << 'EOF' > "${TEMP_DIR}/out/main_entry.cpp"
#include <iostream>
#include <cassert>
#include "MathHelper.h"
#include "Describable.h"
#include "Animal.h"
#include "Dog.h"
#include "SampleApp.h"
#include "Main.h"
#include "System.h"
#include "Array.h"

int main() {
    GC_INIT();

    // Verify static fields and static method invocation
    assert(MathHelper::FACTOR == 7);
    int32_t scaled = MathHelper::multiplyByFactor(6);
    assert(scaled == 42);

    // Verify inheritance, polymorphism, and virtual method override
    Dog* dog = java_new<Dog>(java_new<java::lang::String>("Golden Retriever"));
    assert(dog != nullptr);
    assert(dog->getLegCount() == 4);

    // Dynamic dispatch / polymorphic base class pointer
    Animal* animal = dog;
    java::lang::String* sound = animal->speak();
    assert(sound != nullptr);
    assert(sound->getStdString() == "Woof");

    // Dynamic dispatch / interface pointer
    Describable* desc = dog;
    java::lang::String* description = desc->describe();
    assert(description != nullptr);
    assert(description->getStdString() == "A loyal dog");

    // Verify default and parameterized constructor
    SampleApp* defaultApp = java_new<SampleApp>();
    assert(defaultApp != nullptr);
    assert(defaultApp->getBaseValue() == 100);

    SampleApp* customApp = java_new<SampleApp>(350);
    assert(customApp != nullptr);
    assert(customApp->getBaseValue() == 350);

    int32_t sum = defaultApp->compute(15, 27);
    assert(sum == 42);

    int32_t loopSum = defaultApp->calculateSum(5);
    assert(loopSum == 10);

    int32_t branch1 = defaultApp->checkBranch(20);
    assert(branch1 == 1);

    int32_t branch2 = defaultApp->checkBranch(5);
    assert(branch2 == 0);

    // Verify array creation and indexing
    java::lang::Array<int32_t>* numbers = defaultApp->createArray(4);
    assert(numbers != nullptr);
    assert(numbers->length == 4);
    assert((*numbers)[0] == 0);
    assert((*numbers)[1] == 10);
    assert((*numbers)[2] == 20);
    assert((*numbers)[3] == 30);

    // Verify array bounds checking
    bool boundsCaught = false;
    try {
        int32_t bad = (*numbers)[10];
        (void)bad;
    } catch (const java::lang::ArrayIndexOutOfBoundsException& e) {
        boundsCaught = true;
    }
    assert(boundsCaught);

    defaultApp->testException();

    Main* m = java_new<Main>();
    m->start();

    // Verify System.out.println
    java::lang::System::out->println("Testing System.out.println string");
    java::lang::System::out->println(42);
    java::lang::System::out->println(true);

    // Verify String representation and Object monitor
    java::lang::String* str = defaultApp->toString();
    assert(str != nullptr);
    assert(!str->getStdString().empty());

    std::mutex* mon = defaultApp->getMonitor();
    assert(mon != nullptr);

    // Verify JAVA_NULL_CHECK throws NullPointerException
    bool caught = false;
    SampleApp* nullApp = nullptr;
    try {
        JAVA_NULL_CHECK(nullApp)->compute(1, 2);
    } catch (const java::lang::NullPointerException& e) {
        caught = true;
    }
    assert(caught);

    std::cout << "E2E Transpilation and Native Execution Verification PASSED" << std::endl;
    return 0;
}
EOF

echo "=== Compiling native executable with g++ ==="
g++ -std=c++20 \
    "${TEMP_DIR}/out"/*.cpp \
    -I"${REPO_ROOT}/libjruntime" \
    -I"${REPO_ROOT}/bdwgc/include" \
    -L"${REPO_ROOT}/bdwgc/.libs" \
    -lgc \
    -Wl,-rpath,"${REPO_ROOT}/bdwgc/.libs" \
    -o "${TEMP_DIR}/native_app"

echo "=== Running native executable ==="
"${TEMP_DIR}/native_app"

echo "=== All E2E Integration Tests Passed Successfully ==="
