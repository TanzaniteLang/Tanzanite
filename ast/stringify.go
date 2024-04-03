package ast

import (
    "fmt"
    "strings"
)

func (t *TypeLiteral) Stringify() string {
    return strType(t.Type)
}

func (t *TypeCast) Stringify() string {
    return fmt.Sprintf("(%s) %s", t.Target.Stringify(), strExpr(t.Expr))
}

func (b *BinaryExpr) Stringify() string {
    if b.Operator == "**" {
        return fmt.Sprintf("pow(%s, %s)", strExpr(b.Left), strExpr(b.Right))
    } else if b.Operator == "//" {
        return fmt.Sprintf("floor(%s) / floor(%s)", strExpr(b.Left), strExpr(b.Right))
    }
    return fmt.Sprintf("%s %s %s", strExpr(b.Left), b.Operator, strExpr(b.Right))
}

func (u *UnaryExpr) Stringify() string {
    return fmt.Sprintf("%s%s", u.Operator, strExpr(u.Operand))
}

func (v *VarDeclaration) Stringify() string {
    if v.Value != nil {
        return fmt.Sprintf("%s %s = %s", strType(v.Type), v.Name, strExpr(v.Value))
    }
    return fmt.Sprintf("%s %s", strType(v.Type), v.Name)
}

func (v *VarDeclaration) StringifyAsArg() string {
    return fmt.Sprintf("%s %s", strType(v.Type), v.Name)
}

func (a *AssignExpr) Stringify() string {
    if a.Operator == "**=" {
        return fmt.Sprintf("%s = pow(%s, %s)", strExpr(a.Name), strExpr(a.Name), strExpr(a.Value))
    } else if a.Operator == "//=" {
        return fmt.Sprintf("%s = floor(%s) / floor(%s)", strExpr(a.Name), strExpr(a.Name), strExpr(a.Value))
    }
    return fmt.Sprintf("%s %s %s", strExpr(a.Name), a.Operator, strExpr(a.Value))
}

func (r *ReturnExpr) Stringify() string {
    return fmt.Sprintf("return %s", strExpr(r.Value))
}

func (b *BracketExpr) Stringify() string {
    return fmt.Sprintf("(%s)", strExpr(b.Expr))
}

func (c *ConditionalExpr) Stringify() string {
    return fmt.Sprintf("(%s ? %s : %s)", strExpr(c.Condition), strExpr(c.TrueExpr), strExpr(c.FalseExpr))
}

func (f *ForwardPipeExpr) Stringify() string {
    return fmt.Sprintf("%s(%s)", strExpr(f.Target), strExpr(f.Value))
}

func (f *FunctionDecl) StringifyHead() string {
    text := strType(f.ReturnType) + " " + f.Name + "("

    for _, arg := range f.Arguments {
        if arg.GetKind() == VarDeclarationType {
            expr := arg.(VarDeclaration)
            text += expr.StringifyAsArg()
        } else if arg.GetKind() == VariadicArgType {
            text += "..."
        }

        text += ", "
    }

    return strings.TrimSuffix(text, ", ") + ")"
}

func (f *FunctionDecl) StringifyBody() string {
    body := ""
    /*
    for i, stmt := range f.Body {
        body += f.Debug[i].Stringify() + "\n"
        body += strExpr(stmt) + ";\n"
    }
    */

    return body
}

func (f *FunctionDecl) Stringify() string {
    return f.StringifyHead() + " {\n" + f.StringifyBody() + "}\n"
}

func (f *FunctionCall) StringifyArgs() string {
    args := ""
    for _, arg := range f.Args {
        args += strExpr(arg) + ", "
    }

    return strings.TrimSuffix(args, ", ")
}

func (f *FunctionCall) Stringify() string {
    return fmt.Sprintf("%s(%s)", strExpr(f.Calle), f.StringifyArgs())
}

func (i *IfStatement) Stringify() string {
    code := ""
    if i.Unless {
        code = fmt.Sprintf("if (!(%s)) {\n", strExpr(i.Condition))
    } else {
        code = fmt.Sprintf("if (%s) {\n", strExpr(i.Condition))
    }
/*
    for iter, stmt := range i.Body {
        code += i.Debug[iter].Stringify() + "\n"
        code += strExpr(stmt) + ";\n"
    }
    */
    code += "}"

    if i.Next != nil {
        code += " "
        code += strExpr(i.Next)
    }

    return code
}

