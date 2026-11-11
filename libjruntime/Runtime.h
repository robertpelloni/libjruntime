#pragma once

#include "Exceptions.h"

namespace java {
namespace lang {

template<typename T>
inline T* check_not_null(T* obj) {
    if (obj == nullptr) {
        throw java::lang::NullPointerException();
    }
    return obj;
}

} // namespace lang
} // namespace java

#define JAVA_NULL_CHECK(obj) (java::lang::check_not_null(obj))
