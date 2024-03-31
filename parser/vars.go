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
            p.consume()

            return ast.VarDeclaration{
                Name: ident.Text,
                Type: varType,
                Value: p.parseExpression(),
            }
        } else {
            return ast.VarDeclaration{
                Name: ident.Text,
                Type: varType,
                Value: nil,
            }
        }
    } else if p.current().Info == tokens.Assign {
        c := p.current()

        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column - 1)
        dbg.ThrowHint("Specify the Type until Static Analyzer is present!",
            "Use any of these types: Char, Bool, Int or Float", ": Type", p.warn || p.Dead)
        p.warn = true
        p.skipToNewLine()
        return nil
    }

    c := p.current()
 
    dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
    dbg.ThrowError("Expected : or =, but got " + c.Text + " instead!", p.warn || p.Dead)
    p.Dead = true
    p.skipToNewLine()
    
    // TODO: Throw error of invalid syntax
    return nil
}
