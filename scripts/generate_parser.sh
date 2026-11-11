#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ANTLR_JAR="${REPO_ROOT}/antlr-4.13.1-complete.jar"
GRAMMAR_DIR="${REPO_ROOT}/grammars-v4/java/java"
OUT_DIR="${REPO_ROOT}/compiler-engine/parser"

echo "=== Generating Go ANTLR4 Parser Bindings from grammars-v4 ==="

if [ ! -f "${ANTLR_JAR}" ]; then
    echo "Error: ANTLR jar not found at ${ANTLR_JAR}" >&2
    exit 1
fi

if [ ! -d "${GRAMMAR_DIR}" ]; then
    echo "Error: grammars-v4 submodule not found at ${GRAMMAR_DIR}" >&2
    exit 1
fi

mkdir -p "${OUT_DIR}"

java -jar "${ANTLR_JAR}" \
    -Dlanguage=Go \
    -visitor \
    -package parser \
    -o "${OUT_DIR}" \
    "${GRAMMAR_DIR}/JavaLexer.g4" \
    "${GRAMMAR_DIR}/JavaParser.g4"

# ANTLR Go target translates Java semantic predicates with 'this.' as literal 'this.' in Go.
# Normalize 'this.' to parser receiver 'p.' in generated java_parser.go
if [ -f "${OUT_DIR}/java_parser.go" ]; then
    sed -i 's/this\.IsNotIdentifierAssign/p.IsNotIdentifierAssign/g' "${OUT_DIR}/java_parser.go"
    sed -i 's/this\.DoLastRecordComponent/p.DoLastRecordComponent/g' "${OUT_DIR}/java_parser.go"
fi

echo "Parser code generated successfully in ${OUT_DIR}"
