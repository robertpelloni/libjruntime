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
	outDir       string
	currentClass string
	headerCode   strings.Builder
	cppCode      strings.Builder
}

func NewCppTranspilerVisitor(outDir string) *CppTranspilerVisitor {
	return &CppTranspilerVisitor{
		BaseJavaParserVisitor: &parser.BaseJavaParserVisitor{},
		outDir:                outDir,
	}
}

func mapType(javaType string) string {
	switch javaType {
	case "int":
		return "int32_t"
	case "boolean":
		return "bool"
	case "long":
		return "int64_t"
	case "String":
		return "java::lang::String*"
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

	v.headerCode.WriteString(fmt.Sprintf("#pragma once\n\n#include <cstdint>\n#include \"Object.h\"\n#include \"String.h\"\n#include \"GC.h\"\n#include \"Runtime.h\"\n#include \"FinallyGuard.h\"\n\nclass %s : public java::lang::Object {\npublic:\n", className))
	v.cppCode.WriteString(fmt.Sprintf("#include \"%s.h\"\n\n", className))

	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}

	v.headerCode.WriteString("};\n")

	err := os.WriteFile(filepath.Join(v.outDir, className+".h"), []byte(v.headerCode.String()), 0644)
	if err != nil {
		fmt.Printf("Failed to write header file: %v\n", err)
	}
	err = os.WriteFile(filepath.Join(v.outDir, className+".cpp"), []byte(v.cppCode.String()), 0644)
	if err != nil {
		fmt.Printf("Failed to write cpp file: %v\n", err)
	}

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
	for _, child := range ctx.GetChildren() {
		if payload, ok := child.(antlr.ParseTree); ok {
			payload.Accept(v)
		}
	}
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

func (v *CppTranspilerVisitor) VisitMethodDeclaration(ctx *parser.MethodDeclarationContext) interface{} {
	methodName := ctx.Identifier().GetText()

	fmt.Printf("Discovered method: %s\n", methodName)

	returnType := "void"
	if ctx.TypeTypeOrVoid().TypeType() != nil {
		returnType = mapType(ctx.TypeTypeOrVoid().TypeType().GetText())
	}

	params := ""
	if ctx.FormalParameters() != nil {
        fpListCtx := ctx.FormalParameters().GetChild(1) // FormalParameterList is child 1 if present.
        if fpList, ok := fpListCtx.(*parser.FormalParameterListContext); ok {
            var pList []string
            for _, p := range fpList.AllFormalParameter() {
                pType := mapType(p.TypeType().GetText())
                pName := p.VariableDeclaratorId().GetText()
                pList = append(pList, fmt.Sprintf("%s %s", pType, pName))
            }
            params = strings.Join(pList, ", ")
        }
	}

	v.headerCode.WriteString(fmt.Sprintf("    virtual %s %s(%s);\n", returnType, methodName, params))

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
		v.headerCode.WriteString(fmt.Sprintf("    %s %s;\n", fieldType, fieldName))
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
        typ := mapType(decl.TypeType().GetText())
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

    if ctx.RETURN() != nil {
        v.cppCode.WriteString("    return ")
        if ctx.Expression(0) != nil {
            res := v.Visit(ctx.Expression(0))
            if res != nil {
                v.cppCode.WriteString(fmt.Sprintf("%v", res))
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
             className := creator.CreatedName().GetText()
             // Check if mapped
             if className == "String" {
                 className = "java::lang::String"
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

             if ctx.Identifier() != nil {
                 return fmt.Sprintf("JAVA_NULL_CHECK(%v)->%v", objRes, ctx.Identifier().GetText())
             } else if ctx.MethodCall() != nil {
                 methodCallRes := v.Visit(ctx.MethodCall())
                 if methodCallRes == nil {
                     methodCallRes = ctx.MethodCall().GetText()
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

    // search for ExpressionList manually
    for _, child := range ctx.GetChildren() {
        if exprListCtx, ok := child.(*parser.ExpressionListContext); ok {
             var parts []string
             for _, expr := range exprListCtx.AllExpression() {
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

    return fmt.Sprintf("%s(%s)", methodName, argsStr)
}
