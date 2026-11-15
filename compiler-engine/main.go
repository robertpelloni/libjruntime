package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/user/j2cpp-engine/parser"
)

func main() {
	srcDir := flag.String("src", "", "Source directory containing Java files")
	outDir := flag.String("out", "", "Output directory for C++ files")
	flag.Parse()

	if *srcDir == "" || *outDir == "" {
		fmt.Println("Usage: j2cpp-engine -src=<source_directory> -out=<output_directory>")
		os.Exit(1)
	}

	err := os.MkdirAll(*outDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	err = filepath.WalkDir(*srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".java") {
			processJavaFile(path, *outDir)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking source directory: %v", err)
	}
}

func processJavaFile(javaFile string, outDir string) {
	fmt.Printf("Processing %s\n", javaFile)

	input, err := antlr.NewFileStream(javaFile)
	if err != nil {
		log.Printf("Failed to read file %s: %v", javaFile, err)
		return
	}

	lexer := parser.NewJavaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewJavaParser(stream)

	tree := p.CompilationUnit()

	visitor := NewCppTranspilerVisitor(outDir)
	visitor.Visit(tree)
}
