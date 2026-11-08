#pragma once

#include <mutex>
#include <atomic>
#include "Object.h"
#include "GC.h"

namespace java {
namespace util {
namespace concurrent {
namespace locks {

class ReentrantLock : public ::java::lang::Object {
private:
    std::recursive_mutex mtx;
    std::atomic<bool> locked{false};

public:
    ReentrantLock() = default;

    virtual void lock() {
        mtx.lock();
        locked.store(true);
    }

    virtual void unlock() {
        locked.store(false);
        mtx.unlock();
    }

    virtual bool tryLock() {
        if (mtx.try_lock()) {
            locked.store(true);
            return true;
        }
        return false;
    }

    virtual bool isLocked() const {
        return locked.load();
    }
};

} // namespace locks
} // namespace concurrent
} // namespace util
} // namespace java
