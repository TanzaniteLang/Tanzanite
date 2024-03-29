package main

import (
    "fmt"
    "os"
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/ccg"
)

func main() {
    cmdArgs := os.Args[1:]

    if len(cmdArgs) != 1 && len(cmdArgs) != 3 {
        fmt.Println("Expected only 1 or 3 arguments!")
        fmt.Println("help: tanzanite [file] (-o [output])")
        os.Exit(1)
    }

    code, err := os.ReadFile(cmdArgs[0])
    if err != nil {
        fmt.Print(err)
        os.Exit(1)
    }

    par := parser.NewParser()
    out := par.ProduceAST(string(code))

    output := ""

    if len(cmdArgs) == 3 {
        output = cmdArgs[2]
    }

    src := ccg.NewSource("")

    for _, stmt := range out.Body {
        if stmt.GetKind() == ast.FunctionDeclType {
            src.Functions = append(src.Functions, stmt.(ast.FunctionDecl))
        }
    }

    if len(output) > 0 {
        err := os.WriteFile(output, []byte(src.Generate()), 0666)
        if err != nil {
            fmt.Print(err)
            os.Exit(1)
        }
    } else {
        fmt.Println(src.Generate())
    }
}
