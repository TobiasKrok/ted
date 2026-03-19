package input

import (
	"bufio"
	"fmt"
)

type Key int

const (
	KeyChar = iota
	KeyEsc
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
	if char == '\x1b' {
		if r.Buffered() == 0 {
			return KeyEvent{Type: KeyEsc}, nil
		}
		fmt.Println("RUNE:")
		fmt.Println(char)
		return KeyEvent{Type: KeyEsc}, nil
	}

	return KeyEvent{Type: KeyChar, Char: char}, nil
}
