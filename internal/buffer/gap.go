package buffer

import "fmt"

type GapBuffer struct {
	data   []rune
	start  int
	end    int
	cursor int
}

const NEW_GAP_SIZE = 64
const MAX_SIZE = 1024 * 8

func NewGapBuffer(initialCap int) *GapBuffer {

	d := make([]rune, initialCap)
	return &GapBuffer{
		data:   d,
		start:  0,
		end:    initialCap,
		cursor: 0,
	}
}

func (g *GapBuffer) Insert(r rune) {

	if g.start+1 == g.end {
		g.grow()
	}
	g.data[g.cursor] = r
	g.cursor += 1
	g.start += 1
}

func (g *GapBuffer) DeleteBackwards() {

	if g.start == 0 {
		return
	}
	g.start -= 1
	g.cursor -= 1
}

// Delete a range in the buffer, start and end are logical positions
func (g *GapBuffer) DeleteRange(start, end int) error {
	ps := g.toPhysical(min(start, end))
	pe := g.toPhysical(max(start, end))

	if ps < 0 || pe > len(g.data) {
		return fmt.Errorf("deletion range out of bounds: start=%d end=%d buflen=%d", ps, pe, len(g.data))
	}
	// A B C D E [____]  F G H
	// |                   |
	// A [_________] H
	if ps < g.start && pe > g.start {
		g.start = ps
		g.cursor = ps
		if pe > g.end {
			g.end = pe

		}
	} else if ps < g.start && pe <= g.start {

		// A B C D E [____]  F G H
		//   |   |
		// AE [____] F G H
		bef := g.data[pe:g.start] // chars form the end of the deletion range to the start of the gap

		ns := ps + len(bef) // new gap start
		copy(g.data[ps:ns], bef)
		g.start = ns
		g.cursor = ns
	} else if ps >= g.end {
		//TODO: move cursor
		aft := g.data[g.end:ps]
		newEnd := pe - len(aft)
		copy(g.data[newEnd:], aft)
		g.end = newEnd
	}

	return nil
}

// Moves the cursor and the gap to the position
func (g *GapBuffer) MoveGap(pos int) error {
	p := g.toPhysical(pos)

	if p < 0 || p > len(g.data) {
		return fmt.Errorf("cannot move to index %d (physical %d) because its out of bounds", pos, p)
	}

	if p == g.start && p == g.cursor {
		return nil
	}

	// we're gonna move them before the gap
	if p < g.start {
		after := g.data[p:g.start]
		ne := g.end - len(after)
		copy(g.data[ne:g.end], after)
		g.start = p
		g.cursor = p
		g.end = ne
	} else if p > g.end {
		before := g.data[g.end:p]
		ns := g.start + len(before)
		copy(g.data[g.start:ns], before)
		g.start = ns
		g.cursor = ns
		g.end = p
	}
	return nil
}

// grows the gap and moves the cursor to the start
func (g *GapBuffer) grow() {
	newCap := len(g.data) + NEW_GAP_SIZE
	b := make([]rune, newCap)

	copy(b, g.data[:g.start]) // start of buffer tp gap start

	end := g.data[g.end:]
	copy(b[len(end):], end)

	g.cursor = g.start

	g.end = newCap
	g.data = b

}

// return the character at the position. return -1 if its out of bounds
func (g *GapBuffer) CharAt(pos int) rune {
	p := g.toPhysical(pos)

	if p < 0 || p >= len(g.data) {
		return -1
	}

	return g.data[p]
}

func (g *GapBuffer) String() string {
	return string(g.data[:g.start]) + string(g.data[g.end:])
}

func (g *GapBuffer) Length() int {
	return len(g.data) - g.gapLen()
}

func (g *GapBuffer) gapLen() int {

	return g.end - g.start
}

func (g *GapBuffer) toPhysical(pos int) int {

	p := pos
	gl := g.gapLen()

	if p >= g.start {
		p = pos + gl
	}
	return p
}
