#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "=== Initializing and updating submodules recursively ==="
git -C "${REPO_ROOT}" submodule update --init --recursive

echo "=== Verifying Boehm GC submodule build ==="
if [ ! -f "${REPO_ROOT}/bdwgc/.libs/libgc.so" ]; then
    echo "Building Boehm GC (bdwgc)..."
    (
        cd "${REPO_ROOT}/bdwgc"
        ./autogen.sh
        ./configure --enable-threads=posix --enable-shared --disable-static
        make -j"$(nproc)"
    )
else
    echo "Boehm GC shared library already built at ${REPO_ROOT}/bdwgc/.libs/libgc.so"
fi

echo "=== Verifying Submodule Readiness ==="
if [ -f "${REPO_ROOT}/bdwgc/.libs/libgc.so" ] && [ -f "${REPO_ROOT}/grammars-v4/java/java/JavaParser.g4" ]; then
    echo "Submodules grammars-v4 and bdwgc are initialized and ready for development."
else
    echo "Error: Submodule readiness check failed." >&2
    exit 1
fi
