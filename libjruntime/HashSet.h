#pragma once

#include <cstdint>
#include <new>
#include <algorithm>
#include "Object.h"
#include "HashMap.h"
#include "GC.h"

namespace java {
namespace util {

template <typename E>
class HashSet : public ::java::lang::Object {
private:
    HashMap<E, ::java::lang::Object*>* map;
    E* elementList;
    int32_t listCapacity;
    int32_t listSize;

    void ensureCapacity(int32_t minCap) {
        if (minCap <= listCapacity) return;
        int32_t newCap = listCapacity == 0 ? 10 : listCapacity * 2;
        if (newCap < minCap) newCap = minCap;

        void* mem = GC_MALLOC(sizeof(E) * newCap);
        if (mem == nullptr) {
            throw std::bad_alloc();
        }
        E* newList = static_cast<E*>(mem);
        for (int32_t i = 0; i < listSize; ++i) {
            new (&newList[i]) E(elementList[i]);
        }
        elementList = newList;
        listCapacity = newCap;
    }

public:
    HashSet() : elementList(nullptr), listCapacity(0), listSize(0) {
        map = java_new<HashMap<E, ::java::lang::Object*>>();
        ensureCapacity(10);
    }

    HashSet(int32_t initCap) : elementList(nullptr), listCapacity(0), listSize(0) {
        map = java_new<HashMap<E, ::java::lang::Object*>>(initCap);
        ensureCapacity(initCap > 0 ? initCap : 10);
    }

    virtual bool add(E element) {
        if (map->containsKey(element)) {
            return false;
        }
        map->put(element, reinterpret_cast<::java::lang::Object*>(1));
        ensureCapacity(listSize + 1);
        elementList[listSize++] = element;
        return true;
    }

    virtual bool contains(E element) const {
        return map->containsKey(element);
    }

    virtual bool remove(E element) {
        if (!map->containsKey(element)) {
            return false;
        }
        map->remove(element);
        for (int32_t i = 0; i < listSize; ++i) {
            if (internal::areKeysEqual(elementList[i], element)) {
                for (int32_t j = i; j < listSize - 1; ++j) {
                    elementList[j] = elementList[j + 1];
                }
                listSize--;
                break;
            }
        }
        return true;
    }

    virtual int32_t size() const {
        return map->size();
    }

    virtual bool isEmpty() const {
        return map->isEmpty();
    }

    virtual void clear() {
        map->clear();
        listSize = 0;
    }

    E* begin() noexcept { return elementList; }
    E* end() noexcept { return elementList + listSize; }
    const E* begin() const noexcept { return elementList; }
    const E* end() const noexcept { return elementList + listSize; }
    const E* cbegin() const noexcept { return elementList; }
    const E* cend() const noexcept { return elementList + listSize; }
};

} // namespace util
} // namespace java
