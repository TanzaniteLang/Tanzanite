package ast

import (
    "fmt"
    "strings"
)

type NodeType int

const (
    ProgramType = 0
    IntLiteralType = 1
    FloatLiteralType = 2
    IdentifierType = 3
    StringType = 4
    PointerType = 5
    BinaryExprType = 6
    VarDeclarationType = 7
    AssignExprType = 8

    FunctionArgType = 9
    VariadicArgType = 10
    FunctionDeclType = 11
    FunctionCallType = 12
    ReturnExprType = 13

    ConditionalExprType = 14
    UnaryExprType = 15
    ForwardPipeExprType = 16

    BracketExprType = 17
)

type Statement interface {
    GetKind() NodeType
}

type Program struct {
    Body []Statement
}

type Expression Statement

type ConditionalExpr struct {
    Condition Expression
    TrueExpr Expression
    FalseExpr Expression
}

func (c ConditionalExpr) GetKind() NodeType {
    return ConditionalExprType
}

type UnaryExpr struct {
    Operator string
    Operand Expression
}

func (u UnaryExpr) GetKind() NodeType {
    return UnaryExprType
}

type BinaryExpr struct {
    Left Expression
    Right Expression
    Operator string
}

func (b *BinaryExpr) Stringify() string {
    return fmt.Sprintf("%s %s %s", stringifyExpr(b.Left), b.Operator, stringifyExpr(b.Right))
}

func (b BinaryExpr) GetKind() NodeType {
    return BinaryExprType
}

type ForwardPipeExpr struct {
    Value Expression
    Target Expression
}

func (f ForwardPipeExpr) GetKind() NodeType {
    return ForwardPipeExprType
}

type Identifier struct {
    Symbol string
}

func (i Identifier) GetKind() NodeType {
    return IdentifierType
}

type IntLiteral struct {
    Value int64 
}

func (i IntLiteral) GetKind() NodeType {
    return IntLiteralType
}

type FloatLiteral struct {
    Value float64 
}

func (f FloatLiteral) GetKind() NodeType {
    return FloatLiteralType
}

type String struct {
    Value string
}

func (s String) GetKind() NodeType {
    return StringType
}

type VarDeclaration struct {
    Name string
    Type []Statement
    Value Expression
}

type VariadicArg struct {}

func (v VariadicArg) GetKind() NodeType {
    return VariadicArgType
}

type FunctionArg struct {
    Name string
    Type []Statement
    Value Expression
}

func (f FunctionArg) GetKind() NodeType {
    return FunctionArgType
}

type FunctionDecl struct {
    Name string
    Arguments []Statement
    ReturnType []Statement
    Immutable bool // True if this is FUN function
    Variadic bool
    Body []Statement
}

func (f *FunctionDecl) GenDecl() string {
    decl := strings.ToLower(stringifyType(f.ReturnType)) // TODO: Will break custom types
    decl += " "
    decl += f.Name

    decl += "("
    for _, arg := range f.Arguments {
        if arg == nil { continue }

        if arg.GetKind() == VarDeclarationType {
            vDecl := arg.(VarDeclaration)
            decl += strings.ToLower(stringifyType(vDecl.Type))
            decl += " "
            decl += vDecl.Name
        } else if arg.GetKind() == VariadicArgType {
            decl += "..."
        }
        decl += ", "
    }

    return strings.TrimSuffix(decl, ", ") + ")"
}

func (f *FunctionDecl) GenBody() string {
    body := ""

    for _, stmt := range f.Body {
        if stmt.GetKind() == FunctionCallType {
            fnCall := stmt.(FunctionCall)
            if fnCall.Calle.GetKind() == IdentifierType {
                body += fnCall.Calle.(Identifier).Symbol
            } else {
                panic("Other forms not implemented")
            }

            body += "(" + fnCall.StringifyArgs() + ");\n"
        } else if stmt.GetKind() == ReturnExprType {
            ret := stmt.(ReturnExpr)
            body += fmt.Sprintf("%s;\n", ret.Stringify())
        }
    }

    return body
}

func (f FunctionDecl) GetKind() NodeType {
    return FunctionDeclType
}

type FunctionCall struct {
    Calle Expression
    Args []Expression
}

func (f *FunctionCall) StringifyArgs() string {
    args := ""

    for _, arg := range f.Args {
        if arg == nil { continue }

        args += stringifyExpr(arg) + ", "
    }

    return strings.TrimSuffix(args, ", ")
}

func (f FunctionCall) GetKind() NodeType {
    return FunctionCallType
}

type ReturnExpr struct {
    Value Expression
}

func (r *ReturnExpr) Stringify() string {
    return fmt.Sprintf("return %s", stringifyExpr(r.Value))
}

func (r ReturnExpr) GetKind() NodeType {
    return ReturnExprType
}

type BracketExpr struct {
    Expr Expression
}

func (b BracketExpr) GetKind() NodeType {
    return BracketExprType
}

type AssignExpr struct {
    Name Expression
    Value Expression
    Operator string
}

type Pointer struct {}

func (p Pointer) GetKind() NodeType {
    return PointerType
}

func (a AssignExpr) GetKind() NodeType {
    return AssignExprType
}

func (v VarDeclaration) GetKind() NodeType {
    return VarDeclarationType
}

func stringifyExpr(e Expression) string {
    switch (e.GetKind()) {
    case IdentifierType:
        return e.(Identifier).Symbol
    case IntLiteralType:
        return fmt.Sprintf("%d", e.(IntLiteral).Value)
    case StringType:
        return "\"" + e.(String).Value + "\""
    case BinaryExprType:
        expr := e.(BinaryExpr)
        return expr.Stringify()
    default:
        return ""
    }
}

func stringifyType(t []Statement) string {
    text := ""

    for _, part := range t {
        if part.GetKind() == IdentifierType {
            text += part.(Identifier).Symbol
        } else if part.GetKind() == PointerType {
            text += "*"
        }
    }

    return text
}
