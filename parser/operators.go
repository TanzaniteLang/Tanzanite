package parser

import "codeberg.org/Tanzanite/Tanzanite/ast"

func (p *Parser) parseAdditiveExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseMultiplicativeExpr()

    for p.current().Text == "+" || p.current().Text == "-" {
        operator := p.consume().Text
        right := p.parseMultiplicativeExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseMultiplicativeExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseAccessExpr()

    for p.current().Text == "/" ||
        p.current().Text == "*" ||
        p.current().Text == "**" ||
        p.current().Text == "//" ||
        p.current().Text == "%" {
        operator := p.consume().Text
        right := p.parseAccessExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseAccessExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parsePrimaryExpr()

    for p.current().Text == "." {
        p.consume()
        right := p.parseAccessExpr()
        left = ast.FieldAccess{
            Left: left,
            Right: right,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseShiftExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseAdditiveExpr()

    for p.current().Text == "<<" || p.current().Text == ">>" {
        operator := p.consume().Text
        right := p.parseAdditiveExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseComparativeExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseShiftExpr()

    for p.current().Text == "<" ||
        p.current().Text == "<=" ||
        p.current().Text == ">" ||
        p.current().Text == ">=" {
        operator := p.consume().Text
        right := p.parseShiftExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseEqaulityExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseComparativeExpr()

    for p.current().Text == "==" ||
        p.current().Text == "!=" {
        operator := p.consume().Text
        right := p.parseComparativeExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseSpaceshipExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseEqaulityExpr()

    for p.current().Text == "<=>" {
        operator := p.consume().Text
        right := p.parseEqaulityExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseBitwiseAndExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseSpaceshipExpr()

    for p.current().Text == "&" {
        operator := p.consume().Text
        right := p.parseSpaceshipExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseBitwiseXorExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseBitwiseAndExpr()

    for p.current().Text == "^" {
        operator := p.consume().Text
        right := p.parseBitwiseAndExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseBitwiseOrExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseBitwiseXorExpr()

    for p.current().Text == "|" {
        operator := p.consume().Text
        right := p.parseBitwiseXorExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseLogicalAndExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseBitwiseOrExpr()

    for p.current().Text == "&&" {
        operator := p.consume().Text
        right := p.parseBitwiseOrExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseLogicalOrExpr() ast.Expression {
    start_pos := p.current().Position
    left := p.parseLogicalAndExpr()

    for p.current().Text == "||" {
        operator := p.consume().Text
        right := p.parseLogicalAndExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
            Position: start_pos,
        }
    }

    return left
}

func (p *Parser) parseConditionalExpr() ast.Expression {
    start_pos := p.current().Position
    condition := p.parseLogicalOrExpr()

    if p.current().Text == "?" {
        p.consume()

        trueExpr := p.parseExpression()

        if p.current().Text != ":" {
            panic("Missing :")
        }
        p.consume()

        falseExpr := p.parseConditionalExpr()
        return ast.ConditionalExpr{
            Condition: condition,
            TrueExpr: trueExpr,
            FalseExpr: falseExpr,
            Position: start_pos,
        }
    }

    return condition
}

func (p *Parser) parseForwardPipeExpr() ast.Expression {
    start_pos := p.current().Position
    value := p.parseConditionalExpr()

    for p.current().Text == "|>" {
        p.consume()
        old := p.parsingFn
        p.parsingFn = true
        target := p.parseExpression()
        p.parsingFn = old
        return ast.ForwardPipeExpr{
            Value: value,
            Target: target,
            Position: start_pos,
        }
    }

    return value
}

func (p *Parser) parseAssignExpr() ast.Expression {
    start_pos := p.current().Position
    value := p.parseForwardPipeExpr()

    for p.current().Text == "=" ||
        p.current().Text == "+=" ||
        p.current().Text == "-=" ||
        p.current().Text == "*=" ||
        p.current().Text == "**=" ||
        p.current().Text == "/=" || 
        p.current().Text == "//=" || 
        p.current().Text == "%=" || 
        p.current().Text == "!=" || 
        p.current().Text == "~=" || 
        p.current().Text == "&=" || 
        p.current().Text == "|=" || 
        p.current().Text == "^=" || 
        p.current().Text == "<<=" || 
        p.current().Text == ">>=" {

        operator := p.current().Text
        p.consume()
        target := p.parseAssignExpr()
        return ast.AssignExpr{
            Name: value,
            Value: target,
            Operator: operator,
            Position: start_pos,
        }
    }

    return value
}

func (p *Parser) parseUnaryExpr() ast.Expression {
    start_pos := p.current().Position
    operator := p.consume().Text
    operand := p.parseExpression()

    return ast.UnaryExpr{
        Operator: operator,
        Operand: operand,
        Position: start_pos,
    }
}
