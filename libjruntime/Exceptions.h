#pragma once

#include <exception>
#include <string>
#include "Object.h"
#include "String.h"

namespace java {
namespace lang {

class Exception : public std::exception {
protected:
    std::string message;
public:
    Exception() : message("Exception") {}
    Exception(const std::string& msg) : message(msg) {}
    Exception(String* msg) : message(msg != nullptr ? msg->getStdString() : "null") {}

    virtual const char* what() const noexcept override {
        return message.c_str();
    }
};

class NullPointerException : public Exception {
public:
    NullPointerException() : Exception("NullPointerException") {}
    NullPointerException(const std::string& msg) : Exception("NullPointerException: " + msg) {}
    NullPointerException(String* msg)
        : Exception(msg != nullptr ? "NullPointerException: " + msg->getStdString() : "NullPointerException") {}
};

class ArrayIndexOutOfBoundsException : public Exception {
public:
    ArrayIndexOutOfBoundsException() : Exception("ArrayIndexOutOfBoundsException") {}
    ArrayIndexOutOfBoundsException(const std::string& msg) : Exception("ArrayIndexOutOfBoundsException: " + msg) {}
    ArrayIndexOutOfBoundsException(String* msg)
        : Exception(msg != nullptr ? "ArrayIndexOutOfBoundsException: " + msg->getStdString() : "ArrayIndexOutOfBoundsException") {}
};

} // namespace lang
} // namespace java
