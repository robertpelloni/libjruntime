#pragma once

#include <algorithm>
#include <random>
#include <cstdint>
#include "Object.h"
#include "ArrayList.h"
#include "Exceptions.h"

namespace java {
namespace util {

class Collections : public ::java::lang::Object {
public:
    template<typename T>
    static void sort(ArrayList<T>* list) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.sort: list is null");
        }
        std::sort(list->begin(), list->end());
    }

    template<typename T, typename Compare>
    static void sort(ArrayList<T>* list, Compare comp) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.sort: list is null");
        }
        std::sort(list->begin(), list->end(), comp);
    }

    template<typename T>
    static void reverse(ArrayList<T>* list) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.reverse: list is null");
        }
        std::reverse(list->begin(), list->end());
    }

    template<typename T>
    static void shuffle(ArrayList<T>* list) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.shuffle: list is null");
        }
        std::random_device rd;
        std::mt19937 g(rd());
        std::shuffle(list->begin(), list->end(), g);
    }

    template<typename T>
    static T min(ArrayList<T>* list) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.min: list is null");
        }
        if (list->isEmpty()) {
            throw ::java::lang::Exception("NoSuchElementException: empty collection");
        }
        return *std::min_element(list->begin(), list->end());
    }

    template<typename T>
    static T max(ArrayList<T>* list) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.max: list is null");
        }
        if (list->isEmpty()) {
            throw ::java::lang::Exception("NoSuchElementException: empty collection");
        }
        return *std::max_element(list->begin(), list->end());
    }

    template<typename T>
    static void fill(ArrayList<T>* list, T val) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.fill: list is null");
        }
        std::fill(list->begin(), list->end(), val);
    }

    template<typename T>
    static int32_t frequency(ArrayList<T>* list, T val) {
        if (list == nullptr) {
            throw ::java::lang::NullPointerException("Collections.frequency: list is null");
        }
        int32_t count = 0;
        for (const auto& item : *list) {
            if (item == val) {
                count++;
            }
        }
        return count;
    }
};

} // namespace util
} // namespace java