func (e *ElsifStatement) Stringify() string {
    code := fmt.Sprintf("else if (%s) {\n", strExpr(e.Condition))
    /*
    for i, stmt := range e.Body {
        code += e.Debug[i].Stringify() + "\n"
        code += strExpr(stmt) + ";\n"
    }*/
    code += "}"

    if e.Next != nil {
        code += " "
        code += strExpr(e.Next)
    }

    return code
}

func (e *ElseStatement) Stringify() string {
    code := "else {\n"
    /*
    for i, stmt := range e.Body {
        code += e.Debug[i].Stringify() + "\n"
        code += strExpr(stmt) + ";\n"
    }*/
    code += "}"

    return code
}

func (w *WhileStatement) Stringify() string {
    if !w.DoWhile {
        code := ""
        if w.Until {
            code = fmt.Sprintf("while (!(%s)) {\n", strExpr(w.Condition))
        } else {
            code = fmt.Sprintf("while (%s) {\n", strExpr(w.Condition))
        }

        /*
        for i, stmt := range w.Body {
            code += w.Debug[i].Stringify() + "\n"
            code += strExpr(stmt) + ";\n"
        }
        */
        code += "}"

        return code
    }

    code := "do {\n"
    /*
    for i, stmt := range w.Body {
        code += w.Debug[i].Stringify() + "\n"
        code += strExpr(stmt) + ";\n"
    }
    */
    code += "} "

    if w.Until {
        code += fmt.Sprintf("while (!(%s))", strExpr(w.Condition))
    } else {
        code += fmt.Sprintf("while (%s)", strExpr(w.Condition))
    }

    return code
}

func (l *LoopControlStatement) Stringify() string {
    if l.Break {
        return "break"
    } else {
        return "continue"
    }
}

func strExpr(e Expression) string {
    switch (e.GetKind()) {
        // Types
    case IntLiteralType:
        return fmt.Sprintf("%d", e.(IntLiteral).Value)
    case FloatLiteralType:
        return fmt.Sprintf("%f", e.(FloatLiteral).Value)
    case StringType:
        return "\"" + e.(String).Value + "\""
    case CharType:
        return "'" + e.(Char).Value + "'"
    case BoolType:
        return e.(Bool).Value
    case PointerType:
        return "*"
    case IdentifierType:
        return e.(Identifier).Symbol
    case TypeLiteralType:
        expr := e.(TypeLiteral)
        return expr.Stringify()
    case TypeCastType:
        expr := e.(TypeCast)
        return expr.Stringify()

        // Expressions
    case BinaryExprType:
        expr := e.(BinaryExpr)
        return expr.Stringify()
    case UnaryExprType:
        expr := e.(UnaryExpr)
        return expr.Stringify()
    case VarDeclarationType:
        expr := e.(VarDeclaration)
        return expr.Stringify()
    case AssignExprType:
        expr := e.(AssignExpr)
        return expr.Stringify()
    case ReturnExprType:
        expr := e.(ReturnExpr)
        return expr.Stringify()
    case BracketExprType:
        expr := e.(BracketExpr)
        return expr.Stringify()
    case ConditionalExprType:
        expr := e.(ConditionalExpr)
        return expr.Stringify()
    case ForwardPipeExprType:
        expr := e.(ForwardPipeExpr)
        return expr.Stringify()

        // Functions
    case VariadicArgType:
        return "..."
    case FunctionDeclType:
        expr := e.(FunctionDecl)
        return expr.Stringify()
    case FunctionCallType:
        expr := e.(FunctionCall)
        return expr.Stringify()
    
        // Conditions
    case IfStatementType:
        expr := e.(IfStatement)
        return expr.Stringify()
    case ElsifStatementType:
        expr := e.(ElsifStatement)
        return expr.Stringify()
    case ElseStatementType:
        expr := e.(ElseStatement)
        return expr.Stringify()

        // Loops
    case WhileStatementType:
        expr := e.(WhileStatement)
        return expr.Stringify()
    case LoopControlStatementType:
        expr := e.(LoopControlStatement)
        return expr.Stringify()
    default:
        return ""
    }
}

func strType(t []Statement) string {
    text := ""

    for _, part := range t {
        text += strExpr(part)
    }

    return text
}
