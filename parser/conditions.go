package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) parseElse() ast.Statement {
    p.consume()

    stat := ast.ElseStatement{}

    current := p.current()
    for current.Info != tokens.End {
        stat.Debug = append(stat.Debug, debug.NewSourceLocation(p.source, current.Position.Line, current.Position.Column))
        stat.Body = append(stat.Body, p.parseStatement())
        current = p.current()
    }

    p.consume()

    return stat
}

func (p *Parser) parseElsif() ast.Statement {
    p.consume()

    stat := ast.ElsifStatement{
        Condition: p.parseExpression(),
        Next: nil,
    }

    current := p.current()
    for current.Info != tokens.End && current.Info != tokens.Elsif && current.Info != tokens.Else {
        stat.Debug = append(stat.Debug, debug.NewSourceLocation(p.source, current.Position.Line, current.Position.Column))
        stat.Body = append(stat.Body, p.parseStatement())
        current = p.current()
    }

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
    p.consume()

    stat := ast.IfStatement{
        Condition: p.parseExpression(),
        Unless: unless,
        Next: nil,
    }

    current := p.current()
    for current.Info != tokens.End && current.Info != tokens.Elsif && current.Info != tokens.Else {
        stat.Debug = append(stat.Debug, debug.NewSourceLocation(p.source, current.Position.Line, current.Position.Column))
        stat.Body = append(stat.Body, p.parseStatement())
        current = p.current()
    }

    if current.Info == tokens.Elsif {
        stat.Next = p.parseElsif()
    } else if current.Info == tokens.Else {
        stat.Next = p.parseElse()
    } else {
        p.consume()
    }

    return stat
}
