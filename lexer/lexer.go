package lexer

import (
    "io"
    "unicode"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/reader"
)

type Lexer struct {
    reader *reader.Reader
    pos tokens.Position
}

func InitLexer(code string) *Lexer {
    return &Lexer {
        reader: reader.NewReader(code),
        pos: tokens.Position{Line: 1, Column: 0},
    }
}

func (l *Lexer) Lex() (tokens.Position, tokens.Token, string) {
    for {
        current, _, err := l.reader.ReadRune()

        if err != nil {
            if err == io.EOF {
                return l.pos, tokens.Eof, ""
            }

            panic(err)
        }

        l.pos.Column++

        switch current {
        case '\n':
            l.newLine()
        case '#':
            l.skipComment()
        case '?':
            return l.pos, tokens.QuestionMark, "?"
        case '.':
            return l.pos, tokens.Dot, "."
        case ',':
            return l.pos, tokens.Comma, ","
        case ':':
            return l.pos, tokens.Colon, ":"
        case ';':
            return l.pos, tokens.Semicolon, ";"
        case '(':
            return l.pos, tokens.LBracket, "("
        case ')':
            return l.pos, tokens.RBracket, ")"
        case '[':
            return l.pos, tokens.LSquareBracket, "["
        case ']':
            return l.pos, tokens.RSquareBracket, "]"
        case '{':
            return l.pos, tokens.LSquiglyBracket, "}"
        case '}':
            return l.pos, tokens.RSquiglyBracket, "}"
        case '=':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '+':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '-':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '%':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '!':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '~':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '^':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '&':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '|':
            pos := l.pos
            tok, text := l.twoOperators(current)
            return pos, tok, text
        case '*':
            pos := l.pos
            tok, text := l.threeOperators(current)
            return pos, tok, text
        case '/':
            pos := l.pos
            tok, text := l.threeOperators(current)
            return pos, tok, text
        case '<':
            pos := l.pos
            tok, text := l.threeOperators(current)
            return pos, tok, text
        case '>':
            pos := l.pos
            tok, text := l.threeOperators(current)
            return pos, tok, text
        case '\'':
            pos := l.pos
            text := l.parseChar()
            return pos, tokens.CharVal, text
        case '"':
            pos := l.pos
            text := l.parseString()
            return pos, tokens.StringVal, text
        default:
            if unicode.IsSpace(current) {
                continue
            } else if unicode.IsDigit(current) {
                pos := l.pos
                l.undo()
                tok, text := l.parseNumber()
                return pos, tok, text
            } else if unicode.IsLetter(current) {
                pos := l.pos
                l.undo()
                tok, text := l.parseIdentifier()
                return pos, tok, text
            }

            return l.pos, tokens.Illegal, string(current)
        }
    }
}

func (l *Lexer) newLine() {
    l.pos.Line++
    l.pos.Column = 0
}

func (l *Lexer) undo() {
    if err := l.reader.UnreadRune(); err != nil {
        panic(err)
    }

    l.pos.Column--
}

func (l *Lexer) parseString() string {
    text := ""
    r, _, err := l.reader.ReadRune()
    if err != nil {
        return ""
    }

    l.pos.Column++
    for {
        next, _, err := l.reader.ReadRune()
        if err != nil {
            return text
        }
        l.pos.Column++

        if r == '\\' && next == '"' {
            text += "\\\""
            r, _, err = l.reader.ReadRune()
            if err != nil {
                return text
            }
            l.pos.Column++
            continue
        }

        l.undo()
        text += string(r)

        r, _, err = l.reader.ReadRune()
        if err != nil {
            return text
        }
        l.pos.Column++
        if r == '"' { break }
    }

    return text
}

func (l *Lexer) parseChar() string {
    // TODO: Ensure that it is closed with '
    r, _, err := l.reader.ReadRune()
    if err != nil {
        return ""
    }
    l.pos.Column++

    if r == '\\' {
        next, _, err := l.reader.ReadRune()
        if err != nil {
            return ""
        }
        l.pos.Column++

        _, _, err = l.reader.ReadRune()
        if err != nil {
            return ""
        }
        l.pos.Column++

        return "\\" + string(next)
    }

    _, _, err = l.reader.ReadRune()
    if err != nil {
        return ""
    }
    l.pos.Column++

    return string(r)
}

func (l *Lexer) parseIdentifier() (tokens.Token, string) {
    text := ""
    r, _, err := l.reader.ReadRune()
    if err != nil {
        return tokens.Eof, ""
    }

    l.pos.Column++

    for unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_' || r == '@' {
        l.pos.Column++
        text += string(r)
        r, _, err = l.reader.ReadRune()
        if err != nil {
            if err == io.EOF {
                if tok := tokens.Search(text); tok == tokens.Eof {
                    return tokens.Identifier, text
                } else {
                    return tok, text
                }
            }
            panic(err)
        }
    }

    l.undo()

    if tok := tokens.Search(text); tok == tokens.Eof {
        return tokens.Identifier, text
    } else {
        return tok, text
    }
}

