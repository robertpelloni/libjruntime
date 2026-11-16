#pragma once

#include <cmath>
#include <cstdint>
#include <algorithm>

namespace java {
namespace lang {

class Math {
public:
    inline static constexpr double PI = 3.14159265358979323846;
    inline static constexpr double E  = 2.71828182845904523536;

    static int32_t abs(int32_t a) {
        return a < 0 ? -a : a;
    }

    static int64_t abs(int64_t a) {
        return a < 0 ? -a : a;
    }

    static double abs(double a) {
        return std::abs(a);
    }

    static int32_t max(int32_t a, int32_t b) {
        return std::max(a, b);
    }

    static int64_t max(int64_t a, int64_t b) {
        return std::max(a, b);
    }

    static double max(double a, double b) {
        return std::max(a, b);
    }

    static int32_t min(int32_t a, int32_t b) {
        return std::min(a, b);
    }

    static int64_t min(int64_t a, int64_t b) {
        return std::min(a, b);
    }

    static double min(double a, double b) {
        return std::min(a, b);
    }

    static double sqrt(double a) {
        return std::sqrt(a);
    }

    static double pow(double a, double b) {
        return std::pow(a, b);
    }

    static double floor(double a) {
        return std::floor(a);
    }

    static double ceil(double a) {
        return std::ceil(a);
    }

    static int64_t round(double a) {
        return static_cast<int64_t>(std::round(a));
    }
};

} // namespace lang
} // namespace java
