package buffer

type GapBuffer struct {
	data  []rune
	start int
	end   int
	// length int // ?
	cursor int

	chbef int // chars before cursor
}

const NEW_GAP_SIZE = 64
const MAX_SIZE = 1024 * 4

func NewGapBuffer(initialCap int) *GapBuffer {

	d := make([]rune, initialCap)
	return &GapBuffer{
		data:   d,
		start:  0,
		end:    initialCap,
		cursor: 0,
		// length: 0,
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

func (g *GapBuffer) DeleteRange(start, end int) {
	s := min(start, end)
	e := max(start, end)

	if s < 0 || e > len(g.data) {
		return
	}

	// its possible we delete inside the gap
	// dddd [____] ddf
	//  |  |    

	if e < g.start {


	}
	if s < g.start && e < g.end {
		g.start = s
		g.cursor = s
	} else if {

	}
	// we're deleting backwards
	if s < g.start {
		// just expand the gap, its likely that the user will write back again anyways
		g.start = s
		g.cursor = s
	} else if e > g.start {
		g.end = e
		// g.start = s // ??
	}
}

// Moves the cursor and the gap to the position
func (g *GapBuffer) MoveCursorAt(pos int) {

	// cases:
	// are we moving to the same pos?
	if g.cursor == pos && g.start == g.cursor {
		return
	}

	g.cursor = pos
	g.start = pos
	// g.end =

}

// grows the gap and moves the cursor to the start
func (g *GapBuffer) grow() {
	// .
	newCap := len(g.data) + NEW_GAP_SIZE
	b := make([]rune, newCap)

	copy(b, g.data[:g.start]) // start of buffer tp gap start

	end := g.data[g.end:]
	copy(b[len(end):], end)

	g.cursor = g.start

	g.end = newCap
	g.data = b

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
