package debug

import "fmt"

type SourceLocation struct {
    Line uint64
    File string
}

func NewSourceLocation(file string, line uint64) SourceLocation {
    return SourceLocation{
        File: file,
        Line: line,
    }
}

func (s *SourceLocation) Stringify() string {
    return fmt.Sprintf("#line %d \"%s\"", s.Line, s.File)
}
