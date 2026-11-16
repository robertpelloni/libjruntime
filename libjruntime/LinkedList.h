#pragma once

#include <cstdint>
#include <new>
#include <string>
#include "Object.h"
#include "Exceptions.h"
#include "GC.h"

namespace java {
namespace util {

template<typename T>
class LinkedList : public ::java::lang::Object {
private:
    struct Node {
        T data;
        Node* prev;
        Node* next;

        Node(const T& val, Node* p = nullptr, Node* n = nullptr)
            : data(val), prev(p), next(n) {}
    };

    Node* head;
    Node* tail;
    int32_t currentSize;

    Node* getNode(int32_t index) const {
        if (index < 0 || index >= currentSize) {
            throw ::java::lang::ArrayIndexOutOfBoundsException(
                "Index " + std::to_string(index) + " out of bounds for size " + std::to_string(currentSize));
        }
        if (index < (currentSize / 2)) {
            Node* cur = head;
            for (int32_t i = 0; i < index; ++i) {
                cur = cur->next;
            }
            return cur;
        } else {
            Node* cur = tail;
            for (int32_t i = currentSize - 1; i > index; --i) {
                cur = cur->prev;
            }
            return cur;
        }
    }

public:
    class Iterator {
    private:
        Node* current;
    public:
        Iterator(Node* n) : current(n) {}

        T& operator*() const {
            return current->data;
        }

        Iterator& operator++() {
            if (current != nullptr) {
                current = current->next;
            }
            return *this;
        }

        bool operator!=(const Iterator& other) const {
            return current != other.current;
        }

        bool operator==(const Iterator& other) const {
            return current == other.current;
        }
    };

    LinkedList() : head(nullptr), tail(nullptr), currentSize(0) {}

    virtual bool add(T element) {
        addLast(element);
        return true;
    }

    virtual void addFirst(T element) {
        void* mem = GC_MALLOC(sizeof(Node));
        if (mem == nullptr) {
            throw std::bad_alloc();
        }
        Node* newNode = new (mem) Node(element, nullptr, head);
        if (head == nullptr) {
            head = tail = newNode;
        } else {
            head->prev = newNode;
            head = newNode;
        }
        currentSize++;
    }

    virtual void addLast(T element) {
        void* mem = GC_MALLOC(sizeof(Node));
        if (mem == nullptr) {
            throw std::bad_alloc();
        }
        Node* newNode = new (mem) Node(element, tail, nullptr);
        if (tail == nullptr) {
            head = tail = newNode;
        } else {
            tail->next = newNode;
            tail = newNode;
        }
        currentSize++;
    }

    virtual T get(int32_t index) const {
        return getNode(index)->data;
    }

    virtual T getFirst() const {
        if (head == nullptr) {
            throw ::java::lang::Exception("NoSuchElementException: LinkedList is empty");
        }
        return head->data;
    }

    virtual T getLast() const {
        if (tail == nullptr) {
            throw ::java::lang::Exception("NoSuchElementException: LinkedList is empty");
        }
        return tail->data;
    }

    virtual T removeFirst() {
        if (head == nullptr) {
            throw ::java::lang::Exception("NoSuchElementException: LinkedList is empty");
        }
        Node* oldHead = head;
        T oldVal = oldHead->data;
        head = head->next;
        if (head != nullptr) {
            head->prev = nullptr;
        } else {
            tail = nullptr;
        }
        currentSize--;
        return oldVal;
    }

    virtual T removeLast() {
        if (tail == nullptr) {
            throw ::java::lang::Exception("NoSuchElementException: LinkedList is empty");
        }
        Node* oldTail = tail;
        T oldVal = oldTail->data;
        tail = tail->prev;
        if (tail != nullptr) {
            tail->next = nullptr;
        } else {
            head = nullptr;
        }
        currentSize--;
        return oldVal;
    }

    virtual T remove(int32_t index) {
        Node* node = getNode(index);
        T oldVal = node->data;
        if (node->prev != nullptr) {
            node->prev->next = node->next;
        } else {
            head = node->next;
        }
        if (node->next != nullptr) {
            node->next->prev = node->prev;
        } else {
            tail = node->prev;
        }
        currentSize--;
        return oldVal;
    }

    virtual int32_t size() const {
        return currentSize;
    }

    virtual bool isEmpty() const {
        return currentSize == 0;
    }

    virtual void clear() {
        head = tail = nullptr;
        currentSize = 0;
    }

    virtual bool contains(T element) const {
        for (Node* cur = head; cur != nullptr; cur = cur->next) {
            if (cur->data == element) {
                return true;
            }
        }
        return false;
    }

    Iterator begin() const { return Iterator(head); }
    Iterator end() const { return Iterator(nullptr); }
};

} // namespace util
} // namespace java
