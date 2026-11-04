package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/user/j2cpp-engine/parser"
)

func transpileJavaSnippet(t *testing.T, javaSource string) (string, string) {
	tempDir, err := os.MkdirTemp("", "transpiler_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	input := antlr.NewInputStream(javaSource)
	lexer := parser.NewJavaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewJavaParser(stream)

	tree := p.CompilationUnit()
	visitor := NewCppTranspilerVisitor(tempDir)
	visitor.Visit(tree)

	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read temp dir: %v", err)
	}

	var headerContent, cppContent string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".h") {
			data, _ := os.ReadFile(filepath.Join(tempDir, f.Name()))
			headerContent = string(data)
		} else if strings.HasSuffix(f.Name(), ".cpp") {
			data, _ := os.ReadFile(filepath.Join(tempDir, f.Name()))
			cppContent = string(data)
		}
	}
	return headerContent, cppContent
}

func TestTranspileClassAndMethods(t *testing.T) {
	javaSrc := `
public class Calculator {
    int total;

    public int add(int a, int b) {
        return a + b;
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "class Calculator : public java::lang::Object") {
		t.Errorf("Expected class inheritance in header, got:\n%s", header)
	}
	if !strings.Contains(header, "virtual int32_t add(int32_t a, int32_t b);") {
		t.Errorf("Expected virtual add method in header, got:\n%s", header)
	}
	if !strings.Contains(header, "int32_t total;") {
		t.Errorf("Expected total field declaration in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "int32_t Calculator::add(int32_t a, int32_t b) {") {
		t.Errorf("Expected add method implementation in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "return a + b;") {
		t.Errorf("Expected return statement in cpp, got:\n%s", cpp)
	}
}

func TestTranspileTryCatchFinally(t *testing.T) {
	javaSrc := `
public class TryCatchTest {
    public void execute() {
        try {
            int x = 1;
        } catch (Exception e) {
            int err = 2;
        } finally {
            int done = 3;
        }
    }
}
`
	_, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(cpp, "try {") {
		t.Errorf("Expected try block in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "catch (java::lang::Exception* e) {") {
		t.Errorf("Expected catch block in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_FINALLY(") {
		t.Errorf("Expected JAVA_FINALLY guard in cpp, got:\n%s", cpp)
	}
}

func TestTranspileObjectCreationAndNullCheck(t *testing.T) {
	javaSrc := `
public class ObjectTest {
    public void run() {
        Object obj = new Object();
        obj.toString();
    }
}
`
	_, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(cpp, "java_new<Object>()") {
		t.Errorf("Expected java_new<Object>() allocation, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(obj)->toString()") {
		t.Errorf("Expected JAVA_NULL_CHECK macro wrapping member access, got:\n%s", cpp)
	}
}
