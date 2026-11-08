#pragma once

#include <string>
#include <string_view>
#include "Object.h"
#include "String.h"
#include "GC.h"

namespace java {
namespace lang {

template <typename T>
class Class : public Object {
private:
    String* fullName;
    String* simpleName;

public:
    Class(const char* name) {
        fullName = java_new<String>(name);
        std::string_view sv(name);
        auto lastColon = sv.rfind(':');
        if (lastColon != std::string_view::npos) {
            simpleName = java_new<String>(std::string(sv.substr(lastColon + 1)));
        } else {
            simpleName = fullName;
        }
    }

    virtual String* getName() const {
        return fullName;
    }

    virtual String* getSimpleName() const {
        return simpleName;
    }

    virtual bool isInstance(Object* obj) const {
        if (obj == nullptr) {
            return false;
        }
        if constexpr (std::is_base_of_v<Object, T>) {
            return dynamic_cast<T*>(obj) != nullptr;
        } else {
            return false;
        }
    }

    virtual String* toString() override {
        std::string s = "class " + fullName->getStdString();
        return java_new<String>(s);
    }

    static Class<T>* forType(const char* name) {
        return java_new<Class<T>>(name);
    }
};

inline Class<Object>* Object::getClass() {
    return Class<Object>::forType("java.lang.Object");
}

} // namespace lang
} // namespace java
