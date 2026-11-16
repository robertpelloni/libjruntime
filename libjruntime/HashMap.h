#pragma once

#include <cstdint>
#include <new>
#include <functional>
#include "Object.h"
#include "GC.h"

namespace java {
namespace util {

namespace internal {

template<typename T>
bool areKeysEqual(const T& a, const T& b) {
    if constexpr (std::is_pointer_v<T>) {
        if (a == b) return true;
        if (a == nullptr || b == nullptr) return false;
        if constexpr (std::is_base_of_v<::java::lang::Object, std::remove_pointer_t<T>>) {
            return a->equals(b);
        } else {
            return *a == *b;
        }
    } else {
        return a == b;
    }
}

template<typename T>
std::size_t computeHash(const T& a) {
    if constexpr (std::is_pointer_v<T>) {
        if (a == nullptr) return 0;
        if constexpr (std::is_base_of_v<::java::lang::Object, std::remove_pointer_t<T>>) {
            return static_cast<std::size_t>(a->hashCode());
        } else {
            return std::hash<T>{}(a);
        }
    } else {
        return std::hash<T>{}(a);
    }
}

} // namespace internal

template<typename K, typename V>
class HashMap : public ::java::lang::Object {
private:
    struct Node {
        K key;
        V value;
        Node* next;

        Node(const K& k, const V& v, Node* n = nullptr)
            : key(k), value(v), next(n) {}
    };

    Node** buckets;
    int32_t capacity;
    int32_t currentSize;

    void rehash(int32_t newCap) {
        void* mem = GC_MALLOC(sizeof(Node*) * newCap);
        if (mem == nullptr) {
            throw std::bad_alloc();
        }
        Node** newBuckets = static_cast<Node**>(mem);
        for (int32_t i = 0; i < newCap; ++i) {
            newBuckets[i] = nullptr;
        }

        for (int32_t i = 0; i < capacity; ++i) {
            Node* curr = buckets[i];
            while (curr != nullptr) {
                Node* nextNode = curr->next;
                std::size_t idx = internal::computeHash(curr->key) % newCap;
                curr->next = newBuckets[idx];
                newBuckets[idx] = curr;
                curr = nextNode;
            }
        }

        buckets = newBuckets;
        capacity = newCap;
    }

public:
    HashMap(int32_t initCap = 16) : capacity(initCap > 0 ? initCap : 16), currentSize(0) {
        void* mem = GC_MALLOC(sizeof(Node*) * capacity);
        if (mem == nullptr) {
            throw std::bad_alloc();
        }
        buckets = static_cast<Node**>(mem);
        for (int32_t i = 0; i < capacity; ++i) {
            buckets[i] = nullptr;
        }
    }

    virtual V put(K key, V value) {
        if (currentSize + 1 > capacity * 0.75) {
            rehash(capacity * 2);
        }

        std::size_t idx = internal::computeHash(key) % capacity;
        Node* curr = buckets[idx];
        while (curr != nullptr) {
            if (internal::areKeysEqual(curr->key, key)) {
                V oldVal = curr->value;
                curr->value = value;
                return oldVal;
            }
            curr = curr->next;
        }

        void* mem = GC_MALLOC(sizeof(Node));
        if (mem == nullptr) {
            throw std::bad_alloc();
        }
        Node* newNode = new (mem) Node(key, value, buckets[idx]);
        buckets[idx] = newNode;
        currentSize++;

        if constexpr (std::is_pointer_v<V>) {
            return nullptr;
        } else {
            return V{};
        }
    }

    virtual V get(K key) const {
        if (capacity == 0) {
            if constexpr (std::is_pointer_v<V>) return nullptr;
            else return V{};
        }
        std::size_t idx = internal::computeHash(key) % capacity;
        Node* curr = buckets[idx];
        while (curr != nullptr) {
            if (internal::areKeysEqual(curr->key, key)) {
                return curr->value;
            }
            curr = curr->next;
        }
        if constexpr (std::is_pointer_v<V>) return nullptr;
        else return V{};
    }

    virtual bool containsKey(K key) const {
        if (capacity == 0) return false;
        std::size_t idx = internal::computeHash(key) % capacity;
        Node* curr = buckets[idx];
        while (curr != nullptr) {
            if (internal::areKeysEqual(curr->key, key)) {
                return true;
            }
            curr = curr->next;
        }
        return false;
    }

    virtual V remove(K key) {
        if (capacity == 0) {
            if constexpr (std::is_pointer_v<V>) return nullptr;
            else return V{};
        }
        std::size_t idx = internal::computeHash(key) % capacity;
        Node* curr = buckets[idx];
        Node* prev = nullptr;

        while (curr != nullptr) {
            if (internal::areKeysEqual(curr->key, key)) {
                if (prev == nullptr) {
                    buckets[idx] = curr->next;
                } else {
                    prev->next = curr->next;
                }
                V val = curr->value;
                currentSize--;
                return val;
            }
            prev = curr;
            curr = curr->next;
        }
        if constexpr (std::is_pointer_v<V>) return nullptr;
        else return V{};
    }

    virtual int32_t size() const {
        return currentSize;
    }

    virtual bool isEmpty() const {
        return currentSize == 0;
    }

    virtual void clear() {
        for (int32_t i = 0; i < capacity; ++i) {
            buckets[i] = nullptr;
        }
        currentSize = 0;
    }
};

} // namespace util
} // namespace java
