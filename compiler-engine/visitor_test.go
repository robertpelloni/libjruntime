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

func TestTranspileTernaryExpression(t *testing.T) {
	javaSrc := `
public class TernaryTest {
    public int choose(int x) {
        return (x > 10 ? 100 : 200);
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "virtual int32_t choose(int32_t x);") {
		t.Errorf("Expected choose method in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "return (x > 10 ? 100 : 200);") {
		t.Errorf("Expected ternary expression in cpp, got:\n%s", cpp)
	}
}

func TestTranspileSwitchStatement(t *testing.T) {
	javaSrc := `
public class SwitchTest {
    public int eval(int code) {
        int result = 0;
        switch (code) {
            case 1:
                result = 10;
                break;
            case 2:
                result = 20;
                break;
            default:
                result = -1;
                break;
        }
        return result;
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "virtual int32_t eval(int32_t code);") {
		t.Errorf("Expected eval method declaration in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "switch (code) {") {
		t.Errorf("Expected switch statement in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "case 1:") || !strings.Contains(cpp, "case 2:") {
		t.Errorf("Expected cases in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "default:") {
		t.Errorf("Expected default label in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "break;") {
		t.Errorf("Expected break statement in cpp, got:\n%s", cpp)
	}
}

func TestTranspileCastAndInstanceOf(t *testing.T) {
	javaSrc := `
public class CastTest {
    public boolean checkAndCast(Object obj) {
        if (obj instanceof Dog) {
            Dog d = (Dog) obj;
            return true;
        }
        return false;
    }
}
`
	_, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(cpp, "#include \"Dog.h\"") {
		t.Errorf("Expected #include \"Dog.h\" in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "(dynamic_cast<Dog*>(obj) != nullptr)") {
		t.Errorf("Expected dynamic_cast instanceof check in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "dynamic_cast<Dog*>(obj)") {
		t.Errorf("Expected dynamic_cast in cpp, got:\n%s", cpp)
	}
}

func TestTranspileEnhancedForLoop(t *testing.T) {
	javaSrc := `
public class ForEachTest {
    public int sum(int[] numbers) {
        int total = 0;
        for (int n : numbers) {
            total = total + n;
        }
        return total;
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "virtual int32_t sum(java::lang::Array<int32_t>* numbers);") {
		t.Errorf("Expected sum method declaration with Array parameter, got:\n%s", header)
	}
	if !strings.Contains(cpp, "for (int32_t n : *JAVA_NULL_CHECK(numbers)) {") {
		t.Errorf("Expected range-based for loop in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "total = total + n;") {
		t.Errorf("Expected loop body statement in cpp, got:\n%s", cpp)
	}
}

func TestTranspileHashMap(t *testing.T) {
	javaSrc := `
public class MapTest {
    public void run() {
        HashMap map = new HashMap();
        map.put("key", "val");
        Object v = map.get("key");
        int count = map.size();
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"HashMap.h\"") {
		t.Errorf("Expected #include \"HashMap.h\" in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "java::util::HashMap<java::lang::Object*, java::lang::Object*>* map = java_new<java::util::HashMap<java::lang::Object*, java::lang::Object*>>();") {
		t.Errorf("Expected HashMap creation in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(map)->put(java_new<java::lang::String>(\"key\"), java_new<java::lang::String>(\"val\"));") {
		t.Errorf("Expected put call in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(map)->get(java_new<java::lang::String>(\"key\"));") {
		t.Errorf("Expected get call in cpp, got:\n%s", cpp)
	}
}

func TestTranspileArrayList(t *testing.T) {
	javaSrc := `
public class ListTest {
    public void run() {
        ArrayList list = new ArrayList();
        list.add("Item1");
        int count = list.size();
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"ArrayList.h\"") {
		t.Errorf("Expected #include \"ArrayList.h\" in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "java::util::ArrayList<java::lang::Object*>* list = java_new<java::util::ArrayList<java::lang::Object*>>();") {
		t.Errorf("Expected ArrayList creation, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(list)->add(java_new<java::lang::String>(\"Item1\"));") {
		t.Errorf("Expected add method call, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(list)->size()") {
		t.Errorf("Expected size method call, got:\n%s", cpp)
	}
}

func TestTranspilePrimitiveWrappers(t *testing.T) {
	javaSrc := `
public class BoxTest {
    public void run() {
        Integer i = new Integer(42);
        Boolean b = new Boolean(true);
        Long l = new Long(100);
        int unboxed = i.intValue();
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"Boxes.h\"") {
		t.Errorf("Expected #include \"Boxes.h\" in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "java::lang::Integer* i = java_new<java::lang::Integer>(42);") {
		t.Errorf("Expected Integer creation, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "java::lang::Boolean* b = java_new<java::lang::Boolean>(true);") {
		t.Errorf("Expected Boolean creation, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "java::lang::Long* l = java_new<java::lang::Long>(100);") {
		t.Errorf("Expected Long creation, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(i)->intValue()") {
		t.Errorf("Expected intValue invocation, got:\n%s", cpp)
	}
}

func TestTranspileStaticMembers(t *testing.T) {
	javaSrc := `
public class MathUtil {
    public static int CONSTANT = 100;

    public static int multiply(int a, int b) {
        return a * b;
    }

    public void test() {
        int x = MathUtil.CONSTANT;
        int y = MathUtil.multiply(2, 3);
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "inline static int32_t CONSTANT = 100;") {
		t.Errorf("Expected inline static field in header, got:\n%s", header)
	}
	if !strings.Contains(header, "static int32_t multiply(int32_t a, int32_t b);") {
		t.Errorf("Expected static method declaration in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "MathUtil::CONSTANT") {
		t.Errorf("Expected MathUtil::CONSTANT reference in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "MathUtil::multiply(2, 3)") {
		t.Errorf("Expected MathUtil::multiply invocation in cpp, got:\n%s", cpp)
	}
}

func TestTranspileMathOperations(t *testing.T) {
	javaSrc := `
public class MathTest {
    public void run() {
        int maximum = Math.max(10, 20);
        int absolute = Math.abs(-5);
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"Math.h\"") {
		t.Errorf("Expected #include \"Math.h\" in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "java::lang::Math::max(10, 20)") {
		t.Errorf("Expected java::lang::Math::max invocation, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "java::lang::Math::abs(-5)") {
		t.Errorf("Expected java::lang::Math::abs invocation, got:\n%s", cpp)
	}
}

func TestTranspileStaticInitBlock(t *testing.T) {
	javaSrc := `
public class StaticInitTest {
    public static int VAL;

    static {
        VAL = 500;
    }
}
`
	_, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(cpp, "static inline bool _static_init_StaticInitTest_1 = []() {") {
		t.Errorf("Expected static initialization lambda in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "return true;\n}();") {
		t.Errorf("Expected static lambda self-invocation in cpp, got:\n%s", cpp)
	}
}

func TestTranspileStringBuilder(t *testing.T) {
	javaSrc := `
public class BuilderTest {
    public String build() {
        StringBuilder sb = new StringBuilder();
        sb.append("Hello");
        return sb.toString();
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"StringBuilder.h\"") {
		t.Errorf("Expected #include \"StringBuilder.h\" in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "java::lang::StringBuilder* sb = java_new<java::lang::StringBuilder>();") {
		t.Errorf("Expected StringBuilder allocation in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(sb)->append(java_new<java::lang::String>(\"Hello\"));") {
		t.Errorf("Expected append call in cpp, got:\n%s", cpp)
	}
}

func TestTranspileInheritance(t *testing.T) {
	javaSrc := `
public class Dog extends Animal {
    public void bark() {
        System.out.println("Woof");
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"Animal.h\"") {
		t.Errorf("Expected #include \"Animal.h\" in header, got:\n%s", header)
	}
	if !strings.Contains(header, "class Dog : public Animal") {
		t.Errorf("Expected class Dog : public Animal, got:\n%s", header)
	}
	if !strings.Contains(cpp, "void Dog::bark() {") {
		t.Errorf("Expected bark method in cpp, got:\n%s", cpp)
	}
}

func TestTranspileInterfaceAndImplements(t *testing.T) {
	javaIface := `
public interface Greeter {
    void greet(String name);
}
`
	header, _ := transpileJavaSnippet(t, javaIface)
	if !strings.Contains(header, "class Greeter {") {
		t.Errorf("Expected class Greeter in header, got:\n%s", header)
	}
	if !strings.Contains(header, "virtual ~Greeter() = default;") {
		t.Errorf("Expected virtual destructor in interface header, got:\n%s", header)
	}
	if !strings.Contains(header, "virtual void greet(java::lang::String* name) = 0;") {
		t.Errorf("Expected pure virtual greet method in interface header, got:\n%s", header)
	}

	javaClass := `
public class FriendlyGreeter extends Object implements Greeter {
    public void greet(String name) {
        System.out.println(name);
    }
}
`
	classHeader, classCpp := transpileJavaSnippet(t, javaClass)
	if !strings.Contains(classHeader, "#include \"Greeter.h\"") {
		t.Errorf("Expected #include \"Greeter.h\" in class header, got:\n%s", classHeader)
	}
	if !strings.Contains(classHeader, "class FriendlyGreeter : public java::lang::Object, public Greeter") {
		t.Errorf("Expected multiple inheritance in class header, got:\n%s", classHeader)
	}
	if !strings.Contains(classCpp, "void FriendlyGreeter::greet(java::lang::String* name) {") {
		t.Errorf("Expected greet method implementation in cpp, got:\n%s", classCpp)
	}
}

func TestTranspileArrayOperations(t *testing.T) {
	javaSrc := `
public class ArrayTest {
    public void run() {
        int[] numbers = new int[5];
        int first = numbers[0];
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"Array.h\"") {
		t.Errorf("Expected #include \"Array.h\" in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "java::lang::Array<int32_t>* numbers = java_new<java::lang::Array<int32_t>>(5);") {
		t.Errorf("Expected array creation in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "(*JAVA_NULL_CHECK(numbers))[0]") {
		t.Errorf("Expected array bracket indexing in cpp, got:\n%s", cpp)
	}
}

func TestTranspileConstructor(t *testing.T) {
	javaSrc := `
public class Person {
    String name;
    int age;

    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "Person(java::lang::String* name, int32_t age);") {
		t.Errorf("Expected constructor in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "Person::Person(java::lang::String* name, int32_t age) {") {
		t.Errorf("Expected constructor in cpp, got:\n%s", cpp)
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

func TestTranspileSystemOutPrintln(t *testing.T) {
	javaSrc := `
public class PrintTest {
    public void run() {
        System.out.println("Hello World");
    }
}
`
	header, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(header, "#include \"System.h\"") {
		t.Errorf("Expected System.h include in header, got:\n%s", header)
	}
	if !strings.Contains(cpp, "JAVA_NULL_CHECK(java::lang::System::out)->println(java_new<java::lang::String>(\"Hello World\"));") {
		t.Errorf("Expected java::lang::System::out println call, got:\n%s", cpp)
	}
}

func TestTranspileControlFlow(t *testing.T) {
	javaSrc := `
public class FlowTest {
    public void run(int val) {
        if (val > 0) {
            val = val - 1;
        } else {
            val = 0;
        }

        while (val < 10) {
            val = val + 1;
        }

        for (int i = 0; i < 5; i++) {
            val = val + i;
        }
    }
}
`
	_, cpp := transpileJavaSnippet(t, javaSrc)

	if !strings.Contains(cpp, "if (val > 0) {") {
		t.Errorf("Expected if block in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "else {") {
		t.Errorf("Expected else block in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "while (val < 10) {") {
		t.Errorf("Expected while block in cpp, got:\n%s", cpp)
	}
	if !strings.Contains(cpp, "for (int32_t i = 0; i < 5; i++) {") {
		t.Errorf("Expected for loop in cpp, got:\n%s", cpp)
	}
}
