package main

import (
    "fmt"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/lexer"
)

func main() {
    lex := lexer.InitLexer(`fun hello(msg: *Char = "World")
    puts "Hello, #{msg}!"
end`)

    for {
        pos, tok, text := lex.Lex()

        if tok == tokens.Eof {
            break
        }
        fmt.Println(pos, tok.String(), text)
    }
}
