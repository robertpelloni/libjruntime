#pragma once

#include <mutex>
#include <string>
#include <atomic>

namespace java {
namespace lang {

class String;

class Object {
private:
    std::atomic<std::mutex*> monitor{nullptr};

public:
    virtual ~Object() {
        std::mutex* m = monitor.load();
        if (m != nullptr) {
            delete m;
        }
    }

    virtual int hashCode() {
        return static_cast<int>(reinterpret_cast<std::intptr_t>(this));
    }

    virtual bool equals(Object* obj) {
        return this == obj;
    }

    virtual String* toString();

    std::mutex* getMonitor() {
        std::mutex* m = monitor.load(std::memory_order_acquire);
        if (m == nullptr) {
            std::mutex* new_m = new std::mutex();
            std::mutex* expected = nullptr;
            if (monitor.compare_exchange_strong(expected, new_m, std::memory_order_acq_rel)) {
                m = new_m;
            } else {
                delete new_m;
                m = expected;
            }
        }
        return m;
    }
};

} // namespace lang
} // namespace java
