package buffer

type GapBuffer struct {
	data   []rune
	start  int
	end    int
	length int // ?
	cursor int
}

const GAP_SIZE = 5

func New() *GapBuffer {

	d := make([]rune, GAP_SIZE)
	return &GapBuffer{
		data:   d,
		start:  0,
		end:    GAP_SIZE,
		cursor: 0,
		length: 0,
	}
}

func (g *GapBuffer) Insert(r rune) {
	// we can simplify this by just checking if its not at cursor and move if it needs to
	g.data[g.cursor] = r
	g.cursor += 1
	g.start += 1
}

func (g *GapBuffer) String() string {
	return string(g.data)
}

func (g *GapBuffer) Length() int {
	return len(g.data) - g.gapLen()
}

func (g *GapBuffer) gapLen() int {

	return g.end - g.start
}
