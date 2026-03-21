package buffer

import (
	"strings"
	"testing"
)

func TestNewGapBuffer(t *testing.T) {
	ic := 10
	gb := NewGapBuffer(ic)

	if len(gb.data) != ic {

		t.Errorf("Expected len to be %d, got %d", ic, len(gb.data))
	}

	if gb.String() != "" {
		t.Errorf("Expected empty buffer, got %q", gb.String())
	}

	if gb.start != 0 {
		t.Errorf("Expected gap start at 0, got %d", gb.start)
	}

	if gb.cursor != gb.start {
		t.Errorf("Expected cursor to be at gap start, but is at index %d", gb.cursor)
	}
}

func TestInsert(t *testing.T) {

	ic := 10
	gb := NewGapBuffer(ic)

	gb.Insert('A')
	gb.Insert('B')
	gb.Insert('C')

	if gb.Length() != 3 {

		t.Errorf("Expected content length 3, got %q", gb.String())
	}
	if gb.String() != "ABC" {
		t.Errorf("Expected 'ABC', got %q", gb.String())
	}
	// ABC[_____]
	if gb.start != 3 {
		t.Errorf("Expected gap start at idx 3, got %d", gb.start)
	}
	if gb.end != ic {
		t.Errorf("Expected gap end at %d, got %d", ic, gb.end)
	}

	if gb.cursor != gb.start {
		t.Errorf("Expected cursor to be at gap start at %d, got %d", gb.start, gb.cursor)
	}
}

// simple, the gap end should be at the last index
func TestManualGrowSimple(t *testing.T) {

	ic := 5
	gb := NewGapBuffer(ic)
	gb.Insert('A')
	gb.Insert('B')
	gb.Insert('C')

	gb.grow()

	nc := ic + GAP_SIZE

	if len(gb.data) != nc {
		t.Errorf("Expected buffer size to be %d, got %d", nc, len(gb.data))
	}

	if gb.start != 3 {
		t.Errorf("Expected gap start at idx 3, got %d", gb.start)
	}
	if gb.end != nc {

		t.Errorf("Expected gap end at %d, got %d", nc, gb.end)
	}
	// to make sure the gap pointers are updated
	if gb.String() != "ABC" {
		t.Errorf("Expected 'ABC', got %q", gb.String())
	}
}

// its probably a good tedt case? idk
func TestManualGrowEmpty(t *testing.T) {

	ic := 5
	gb := NewGapBuffer(ic)
	gb.grow()
	nc := ic + GAP_SIZE
	if len(gb.data) != nc {
		t.Errorf("Expected buffer size to be %d, got %d", nc, len(gb.data))
	}

	if gb.start != 0 {
		t.Errorf("Expected gap start at idx 0, got %d", gb.start)
	}
	if gb.end != nc {
		t.Errorf("Expected gap end at %d, got %d", nc, gb.end)
	}
}

// test if the buffer properly grows dynamically during insert
func TestGrowDynamic(t *testing.T) {
	ic := 5
	gb := NewGapBuffer(ic)
	c := 10

	var sb strings.Builder
	for range c {
		sb.WriteString("A")
		gb.Insert('A')
	}

	nc := ic + GAP_SIZE
	if len(gb.data) != nc {
		t.Errorf("Expected buffer size to be %d, got %d", nc, len(gb.data))
	}
	if gb.start != c {
		t.Errorf("Expected gap start at idx %d, got %d", c, gb.start)
	}

	if gb.cursor != gb.start {
		t.Errorf("Expected cursor at idx %d, got %d", gb.cursor, gb.start)
	}

	if gb.String() != sb.String() {
		t.Errorf("Expected '%q', got %q", sb.String(), gb.String())
	}
}
