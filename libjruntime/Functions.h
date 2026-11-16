#pragma once

#include <functional>
#include "Object.h"
#include "GC.h"

namespace java {
namespace util {
namespace function {

template <typename T, typename R>
class Function : public virtual java::lang::Object {
public:
    virtual ~Function() = default;
    virtual R apply(T t) = 0;
};

template <typename T>
class Consumer : public virtual java::lang::Object {
public:
    virtual ~Consumer() = default;
    virtual void accept(T t) = 0;
};

template <typename T>
class Supplier : public virtual java::lang::Object {
public:
    virtual ~Supplier() = default;
    virtual T get() = 0;
};

template <typename T>
class Predicate : public virtual java::lang::Object {
public:
    virtual ~Predicate() = default;
    virtual bool test(T t) = 0;
};

// Helper adapters to convert C++ lambdas / callables into GC-tracked functional interfaces
template <typename T, typename R, typename F>
class LambdaFunction : public Function<T, R> {
private:
    F fn;
public:
    LambdaFunction(F f) : fn(f) {}
    R apply(T t) override {
        return fn(t);
    }
};

template <typename T, typename F>
class LambdaConsumer : public Consumer<T> {
private:
    F fn;
public:
    LambdaConsumer(F f) : fn(f) {}
    void accept(T t) override {
        fn(t);
    }
};

template <typename T, typename F>
class LambdaSupplier : public Supplier<T> {
private:
    F fn;
public:
    LambdaSupplier(F f) : fn(f) {}
    T get() override {
        return fn();
    }
};

template <typename T, typename F>
class LambdaPredicate : public Predicate<T> {
private:
    F fn;
public:
    LambdaPredicate(F f) : fn(f) {}
    bool test(T t) override {
        return fn(t);
    }
};

template <typename T, typename R, typename F>
Function<T, R>* make_function(F&& f) {
    return java_new<LambdaFunction<T, R, std::decay_t<F>>>(std::forward<F>(f));
}

template <typename T, typename F>
Consumer<T>* make_consumer(F&& f) {
    return java_new<LambdaConsumer<T, std::decay_t<F>>>(std::forward<F>(f));
}

template <typename T, typename F>
Supplier<T>* make_supplier(F&& f) {
    return java_new<LambdaSupplier<T, std::decay_t<F>>>(std::forward<F>(f));
}

template <typename T, typename F>
Predicate<T>* make_predicate(F&& f) {
    return java_new<LambdaPredicate<T, std::decay_t<F>>>(std::forward<F>(f));
}

} // namespace function
} // namespace util
} // namespace java
