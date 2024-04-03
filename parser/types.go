package parser

import (
    "strconv"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/debug"
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

func (p *Parser) parseType() []ast.Statement {
    typeConstruct := make([]ast.Statement, 0)

    if !p.checkType() {
        c := p.current()

        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
        dbg.ThrowError("Specify the Type until Static Analyzer is present!", p.warn || p.Dead, &debug.Hint{
            Msg: "Use any of these types: Char, Bool, Int or Float", 
            Code: "Type ",
        })
        p.Dead = true
        p.skipToNewLine()

        return typeConstruct
    }

    typeConstruct = append(typeConstruct, ast.Identifier{Symbol: p.consume().Text})

    current := p.current().Info
    for current == tokens.Asterisk || current == tokens.DoubleAsterisk {
        if current == tokens.DoubleAsterisk {
            typeConstruct = append(typeConstruct, ast.Pointer{})
            typeConstruct = append(typeConstruct, ast.Pointer{})
        } else {
            typeConstruct = append(typeConstruct, ast.Pointer{})
        }
        p.consume()

        current = p.current().Info
    }

    return typeConstruct
}

func (p *Parser) checkType() bool {
    tok := p.current().Info

    switch tok {
    case tokens.Char:
        return true
    case tokens.Int:
        return true
    case tokens.Float:
        return true
    case tokens.Bool:
        return true
    case tokens.Void:
        return true
    case tokens.Identifier:
        return true
    }
    return false
}

func (p *Parser) parsePrimaryExpr() ast.Expression {
    tok := p.current().Info

    switch tok {
    case tokens.Identifier:
        fn, ok := p.Globals.Scope[p.current().Text]
        if ok && !p.parsingFn { // This is a function call
            return p.parseFnCall(fn)
        }
        return ast.Identifier{ Symbol: p.consume().Text }
    case tokens.CharVal:
        return ast.Char{ Value: p.consume().Text }
    case tokens.StringVal:
        return ast.String{ Value: p.consume().Text }
    case tokens.FloatVal:
        val, _ := strconv.ParseFloat(p.consume().Text, 64)
        return ast.FloatLiteral{
            Value: val,
        }
    case tokens.IntVal:
        val, _ := strconv.ParseInt(p.consume().Text, 10, 64)
        return ast.IntLiteral{
            Value: val,
        }
    case tokens.BoolVal:
        return ast.Bool{ Value: p.consume().Text }
    case tokens.LBracket:
        p.consume()

        if p.checkType() {
            val := ast.TypeCast{
                Target: ast.TypeLiteral{
                    Type: p.parseType(),
                },
                Expr: nil,
            }
            p.consume()

            val.Expr = p.parseExpression()

            return val
        } else {
            val := ast.BracketExpr{
                Expr: p.parseExpression(),
            }
            p.consume()

            return val
        }
    case tokens.Nil:
        p.consume()
        return ast.Identifier{ Symbol: "nil" }
    case tokens.Plus:
        return p.parseUnaryExpr()
    case tokens.Minus:
        return p.parseUnaryExpr()
    case tokens.Bang:
        return p.parseUnaryExpr()
    case tokens.Tilda:
        return p.parseUnaryExpr()
    case tokens.Ampersand:
        return p.parseUnaryExpr()
    case tokens.Asterisk:
        return p.parseUnaryExpr()
    case tokens.Sizeof:
        return p.parseUnaryExpr()
    }

    return nil
}
