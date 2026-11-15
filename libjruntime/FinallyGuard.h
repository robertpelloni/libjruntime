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

#define JAVA_FINALLY(code) auto finally_guard = make_finally_guard([&]() { code });
