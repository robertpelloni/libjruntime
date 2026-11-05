#pragma once

#include <cstdint>
#include <string>
#include "Object.h"
#include "String.h"
#include "GC.h"

namespace java {
namespace lang {

class Integer : public Object {
private:
    int32_t value;

public:
    Integer(int32_t val) : value(val) {}

    static Integer* valueOf(int32_t val) {
        return java_new<Integer>(val);
    }

    int32_t intValue() const {
        return value;
    }

    virtual int hashCode() override {
        return value;
    }

    virtual bool equals(Object* obj) override {
        Integer* other = dynamic_cast<Integer*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return java_new<String>(std::to_string(value));
    }
};

class Boolean : public Object {
private:
    bool value;

public:
    Boolean(bool val) : value(val) {}

    static Boolean* valueOf(bool val) {
        return java_new<Boolean>(val);
    }

    bool booleanValue() const {
        return value;
    }

    virtual int hashCode() override {
        return value ? 1231 : 1237;
    }

    virtual bool equals(Object* obj) override {
        Boolean* other = dynamic_cast<Boolean*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return java_new<String>(value ? "true" : "false");
    }
};

class Long : public Object {
private:
    int64_t value;

public:
    Long(int64_t val) : value(val) {}

    static Long* valueOf(int64_t val) {
        return java_new<Long>(val);
    }

    int64_t longValue() const {
        return value;
    }

    int32_t intValue() const {
        return static_cast<int32_t>(value);
    }

    double doubleValue() const {
        return static_cast<double>(value);
    }

    virtual int hashCode() override {
        return static_cast<int>(value ^ (value >> 32));
    }

    virtual bool equals(Object* obj) override {
        Long* other = dynamic_cast<Long*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return java_new<String>(std::to_string(value));
    }
};

} // namespace lang
} // namespace java
