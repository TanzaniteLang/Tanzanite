package ast

import "codeberg.org/Tanzanite/Tanzanite/tokens"

type NodeType int

const (
    ProgramType = 0
    BodyType = 1

    // Types
    IntLiteralType = 2
    FloatLiteralType = 3
    StringType = 4
    CharType = 5
    BoolType = 6
    PointerType = 7
    IdentifierType = 8
    TypeLiteralType = 9
    TypeCastType = 10

    // Expressions
    BinaryExprType = 11
    FieldAccessType = 12
    UnaryExprType = 13
    VarDeclarationType = 14
    AssignExprType = 15
    ReturnExprType = 16
    BracketExprType = 17
    ConditionalExprType = 18
    ForwardPipeExprType = 19

    // Functions
    VariadicArgType = 20
    FunctionDeclType = 21
    DefFunctionDecl = 22
    FunctionCallType = 23

    // Conditions
    IfStatementType = 24
    ElsifStatementType = 25
    ElseStatementType = 26

    // Loops
    WhileStatementType = 27
    LoopControlStatementType = 28
)

type Statement interface {
    GetKind() NodeType
}

type Body struct {
    Scope map[string]*VarDeclaration
    Body []Statement
    // TODO: Debug
}

func (b Body) GetKind() NodeType {
    return BodyType
}

func (b *Body) RegisterVar(name string, val *VarDeclaration) {
    b.Scope[name] = val    
}

func (b *Body) HasVar(name string) bool {
    _, ok := b.Scope[name]

    return ok
}

func (b *Body) Append(stat Statement) {
    b.Body = append(b.Body, stat)
}

type Program struct {
    Body Body
}

type Expression Statement

type ConditionalExpr struct {
    Position tokens.Position
    Condition Expression
    TrueExpr Expression
    FalseExpr Expression
}

func (c ConditionalExpr) GetKind() NodeType {
    return ConditionalExprType
}

type UnaryExpr struct {
    Position tokens.Position
    Operator string
    Operand Expression
}

func (u UnaryExpr) GetKind() NodeType {
    return UnaryExprType
}

type BinaryExpr struct {
    Position tokens.Position
    Left Expression
    Right Expression
    Operator string
}

type FieldAccess struct {
    Position tokens.Position
    Left Expression
    Right Expression
}

func (f FieldAccess) GetKind() NodeType {
    return FieldAccessType
}

type TypeLiteral struct {
    Type []Statement
}

func (t TypeLiteral) GetKind() NodeType {
    return TypeLiteralType
}

type TypeCast struct {
    Position tokens.Position
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
    Position tokens.Position
    Value Expression
    Target Expression
}

func (f ForwardPipeExpr) GetKind() NodeType {
    return ForwardPipeExprType
}

type Identifier struct {
    Position tokens.Position
    Symbol string
}

func (i Identifier) GetKind() NodeType {
    return IdentifierType
}

type IntLiteral struct {
    Position tokens.Position
    Value int64 
}

func (i IntLiteral) GetKind() NodeType {
    return IntLiteralType
}

type FloatLiteral struct {
    Position tokens.Position
    Value float64 
}

func (f FloatLiteral) GetKind() NodeType {
    return FloatLiteralType
}

type String struct {
    Position tokens.Position
    Value string
}

func (s String) GetKind() NodeType {
    return StringType
}

type Char struct {
    Position tokens.Position
    Value string
}

func (c Char) GetKind() NodeType {
    return CharType
}

type Bool struct {
    Position tokens.Position
    Value string
}

func (b Bool) GetKind() NodeType {
    return BoolType
}

type VarDeclaration struct {
    Position tokens.Position
    Name string
    Type []Statement
    Value Expression
}

type VariadicArg struct {}

func (v VariadicArg) GetKind() NodeType {
    return VariadicArgType
}

type FunctionDecl struct {
    Position tokens.Position
    Name string
    Arguments []Statement
    ReturnType []Statement
    Immutable bool // True if this is FUN function
    Variadic bool
    Failed bool
    Body Body
}

func (f FunctionDecl) GetKind() NodeType {
    return FunctionDeclType
}

type FunctionCall struct {
    Position tokens.Position
    Calle Expression
    Args []Expression
}

func (f FunctionCall) GetKind() NodeType {
    return FunctionCallType
}

type ReturnExpr struct {
    Position tokens.Position
    Value Expression
}

func (r ReturnExpr) GetKind() NodeType {
    return ReturnExprType
}

type BracketExpr struct {
    Position tokens.Position
    Expr Expression
}

func (b BracketExpr) GetKind() NodeType {
    return BracketExprType
}

type AssignExpr struct {
    Position tokens.Position
    Name Expression
    Value Expression
    Operator string
}


type IfStatement struct {
    Position tokens.Position
    Condition Expression
    Unless bool
    Next Statement
    Body Body
}

func (i IfStatement) GetKind() NodeType {
    return IfStatementType
}

type ElsifStatement struct {
    Position tokens.Position
    Condition Expression
    Next Statement
    Body Body
}

func (e ElsifStatement) GetKind() NodeType {
    return ElsifStatementType
}

type ElseStatement struct {
    Position tokens.Position
    Body Body
}

func (e ElseStatement) GetKind() NodeType {
    return ElseStatementType
}

type WhileStatement struct {
    Position tokens.Position
    Condition Expression
    Until bool
    DoWhile bool
    Body Body
}

func (w WhileStatement) GetKind() NodeType {
    return WhileStatementType
}

type LoopControlStatement struct {
    Position tokens.Position
    Break bool
}

func (l LoopControlStatement) GetKind() NodeType {
    return LoopControlStatementType
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
