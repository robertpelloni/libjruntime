#pragma once

#include <exception>
#include <string>

namespace java {
namespace lang {

class Exception : public std::exception {
protected:
    std::string message;
public:
    Exception() : message("Exception") {}
    Exception(const std::string& msg) : message(msg) {}

    virtual const char* what() const noexcept override {
        return message.c_str();
    }
};

class NullPointerException : public Exception {
public:
    NullPointerException() : Exception("NullPointerException") {}
    NullPointerException(const std::string& msg) : Exception("NullPointerException: " + msg) {}
};

class ArrayIndexOutOfBoundsException : public Exception {
public:
    ArrayIndexOutOfBoundsException() : Exception("ArrayIndexOutOfBoundsException") {}
    ArrayIndexOutOfBoundsException(const std::string& msg) : Exception("ArrayIndexOutOfBoundsException: " + msg) {}
};

} // namespace lang
} // namespace java
