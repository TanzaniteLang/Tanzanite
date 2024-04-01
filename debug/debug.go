package debug

import (
    "path/filepath"
    "strings"
    "os"
    "fmt"
)

var fileCache map[string][]byte 

type SourceLocation struct {
    Line uint64
    Column uint64
    File string
}

type Hint struct {
    Msg string
    Code string
}

func NewSourceLocation(file string, line uint64, col uint64) SourceLocation {
    if fileCache == nil {
        fileCache = map[string][]byte{}
    }
    return SourceLocation{
        File: file,
        Line: line,
        Column: col,
    }
}

func (s *SourceLocation) Stringify() string {
    return fmt.Sprintf("#line %d \"%s\"", s.Line, s.File)
}

func (s *SourceLocation) ThrowError(msg string, top bool, hint *Hint) {
    if !top {
        fmt.Fprintln(os.Stderr, "──────────────────────────────────────")
    }

    if hint != nil {
        n, _ := fmt.Fprintf(os.Stderr, "[%s:%d:%d]", filepath.Base(s.File), s.Line, s.Column)
        fmt.Fprintf(os.Stderr, " \x1b[31;1mError\x1b[0m: %s\n", msg)
        fmt.Fprintf(os.Stderr, "%*c \x1b[34;1mHint\x1b[0m: %s\n", n, ' ', hint.Msg)
    } else {
        fmt.Fprintf(os.Stderr, "[%s:%d:%d] \x1b[31;1mError\x1b[0m: %s\n", filepath.Base(s.File), s.Line, s.Column, msg)
    }

    code, ok := fileCache[s.File]

    if !ok {
        code, _ = os.ReadFile(s.File)
        fileCache[s.File] = code
    }

    lines := strings.Split(string(code), "\n")

    fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━ code ━━━━━━━━━━━━━━━━")

    startLine := int64(s.Line - 5)
    if startLine < 0 {
        startLine = 0
    }

    width := len(fmt.Sprintf("%d", s.Line))

    if hint != nil {
        for i := uint64(startLine); i < s.Line - 1; i++ {
            fmt.Fprintf(os.Stderr, "%*d | %s\n", width, i + 1, lines[i])
        }
        fmt.Fprintf(os.Stderr, "%*d | %s\x1b[32;1m%s\x1b[0m%s\n", width, s.Line, lines[s.Line - 1][:s.Column - 1], hint.Code, lines[s.Line - 1][s.Column - 1:])

        fmt.Fprintf(os.Stderr, "%*c |%*c", width, ' ', s.Column, ' ')
        fmt.Fprintf(os.Stderr, "\x1b[32;1m%s\x1b[0m\n", strings.Repeat("+", len(hint.Code)))
    } else {
        for i := uint64(startLine); i < s.Line; i++ {
            fmt.Fprintf(os.Stderr, "%*d | %s\n", width, i + 1, lines[i])
        }
        fmt.Fprintf(os.Stderr, "%*c |%*c", width, ' ', s.Column, ' ')
        fmt.Fprintf(os.Stderr, "^\n")
    }
    fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func (s *SourceLocation) ThrowWarning(msg string, top bool, hint *Hint) {
    if !top {
        fmt.Fprintln(os.Stderr, "──────────────────────────────────────")
    }

    if hint != nil {
        n, _ := fmt.Fprintf(os.Stderr, "[%s:%d:%d]", filepath.Base(s.File), s.Line, s.Column)
        fmt.Fprintf(os.Stderr, " \x1b[33;1mWarning\x1b[0m: %s\n", msg)
        fmt.Fprintf(os.Stderr, "%*c \x1b[34;1mHint\x1b[0m: %s\n", n, ' ', hint.Msg)
    } else {
        fmt.Fprintf(os.Stderr, "[%s:%d:%d] \x1b[33;1Warning\x1b[0m: %s\n", filepath.Base(s.File), s.Line, s.Column, msg)
    }

    code, ok := fileCache[s.File]

    if !ok {
        code, _ = os.ReadFile(s.File)

        fileCache[s.File] = code
    }

    lines := strings.Split(string(code), "\n")

    fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━ code ━━━━━━━━━━━━━━━━")

    startLine := int64(s.Line - 5)
    if startLine < 0 {
        startLine = 0
    }

    width := len(fmt.Sprintf("%d", s.Line))

    if hint != nil {
        for i := uint64(startLine); i < s.Line - 1; i++ {
            fmt.Fprintf(os.Stderr, "%*d | %s\n", width, i + 1, lines[i])
        }
        fmt.Fprintf(os.Stderr, "%*d | %s\x1b[32;1m%s\x1b[0m%s\n", width, s.Line, lines[s.Line - 1][:s.Column - 1], hint.Code, lines[s.Line - 1][s.Column - 1:])

        fmt.Fprintf(os.Stderr, "%*c |%*c", width, ' ', s.Column, ' ')
        fmt.Fprintf(os.Stderr, "\x1b[32;1m%s\x1b[0m\n", strings.Repeat("+", len(hint.Code)))
    } else {
        for i := uint64(startLine); i < s.Line; i++ {
            fmt.Fprintf(os.Stderr, "%*d | %s\n", width, i + 1, lines[i])
        }
        fmt.Fprintf(os.Stderr, "%*c |%*c", width, ' ', s.Column, ' ')
        fmt.Fprintf(os.Stderr, "^\n")
    }
    fmt.Fprintln(os.Stderr, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
