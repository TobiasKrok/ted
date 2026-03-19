package editor

import (
	"fmt"
)

type CursorStyle int
type CursorMoveDir int

const (
	CursorMoveRight CursorMoveDir = iota
	CursorMoveLeft
	CursorMoveUp
	CursorMoveDown
)
const (
	CursorStyleInsert CursorStyle = iota
	CursorStyleNormal
)

type Cursor struct {
	Row   int
	Col   int
	style CursorStyle
}

func NewCursor(tRow, tCol int) *Cursor {

	return &Cursor{
		Row:   0,
		Col:   0,
		style: CursorStyleNormal,
	}
}

// https://vt100.net/docs/vt510-rm/DECSCUSR.html
func (c *Cursor) setStyle(s CursorStyle) {
	var sequence string
	switch s {
	case CursorStyleInsert:
		sequence = "\x1b[5 q"
	case CursorStyleNormal:
		sequence = "\x1b[2 q"
	}
	fmt.Print(sequence)
}

// TODO: limits?
func (t *Cursor) move(dir CursorMoveDir, count int) {

	switch dir {
	case CursorMoveLeft:
		t.Col -= count
	case CursorMoveRight:
		t.Col += count
	case CursorMoveDown:
		t.Row += count
	case CursorMoveUp:
		t.Row -= count
	}
}

func (t *Cursor) moveHome() {
	fmt.Print("\x1b[0;2H")
}
