package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/lexer"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

type GlobalScope struct {
    Scope map[string]*ast.FunctionDecl
}

func (g *GlobalScope) RegisterFunction(name string, fn *ast.FunctionDecl) {
    g.Scope[name] = fn
}

func (g *GlobalScope) HasFunction(name string) bool {
    _, ok := g.Scope[name]

    return ok
}

type Token struct {
    Info tokens.Token
    Position tokens.Position
    Text string
}

type Parser struct {
    tokens []Token
    Globals GlobalScope

    scopes []*ast.Body

    pos int
    source string
    parsingFn bool // Breaks reccursion
    requireBrackets bool

    // Programmer info
    warn bool
    Dead bool
}

func NewParser(file string) *Parser {
    return &Parser{
        tokens: []Token{},
        Globals: GlobalScope{
            Scope: map[string]*ast.FunctionDecl{},
        },
        source: file,
        scopes: []*ast.Body{},
        pos: 0,
        parsingFn: false,
        requireBrackets: false,
        warn: false,
        Dead: false,
    }
}

func (p *Parser) findVariable(name string) *ast.VarDeclaration {
    last := len(p.scopes) - 1
    for last >= 0 {
        decl, ok := p.scopes[last].Scope[name]
        if ok {
            return decl
        }
        last--
    }
    return nil
}

func (p *Parser) RegisterVar(name string, decl *ast.VarDeclaration) {
    last := p.scopes[len(p.scopes) - 1]
    last.RegisterVar(name, decl)
}

func (p *Parser) notEof() bool {
    return p.tokens[p.pos].Info != tokens.Eof
}

func (p *Parser) current() Token {
    return p.tokens[p.pos]
}

func (p *Parser) consume() Token {
    p.pos++
    return p.tokens[p.pos - 1]
}

func (p *Parser) previous() Token {
    if p.pos == 0 {
        return p.tokens[0]
    }

    return p.tokens[p.pos - 1]
}

func (p *Parser) skipToNewLine() {
    pos := p.previous().Position.Line

    for p.current().Position.Line == pos {
        p.pos++
    }
}

func (p *Parser) AppendScope(scope *ast.Body) {
    p.scopes = append(p.scopes, scope)
}

func (p *Parser) PopScope() {
    p.scopes = p.scopes[:len(p.scopes) - 1]
}

func (p *Parser) ProduceAST(code string) ast.Program {
    lex := lexer.InitLexer(code)

    for {
        pos, tok, text := lex.Lex()

        p.tokens = append(p.tokens, Token { Info: tok, Position: pos, Text: text})

        if tok == tokens.Eof {
            break
        }
    }

    prog := ast.Program {
        Body: ast.Body{
            Scope: map[string]*ast.VarDeclaration{},
            Body: []ast.Statement{},
        },
    }

    p.AppendScope(&prog.Body)

    for p.notEof() {
        prog.Body.Append(p.parseStatement())
    }

    return prog
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.current().Info {
    case tokens.Def:
        panic("Def functions are not yet implemented!")
    case tokens.Fun:
        p.consume()
        fn := p.parseFunction(true).(ast.FunctionDecl)
        p.Globals.RegisterFunction(fn.Name, &fn)

        return fn
    case tokens.Break:
        p.consume()
        return ast.LoopControlStatement{
            Break: true,
        }
    case tokens.Next:
        p.consume()
        return ast.LoopControlStatement{
            Break: false,
        }
    case tokens.If:
        return p.parseIf(false)
    case tokens.Unless:
        return p.parseIf(true)
    case tokens.While:
        return p.parseWhile(false)
    case tokens.Until:
        return p.parseWhile(true)
    case tokens.Begin:
        return p.parseBegin()
    case tokens.Return:
        p.consume()
        expr := p.parseExpression()
        if expr == nil {
            c := p.previous()
            dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1 + uint64(len(c.Text)))
            dbg.ThrowError("Return statement is missing a value!", p.warn || p.Dead, nil)
            p.Dead = true
            p.skipToNewLine()
        }

        return ast.ReturnExpr{
            Value: expr,
        }
    case tokens.Identifier:
        if p.Globals.HasFunction(p.current().Text) {
            fn, _ := p.Globals.Scope[p.current().Text]
            if fn.Failed {
                c := p.current()
                p.consume()

                dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
                dbg.ThrowError("Function \"" + c.Text + "\" failed to parse!", p.warn || p.Dead, &debug.Hint{
                    Msg: "Fix the function before using it", 
                    Code: "",
                })
                p.Dead = true
                p.skipToNewLine()
                return nil
            }
            return p.parseFnCall(fn)
        }

        decl := p.findVariable(p.current().Text)

        if decl == nil {
            possiblyVar := p.parseVarDeclaration()
            if possiblyVar == nil {
                return nil
            }
            stmt := possiblyVar.(ast.VarDeclaration)
            p.RegisterVar(stmt.Name, &stmt)
            return stmt
        }
        return p.parseAssignExpr()
    default:
        return p.parseExpression()
    }
}

func (p *Parser) parseExpression() ast.Expression {
    return p.parseAssignExpr()
}