func (l *Lexer) parseNumber() (tokens.Token, string) {
    text := ""
    is_float := false

    r, _, err := l.reader.ReadRune()
    if err != nil {
        return tokens.Eof, ""
    }

    l.pos.Column++

    for unicode.IsDigit(r) || r == '.' && !is_float || r == '_' {
        l.pos.Column++
        text += string(r)
        if r == '.' {
            is_float = true
        }

        r, _, err = l.reader.ReadRune()
        if err != nil {
            if err == io.EOF {
                if is_float {
                    return tokens.Float, text
                } else {
                    return tokens.Int, text
                }
            }
            panic(err)
        }
    }

    l.undo()
    if is_float {
        return tokens.FloatVal, text
    } else {
        return tokens.IntVal, text
    }
}

func (l *Lexer) skipComment() {
    r, _, err := l.reader.ReadRune()
    if err != nil {
        return
    }
    l.pos.Column++

    for r != '\n' {
        l.pos.Column++
        r, _, err = l.reader.ReadRune()
        if err != nil {
            return
        }
    }
    l.newLine()
}

func (l *Lexer) threeOperators(current rune) (tokens.Token, string) {
    r, _, err := l.reader.ReadRune()
    if err != nil {
        if err == io.EOF {
            return tokens.Search(string(current)), string(current)
        }
    } 

    r2, _, err := l.reader.ReadRune()
    if err != nil {
        if err == io.EOF {
            return tokens.Search(string(current)), string(current)
        }
    } 

    l.pos.Column += 2

    switch current {
    case '*':
        if r == '=' {
            l.undo();
            return tokens.AsteriskAssign, "*="
        } else if r == '*' {
            if r2 == '=' {
                return tokens.DoubleAsteriskAssign, "**="
            }
            l.undo();
            return tokens.DoubleAsterisk, "**"
        }
        l.undo();
        l.undo();
        return tokens.Asterisk, "*"
    case '/':
        if r == '=' {
            l.undo();
            return tokens.SlashAssign, "/="
        } else if r == '/' {
            if r2 == '=' {
                return tokens.DoubleSlashAssign, "//="
            }
            l.undo();
            return tokens.DoubleSlash, "//"
        }
        l.undo();
        l.undo();
        return tokens.Slash, "/"
    case '>':
        if r == '=' {
            l.undo();
            return tokens.GreaterEquals, ">="
        } else if r == '>' {
            if r2 == '=' {
                return tokens.RightShiftAssign, ">>="
            }
            l.undo();
            return tokens.RightShift, ">>"
        }
        l.undo();
        l.undo();
        return tokens.Greater, ">"
    case '<':
        if r == '=' {
            l.undo();
            return tokens.LessEquals, "<="
        } else if r == '<' {
            if r2 == '=' {
                return tokens.LeftShiftAssign, "<<="
            }
            l.undo();
            return tokens.LeftShift, "<<"
        }
        l.undo();
        l.undo();
        return tokens.Less, "<"
    }

    return tokens.Eof, ""
}

func (l *Lexer) twoOperators(current rune) (tokens.Token, string) {
    r, _, err := l.reader.ReadRune()
    if err != nil {
        if err == io.EOF {
            return tokens.Search(string(current)), string(current)
        }
    } 

    l.pos.Column++

    switch current {
    case '=':
        if r == '=' {
            return tokens.Equals, "=="
        }
        l.undo()
        return tokens.Assign, "="
    case '+':
        if r == '=' {
            return tokens.PlusAssign, "+="
        }
        l.undo()
        return tokens.Plus, "+"
    case '-':
        if r == '=' {
            return tokens.MinusAssign, "-="
        }
        l.undo()
        return tokens.Minus, "-"
    case '%':
        if r == '=' {
            return tokens.ModuloAssign, "%="
        }
        l.undo()
        return tokens.Modulo, "%"
    case '!':
        if r == '=' {
            return tokens.NotEquals, "!="
        }
        l.undo()
        return tokens.Bang, "!"
    case '~':
        if r == '=' {
            return tokens.TildaAssign, "~="
        }
        l.undo()
        return tokens.Tilda, "~"
    case '^':
        if r == '=' {
            return tokens.CaretAssign, "^="
        }
        l.undo()
        return tokens.Caret, "^"
    case '&':
        if r == '=' {
            return tokens.AmpersandAssign, "&="
        } else if r == '&' {
            return tokens.And, "&&"
        }
        l.undo()
        return tokens.Ampersand, "&"
    case '|':
        if r == '=' {
            return tokens.PipeAssign, "|="
        } else if r == '|' {
            return tokens.Or, "||"
        } else if r == '>' {
            return tokens.PipeTo, "|>"
        }
        l.undo()
        return tokens.Pipe, "|"
    }

    return tokens.Eof, ""
}
