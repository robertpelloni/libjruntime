package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/user/j2cpp-engine/parser"
)

type CppTranspilerVisitor struct {
	*parser.BaseJavaParserVisitor
	outDir          string
	currentClass    string
	currentSuper    string
	headerCode      strings.Builder
	cppCode         strings.Builder
	included        map[string]bool
	currentIsStatic bool
	staticInitCount int
	staticFields    map[string]bool
}

func NewCppTranspilerVisitor(outDir string) *CppTranspilerVisitor {
	return &CppTranspilerVisitor{
		BaseJavaParserVisitor: &parser.BaseJavaParserVisitor{},
		outDir:                outDir,
		included:              make(map[string]bool),
		staticFields:          make(map[string]bool),
	}
}

func mapType(javaType string) string {
	if strings.HasSuffix(javaType, "[]") {
		elemType := mapType(strings.TrimSuffix(javaType, "[]"))
		return fmt.Sprintf("java::lang::Array<%s>*", elemType)
	}
	if strings.HasPrefix(javaType, "ArrayList<") && strings.HasSuffix(javaType, ">") {
		innerType := javaType[10 : len(javaType)-1]
		mappedInner := mapType(innerType)
		return fmt.Sprintf("java::util::ArrayList<%s>*", mappedInner)
	}
	if strings.HasPrefix(javaType, "HashMap<") && strings.HasSuffix(javaType, ">") {
		innerTypes := javaType[8 : len(javaType)-1]
		parts := strings.Split(innerTypes, ",")
		if len(parts) == 2 {
			k := mapType(strings.TrimSpace(parts[0]))
			v := mapType(strings.TrimSpace(parts[1]))
			return fmt.Sprintf("java::util::HashMap<%s, %s>*", k, v)
		}
		return "java::util::HashMap<java::lang::Object*, java::lang::Object*>*"
	}
	switch javaType {
	case "int":
		return "int32_t"
	case "boolean":
		return "bool"
	case "long":
		return "int64_t"
	case "String":
		return "java::lang::String*"
	case "StringBuilder":
		return "java::lang::StringBuilder*"
	case "Integer":
		return "java::lang::Integer*"
	case "Boolean":
		return "java::lang::Boolean*"
	case "Long":
		return "java::lang::Long*"
	case "ArrayList":
		return "java::util::ArrayList<java::lang::Object*>*"
	case "HashMap":
		return "java::util::HashMap<java::lang::Object*, java::lang::Object*>*"
	case "Object":
		return "java::lang::Object*"
	case "void":
		return "void"
	case "Exception":
	    return "java::lang::Exception*"
	default:
		return javaType + "*" // default to pointer for reference types
	}
}

