#pragma once

#include <cstdint>
#include "Object.h"
#include "ArrayList.h"
#include "Functions.h"
#include "GC.h"

namespace java {
namespace util {
namespace stream {

template <typename T>
class Stream : public ::java::lang::Object {
    template <typename> friend class Stream;
private:
    ::java::util::ArrayList<T>* elements;

public:
    Stream() {
        elements = java_new<::java::util::ArrayList<T>>();
    }

    Stream(::java::util::ArrayList<T>* list) {
        elements = java_new<::java::util::ArrayList<T>>();
        if (list != nullptr) {
            for (int32_t i = 0; i < list->size(); ++i) {
                elements->add(list->get(i));
            }
        }
    }

    virtual Stream<T>* filter(::java::util::function::Predicate<T>* predicate) {
        Stream<T>* result = java_new<Stream<T>>();
        if (predicate != nullptr) {
            for (int32_t i = 0; i < elements->size(); ++i) {
                T item = elements->get(i);
                if (predicate->test(item)) {
                    result->elements->add(item);
                }
            }
        }
        return result;
    }

    template <typename R>
    Stream<R>* map(::java::util::function::Function<T, R>* mapper) {
        Stream<R>* result = java_new<Stream<R>>();
        if (mapper != nullptr) {
            for (int32_t i = 0; i < elements->size(); ++i) {
                result->elements->add(mapper->apply(elements->get(i)));
            }
        }
        return result;
    }

    virtual Stream<T>* limit(int32_t maxSize) {
        Stream<T>* result = java_new<Stream<T>>();
        int32_t n = elements->size() < maxSize ? elements->size() : maxSize;
        for (int32_t i = 0; i < n; ++i) {
            result->elements->add(elements->get(i));
        }
        return result;
    }

    virtual Stream<T>* skip(int32_t n) {
        Stream<T>* result = java_new<Stream<T>>();
        for (int32_t i = n; i < elements->size(); ++i) {
            result->elements->add(elements->get(i));
        }
        return result;
    }

    virtual int64_t count() const {
        return static_cast<int64_t>(elements->size());
    }

    virtual void forEach(::java::util::function::Consumer<T>* action) {
        if (action != nullptr) {
            for (int32_t i = 0; i < elements->size(); ++i) {
                action->accept(elements->get(i));
            }
        }
    }

    virtual ::java::util::ArrayList<T>* toList() {
        return elements;
    }
};

} // namespace stream
} // namespace util
} // namespace java

namespace java {
namespace util {

template <typename T>
::java::util::stream::Stream<T>* ArrayList<T>::stream() {
    return java_new<::java::util::stream::Stream<T>>(this);
}

} // namespace util
} // namespace java
