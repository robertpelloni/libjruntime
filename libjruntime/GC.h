#pragma once

#include <new>
#include <gc.h>

template<typename T, typename... Args>
inline T* java_new(Args&&... args) {
    void* mem = GC_MALLOC(sizeof(T));
    if (mem == nullptr) {
        throw std::bad_alloc();
    }
    return new(mem) T(std::forward<Args>(args)...);
}
