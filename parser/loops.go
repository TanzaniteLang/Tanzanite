package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) parseWhile(until bool) ast.Statement {
    p.consume()

    stat := ast.WhileStatement{
        Condition: p.parseExpression(),
        Until: until,
    }

    current := p.current()
    for current.Info != tokens.End {
        stat.Debug = append(stat.Debug, debug.NewSourceLocation(p.source, current.Position.Line, current.Position.Column))
        stat.Body = append(stat.Body, p.parseStatement())
        current = p.current()
    }

    p.consume()

    return stat
}
