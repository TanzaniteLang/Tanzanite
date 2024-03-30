package parser

import (
    "strconv"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

func (p *Parser) parseType() []ast.Statement {
    typeConstruct := make([]ast.Statement, 0)
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

func (p *Parser) parsePrimaryExpr() ast.Expression {
    tok := p.current().Info;

    switch tok {
    case tokens.Identifier:
        fn, ok := p.env.Fns[p.current().Text]
        if ok && !p.parsingFn { // This is a function call
            return p.parseFnCall(fn)
        }
        return ast.Identifier{ Symbol: p.consume().Text }
    case tokens.Char:
        return ast.TypeLiteral{ Type: p.consume().Text }
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
    case tokens.LBracket:
        p.consume()
 
        val := ast.BracketExpr{
            Expr: p.parseExpression(),
        }
        p.consume()

        return val
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
