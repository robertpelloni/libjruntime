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

cat << 'EOF' > "${TEMP_DIR}/src/Config.java"
public class Config {
    public static int DEFAULT_PORT;
    public static String APP_NAME;

    static {
        DEFAULT_PORT = 8080;
        APP_NAME = "LibJRuntimeApp";
    }
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/ListHelper.java"
public class ListHelper {
    public static ArrayList createList() {
        ArrayList list = new ArrayList();
        list.add("First");
        list.add("Second");
        list.add("Third");
        return list;
    }

    public static int sumArray(int[] numbers) {
        int total = 0;
        for (int n : numbers) {
            total = total + n;
        }
        return total;
    }

    public static int countListLength(ArrayList list) {
        int count = 0;
        for (Object item : list) {
            count = count + 1;
        }
        return count;
    }

    public static boolean isDog(Object obj) {
        return obj instanceof Dog;
    }

    public static Dog castToDog(Object obj) {
        if (obj instanceof Dog) {
            return (Dog) obj;
        }
        return null;
    }

    public static int switchBranch(int code) {
        int res = 0;
        switch (code) {
            case 10:
                res = 100;
                break;
            case 20:
                res = 200;
                break;
            default:
                res = -1;
                break;
        }
        return res;
    }

    public static int ternarySelect(int val) {
        return (val > 50 ? 999 : 111);
    }

    public static HashMap createMap() {
        HashMap map = new HashMap();
        map.put("one", "1");
        map.put("two", "2");
        map.put("three", "3");
        return map;
    }

    public static void checkPositive(int val) {
        if (val < 0) {
            throw new Exception("Value cannot be negative");
        }
    }

    public static int testFileRoundTrip(String filepath) {
        File file = new File(filepath);
        FileOutputStream out = new FileOutputStream(file);
        out.write(72);
        out.write(105);
        out.close();

        FileInputStream in = new FileInputStream(file);
        int b1 = in.read();
        int b2 = in.read();
        in.close();

        file.deleteFile();
        return b1 + b2;
    }

    public static int syncExecution(Object lock) {
        int count = 0;
        synchronized (lock) {
            count = count + 42;
        }
        return count;
    }
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/MathHelper.java"
public class MathHelper {
    public static int FACTOR = 7;

    public static int multiplyByFactor(int n) {
        return n * FACTOR;
    }

    public static int getClampedValue(int val, int low, int high) {
        return Math.max(low, Math.min(val, high));
    }
}
EOF

cat << 'EOF' > "${TEMP_DIR}/src/BoxHelper.java"
public class BoxHelper {
    public static Integer boxInt(int val) {
        return new Integer(val);
    }

    public static int unboxInt(Integer boxed) {
        return boxed.intValue();
    }

    public static Boolean boxBool(boolean val) {
        return new Boolean(val);
    }

    public static boolean unboxBool(Boolean boxed) {
        return boxed.booleanValue();
    }

    public static Double boxDouble(double val) {
        return new Double(val);
    }

    public static double unboxDouble(Double boxed) {
        return boxed.doubleValue();
    }

    public static Float boxFloat(float val) {
        return new Float(val);
    }

