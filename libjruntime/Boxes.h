#pragma once

#include <cstdint>
#include <string>
#include <bit>
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

class Double : public Object {
private:
    double value;

public:
    Double(double val) : value(val) {}

    static Double* valueOf(double val) {
        return java_new<Double>(val);
    }

    double doubleValue() const {
        return value;
    }

    float floatValue() const {
        return static_cast<float>(value);
    }

    int32_t intValue() const {
        return static_cast<int32_t>(value);
    }

    int64_t longValue() const {
        return static_cast<int64_t>(value);
    }

    virtual int hashCode() override {
        int64_t bits = std::bit_cast<int64_t>(value);
        return static_cast<int>(bits ^ (bits >> 32));
    }

    virtual bool equals(Object* obj) override {
        Double* other = dynamic_cast<Double*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return java_new<String>(std::to_string(value));
    }
};

class Float : public Object {
private:
    float value;

public:
    Float(float val) : value(val) {}

    static Float* valueOf(float val) {
        return java_new<Float>(val);
    }

    float floatValue() const {
        return value;
    }

    double doubleValue() const {
        return static_cast<double>(value);
    }

    int32_t intValue() const {
        return static_cast<int32_t>(value);
    }

    int64_t longValue() const {
        return static_cast<int64_t>(value);
    }

    virtual int hashCode() override {
        int32_t bits = std::bit_cast<int32_t>(value);
        return bits;
    }

    virtual bool equals(Object* obj) override {
        Float* other = dynamic_cast<Float*>(obj);
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

class Byte : public Object {
private:
    int8_t value;

public:
    Byte(int8_t val) : value(val) {}

    static Byte* valueOf(int8_t val) {
        return java_new<Byte>(val);
    }

    int8_t byteValue() const {
        return value;
    }

    int32_t intValue() const {
        return static_cast<int32_t>(value);
    }

    virtual int hashCode() override {
        return static_cast<int>(value);
    }

    virtual bool equals(Object* obj) override {
        Byte* other = dynamic_cast<Byte*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return java_new<String>(std::to_string(static_cast<int32_t>(value)));
    }
};

class Character : public Object {
private:
    char16_t value;

public:
    Character(char16_t val) : value(val) {}

    static Character* valueOf(char16_t val) {
        return java_new<Character>(val);
    }

    char16_t charValue() const {
        return value;
    }

    virtual int hashCode() override {
        return static_cast<int>(value);
    }

    virtual bool equals(Object* obj) override {
        Character* other = dynamic_cast<Character*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return java_new<String>(std::string(1, static_cast<char>(value)));
    }
};

class Short : public Object {
private:
    int16_t value;

public:
    Short(int16_t val) : value(val) {}

    static Short* valueOf(int16_t val) {
        return java_new<Short>(val);
    }

    int16_t shortValue() const {
        return value;
    }

    int32_t intValue() const {
        return static_cast<int32_t>(value);
    }

    double doubleValue() const {
        return static_cast<double>(value);
    }

    virtual int hashCode() override {
        return static_cast<int>(value);
    }

    virtual bool equals(Object* obj) override {
        Short* other = dynamic_cast<Short*>(obj);
        if (other == nullptr) return false;
        return this->value == other->value;
    }

    virtual String* toString() override {
        return java_new<String>(std::to_string(static_cast<int32_t>(value)));
    }
};

} // namespace lang
} // namespace java
