package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) parseVarDeclaration() ast.Statement {
    ident := p.consume()

    if p.current().Info == tokens.Colon {
        p.consume()

        varType := p.parseType()

        if p.current().Info == tokens.Assign {
            c := p.consume()

            expr := p.parseExpression()

            if expr == nil {
                dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1)
                dbg.ThrowError("Missign expression!", p.warn || p.Dead, nil)
                p.Dead = true
                p.skipToNewLine()
            }

            return ast.VarDeclaration{
                Name: ident.Text,
                Type: varType,
                Value: expr,
            }
        } else {
            return ast.VarDeclaration{
                Name: ident.Text,
                Type: varType,
                Value: ast.IntLiteral{
                    Value: 0,
                },
            }
        }
    } else if p.current().Info == tokens.Assign {
        c := p.current()

        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column - 1)
        dbg.ThrowWarning("Specify the Type until Static Analyzer is present!", p.warn || p.Dead, &debug.Hint{
            Msg: "Use any of these types: Char, Bool, Int or Float", 
            Code: ": Type",
        })
        p.warn = true
        p.skipToNewLine()
        return nil
    }

    c := p.current()
 
    dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
    dbg.ThrowError("Expected : or =, but got " + c.Text + " instead!", p.warn || p.Dead, nil)
    p.Dead = true
    p.skipToNewLine()
    
    return nil
}
