#pragma once

#include <iostream>
#include <cstdint>
#include <string>
#include "Object.h"
#include "String.h"

namespace java {
namespace lang {

class PrintStream : public Object {
public:
    void print(String* s) {
        if (s != nullptr) {
            std::cout << s->getStdString();
        } else {
            std::cout << "null";
        }
    }

    void print(const char* s) {
        if (s != nullptr) {
            std::cout << s;
        } else {
            std::cout << "null";
        }
    }

    void print(int32_t val) {
        std::cout << val;
    }

    void print(int64_t val) {
        std::cout << val;
    }

    void print(bool val) {
        std::cout << (val ? "true" : "false");
    }

    void print(Object* obj) {
        if (obj != nullptr) {
            String* s = obj->toString();
            if (s != nullptr) {
                std::cout << s->getStdString();
            } else {
                std::cout << "null";
            }
        } else {
            std::cout << "null";
        }
    }

    void println() {
        std::cout << std::endl;
    }

    void println(String* s) {
        print(s);
        println();
    }

    void println(const char* s) {
        print(s);
        println();
    }

    void println(int32_t val) {
        print(val);
        println();
    }

    void println(int64_t val) {
        print(val);
        println();
    }

    void println(bool val) {
        print(val);
        println();
    }

    void println(Object* obj) {
        print(obj);
        println();
    }
};

class System {
private:
    static PrintStream* initOut() {
        static PrintStream instance;
        return &instance;
    }

public:
    inline static PrintStream* out = initOut();
};

} // namespace lang
} // namespace java
