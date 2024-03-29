package ast

type NodeType int

const (
    ProgramType = 0
    NumericLiteralType
    IdentifierType
    StringType
    PointerType
    BinaryExprType
    VarDeclarationType
    AssignExprType

    FunctionArgType
    VariadicArgType
    FunctionDeclType
    FunctionCallType
    ReturnExprType
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
    Body []Statement
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

type AssignExpr struct {
    Name Expression
    Value Expression
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
