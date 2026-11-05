#pragma once

#include <cstdint>
#include <new>
#include <string>
#include <algorithm>
#include "Object.h"
#include "Exceptions.h"
#include "GC.h"

namespace java {
namespace util {

template<typename T>
class ArrayList : public ::java::lang::Object {
private:
    T* data;
    int32_t capacity;
    int32_t currentSize;

    void ensureCapacity(int32_t minCap) {
        if (minCap <= capacity) {
            return;
        }
        int32_t newCap = capacity == 0 ? 10 : capacity * 2;
        if (newCap < minCap) {
            newCap = minCap;
        }

        void* mem = GC_MALLOC(sizeof(T) * newCap);
        if (mem == nullptr) {
            throw std::bad_alloc();
        }
        T* newData = static_cast<T*>(mem);

        for (int32_t i = 0; i < currentSize; ++i) {
            new (&newData[i]) T(data[i]);
        }

        data = newData;
        capacity = newCap;
    }

public:
    ArrayList() : data(nullptr), capacity(0), currentSize(0) {
        ensureCapacity(10);
    }

    ArrayList(int32_t initialCapacity) : data(nullptr), capacity(0), currentSize(0) {
        if (initialCapacity < 0) {
            throw ::java::lang::Exception("IllegalArgumentException: Negative initial capacity");
        }
        ensureCapacity(initialCapacity > 0 ? initialCapacity : 10);
    }

    virtual bool add(T element) {
        ensureCapacity(currentSize + 1);
        data[currentSize++] = element;
        return true;
    }

    virtual void add(int32_t index, T element) {
        if (index < 0 || index > currentSize) {
            throw ::java::lang::ArrayIndexOutOfBoundsException(
                "Index " + std::to_string(index) + " out of bounds for size " + std::to_string(currentSize));
        }
        ensureCapacity(currentSize + 1);
        for (int32_t i = currentSize; i > index; --i) {
            data[i] = data[i - 1];
        }
        data[index] = element;
        currentSize++;
    }

    virtual T get(int32_t index) const {
        if (index < 0 || index >= currentSize) {
            throw ::java::lang::ArrayIndexOutOfBoundsException(
                "Index " + std::to_string(index) + " out of bounds for size " + std::to_string(currentSize));
        }
        return data[index];
    }

    virtual T set(int32_t index, T element) {
        if (index < 0 || index >= currentSize) {
            throw ::java::lang::ArrayIndexOutOfBoundsException(
                "Index " + std::to_string(index) + " out of bounds for size " + std::to_string(currentSize));
        }
        T old = data[index];
        data[index] = element;
        return old;
    }

    virtual T remove(int32_t index) {
        if (index < 0 || index >= currentSize) {
            throw ::java::lang::ArrayIndexOutOfBoundsException(
                "Index " + std::to_string(index) + " out of bounds for size " + std::to_string(currentSize));
        }
        T old = data[index];
        for (int32_t i = index; i < currentSize - 1; ++i) {
            data[i] = data[i + 1];
        }
        currentSize--;
        return old;
    }

    virtual int32_t size() const {
        return currentSize;
    }

    virtual bool isEmpty() const {
        return currentSize == 0;
    }

    virtual void clear() {
        currentSize = 0;
    }

    virtual bool contains(T element) const {
        for (int32_t i = 0; i < currentSize; ++i) {
            if (data[i] == element) {
                return true;
            }
        }
        return false;
    }
};

} // namespace util
} // namespace java
