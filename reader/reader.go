package reader

import (
    "io"
    "errors"
)

type Reader struct {
    runes []rune
    length int
    pos int
}

func NewReader(text string) *Reader {
    return &Reader {
        runes: []rune(text),
        length: len(text),
        pos: 0,
    }
}

func (r *Reader) ReadRune() (rune, int, error) {
    if r.pos == r.length {
        return ' ', 0, io.EOF
    }

    cur := r.runes[r.pos]
    r.pos++

    return cur, 4, nil
}

func (r *Reader) UnreadRune() error {
    r.pos--
    if r.pos < 0 {
        return errors.New("Cannot unread rune before start!")
    }

    return nil
}
