package parser

import "codeberg.org/Tanzanite/Tanzanite/ast"

func (p *Parser) parseAdditiveExpr() ast.Expression {
    left := p.parseMultiplicativeExpr()

    for p.current().Text == "+" || p.current().Text == "-" {
        operator := p.consume().Text
        right := p.parseMultiplicativeExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseMultiplicativeExpr() ast.Expression {
    left := p.parsePrimaryExpr()

    for p.current().Text == "/" ||
        p.current().Text == "*" ||
        p.current().Text == "**" ||
        p.current().Text == "//" ||
        p.current().Text == "%" {
        operator := p.consume().Text
        right := p.parsePrimaryExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseShiftExpr() ast.Expression {
    left := p.parseAdditiveExpr()

    for p.current().Text == "<<" || p.current().Text == ">>" {
        operator := p.consume().Text
        right := p.parseAdditiveExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseComparativeExpr() ast.Expression {
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
        }
    }

    return left
}

func (p *Parser) parseEqaulityExpr() ast.Expression {
    left := p.parseComparativeExpr()

    for p.current().Text == "==" ||
        p.current().Text == "!=" {
        operator := p.consume().Text
        right := p.parseComparativeExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseSpaceshipExpr() ast.Expression {
    left := p.parseEqaulityExpr()

    for p.current().Text == "<=>" {
        operator := p.consume().Text
        right := p.parseEqaulityExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseBitwiseAndExpr() ast.Expression {
    left := p.parseSpaceshipExpr()

    for p.current().Text == "&" {
        operator := p.consume().Text
        right := p.parseSpaceshipExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseBitwiseXorExpr() ast.Expression {
    left := p.parseBitwiseAndExpr()

    for p.current().Text == "^" {
        operator := p.consume().Text
        right := p.parseBitwiseAndExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseBitwiseOrExpr() ast.Expression {
    left := p.parseBitwiseXorExpr()

    for p.current().Text == "|" {
        operator := p.consume().Text
        right := p.parseBitwiseXorExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseLogicalAndExpr() ast.Expression {
    left := p.parseBitwiseOrExpr()

    for p.current().Text == "&&" {
        operator := p.consume().Text
        right := p.parseBitwiseOrExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseLogicalOrExpr() ast.Expression {
    left := p.parseLogicalAndExpr()

    for p.current().Text == "||" {
        operator := p.consume().Text
        right := p.parseLogicalAndExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseConditionalExpr() ast.Expression {
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
        }
    }

    return condition
}

func (p *Parser) parseForwardPipeExpr() ast.Expression {
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
        }
    }

    return value
}

func (p *Parser) parseAssignExpr() ast.Expression {
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
        }
    }

    return value
}

func (p *Parser) parseUnaryExpr() ast.Expression {
    operator := p.consume().Text
    operand := p.parseExpression()

    return ast.UnaryExpr{
        Operator: operator,
        Operand: operand,
    }
}
