package ast

import "codeberg.org/Tanzanite/Tanzanite/debug"

type NodeType int

const (
    ProgramType = 0

    // Types
    IntLiteralType = 1
    FloatLiteralType = 2
    StringType = 4
    BoolType = 5
    PointerType = 6
    IdentifierType = 7
    TypeLiteralType = 8

    // Expressions
    BinaryExprType = 9
    UnaryExprType = 10
    VarDeclarationType = 11
    AssignExprType = 12
    ReturnExprType = 13
    BracketExprType = 14
    ConditionalExprType = 15
    ForwardPipeExprType = 16

    // Functions
    VariadicArgType = 17
    FunctionDeclType = 18
    FunctionCallType = 19
)

type Statement interface {
    GetKind() NodeType
}

type Program struct {
    Body []Statement
    Debug []debug.SourceLocation
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

type TypeLiteral struct {
    Type string
}

func (t TypeLiteral) GetKind() NodeType {
    return TypeLiteralType
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

type FunctionDecl struct {
    Name string
    Arguments []Statement
    ReturnType []Statement
    Immutable bool // True if this is FUN function
    Variadic bool
    Body []Statement
    Debug []debug.SourceLocation
}

func (f FunctionDecl) GetKind() NodeType {
    return FunctionDeclType
}

type FunctionCall struct {
    Calle Expression
    Args []Expression
}

func (f FunctionCall) GetKind() NodeType {
    return FunctionCallType
}

type ReturnExpr struct {
    Value Expression
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
