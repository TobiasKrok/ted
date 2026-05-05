package input

import (
	"bufio"
)

type Key int

const (
	KeyChar = iota
	KeyEsc
	KeyBackspace
	KeyEnter
)

type KeyEvent struct {
	Type Key
	Char rune
}

func ReadInput(r *bufio.Reader) (KeyEvent, error) {

	char, _, err := r.ReadRune()
	if err != nil {
		return KeyEvent{}, err
	}

	// escape
	if char == '\x1b' && r.Buffered() == 0 {
		return KeyEvent{Type: KeyEsc}, nil
	} else if char == 127 || char == 8 {
		return KeyEvent{Type: KeyBackspace}, nil
	} else if char == '\r' || char == '\n' {
		return KeyEvent{Type: KeyEnter}, nil
	}

	return KeyEvent{Type: KeyChar, Char: char}, nil
}
