package ast

import "codeberg.org/Tanzanite/Tanzanite/debug"

type NodeType int

const (
    ProgramType = 0

    // Types
    IntLiteralType = 1
    FloatLiteralType = 2
    StringType = 4
    CharType = 5
    BoolType = 6
    PointerType = 7
    IdentifierType = 8
    TypeLiteralType = 9
    TypeCastType = 10

    // Expressions
    BinaryExprType = 11
    UnaryExprType = 12
    VarDeclarationType = 13
    AssignExprType = 14
    ReturnExprType = 15
    BracketExprType = 16
    ConditionalExprType = 17
    ForwardPipeExprType = 18

    // Functions
    VariadicArgType = 19
    FunctionDeclType = 20
    FunctionCallType = 21

    // Conditions
    IfStatementType = 22
    ElsifStatementType = 23
    ElseStatementType = 24
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
    Type []Statement
}

func (t TypeLiteral) GetKind() NodeType {
    return TypeLiteralType
}

type TypeCast struct {
    Target TypeLiteral
    Expr Expression
}

func (t TypeCast) GetKind() NodeType {
    return TypeCastType
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

type Char struct {
    Value string
}

func (c Char) GetKind() NodeType {
    return CharType
}

type Bool struct {
    Value string
}

func (b Bool) GetKind() NodeType {
    return BoolType
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


type IfStatement struct {
    Condition Expression
    Body []Statement
    Next Statement
    Debug []debug.SourceLocation
}

func (i IfStatement) GetKind() NodeType {
    return IfStatementType
}

type ElsifStatement struct {
    Condition Expression
    Body []Statement
    Debug []debug.SourceLocation
    Next Statement
}

func (e ElsifStatement) GetKind() NodeType {
    return ElsifStatementType
}

type ElseStatement struct {
    Body []Statement
    Debug []debug.SourceLocation
}

func (e ElseStatement) GetKind() NodeType {
    return ElseStatementType
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
