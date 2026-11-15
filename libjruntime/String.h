#pragma once

#include "Object.h"
#include <string>
#include <string_view>

namespace java {
namespace lang {

class String : public Object {
private:
    std::string value;

public:
    String() : value("") {}
    String(const char* str) : value(str) {}
    String(std::string_view str) : value(str) {}

    virtual int hashCode() override {
        int hash = 0;
        for (char c : value) {
            hash = 31 * hash + c;
        }
        return hash;
    }

    virtual bool equals(Object* obj) override {
        String* other = dynamic_cast<String*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return this;
    }

    const std::string& getStdString() const {
        return value;
    }

    String* operator+(const String& other) {
        return new String(this->value + other.value);
    }
};

inline String* Object::toString() {
    return new String("Object@" + std::to_string(reinterpret_cast<long long>(this)));
}

} // namespace lang
} // namespace java