    public static float unboxFloat(Float boxed) {
        return boxed.floatValue();
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

    public Dog() {
        this("Unknown Breed");
    }

    public Dog(String breed) {
        super();
        this.legs = 4;
        this.breed = breed;
    }

    public String getBreed() {
        return this.breed;
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

    public String buildMessage(String prefix, int value) {
        StringBuilder sb = new StringBuilder();
        sb.append(prefix);
        sb.append(" : ");
        sb.append(value);
        return sb.toString();
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
#include "Config.h"
#include "MathHelper.h"
#include "BoxHelper.h"
#include "Describable.h"
#include "Animal.h"
#include "Dog.h"
#include "SampleApp.h"
#include "Main.h"
#include "System.h"
#include "Array.h"
#include "StringBuilder.h"
#include "Math.h"
#include "Boxes.h"
#include "ArrayList.h"
#include "HashMap.h"
#include "FileIO.h"
#include "ListHelper.h"

int main() {
    GC_INIT();

    // Verify static initialization block execution
    assert(Config::DEFAULT_PORT == 8080);
    assert(Config::APP_NAME != nullptr);
    assert(Config::APP_NAME->getStdString() == "LibJRuntimeApp");

    // Verify java.lang.Math constants and operations
    assert(java::lang::Math::abs(-42) == 42);
    assert(java::lang::Math::max(10, 50) == 50);
    assert(java::lang::Math::min(10, 50) == 10);
    assert(java::lang::Math::PI > 3.14 && java::lang::Math::PI < 3.15);

    // Verify MathHelper static method that uses Math.max and Math.min
    int32_t clamped = MathHelper::getClampedValue(150, 0, 100);
    assert(clamped == 100);
    int32_t clampedLow = MathHelper::getClampedValue(-50, 0, 100);
    assert(clampedLow == 0);

    // Verify static fields and static method invocation
    assert(MathHelper::FACTOR == 7);
    int32_t scaled = MathHelper::multiplyByFactor(6);
    assert(scaled == 42);

    // Verify ListHelper and ArrayList functionality
    java::util::ArrayList<java::lang::Object*>* myList = ListHelper::createList();
    assert(myList != nullptr);
    assert(myList->size() == 3);
    assert(!myList->isEmpty());

    java::lang::String* elem0 = dynamic_cast<java::lang::String*>(myList->get(0));
    assert(elem0 != nullptr);
    assert(elem0->getStdString() == "First");

    java::lang::String* elem1 = dynamic_cast<java::lang::String*>(myList->get(1));
    assert(elem1 != nullptr);
    assert(elem1->getStdString() == "Second");

    java::lang::Object* removed = myList->remove(1);
    assert(removed == elem1);
    assert(myList->size() == 2);

    bool listBoundsCaught = false;
    try {
        myList->get(10);
    } catch (const java::lang::ArrayIndexOutOfBoundsException& e) {
        listBoundsCaught = true;
    }
    assert(listBoundsCaught);

    // Verify switch statements and ternary expressions
    assert(ListHelper::switchBranch(10) == 100);
    assert(ListHelper::switchBranch(20) == 200);
    assert(ListHelper::switchBranch(999) == -1);

    assert(ListHelper::ternarySelect(100) == 999);
    assert(ListHelper::ternarySelect(10) == 111);

    // Verify HashMap functionality
    java::util::HashMap<java::lang::Object*, java::lang::Object*>* myMap = ListHelper::createMap();
    assert(myMap != nullptr);
    assert(myMap->size() == 3);
    assert(!myMap->isEmpty());

    java::lang::String* kOne = java_new<java::lang::String>("one");
    assert(myMap->containsKey(kOne));
    java::lang::String* vOne = dynamic_cast<java::lang::String*>(myMap->get(kOne));
    assert(vOne != nullptr);
    assert(vOne->getStdString() == "1");

    java::lang::String* kTwo = java_new<java::lang::String>("two");
    java::lang::String* vTwo = dynamic_cast<java::lang::String*>(myMap->remove(kTwo));
    assert(vTwo != nullptr);
    assert(vTwo->getStdString() == "2");
    assert(myMap->size() == 2);
    assert(!myMap->containsKey(kTwo));

    // Verify throw statement and exception handling
    bool customThrown = false;
    try {
        ListHelper::checkPositive(-42);
    } catch (java::lang::Exception* e) {
        customThrown = true;
        assert(e != nullptr);
    }
    assert(customThrown);

    // Verify synchronized statement block execution with monitor locking
    java::lang::Object* syncLock = java_new<java::lang::Object>();
    int32_t syncResult = ListHelper::syncExecution(syncLock);
    assert(syncResult == 42);

    // Verify File and FileOutputStream / FileInputStream round-trip
    java::lang::String* tempFilePath = java_new<java::lang::String>("/tmp/j2cpp_test_io.txt");
    int32_t byteSum = ListHelper::testFileRoundTrip(tempFilePath);
    assert(byteSum == 72 + 105); // 'H' (72) + 'i' (105) = 177

    // Verify enhanced for-each loops over Array and ArrayList
    int32_t totalListItems = ListHelper::countListLength(myList);
    assert(totalListItems == 2);

    java::lang::Array<int32_t>* testArr = java_new<java::lang::Array<int32_t>>(3);
    (*testArr)[0] = 5;
    (*testArr)[1] = 15;
    (*testArr)[2] = 25;
    int32_t arrSum = ListHelper::sumArray(testArr);
    assert(arrSum == 45);

    // Verify BoxHelper and java.lang.Integer / Boolean / Double / Float boxing & unboxing
    java::lang::Integer* boxedInt = BoxHelper::boxInt(12345);
    assert(boxedInt != nullptr);
    assert(boxedInt->intValue() == 12345);
    assert(BoxHelper::unboxInt(boxedInt) == 12345);
    assert(boxedInt->toString()->getStdString() == "12345");

    java::lang::Boolean* boxedBool = BoxHelper::boxBool(true);
    assert(boxedBool != nullptr);
    assert(boxedBool->booleanValue() == true);
    assert(BoxHelper::unboxBool(boxedBool) == true);
    assert(boxedBool->toString()->getStdString() == "true");

    java::lang::Long* boxedLong = java::lang::Long::valueOf(9876543210LL);
    assert(boxedLong != nullptr);
    assert(boxedLong->longValue() == 9876543210LL);
    assert(boxedLong->intValue() == (int32_t)9876543210LL);

    java::lang::Double* boxedD = BoxHelper::boxDouble(3.14159);
    assert(boxedD != nullptr);
    assert(BoxHelper::unboxDouble(boxedD) > 3.14 && BoxHelper::unboxDouble(boxedD) < 3.15);
    assert(boxedD->intValue() == 3);

    java::lang::Float* boxedF = BoxHelper::boxFloat(2.5f);
    assert(boxedF != nullptr);
    assert(BoxHelper::unboxFloat(boxedF) == 2.5f);
    assert(boxedF->intValue() == 2);

    // Verify inheritance, polymorphism, and virtual method override
    // Verify constructor chaining with super() and this(...)
    Dog* defaultDog = java_new<Dog>();
    assert(defaultDog != nullptr);
    assert(defaultDog->getLegCount() == 4);
    assert(defaultDog->getBreed()->getStdString() == "Unknown Breed");

    Dog* dog = java_new<Dog>(java_new<java::lang::String>("Golden Retriever"));
    assert(dog != nullptr);
    assert(dog->getLegCount() == 4);
    assert(dog->getBreed()->getStdString() == "Golden Retriever");

    // Dynamic dispatch / polymorphic base class pointer
    Animal* animal = dog;
    java::lang::String* sound = animal->speak();
    assert(sound != nullptr);
    assert(sound->getStdString() == "Woof");

    // Verify instanceof and explicit type casts
    assert(ListHelper::isDog(dog) == true);
    assert(ListHelper::isDog(myList) == false);

    Dog* castedDog = ListHelper::castToDog(animal);
    assert(castedDog != nullptr);
    assert(castedDog->getLegCount() == 4);

    Dog* castedFail = ListHelper::castToDog(myList);
    assert(castedFail == nullptr);

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

    // Verify StringBuilder string concatenation
    java::lang::String* built = defaultApp->buildMessage(java_new<java::lang::String>("Score"), 100);
    assert(built != nullptr);
    assert(built->getStdString() == "Score : 100");

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
