#pragma once

template<typename F>
class FinallyGuard {
    F f;
public:
    FinallyGuard(F f) : f(f) {}
    ~FinallyGuard() { f(); }
};

template<typename F>
FinallyGuard<F> make_finally_guard(F f) {
    return FinallyGuard<F>(f);
}

#define _JAVA_FINALLY_CONCAT_INNER(a, b) a ## b
#define _JAVA_FINALLY_CONCAT(a, b) _JAVA_FINALLY_CONCAT_INNER(a, b)
#define JAVA_FINALLY(...) [[maybe_unused]] auto _JAVA_FINALLY_CONCAT(finally_guard_, __LINE__) = make_finally_guard([&]() { __VA_ARGS__ });
