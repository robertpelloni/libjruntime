#pragma once

#include <string>
#include <cstdint>
#include "Object.h"
#include "String.h"
#include "GC.h"

namespace java {
namespace lang {

class StringBuilder : public Object {
private:
    std::string buffer;

public:
    StringBuilder() : buffer("") {}
    StringBuilder(int32_t capacity) : buffer("") {
        if (capacity > 0) {
            buffer.reserve(capacity);
        }
    }
    StringBuilder(String* str) : buffer(str != nullptr ? str->getStdString() : "null") {}
    StringBuilder(const char* str) : buffer(str != nullptr ? str : "null") {}

    StringBuilder* append(String* str) {
        if (str != nullptr) {
            buffer += str->getStdString();
        } else {
            buffer += "null";
        }
        return this;
    }

    StringBuilder* append(const char* str) {
        if (str != nullptr) {
            buffer += str;
        } else {
            buffer += "null";
        }
        return this;
    }

    StringBuilder* append(int32_t val) {
        buffer += std::to_string(val);
        return this;
    }

    StringBuilder* append(int64_t val) {
        buffer += std::to_string(val);
        return this;
    }

    StringBuilder* append(bool val) {
        buffer += (val ? "true" : "false");
        return this;
    }

    StringBuilder* append(Object* obj) {
        if (obj != nullptr) {
            String* s = obj->toString();
            if (s != nullptr) {
                buffer += s->getStdString();
            } else {
                buffer += "null";
            }
        } else {
            buffer += "null";
        }
        return this;
    }

    int32_t length() const {
        return static_cast<int32_t>(buffer.length());
    }

    virtual String* toString() override {
        return java_new<String>(buffer);
    }
};

} // namespace lang
} // namespace java
