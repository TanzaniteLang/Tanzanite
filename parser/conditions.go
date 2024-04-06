package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) parseElse() ast.Statement {
    start := p.consume()

    stat := ast.ElseStatement{
        Body: ast.Body{
            Scope: map[string]*ast.VarDeclaration{},
            Body: []ast.Statement{},
        },
        Position: start.Position,
    }

    p.AppendScope(&stat.Body)
    current := p.current()
    for current.Info != tokens.End {
        stat.Body.Append(p.parseStatement())
        current = p.current()
    }
    p.consume()

    p.PopScope()

    return stat
}

func (p *Parser) parseElsif() ast.Statement {
    c := p.consume()
    index := p.pos
    start_line := c.Position.Line
    expr_pos := p.current().Position

    expr := p.parseExpression()

    if expr == nil || start_line != p.tokens[index].Position.Line {
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1 + uint64(len(c.Text)))
        dbg.ThrowError("Missign expression!", p.warn || p.Dead, nil)
        p.Dead = true
        p.skipToNewLine()
    }

    if p.current().Info == tokens.Then {
        p.consume()
    }

    stat := ast.ElsifStatement{
        Condition: expr,
        Body: ast.Body{
            Scope: map[string]*ast.VarDeclaration{},
            Body: []ast.Statement{},
        },
        Position: expr_pos,
    }

    p.AppendScope(&stat.Body)
    current := p.current()
    for current.Info != tokens.End && current.Info != tokens.Elsif && current.Info != tokens.Else {
        stat.Body.Append(p.parseStatement())
        current = p.current()
    }
    p.PopScope()

    if current.Info == tokens.Elsif {
        stat.Next = p.parseElsif()
    } else if current.Info == tokens.Else {
        stat.Next = p.parseElse()
    } else {
        p.consume()
    }

    return stat
}

func (p *Parser) parseIf(unless bool) ast.Statement {
    c := p.consume()
    index := p.pos
    start_line := c.Position.Line
    expr_pos := p.current().Position
    
    expr := p.parseExpression()

    if expr == nil || start_line != p.tokens[index].Position.Line {
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1 + uint64(len(c.Text)))
        dbg.ThrowError("Missign expression!", p.warn || p.Dead, nil)
        p.Dead = true
        p.skipToNewLine()
    }

    if p.current().Info == tokens.Then {
        p.consume()
    }

    stat := ast.IfStatement{
        Condition: expr,
        Unless: unless,
        Body: ast.Body{
            Scope: map[string]*ast.VarDeclaration{},
            Body: []ast.Statement{},
        },
        Position: expr_pos,
    }

    p.AppendScope(&stat.Body)
    current := p.current()
    for current.Info != tokens.End && current.Info != tokens.Elsif && current.Info != tokens.Else {
        stat.Body.Append(p.parseStatement())
        current = p.current()
    }

    p.PopScope()

    if current.Info == tokens.Elsif {
        stat.Next = p.parseElsif()
    } else if current.Info == tokens.Else {
        stat.Next = p.parseElse()
    } else {
        p.consume()
    }

    return stat
}
