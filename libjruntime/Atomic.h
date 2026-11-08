#pragma once

#include <atomic>
#include <cstdint>
#include "Object.h"
#include "GC.h"

namespace java {
namespace util {
namespace concurrent {
namespace atomic {

class AtomicInteger : public ::java::lang::Object {
private:
    std::atomic<int32_t> val;

public:
    AtomicInteger() : val(0) {}
    AtomicInteger(int32_t initialValue) : val(initialValue) {}

    int32_t get() const noexcept {
        return val.load();
    }

    void set(int32_t newValue) noexcept {
        val.store(newValue);
    }

    int32_t getAndSet(int32_t newValue) noexcept {
        return val.exchange(newValue);
    }

    bool compareAndSet(int32_t expect, int32_t update) noexcept {
        return val.compare_exchange_strong(expect, update);
    }

    int32_t getAndIncrement() noexcept {
        return val.fetch_add(1);
    }

    int32_t getAndDecrement() noexcept {
        return val.fetch_sub(1);
    }

    int32_t getAndAdd(int32_t delta) noexcept {
        return val.fetch_add(delta);
    }

    int32_t incrementAndGet() noexcept {
        return val.fetch_add(1) + 1;
    }

    int32_t decrementAndGet() noexcept {
        return val.fetch_sub(1) - 1;
    }

    int32_t addAndGet(int32_t delta) noexcept {
        return val.fetch_add(delta) + delta;
    }
};

class AtomicBoolean : public ::java::lang::Object {
private:
    std::atomic<bool> val;

public:
    AtomicBoolean() : val(false) {}
    AtomicBoolean(bool initialValue) : val(initialValue) {}

    bool get() const noexcept {
        return val.load();
    }

    void set(bool newValue) noexcept {
        val.store(newValue);
    }

    bool getAndSet(bool newValue) noexcept {
        return val.exchange(newValue);
    }

    bool compareAndSet(bool expect, bool update) noexcept {
        return val.compare_exchange_strong(expect, update);
    }
};

} // namespace atomic
} // namespace concurrent
} // namespace util
} // namespace java
