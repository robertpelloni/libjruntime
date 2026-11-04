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
cat << 'EOF' > "${TEMP_DIR}/src/SampleApp.java"
public class SampleApp {
    int counter;

    public int compute(int a, int b) {
        return a + b;
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
        SampleApp app = new SampleApp();
        app.compute(10, 20);
        app.testException();
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
#include "SampleApp.h"
#include "Main.h"

int main() {
    GC_INIT();

    SampleApp* app = java_new<SampleApp>();
    assert(app != nullptr);

    int32_t sum = app->compute(15, 27);
    assert(sum == 42);

    app->testException();

    Main* m = java_new<Main>();
    m->start();

    // Verify String representation and Object monitor
    java::lang::String* str = app->toString();
    assert(str != nullptr);
    assert(!str->getStdString().empty());

    std::mutex* mon = app->getMonitor();
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
