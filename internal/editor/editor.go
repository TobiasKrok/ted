package editor

import (
	"fmt"
	"strings"
	"ted/internal/buffer"
	"ted/internal/input"
	"ted/internal/term"
)

type mode int

const (
	ModeNormal mode = iota
	ModeInsert
)

type lines struct {
	chars []rune
}

type TedEditor struct {
	term   *term.Terminal
	cursor *Cursor
	mode   mode
	buf    *buffer.GapBuffer
	lines  int
	debug  strings.Builder
}

func NewTedEditor(term *term.Terminal) *TedEditor {
	cursor := NewCursor()
	// buf := buffer.NewGapBuffer(1024 * 4)

	buf := buffer.NewGapBuffer(10)
	return &TedEditor{
		term:   term,
		cursor: cursor,
		mode:   ModeNormal,
		buf:    buf,
	}
}

func (t *TedEditor) Render() error {
	t.term.Clear()
	t.cursor.Hide()
	t.cursor.moveHome()
	fmt.Println(t.buf.String())
	fmt.Printf("\n\nDEBUG: %s", t.debug.String())
	fmt.Printf("\x1b[%d;%dH", t.cursor.Row+1, t.cursor.Col+1)
	t.cursor.Show()

	// var sb strings.Builder
	// for i := range t.term.Size.Cols {
	// 	if i == t.term.Size.Cols-1 {
	// 		sb.WriteString("~")
	// 	} else {
	// 		sb.WriteString("~\r\n")
	// 	}
	// }
	// fmt.Print(sb.String())
	// t.cursor.moveHome()
	return nil
}

func (t *TedEditor) HandleKey(key input.KeyEvent) int {

	if t.mode == ModeNormal {
		return t.handleNormalMode(key)
	} else if t.mode == ModeInsert {
		return t.handleInsertMode(key)
	}
	return 0
}

func (t *TedEditor) handleInsertMode(key input.KeyEvent) int {

	switch key.Type {
	case input.KeyBackspace:
		t.buf.DeleteBackwards()
		t.cursor.move(CursorMoveLeft, 1)
	case input.KeyEsc:
		t.mode = ModeNormal
		t.cursor.setStyle(CursorStyleNormal)
	case input.KeyChar:

		t.buf.Insert(key.Char)
		//TODO: scroll x
		t.cursor.move(CursorMoveRight, 1)
	}

	return 0
}

func (t *TedEditor) handleNormalMode(key input.KeyEvent) int {

	switch key.Char {
	case 'q':
		return 1 // quit
	case 'i':
		t.mode = ModeInsert
		t.cursor.setStyle(CursorStyleInsert)
		err := t.buf.MoveGap(t.cursor.Col)
		if err != nil {
			//TODO: handle error
		}

	case 'a':
		t.mode = ModeInsert
		t.cursor.setStyle(CursorStyleInsert)
		if t.cursor.Col < t.term.Size.Cols {
			t.cursor.move(CursorMoveRight, 1)
		}
		// movement
	case 'l':
		newCol := t.cursor.Col + 1

		if newCol < t.term.Size.Cols && t.buf.CharAt(newCol) != -1 {
			t.cursor.move(CursorMoveRight, 1)
		}
	case 'h':
		if t.cursor.Col-1 >= 0 {
			t.cursor.move(CursorMoveLeft, 1)
		}
	case 'k':
		if t.cursor.Row-1 >= 0 {
			t.cursor.move(CursorMoveUp, 1)
		}
	case 'j':
		if t.cursor.Row+1 < t.term.Size.Rows {

			t.cursor.move(CursorMoveDown, 1)
		}
	}
	return 0
}
