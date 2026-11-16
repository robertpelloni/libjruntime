#pragma once

#include <cstdint>
#include <new>
#include "Object.h"
#include "Exceptions.h"
#include "GC.h"

namespace java {
namespace lang {

template<typename T>
class Array : public Object {
private:
    T* data;
    int32_t len;

public:
    int32_t length;

    Array(int32_t size) : len(size < 0 ? 0 : size), length(size < 0 ? 0 : size) {
        if (size < 0) {
            throw ArrayIndexOutOfBoundsException("Negative array size");
        }
        if (len > 0) {
            void* mem = GC_MALLOC(sizeof(T) * len);
            if (mem == nullptr) {
                throw std::bad_alloc();
            }
            data = static_cast<T*>(mem);
            for (int32_t i = 0; i < len; ++i) {
                new (&data[i]) T();
            }
        } else {
            data = nullptr;
        }
    }

    T& operator[](int32_t index) {
        if (index < 0 || index >= len) {
            throw ArrayIndexOutOfBoundsException("Index " + std::to_string(index) + " out of bounds for length " + std::to_string(len));
        }
        return data[index];
    }

    const T& operator[](int32_t index) const {
        if (index < 0 || index >= len) {
            throw ArrayIndexOutOfBoundsException("Index " + std::to_string(index) + " out of bounds for length " + std::to_string(len));
        }
        return data[index];
    }

    int32_t getLength() const {
        return len;
    }

    T* begin() noexcept { return data; }
    T* end() noexcept { return data + len; }
    const T* begin() const noexcept { return data; }
    const T* end() const noexcept { return data + len; }
    const T* cbegin() const noexcept { return data; }
    const T* cend() const noexcept { return data + len; }
};

} // namespace lang
} // namespace java