func (v *CppTranspilerVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *CppTranspilerVisitor) VisitCompilationUnit(ctx *parser.CompilationUnitContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitTypeDeclaration(ctx *parser.TypeDeclarationContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitClassDeclaration(ctx *parser.ClassDeclarationContext) interface{} {
	className := ctx.Identifier().GetText()
	v.currentClass = className

    fmt.Printf("Discovered class: %s\n", className)

	v.headerCode.Reset()
	v.cppCode.Reset()
	v.included = make(map[string]bool)
	v.staticFields = make(map[string]bool)
	v.staticInitCount = 0

	superClass := "java::lang::Object"
	var extraIncludes []string
	if ctx.EXTENDS() != nil && ctx.TypeType() != nil {
		superClass = ctx.TypeType().GetText()
		if superClass == "Object" {
			superClass = "java::lang::Object"
		} else if superClass != "java::lang::Object" {
			extraIncludes = append(extraIncludes, fmt.Sprintf("#include \"%s.h\"\n", superClass))
			v.included[superClass] = true
		}
	}
	v.currentSuper = superClass

	var baseClasses []string
	baseClasses = append(baseClasses, fmt.Sprintf("public %s", superClass))

	if ctx.IMPLEMENTS() != nil {
		for _, typeList := range ctx.AllTypeList() {
			for _, tt := range typeList.AllTypeType() {
				ifaceName := tt.GetText()
				baseClasses = append(baseClasses, fmt.Sprintf("public %s", ifaceName))
				extraIncludes = append(extraIncludes, fmt.Sprintf("#include \"%s.h\"\n", ifaceName))
				v.included[ifaceName] = true
			}
		}
	}

	extraIncludeStr := strings.Join(extraIncludes, "")
	baseClassesStr := strings.Join(baseClasses, ", ")

	var bodyCode strings.Builder
	// swap headerCode to temporary builder for class body
	origHeader := v.headerCode
	v.headerCode = bodyCode

	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}

	classBodyStr := v.headerCode.String()
	v.headerCode = origHeader

	primitives := map[string]bool{
		"int": true, "boolean": true, "long": true, "byte": true,
		"short": true, "char": true, "float": true, "double": true,
		"void": true, "String": true, "StringBuilder": true, "Math": true, "Integer": true, "Boolean": true, "Long": true, "ArrayList": true, "HashMap": true, "Object": true, "Exception": true,
	}

	for inc := range v.included {
		cleanInc := strings.TrimSuffix(inc, "[]")
		if cleanInc != className && !primitives[cleanInc] && !strings.Contains(cleanInc, "[") && !strings.Contains(cleanInc, "<") {
			extraIncludes = append(extraIncludes, fmt.Sprintf("#include \"%s.h\"\n", cleanInc))
		}
	}

	extraIncludeStr = strings.Join(extraIncludes, "")
	v.headerCode.WriteString(fmt.Sprintf("#pragma once\n\n#include <cstdint>\n#include \"Object.h\"\n#include \"String.h\"\n#include \"GC.h\"\n#include \"Runtime.h\"\n#include \"FinallyGuard.h\"\n#include \"System.h\"\n#include \"Array.h\"\n#include \"ArrayList.h\"\n#include \"HashMap.h\"\n#include \"StringBuilder.h\"\n#include \"Math.h\"\n#include \"Boxes.h\"\n%s\nclass %s : %s {\npublic:\n", extraIncludeStr, className, baseClassesStr))
	v.headerCode.WriteString(classBodyStr)
	v.headerCode.WriteString("};\n")

	var cppFinal strings.Builder
	cppFinal.WriteString(fmt.Sprintf("#include \"%s.h\"\n", className))
	for inc := range v.included {
		cleanInc := strings.TrimSuffix(inc, "[]")
		if cleanInc != className && !primitives[cleanInc] && !strings.Contains(cleanInc, "[") && !strings.Contains(cleanInc, "<") {
			cppFinal.WriteString(fmt.Sprintf("#include \"%s.h\"\n", cleanInc))
		}
	}
	cppFinal.WriteString("\n")
	cppFinal.WriteString(v.cppCode.String())

	err := os.WriteFile(filepath.Join(v.outDir, className+".h"), []byte(v.headerCode.String()), 0644)
	if err != nil {
		fmt.Printf("Failed to write header file: %v\n", err)
	}
	err = os.WriteFile(filepath.Join(v.outDir, className+".cpp"), []byte(cppFinal.String()), 0644)
	if err != nil {
		fmt.Printf("Failed to write cpp file: %v\n", err)
	}

	return nil
}


func (v *CppTranspilerVisitor) VisitInterfaceDeclaration(ctx *parser.InterfaceDeclarationContext) interface{} {
	ifaceName := ctx.Identifier().GetText()
	v.currentClass = ifaceName

	fmt.Printf("Discovered interface: %s\n", ifaceName)

	v.headerCode.Reset()
	v.cppCode.Reset()
	v.included = make(map[string]bool)

	v.headerCode.WriteString(fmt.Sprintf("#pragma once\n\n#include <cstdint>\n#include \"Object.h\"\n#include \"String.h\"\n#include \"GC.h\"\n#include \"Runtime.h\"\n#include \"FinallyGuard.h\"\n#include \"System.h\"\n#include \"Array.h\"\n#include \"ArrayList.h\"\n#include \"HashMap.h\"\n#include \"StringBuilder.h\"\n#include \"Math.h\"\n#include \"Boxes.h\"\n\nclass %s {\npublic:\n    virtual ~%s() = default;\n", ifaceName, ifaceName))

	if ctx.InterfaceBody() != nil {
		v.Visit(ctx.InterfaceBody())
	}

	v.headerCode.WriteString("};\n")

	err := os.WriteFile(filepath.Join(v.outDir, ifaceName+".h"), []byte(v.headerCode.String()), 0644)
	if err != nil {
		fmt.Printf("Failed to write interface header file: %v\n", err)
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitInterfaceBody(ctx *parser.InterfaceBodyContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitInterfaceBodyDeclaration(ctx *parser.InterfaceBodyDeclarationContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitInterfaceMemberDeclaration(ctx *parser.InterfaceMemberDeclarationContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitInterfaceMethodDeclaration(ctx *parser.InterfaceMethodDeclarationContext) interface{} {
	cb := ctx.InterfaceCommonBodyDeclaration()
	if cb == nil {
		return nil
	}

	methodName := cb.Identifier().GetText()
	returnType := "void"
	if cb.TypeTypeOrVoid() != nil && cb.TypeTypeOrVoid().TypeType() != nil {
		returnType = mapType(cb.TypeTypeOrVoid().TypeType().GetText())
	}

	var pList []string
	if cb.FormalParameters() != nil {
		if fp := cb.FormalParameters().FormalParameter(); fp != nil {
			pType := mapType(fp.TypeType().GetText())
			pName := fp.VariableDeclaratorId().GetText()
			pList = append(pList, fmt.Sprintf("%s %s", pType, pName))
		}
		for _, fpList := range cb.FormalParameters().AllFormalParameterList() {
			for _, p := range fpList.AllFormalParameter() {
				pType := mapType(p.TypeType().GetText())
				pName := p.VariableDeclaratorId().GetText()
				pList = append(pList, fmt.Sprintf("%s %s", pType, pName))
			}
		}
	}
	params := strings.Join(pList, ", ")

	v.headerCode.WriteString(fmt.Sprintf("    virtual %s %s(%s) = 0;\n", returnType, methodName, params))
	return nil
}

func (v *CppTranspilerVisitor) VisitClassBody(ctx *parser.ClassBodyContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitClassBodyDeclaration(ctx *parser.ClassBodyDeclarationContext) interface{} {
	v.currentIsStatic = false
	if ctx.STATIC() != nil {
		v.currentIsStatic = true
	}
	for _, mod := range ctx.AllModifier() {
		if mod.ClassOrInterfaceModifier() != nil && mod.ClassOrInterfaceModifier().STATIC() != nil {
			v.currentIsStatic = true
		}
	}

	// Check for static initialization block: static { ... }
	if ctx.STATIC() != nil && ctx.Block() != nil {
		v.staticInitCount++
		v.cppCode.WriteString(fmt.Sprintf("static inline bool _static_init_%s_%d = []() {\n", v.currentClass, v.staticInitCount))
		v.Visit(ctx.Block())
		v.cppCode.WriteString("    return true;\n}();\n\n")
		v.currentIsStatic = false
		return nil
	}

	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	v.currentIsStatic = false
	return nil
}

func (v *CppTranspilerVisitor) VisitMemberDeclaration(ctx *parser.MemberDeclarationContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}


func (v *CppTranspilerVisitor) VisitConstructorDeclaration(ctx *parser.ConstructorDeclarationContext) interface{} {
	constructorName := ctx.Identifier().GetText()
	fmt.Printf("Discovered constructor: %s\n", constructorName)

	var pList []string
	if ctx.FormalParameters() != nil {
		if fp := ctx.FormalParameters().FormalParameter(); fp != nil {
			pType := mapType(fp.TypeType().GetText())
			pName := fp.VariableDeclaratorId().GetText()
			pList = append(pList, fmt.Sprintf("%s %s", pType, pName))
		}
		for _, fpList := range ctx.FormalParameters().AllFormalParameterList() {
			for _, p := range fpList.AllFormalParameter() {
				pType := mapType(p.TypeType().GetText())
				pName := p.VariableDeclaratorId().GetText()
				pList = append(pList, fmt.Sprintf("%s %s", pType, pName))
			}
		}
	}
	params := strings.Join(pList, ", ")

	body := ctx.GetConstructorBody()
	if body == nil {
		body = ctx.Block()
	}

	initListStr := ""
	skipFirstStatement := false

	if body != nil && len(body.AllBlockStatement()) > 0 {
		firstStmt := body.AllBlockStatement()[0]
		if firstStmt.Statement() != nil && firstStmt.Statement().GetStatementExpression() != nil {
			expr := firstStmt.Statement().GetStatementExpression()
			if mcExpr, ok := expr.(*parser.MethodCallExpressionContext); ok && mcExpr.MethodCall() != nil {
				mc := mcExpr.MethodCall()
				var callTarget string
				if mc.SUPER() != nil {
					callTarget = v.currentSuper
				} else if mc.THIS() != nil {
					callTarget = v.currentClass
				}
				if callTarget != "" {
					argsStr := ""
					if mc.Arguments() != nil && mc.Arguments().(*parser.ArgumentsContext).ExpressionList() != nil {
						exprList := mc.Arguments().(*parser.ArgumentsContext).ExpressionList().(*parser.ExpressionListContext)
						var parts []string
						for _, a := range exprList.AllExpression() {
							res := v.Visit(a)
							if res != nil {
								parts = append(parts, fmt.Sprintf("%v", res))
							} else {
								parts = append(parts, a.GetText())
							}
						}
						argsStr = strings.Join(parts, ", ")
					}
					initListStr = fmt.Sprintf(" : %s(%s)", callTarget, argsStr)
					skipFirstStatement = true
				}
			}
		}
	}

	v.headerCode.WriteString(fmt.Sprintf("    %s(%s);\n", constructorName, params))
	v.cppCode.WriteString(fmt.Sprintf("%s::%s(%s)%s {\n", v.currentClass, constructorName, params, initListStr))

	if body != nil {
		for i, bs := range body.AllBlockStatement() {
			if i == 0 && skipFirstStatement {
				continue
			}
			v.Visit(bs)
		}
	}

	v.cppCode.WriteString("}\n\n")
	return nil
}

func (v *CppTranspilerVisitor) VisitMethodDeclaration(ctx *parser.MethodDeclarationContext) interface{} {
	methodName := ctx.Identifier().GetText()

	fmt.Printf("Discovered method: %s\n", methodName)

	returnType := "void"
	if ctx.TypeTypeOrVoid().TypeType() != nil {
		returnType = mapType(ctx.TypeTypeOrVoid().TypeType().GetText())
	}

	var pList []string
	if ctx.FormalParameters() != nil {
		if fp := ctx.FormalParameters().FormalParameter(); fp != nil {
			pType := mapType(fp.TypeType().GetText())
			pName := fp.VariableDeclaratorId().GetText()
			pList = append(pList, fmt.Sprintf("%s %s", pType, pName))
		}
		for _, fpList := range ctx.FormalParameters().AllFormalParameterList() {
			for _, p := range fpList.AllFormalParameter() {
				pType := mapType(p.TypeType().GetText())
				pName := p.VariableDeclaratorId().GetText()
				pList = append(pList, fmt.Sprintf("%s %s", pType, pName))
			}
		}
	}
	params := strings.Join(pList, ", ")

	if v.currentIsStatic {
		v.headerCode.WriteString(fmt.Sprintf("    static %s %s(%s);\n", returnType, methodName, params))
	} else {
		v.headerCode.WriteString(fmt.Sprintf("    virtual %s %s(%s);\n", returnType, methodName, params))
	}

	v.cppCode.WriteString(fmt.Sprintf("%s %s::%s(%s) {\n", returnType, v.currentClass, methodName, params))

	if ctx.MethodBody().Block() != nil {
		v.Visit(ctx.MethodBody().Block())
	}

	v.cppCode.WriteString("}\n\n")

	return nil
}

func (v *CppTranspilerVisitor) VisitFieldDeclaration(ctx *parser.FieldDeclarationContext) interface{} {
    fmt.Printf("Discovered field declaration\n")

	fieldType := mapType(ctx.TypeType().GetText())
	for _, decl := range ctx.VariableDeclarators().AllVariableDeclarator() {
		fieldName := decl.VariableDeclaratorId().GetText()
		initVal := ""
		if decl.VariableInitializer() != nil {
			res := v.Visit(decl.VariableInitializer())
			if res != nil {
				initVal = fmt.Sprintf(" = %v", res)
			} else {
				initVal = fmt.Sprintf(" = %s", decl.VariableInitializer().GetText())
			}
		}
		if v.currentIsStatic {
			v.staticFields[fieldName] = true
			v.headerCode.WriteString(fmt.Sprintf("    inline static %s %s%s;\n", fieldType, fieldName, initVal))
		} else {
			v.headerCode.WriteString(fmt.Sprintf("    %s %s%s;\n", fieldType, fieldName, initVal))
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitBlock(ctx *parser.BlockContext) interface{} {
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitBlockStatement(ctx *parser.BlockStatementContext) interface{} {
    if ctx.LocalVariableDeclaration() != nil {
        decl := ctx.LocalVariableDeclaration()
        typ := "auto"
        if decl.TypeType() != nil {
            rawType := decl.TypeType().GetText()
            v.included[rawType] = true
            typ = mapType(rawType)
        }
        for _, vd := range decl.VariableDeclarators().AllVariableDeclarator() {
            name := vd.VariableDeclaratorId().GetText()
            v.cppCode.WriteString(fmt.Sprintf("    %s %s", typ, name))
            if vd.VariableInitializer() != nil {
                v.cppCode.WriteString(" = ")
                res := v.Visit(vd.VariableInitializer())
                if res != nil {
                    v.cppCode.WriteString(fmt.Sprintf("%v", res))
                }
            }
            v.cppCode.WriteString(";\n")
        }
        return nil
    }


    for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
	return nil
}

func (v *CppTranspilerVisitor) VisitVariableInitializer(ctx *parser.VariableInitializerContext) interface{} {
    if ctx.Expression() != nil {
        return v.Visit(ctx.Expression())
    }
    return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitStatement(ctx *parser.StatementContext) interface{} {
    // Check for try block
    if ctx.TRY() != nil {
        v.cppCode.WriteString("    try {\n")
        if ctx.Block() != nil {
            v.Visit(ctx.Block()) // try block
        }
        v.cppCode.WriteString("    }\n")


        catchIdx := 0
        for ctx.CatchClause(catchIdx) != nil {
            v.cppCode.WriteString("    catch (")
            if ctx.CatchClause(catchIdx).CatchType() != nil {
                // assume one type for now
                catchType := mapType(ctx.CatchClause(catchIdx).CatchType().GetText())
                catchVar := ctx.CatchClause(catchIdx).Identifier().GetText()
                v.cppCode.WriteString(fmt.Sprintf("%s %s", catchType, catchVar))
            } else {
                 v.cppCode.WriteString("java::lang::Exception* e")
            }
            v.cppCode.WriteString(") {\n")
            v.Visit(ctx.CatchClause(catchIdx).Block())
            v.cppCode.WriteString("    }\n")
            catchIdx++
        }


        if ctx.FinallyBlock() != nil {
            v.cppCode.WriteString("    JAVA_FINALLY(\n")
            v.Visit(ctx.FinallyBlock().Block()) // finally block
            v.cppCode.WriteString("    );\n")
        }
        return nil
    }

    // IF statement
    if ctx.IF() != nil {
        v.cppCode.WriteString("    if (")
        if ctx.Expression(0) != nil {
            res := v.Visit(ctx.Expression(0))
            if res != nil {
                v.cppCode.WriteString(fmt.Sprintf("%v", res))
            } else {
                v.cppCode.WriteString(ctx.Expression(0).GetText())
            }
        }
        v.cppCode.WriteString(") {\n")
        if ctx.Statement(0) != nil {
            if ctx.Statement(0).Block() != nil {
                v.Visit(ctx.Statement(0).Block())
            } else {
                v.Visit(ctx.Statement(0))
            }
        }
        v.cppCode.WriteString("    }")
        if ctx.ELSE() != nil && ctx.Statement(1) != nil {
            v.cppCode.WriteString(" else {\n")
            if ctx.Statement(1).Block() != nil {
                v.Visit(ctx.Statement(1).Block())
            } else {
                v.Visit(ctx.Statement(1))
            }
            v.cppCode.WriteString("    }")
        }
        v.cppCode.WriteString("\n")
        return nil
    }

    // WHILE statement
    if ctx.WHILE() != nil && ctx.DO() == nil {
        v.cppCode.WriteString("    while (")
        if ctx.Expression(0) != nil {
            res := v.Visit(ctx.Expression(0))
            if res != nil {
                v.cppCode.WriteString(fmt.Sprintf("%v", res))
            } else {
                v.cppCode.WriteString(ctx.Expression(0).GetText())
            }
        }
        v.cppCode.WriteString(") {\n")
        if ctx.Statement(0) != nil {
            if ctx.Statement(0).Block() != nil {
                v.Visit(ctx.Statement(0).Block())
            } else {
                v.Visit(ctx.Statement(0))
            }
        }
        v.cppCode.WriteString("    }\n")
        return nil
    }

    // FOR statement
    if ctx.FOR() != nil && ctx.ForControl() != nil {
        fc := ctx.ForControl()

        // Enhanced for-each loop: for (Type varName : iterableExpr)
        if fc.EnhancedForControl() != nil {
            efc := fc.EnhancedForControl()
            elemVar := ""
            if efc.VariableDeclaratorId() != nil {
                elemVar = efc.VariableDeclaratorId().GetText()
            }
            elemType := "auto"
            if efc.TypeType() != nil {
                elemType = mapType(efc.TypeType().GetText())
                v.included[efc.TypeType().GetText()] = true
            }

            iterableStr := ""
            if efc.Expression() != nil {
                res := v.Visit(efc.Expression())
                if res != nil {
                    iterableStr = fmt.Sprintf("%v", res)
                } else {
                    iterableStr = efc.Expression().GetText()
                }
            }

            v.cppCode.WriteString(fmt.Sprintf("    for (%s %s : *JAVA_NULL_CHECK(%s)) {\n", elemType, elemVar, iterableStr))
            if ctx.Statement(0) != nil {
                if ctx.Statement(0).Block() != nil {
                    v.Visit(ctx.Statement(0).Block())
                } else {
                    v.Visit(ctx.Statement(0))
                }
            }
            v.cppCode.WriteString("    }\n")
            return nil
        }

        // Standard 3-part for loop: for (init; cond; update)
        initStr := ""
        condStr := ""
        updStr := ""
        if fc.ForInit() != nil {
            initStr = fc.ForInit().GetText()
            if fc.ForInit().LocalVariableDeclaration() != nil {
                decl := fc.ForInit().LocalVariableDeclaration()
                typ := mapType(decl.TypeType().GetText())
                var inits []string
                for _, vd := range decl.VariableDeclarators().AllVariableDeclarator() {
                    name := vd.VariableDeclaratorId().GetText()
                    val := ""
                    if vd.VariableInitializer() != nil {
                        val = fmt.Sprintf(" = %v", vd.VariableInitializer().GetText())
                    }
                    inits = append(inits, fmt.Sprintf("%s %s%s", typ, name, val))
                }
                initStr = strings.Join(inits, ", ")
            }
        }
        if fc.Expression() != nil {
            res := v.Visit(fc.Expression())
            if res != nil {
                condStr = fmt.Sprintf("%v", res)
            } else {
                condStr = fc.Expression().GetText()
            }
        }
        if fc.GetForUpdate() != nil {
            updStr = fc.GetForUpdate().GetText()
        }
        v.cppCode.WriteString(fmt.Sprintf("    for (%s; %s; %s) {\n", initStr, condStr, updStr))
        if ctx.Statement(0) != nil {
            if ctx.Statement(0).Block() != nil {
                v.Visit(ctx.Statement(0).Block())
            } else {
                v.Visit(ctx.Statement(0))
            }
        }
        v.cppCode.WriteString("    }\n")
        return nil
    }

    // SWITCH statement
    if ctx.SWITCH() != nil {
        exprStr := ""
        if ctx.Expression(0) != nil {
            res := v.Visit(ctx.Expression(0))
            if res != nil {
                exprStr = fmt.Sprintf("%v", res)
            } else {
                exprStr = ctx.Expression(0).GetText()
            }
        }
        v.cppCode.WriteString(fmt.Sprintf("    switch (%s) {\n", exprStr))
        for _, group := range ctx.AllSwitchBlockStatementGroup() {
            for _, label := range group.AllSwitchLabel() {
                if label.DEFAULT() != nil {
                    v.cppCode.WriteString("    default:\n")
                } else if label.CASE() != nil {
                    caseExpr := ""
                    if label.Expression() != nil {
                        res := v.Visit(label.Expression())
                        if res != nil {
                            caseExpr = fmt.Sprintf("%v", res)
                        } else {
                            caseExpr = label.Expression().GetText()
                        }
                    }
                    v.cppCode.WriteString(fmt.Sprintf("    case %s:\n", caseExpr))
                }
            }
            for _, bs := range group.AllBlockStatement() {
                v.Visit(bs)
            }
        }
        v.cppCode.WriteString("    }\n")
        return nil
    }

    if ctx.BREAK() != nil {
        v.cppCode.WriteString("    break;\n")
        return nil
    }

    if ctx.CONTINUE() != nil {
        v.cppCode.WriteString("    continue;\n")
        return nil
    }

    // THROW statement
    if ctx.THROW() != nil {
        exprStr := ""
        if ctx.Expression(0) != nil {
            res := v.Visit(ctx.Expression(0))
            if res != nil {
                exprStr = fmt.Sprintf("%v", res)
            } else {
                exprStr = ctx.Expression(0).GetText()
            }
        }
        v.cppCode.WriteString(fmt.Sprintf("    throw %s;\n", exprStr))
        return nil
    }

    if ctx.RETURN() != nil {
        v.cppCode.WriteString("    return")
        if ctx.Expression(0) != nil {
            v.cppCode.WriteString(" ")
            res := v.Visit(ctx.Expression(0))
            if res != nil {
                v.cppCode.WriteString(fmt.Sprintf("%v", res))
            } else {
                v.cppCode.WriteString(ctx.Expression(0).GetText())
            }
        }
        v.cppCode.WriteString(";\n")
        return nil
    }

    if ctx.GetStatementExpression() != nil {
        res := v.Visit(ctx.GetStatementExpression())
        if res != nil {
            v.cppCode.WriteString(fmt.Sprintf("    %v;\n", res))
        }
        return nil
    }

    // Unhandled statements
    v.cppCode.WriteString("    // Unhandled statement: " + ctx.GetText() + "\n")
    return nil
}

// Expression AST handling

func (v *CppTranspilerVisitor) VisitObjectCreationExpression(ctx *parser.ObjectCreationExpressionContext) interface{} {
    if ctx.NEW() != nil && ctx.Creator() != nil {
        creator := ctx.Creator()
        if creator.CreatedName() != nil {
             rawClassName := creator.CreatedName().GetText()
             v.included[rawClassName] = true
             className := rawClassName
             // Check if mapped
             if className == "String" {
                 className = "java::lang::String"
             } else if className == "StringBuilder" {
                 className = "java::lang::StringBuilder"
             } else if className == "Integer" {
                 className = "java::lang::Integer"
             } else if className == "Boolean" {
                 className = "java::lang::Boolean"
             } else if className == "Long" {
                 className = "java::lang::Long"
             } else if className == "Exception" {
                 className = "java::lang::Exception"
             } else if className == "NullPointerException" {
                 className = "java::lang::NullPointerException"
             } else if className == "ArrayIndexOutOfBoundsException" {
                 className = "java::lang::ArrayIndexOutOfBoundsException"
             } else if className == "ArrayList" || className == "ArrayList<>" {
                 className = "java::util::ArrayList<java::lang::Object*>"
             } else if strings.HasPrefix(className, "ArrayList<") && strings.HasSuffix(className, ">") {
                 innerType := className[10 : len(className)-1]
                 if innerType == "" {
                     innerType = "Object"
                 }
                 mappedInner := mapType(innerType)
                 className = fmt.Sprintf("java::util::ArrayList<%s>", mappedInner)
             } else if className == "HashMap" || className == "HashMap<>" {
                 className = "java::util::HashMap<java::lang::Object*, java::lang::Object*>"
             } else if strings.HasPrefix(className, "HashMap<") && strings.HasSuffix(className, ">") {
                 innerTypes := className[8 : len(className)-1]
                 parts := strings.Split(innerTypes, ",")
                 if len(parts) == 2 {
                     k := mapType(strings.TrimSpace(parts[0]))
                     v := mapType(strings.TrimSpace(parts[1]))
                     className = fmt.Sprintf("java::util::HashMap<%s, %s>", k, v)
                 } else {
                     className = "java::util::HashMap<java::lang::Object*, java::lang::Object*>"
                 }
             }

             // Array creation: new T[size]
             if creator.ArrayCreatorRest() != nil {
                 acr := creator.ArrayCreatorRest()
                 sizeExprStr := "0"
                 if len(acr.AllExpression()) > 0 {
                     res := v.Visit(acr.AllExpression()[0])
                     if res != nil {
                         sizeExprStr = fmt.Sprintf("%v", res)
                     } else {
                         sizeExprStr = acr.AllExpression()[0].GetText()
                     }
                 }
                 elemType := mapType(rawClassName)
                 return fmt.Sprintf("java_new<java::lang::Array<%s>>(%s)", elemType, sizeExprStr)
             }

             argsStr := ""
             if creator.ClassCreatorRest() != nil && creator.ClassCreatorRest().Arguments() != nil {
                 args := creator.ClassCreatorRest().Arguments().(*parser.ArgumentsContext)
                 if args.ExpressionList() != nil {
                     exprList := args.ExpressionList().(*parser.ExpressionListContext)
                     var parts []string
                     for _, expr := range exprList.AllExpression() {
                         res := v.Visit(expr)
                         if res != nil {
                             parts = append(parts, fmt.Sprintf("%v", res))
                         } else {
                              parts = append(parts, expr.GetText())
                         }
                     }
                     argsStr = strings.Join(parts, ", ")
                 }
             }
             return fmt.Sprintf("java_new<%s>(%s)", className, argsStr)
        }
    }
    return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitMethodCallExpression(ctx *parser.MethodCallExpressionContext) interface{} {
    // obj.method() -> JAVA_NULL_CHECK(obj)->method()
    if ctx.GetChildCount() >= 3 && ctx.GetChild(1).GetPayload() != nil {
         if token, ok := ctx.GetChild(1).GetPayload().(*antlr.CommonToken); ok && token.GetText() == "." {
             objRes := v.Visit(ctx.GetChild(0).(antlr.ParseTree))
             if objRes == nil {
                  objRes = ctx.GetChild(0).(antlr.ParseTree).GetText()
             }


             if ctx.MethodCall() != nil {
                 methodCallRes := v.Visit(ctx.MethodCall())
                 if methodCallRes == nil {
                     methodCallRes = ctx.MethodCall().GetText()
                 }
                 objStr := fmt.Sprintf("%v", objRes)
                 if objStr == "Integer" || objStr == "Boolean" || objStr == "Long" {
                     return fmt.Sprintf("java::lang::%s::%v", objStr, methodCallRes)
                 }
                 if objStr == "Math" {
                     return fmt.Sprintf("java::lang::Math::%v", methodCallRes)
                 }
                 if len(objStr) > 0 && objStr[0] >= 'A' && objStr[0] <= 'Z' && !strings.Contains(objStr, "->") && !strings.Contains(objStr, "::") {
                     return fmt.Sprintf("%s::%v", objStr, methodCallRes)
                 }
                 return fmt.Sprintf("JAVA_NULL_CHECK(%v)->%v", objRes, methodCallRes)
             }
         }
    }
    // simple method call without object
    if ctx.GetChildCount() == 1 {
         if ctx.MethodCall() != nil {
             return v.Visit(ctx.MethodCall())
         }
    }


    return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitMemberReferenceExpression(ctx *parser.MemberReferenceExpressionContext) interface{} {
    // obj.field -> JAVA_NULL_CHECK(obj)->field
    if ctx.GetChildCount() >= 3 && ctx.GetChild(1).GetPayload() != nil {
         if token, ok := ctx.GetChild(1).GetPayload().(*antlr.CommonToken); ok && token.GetText() == "." {
             objRes := v.Visit(ctx.GetChild(0).(antlr.ParseTree))
             if objRes == nil {
                  objRes = ctx.GetChild(0).(antlr.ParseTree).GetText()
             }


             leftStr := fmt.Sprintf("%v", objRes)
             if leftStr == "System" && ctx.Identifier() != nil && ctx.Identifier().GetText() == "out" {
                 return "java::lang::System::out"
             }
             objStr := fmt.Sprintf("%v", objRes)
             if objStr == "Math" {
                 if ctx.Identifier() != nil {
                     return fmt.Sprintf("java::lang::Math::%v", ctx.Identifier().GetText())
                 } else if ctx.MethodCall() != nil {
                     methodCallRes := v.Visit(ctx.MethodCall())
                     if methodCallRes == nil {
                         methodCallRes = ctx.MethodCall().GetText()
                     }
                     return fmt.Sprintf("java::lang::Math::%v", methodCallRes)
                 }
             }
             isClass := len(objStr) > 0 && objStr[0] >= 'A' && objStr[0] <= 'Z' && !strings.Contains(objStr, "->") && !strings.Contains(objStr, "::")

             if ctx.Identifier() != nil {
                 if isClass {
                     return fmt.Sprintf("%s::%v", objStr, ctx.Identifier().GetText())
                 }
                 return fmt.Sprintf("JAVA_NULL_CHECK(%v)->%v", objRes, ctx.Identifier().GetText())
             } else if ctx.MethodCall() != nil {
                 methodCallRes := v.Visit(ctx.MethodCall())
                 if methodCallRes == nil {
                     methodCallRes = ctx.MethodCall().GetText()
                 }
                 if isClass {
                     return fmt.Sprintf("%s::%v", objStr, methodCallRes)
                 }
                 return fmt.Sprintf("JAVA_NULL_CHECK(%v)->%v", objRes, methodCallRes)
             }
         }
    }
    return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitPrimaryExpression(ctx *parser.PrimaryExpressionContext) interface{} {
    if ctx.Primary() != nil {
        return v.Visit(ctx.Primary())
    }
    return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitPrimary(ctx *parser.PrimaryContext) interface{} {
    if ctx.Identifier() != nil {
        idName := ctx.Identifier().GetText()
        if v.staticFields[idName] {
            return fmt.Sprintf("%s::%s", v.currentClass, idName)
        }
        return idName
    }
    if ctx.Expression() != nil {
        res := v.Visit(ctx.Expression())
        if res != nil {
            return fmt.Sprintf("(%v)", res)
        }
    }
    if ctx.GetChildCount() == 1 {
        child := ctx.GetChild(0)
        if payload, ok := child.(antlr.ParseTree); ok {
             res := v.Visit(payload)
             if res != nil {
                 return res
             }
        }
    }
    return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitMethodCall(ctx *parser.MethodCallContext) interface{} {
	methodName := ""
	if ctx.Identifier() != nil {
		methodName = ctx.Identifier().GetText()
	}

	argsStr := ""
	if ctx.Arguments() != nil && ctx.Arguments().ExpressionList() != nil {
		var parts []string
		for _, expr := range ctx.Arguments().ExpressionList().AllExpression() {
			res := v.Visit(expr)
			if res != nil {
				parts = append(parts, fmt.Sprintf("%v", res))
			} else {
				parts = append(parts, expr.GetText())
			}
		}
		argsStr = strings.Join(parts, ", ")
	}

	return fmt.Sprintf("%s(%s)", methodName, argsStr)
}

func (v *CppTranspilerVisitor) VisitBinaryOperatorExpression(ctx *parser.BinaryOperatorExpressionContext) interface{} {
	exprs := ctx.AllExpression()
	if len(exprs) == 2 {
		left := v.Visit(exprs[0])
		right := v.Visit(exprs[1])
		op := "+"
		if ctx.GetBop() != nil {
			op = ctx.GetBop().GetText()
		}
		return fmt.Sprintf("%v %s %v", left, op, right)
	}
	return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitInstanceOfOperatorExpression(ctx *parser.InstanceOfOperatorExpressionContext) interface{} {
	exprStr := ""
	if ctx.Expression() != nil {
		res := v.Visit(ctx.Expression())
		if res != nil {
			exprStr = fmt.Sprintf("%v", res)
		} else {
			exprStr = ctx.Expression().GetText()
		}
	}

	targetType := "java::lang::Object*"
	if ctx.TypeType() != nil {
		rawType := ctx.TypeType().GetText()
		v.included[rawType] = true
		targetType = mapType(rawType)
	}

	return fmt.Sprintf("(dynamic_cast<%s>(%s) != nullptr)", targetType, exprStr)
}

func (v *CppTranspilerVisitor) VisitCastExpression(ctx *parser.CastExpressionContext) interface{} {
	exprStr := ""
	if ctx.Expression() != nil {
		res := v.Visit(ctx.Expression())
		if res != nil {
			exprStr = fmt.Sprintf("%v", res)
		} else {
			exprStr = ctx.Expression().GetText()
		}
	}

	if len(ctx.AllTypeType()) > 0 {
		rawType := ctx.AllTypeType()[0].GetText()
		v.included[rawType] = true
		mapped := mapType(rawType)
		if strings.HasSuffix(mapped, "*") {
			return fmt.Sprintf("dynamic_cast<%s>(%s)", mapped, exprStr)
		}
		return fmt.Sprintf("static_cast<%s>(%s)", mapped, exprStr)
	}
	return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitPostIncrementDecrementOperatorExpression(ctx *parser.PostIncrementDecrementOperatorExpressionContext) interface{} {
	exprStr := ctx.Expression().GetText()
	res := v.Visit(ctx.Expression())
	if res != nil {
		exprStr = fmt.Sprintf("%v", res)
	}
	op := "++"
	if ctx.GetPostfix() != nil {
		op = ctx.GetPostfix().GetText()
	}
	return fmt.Sprintf("%s%s", exprStr, op)
}

func (v *CppTranspilerVisitor) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	if ctx.STRING_LITERAL() != nil || ctx.TEXT_BLOCK() != nil {
		return fmt.Sprintf("java_new<java::lang::String>(%s)", ctx.GetText())
	}
	if ctx.NULL_LITERAL() != nil {
		return "nullptr"
	}
	return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitTernaryExpression(ctx *parser.TernaryExpressionContext) interface{} {
	exprs := ctx.AllExpression()
	if len(exprs) == 3 {
		cond := v.Visit(exprs[0])
		if cond == nil {
			cond = exprs[0].GetText()
		}
		thenBranch := v.Visit(exprs[1])
		if thenBranch == nil {
			thenBranch = exprs[1].GetText()
		}
		elseBranch := v.Visit(exprs[2])
		if elseBranch == nil {
			elseBranch = exprs[2].GetText()
		}
		return fmt.Sprintf("%v ? %v : %v", cond, thenBranch, elseBranch)
	}
	return ctx.GetText()
}

func (v *CppTranspilerVisitor) VisitUnaryOperatorExpression(ctx *parser.UnaryOperatorExpressionContext) interface{} {
	prefix := ""
	if ctx.GetPrefix() != nil {
		prefix = ctx.GetPrefix().GetText()
	}
	exprStr := ""
	if ctx.Expression() != nil {
		res := v.Visit(ctx.Expression())
		if res != nil {
			exprStr = fmt.Sprintf("%v", res)
		} else {
			exprStr = ctx.Expression().GetText()
		}
	}
	return fmt.Sprintf("%s%s", prefix, exprStr)
}

func (v *CppTranspilerVisitor) VisitSquareBracketExpression(ctx *parser.SquareBracketExpressionContext) interface{} {
	exprs := ctx.AllExpression()
	if len(exprs) == 2 {
		arrRes := v.Visit(exprs[0])
		if arrRes == nil {
			arrRes = exprs[0].GetText()
		}
		idxRes := v.Visit(exprs[1])
		if idxRes == nil {
			idxRes = exprs[1].GetText()
		}
		return fmt.Sprintf("(*JAVA_NULL_CHECK(%v))[%v]", arrRes, idxRes)
	}
	return ctx.GetText()
}
