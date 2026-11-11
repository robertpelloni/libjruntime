#pragma once

#include <cstdint>
#include <thread>
#include <chrono>
#include <atomic>
#include "Object.h"
#include "Exceptions.h"
#include "GC.h"

namespace java {
namespace lang {

class Runnable {
public:
    virtual ~Runnable() = default;
    virtual void run() = 0;
};

class Thread : public Object, public Runnable {
private:
    Runnable* target;
    std::thread nativeThread;
    std::atomic<bool> alive{false};
    std::atomic<bool> started{false};

    void threadEntryPoint() {
        GC_allow_register_threads();
        struct GC_stack_base sb;
        GC_get_stack_base(&sb);
        GC_register_my_thread(&sb);

        try {
            alive.store(true);
            this->run();
        } catch (...) {
            alive.store(false);
            GC_unregister_my_thread();
            throw;
        }

        alive.store(false);
        GC_unregister_my_thread();
    }

public:
    Thread() : target(nullptr) {}

    Thread(Runnable* r) : target(r) {}

    virtual ~Thread() {
        if (nativeThread.joinable()) {
            nativeThread.join();
        }
    }

    virtual void run() override {
        if (target != nullptr) {
            target->run();
        }
    }

    virtual void start() {
        bool expected = false;
        if (!started.compare_exchange_strong(expected, true)) {
            throw Exception("IllegalThreadStateException: Thread already started");
        }
        GC_allow_register_threads();
        nativeThread = std::thread(&Thread::threadEntryPoint, this);
    }

    virtual void join() {
        if (nativeThread.joinable()) {
            nativeThread.join();
        }
    }

    virtual bool isAlive() const {
        return alive.load();
    }

    static void sleep(int64_t millis) {
        if (millis > 0) {
            std::this_thread::sleep_for(std::chrono::milliseconds(millis));
        }
    }
};

} // namespace lang
} // namespace java
