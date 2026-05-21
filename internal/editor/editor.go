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

type TedEditor struct {
	term   *term.Terminal
	cursor *Cursor
	mode   mode
	buf    *buffer.GapBuffer
	lines  []int
	debug  strings.Builder
}

func NewTedEditor(term *term.Terminal) *TedEditor {
	cursor := NewCursor()
	// buf := buffer.NewGapBuffer(1024 * 4)

	buf := buffer.NewGapBuffer(100)
	return &TedEditor{
		term:   term,
		cursor: cursor,
		mode:   ModeNormal,
		buf:    buf,
		lines:  []int{0}, // it contains the row/col position of the 
	}
}

func (t *TedEditor) Render() error {
	t.term.Clear()
	t.cursor.Hide()
	t.cursor.moveHome()
	fmt.Println(t.buf.String())
	fmt.Printf("\n\nDEBUG: %s", t.debug.String())
	t.debug.Reset()
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

// this translates the cursor position into the correct buffer index
// KEEP IN MIND: the order of calls matter since it uses the alreayd calculated coordinates of the cursor. maybe we need a peek at some point
func (t *TedEditor) cursorToBufIdx() int {
	//TODO: this probably needs change when we have scroll implemented
	offset := t.lines[t.cursor.Row]

	return offset + 

}

//TODO: recalculate newlines

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
		if t.cursor.Col+1 >= t.term.Size.Cols {
			t.cursor.move(CursorMoveDown, 1)
			t.cursor.setPos(t.cursor.Row, 0) // at first col
		} else {
			t.cursor.move(CursorMoveRight, 1)
		}
	case input.KeyEnter:
		t.buf.Insert('\n')
		t.lines = append(t.lines, t.cursor.Row+1)
		t.cursor.setPos(t.cursor.Row+1, 0)
	}
	t.debug.WriteString(fmt.Sprintf("%v %v", t.term.Size.Cols, t.buf.InternalLength()))
	return 0
}

func (t *TedEditor) handleNormalMode(key input.KeyEvent) int {

	switch key.Char {
	case 'q':
		return 1 // quit
	case 'i':
		t.mode = ModeInsert
		t.cursor.setStyle(CursorStyleInsert)
		err := t.buf.MoveGap(t.cursorToBufIdx())
		if err != nil {
			//TODO: handle error
		}

	case 'a':
		t.mode = ModeInsert
		t.cursor.setStyle(CursorStyleInsert)
		if t.cursor.Col+1 < t.term.Size.Cols {
			err := t.buf.MoveGap(t.cursorToBufIdx() + 1)
			if err != nil {
				//TODO: handle
			} else {
				t.cursor.move(CursorMoveRight, 1)
			}
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
