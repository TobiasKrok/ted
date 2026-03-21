package buffer

type GapBuffer struct {
	data  []rune
	start int
	end   int
	// length int // ?
	cursor int

	chbef int // chars before cursor
}

const GAP_SIZE = 64
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
	// we can simplify this by just checking if its not at cursor and move if it needs to
	g.data[g.cursor] = r
	g.cursor += 1
	g.start += 1

}

// func (g *GapBuffer) moveGap() {
//
// }

// grows the gap and moves the cursor to the start
func (g *GapBuffer) grow() {
	// .
	newCap := len(g.data) + GAP_SIZE
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
