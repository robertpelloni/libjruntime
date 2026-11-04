#pragma once

#include "Exceptions.h"

#define JAVA_NULL_CHECK(obj) \
    ((obj) == nullptr ? throw java::lang::NullPointerException() : (obj))
