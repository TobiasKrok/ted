package main

import (
	"bufio"
	"fmt"
	"os"
	"ted/internal/editor"
	"ted/internal/input"
	"ted/internal/term"
)

func main() {
	term, err := term.NewTerminal()
	if err != nil {
		panic(err)
	}
	err = term.EnableRawMode()
	if err != nil {
		panic(err)
	}
	editor := editor.NewTedEditor(term)
	defer func() {
		term.DisableRawMode()
	}()
	r := bufio.NewReader(os.Stdin)
	term.Clear()

	for {

		editor.Render()
		key, err := input.ReadInput(r)
		if err != nil {
			fmt.Errorf("%w", err)
			continue
		}
		q := editor.HandleKey(key)
		if q == 1 {
			break
		}
	}

}
