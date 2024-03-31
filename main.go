package main

import (
    "fmt"
    "os"
    "os/exec"
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/ccg"
    "codeberg.org/Tanzanite/Tanzanite/debug"
    "github.com/gookit/goutil/dump"
)

func main() {
    cmdArgs := os.Args[1:]

    if len(cmdArgs) != 1 && len(cmdArgs) != 3 {
        fmt.Fprintln(os.Stderr, "Expected only 1 or 3 arguments!")
        fmt.Fprintln(os.Stderr, "help: tanzanite [file] (-o [output])")
        fmt.Fprintln(os.Stderr, "To see AST, set TZN_DBG=1 env variable, C code will be ommited")
        os.Exit(1)
    }

    code, err := os.ReadFile(cmdArgs[0])
    if err != nil {
        fmt.Print(err)
        os.Exit(1)
    }

    par := parser.NewParser(cmdArgs[0])
    out := par.ProduceAST(string(code))

    dump.Config(func (o *dump.Options) {
        o.MaxDepth = 100
    })

    dbg, ok := os.LookupEnv("TZN_DBG")
    if ok && dbg == "1" {
        dump.Println(out)
        dbg := debug.NewSourceLocation(cmdArgs[0], 7, 1)
        dbg.ThrowError("Oops")
        os.Exit(0)
    }

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
        f, err := os.Create(output + ".c")
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

        cmd := exec.Command("tcc", f.Name(), "-g", "-o", output)
        err = cmd.Run()

        if err != nil {
            fmt.Print(err)
            os.Exit(1)
        }
    } else {
        fmt.Println(src.Generate())
    }
}
