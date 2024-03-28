package tokens

import (
    "slices"
)

type Position struct {
    Line uint64
    Column uint64
}

type Token int

const (
	Eof = iota
	Whitespace
	Illegal

	// Identifier and literals
	Identifier
	String
	Char
	Int
	Float
	Bool
	Nil
	Void
	Command

	// Statement
	Assign
	If
	Else
	Elsif
	For
	While
	Unless
	Until
	Break
	Next
	Case
	When
	With
	Begin
	In

	// Operators
	Equals
	Plus
	PlusAssign
	Minus
	MinusAssign
	Asterisk
	AsteriskAssign
	DoubleAsterisk
	DoubleAsteriskAssign
	Slash
	SlashAssign
	DoubleSlash
	DoubleSlashAssign
	Modulo
	ModuloAssign
	Bang
	NotEquals
	Tilda
	TildaAssign
	Ampersand
	AmpersandAssign
	And
	Pipe
	PipeAssign
	Or
	PipeTo
	Caret
	CaretAssign
	Less
	LessEquals
	LeftShift
	LeftShiftAssign
	Spaceship
	Greater
	GreaterEquals
	RightShift
	RightShiftAssign

	// Delimiters
	QuestionMark
	Dot
	Comma
	Colon
	Semicolon
	LBracket
	RBracket
	LSquareBracket
	RSquareBracket
	LSquiglyBracket
	RSquiglyBracket
	Do
	End
	Then

	// Reserved keywords
	Def
	Fun
	Return
	Sizeof
)

var tokens = []string{
	Assign:               "=",
	Equals:               "==",
	Plus:                 "+",
	PlusAssign:           "+=",
	Minus:                "-",
	MinusAssign:          "-=",
	Asterisk:             "*",
	AsteriskAssign:       "*=",
	DoubleAsterisk:       "**",
	DoubleAsteriskAssign: "**=",
	Slash:                "/",
	SlashAssign:          "/=",
	DoubleSlash:          "//",
	DoubleSlashAssign:    "//=",
	Modulo:               "%",
	ModuloAssign:         "%=",
	Bang:                 "!",
	NotEquals:            "!=",
	Tilda:                "~",
	TildaAssign:          "~=",
	Ampersand:            "&",
	AmpersandAssign:      "&=",
	And:                  "&&",
	Pipe:                 "|",
	PipeAssign:           "|=",
	Or:                   "||",
	PipeTo:               "|>",
	Caret:                "^",
	CaretAssign:          "^=",
	Less:                 "<",
	LessEquals:           "<=",
	LeftShift:            "<<",
	LeftShiftAssign:      "<<=",
	Spaceship:            "<=>",
	Greater:              ">",
	GreaterEquals:        ">=",
	RightShift:           ">>",
	RightShiftAssign:     ">>=",
	If:                   "if",
	Else:                 "else",
	Elsif:                "elsif",
	For:                  "for",
	While:                "while",
	Until:                "until",
	Unless:               "unless",
	Break:                "break",
	Next:                 "next",
	Case:                 "case",
	When:                 "when",
	Def:                  "def",
	Fun:                  "fun",
	Nil:                  "nil",
	Return:               "return",
	Sizeof:               "sizeof",
	With:                 "with",
	Begin:                "begin",
	In:                   "in",
	String:               "String",
    Identifier:           "Identifier",
	Char:                 "Char",
	Int:                  "Int",
	Float:                "Float",
	Bool:                 "Bool",
	Void:                 "Void",
	QuestionMark:         "?",
	Dot:                  ".",
	Comma:                ",",
	Colon:                ":",
	Semicolon:            ";",
	LBracket:             "(",
	RBracket:             ")",
	LSquareBracket:       "[",
	RSquareBracket:       "]",
	LSquiglyBracket:      "{",
	RSquiglyBracket:      "}",
	Do:                   "do",
	End:                  "end",
	Then:                 "then",
}

func (t Token) String() string {
	return tokens[t]
}

func Search(text string) Token {
    if text == "true" || text == "false" {
        return Bool
    }

    index := slices.Index(tokens, text)

    if index == -1 {
        return Eof
    }

    return Token(index)
}
