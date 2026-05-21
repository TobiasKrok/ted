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
	Row    int
	Col    int
	style  CursorStyle
	Hidden bool
}

func NewCursor() *Cursor {

	return &Cursor{
		Row:    0,
		Col:    0,
		Hidden: false,
		style:  CursorStyleNormal,
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
func (c *Cursor) move(dir CursorMoveDir, count int) {

	switch dir {
	case CursorMoveLeft:
		c.Col -= count
	case CursorMoveRight:
		c.Col += count
	case CursorMoveDown:
		c.Row += count
	case CursorMoveUp:
		c.Row -= count
	}
}

func (c *Cursor) setPos(row, col int) {
	c.Col = col
	c.Row = row
}

func (c *Cursor) moveHome() {
	fmt.Print("\x1b[H")
}

func (c *Cursor) Show() {
	c.Hidden = false
	fmt.Print("\x1b[?25h")
}
func (c *Cursor) Hide() {
	c.Hidden = true
	fmt.Print("\x1b[?25l")
}
