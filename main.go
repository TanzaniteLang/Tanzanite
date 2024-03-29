package main

import (
    "fmt"
    "os"
    "os/exec"
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/ccg"
    "github.com/gookit/goutil/dump"
)

func main() {
    cmdArgs := os.Args[1:]

    if len(cmdArgs) != 1 && len(cmdArgs) != 3 {
        fmt.Fprintln(os.Stderr, "Expected only 1 or 3 arguments!")
        fmt.Fprintln(os.Stderr, "help: tanzanite [file] (-o [output])")
        os.Exit(1)
    }

    code, err := os.ReadFile(cmdArgs[0])
    if err != nil {
        fmt.Print(err)
        os.Exit(1)
    }

    par := parser.NewParser()
    out := par.ProduceAST(string(code))

    dump.Config(func (o *dump.Options) {
        o.MaxDepth = 100
    })

    // dump.Println(out)

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
        f, err := os.CreateTemp("/tmp/", output + ".*.c")
        if err != nil {
            fmt.Print(err)
            os.Exit(1)
        }
        defer os.Remove(f.Name())

        if _, err := f.Write([]byte(src.Generate())); err != nil {
            fmt.Print(err)
            os.Exit(1)
        }

        f.Close()

        cmd := exec.Command("tcc", f.Name(), "-o", output)
        err = cmd.Run()

        if err != nil {
            fmt.Print(err)
            os.Exit(1)
        }
    } else {
        fmt.Println(src.Generate())
    }
}
