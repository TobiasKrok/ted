package editor

import (
	"fmt"
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
	fmt.Printf("\x1b[%d;%dH", t.cursor.Row, t.cursor.Col)
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

func (t *TedEditor) HandleKey(key input.KeyEvent) {

	if t.mode == ModeNormal {
		t.handleNormalMode(key)
	} else if t.mode == ModeInsert {
		t.handleInsertMode(key)
	}
}

func (t *TedEditor) handleInsertMode(key input.KeyEvent) {

	if key.Type == input.KeyEsc {
		t.mode = ModeNormal
		t.cursor.setStyle(CursorStyleNormal)
		return
	}

	fmt.Printf("HEL: %v", key.Char)
	switch key.Char {
	case '\b':
		t.buf.DeleteBackwards()
	default:
		t.buf.Insert(key.Char)
	}
}

func (t *TedEditor) handleNormalMode(key input.KeyEvent) {

	switch key.Char {
	case 'i':
		t.mode = ModeInsert
		t.cursor.setStyle(CursorStyleInsert)
	case 'a':
		t.mode = ModeInsert

		if t.cursor.Col < t.term.Size.Cols {
			t.cursor.move(CursorMoveRight, 1)
		}
		// movement
	case 'l':
		if t.cursor.Col+1 < t.term.Size.Cols {
			t.cursor.move(CursorMoveRight, 1)
		}
	case 'h':
		if t.cursor.Col-1 > 1 {
			t.cursor.move(CursorMoveLeft, 1)
		}
	case 'k':
		if t.cursor.Row-1 > 1 {
			t.cursor.move(CursorMoveUp, 1)
		}
	case 'j':
		if t.cursor.Row+1 < t.term.Size.Rows {

			t.cursor.move(CursorMoveDown, 1)
		}
	}
}
