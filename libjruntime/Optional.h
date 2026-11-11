#pragma once

#include <optional>
#include "Object.h"
#include "Exceptions.h"
#include "Functions.h"
#include "GC.h"

namespace java {
namespace util {

template <typename T>
class Optional : public ::java::lang::Object {
    template <typename> friend class Optional;
private:
    bool present;
    T val;

public:
    Optional() : present(false), val() {}
    Optional(T value) : present(true), val(value) {}

    static Optional<T>* empty() {
        return java_new<Optional<T>>();
    }

    static Optional<T>* of(T value) {
        if constexpr (std::is_pointer_v<T>) {
            if (value == nullptr) {
                throw ::java::lang::NullPointerException("Optional.of cannot take null");
            }
        }
        return java_new<Optional<T>>(value);
    }

    static Optional<T>* ofNullable(T value) {
        if constexpr (std::is_pointer_v<T>) {
            if (value == nullptr) {
                return empty();
            }
        }
        return java_new<Optional<T>>(value);
    }

    virtual bool isPresent() const {
        return present;
    }

    virtual bool isEmpty() const {
        return !present;
    }

    virtual T get() const {
        if (!present) {
            throw ::java::lang::Exception("NoSuchElementException: No value present");
        }
        return val;
    }

    virtual T orElse(T other) const {
        return present ? val : other;
    }

    virtual void ifPresent(::java::util::function::Consumer<T>* consumer) const {
        if (present && consumer != nullptr) {
            consumer->accept(val);
        }
    }

    virtual Optional<T>* filter(::java::util::function::Predicate<T>* predicate) const {
        if (!present) {
            return const_cast<Optional<T>*>(this);
        }
        if (predicate != nullptr && predicate->test(val)) {
            return const_cast<Optional<T>*>(this);
        }
        return empty();
    }

    template <typename R>
    Optional<R>* map(::java::util::function::Function<T, R>* mapper) const {
        if (!present || mapper == nullptr) {
            return Optional<R>::empty();
        }
        return Optional<R>::ofNullable(mapper->apply(val));
    }
};

} // namespace util
} // namespace java
