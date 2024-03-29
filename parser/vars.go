package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
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
        p.consume()
        return ast.VarDeclaration{
            Name: ident.Text,
            Type: make([]ast.Statement, 0),
            Value: p.parseExpression(),
        }
    }
    
    // TODO: Throw error of invalid syntax
    return nil
}
