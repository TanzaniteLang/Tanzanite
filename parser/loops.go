package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) parseBegin() ast.Statement {
    p.consume()

    stat := ast.WhileStatement{
        Condition: nil,
        Until: false,
        DoWhile: true,
        Body: ast.Body{
            Scope: map[string]*ast.VarDeclaration{},
            Body: []ast.Statement{},
        },
    }

    current := p.current()
    p.AppendScope(&stat.Body)
    for current.Info != tokens.End {
        stat.Body.Append(p.parseStatement())
        current = p.current()
    }
    p.PopScope()

    p.consume()

    c := p.consume()
    index := p.pos
    start_line := c.Position.Line

    expr := p.parseExpression()

    if expr == nil || start_line != p.tokens[index].Position.Line {
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1 + uint64(len(c.Text)))
        dbg.ThrowError("Missign expression!", p.warn || p.Dead, nil)
        p.Dead = true
        p.skipToNewLine()
    }

    if c.Info == tokens.Unless {
        stat.Until = true
    }

    stat.Condition = expr

    return stat
}

func (p *Parser) parseWhile(until bool) ast.Statement {
    c := p.consume()
    index := p.pos
    start_line := c.Position.Line

    expr := p.parseExpression()

    if expr == nil || start_line != p.tokens[index].Position.Line {
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1 + uint64(len(c.Text)))
        dbg.ThrowError("Missign expression!", p.warn || p.Dead, nil)
        p.Dead = true
        p.skipToNewLine()
    }

    if p.current().Info == tokens.Do {
        p.consume()
    }

    stat := ast.WhileStatement{
        Condition: expr,
        Until: until,
        DoWhile: false,
        Body: ast.Body{
            Scope: map[string]*ast.VarDeclaration{},
            Body: []ast.Statement{},
        },
    }

    p.AppendScope(&stat.Body)
    current := p.current()
    for current.Info != tokens.End {
        stat.Body.Append(p.parseStatement())
        current = p.current()
    }
    p.PopScope()

    p.consume()

    return stat
}
