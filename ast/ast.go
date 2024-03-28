package ast

type NodeType int

const (
    ProgramType = 0
    NumericLiteralType
    IdentifierType
    BinaryExprType
    VarDeclarationType
    AssignExprType
)

type Statement interface {
    GetKind() NodeType
}

type Program struct {
    Body []Statement
}

type Expression Statement

type BinaryExpr struct {
    Left Expression
    Right Expression
    Operator string
}

func (b BinaryExpr) GetKind() NodeType {
    return BinaryExprType
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
    return NumericLiteralType
}

type FloatLiteral struct {
    Value float64 
}

func (f FloatLiteral) GetKind() NodeType {
    return NumericLiteralType
}

type VarDeclaration struct {
    Name string
    Type string
    Value Expression
}

type AssignExpr struct {
    Name Expression
    Value Expression
}

func (a AssignExpr) GetKind() NodeType {
    return AssignExprType
}

func (v VarDeclaration) GetKind() NodeType {
    return VarDeclarationType
}
