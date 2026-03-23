package buffer

import (
	"strings"
	"testing"
)

func newFilledBuffer(s string) *GapBuffer {
	gb := NewGapBuffer(len(s) + NEW_GAP_SIZE)
	for _, r := range s {
		gb.Insert(r)
	}
	return gb
}

func eq[T comparable](t *testing.T, label string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", label, got, want)
	}
}

func checkInvariant(t *testing.T, gb *GapBuffer, label string) {
	t.Helper()
	if total := gb.Length() + gb.gapLen(); total != len(gb.data) {
		t.Errorf("invariant broken at %s: Length(%d) + gapLen(%d) != len(data)(%d)",
			label, gb.Length(), gb.gapLen(), len(gb.data))
	}
}

func TestNewGapBuffer(t *testing.T) {
	ic := 10
	gb := NewGapBuffer(ic)

	eq(t, "String", gb.String(), "")
	eq(t, "Length", gb.Length(), 0)
	eq(t, "start", gb.start, 0)
	eq(t, "end", gb.end, ic)
	eq(t, "cursor", gb.cursor, gb.start)
	eq(t, "gapLen", gb.gapLen(), ic)
	eq(t, "data len", len(gb.data), ic)
}

func TestInsert(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"A"},
		{"ABC"},
		{"héllo wörld 🌍"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			gb := newFilledBuffer(tc.input)
			eq(t, "string", gb.String(), tc.input)
			eq(t, "oength", gb.Length(), len([]rune(tc.input)))
			eq(t, "cursor==start", gb.cursor, gb.start)
		})
	}
}

func TestInsertGapPointers(t *testing.T) {
	ic := 10
	gb := NewGapBuffer(ic)
	gb.Insert('A')
	gb.Insert('B')
	gb.Insert('C')

	eq(t, "start", gb.start, 3)
	eq(t, "end", gb.end, ic)
	eq(t, "cursor", gb.cursor, gb.start)
}

func TestGrow(t *testing.T) {
	cases := []struct {
		name    string
		initial string
	}{
		{"empty", ""},
		{"partial", "ABC"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ic := 5
			gb := NewGapBuffer(ic)
			for _, r := range tc.initial {
				gb.Insert(r)
			}
			gb.grow()

			nc := ic + NEW_GAP_SIZE
			eq(t, "data len", len(gb.data), nc)
			eq(t, "end", gb.end, nc)
			eq(t, "string", gb.String(), tc.initial)
		})
	}
}

func TestGrowDynamic(t *testing.T) {
	cases := []struct {
		name  string
		count int
	}{
		{"single grow", 10},
		{"double grow", NEW_GAP_SIZE + 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gb := NewGapBuffer(5)
			want := strings.Repeat("A", tc.count)
			for range tc.count {
				gb.Insert('A')
			}

			eq(t, "string", gb.String(), want)
			eq(t, "length", gb.Length(), tc.count)
			eq(t, "cursor==start", gb.cursor, gb.start)
			checkInvariant(t, gb, "after inserts")
		})
	}
}

func TestDeleteBackwards(t *testing.T) {
	cases := []struct {
		name    string
		initial string
		deletes int
		want    string
	}{
		{"delete one", "ABC", 1, "AB"},
		{"delete two", "ABC", 2, "A"},
		{"delete all", "AB", 2, ""},
		{"delete on empty", "", 1, ""},
		{"over-delete", "A", 5, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gb := newFilledBuffer(tc.initial)
			for range tc.deletes {
				gb.DeleteBackwards()
			}

			eq(t, "String", gb.String(), tc.want)
			eq(t, "Length", gb.Length(), len([]rune(tc.want)))
			checkInvariant(t, gb, "after deletes")
		})
	}
}

func TestDeleteBackwardsNeverUnderflows(t *testing.T) {
	gb := NewGapBuffer(10)
	gb.DeleteBackwards()

	eq(t, "start", gb.start, 0)
	eq(t, "cursor", gb.cursor, 0)
}

// func TestMoveCursorAt(t *testing.T) {
// 	cases := []struct {
// 		name    string
// 		initial string
// 		pos     int
// 		insert  rune
// 		want    string
// 	}{
// 		{"move to beginning then insert", "ACD", 0, 'X', "XACD"},
// 		{"move to middle then insert", "ACD", 1, 'B', "ABCD"},
// 		{"move to end then insert", "AB", 2, 'C', "ABC"},
// 	}
// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			gb := newFilledBuffer(tc.initial)
// 			gb.MoveCursorAt(tc.pos)
//
// 			eq(t, "String after move", gb.String(), tc.initial)
// 			eq(t, "cursor", gb.cursor, tc.pos)
//
// 			gb.Insert(tc.insert)
// 			eq(t, "String after insert", gb.String(), tc.want)
// 		})
// 	}
// }

func TestDeleteRange(t *testing.T) {
	cases := []struct {
		name    string
		initial string
		s, e    int
		want    string
	}{
		{"delete from start", "ABCDE", 0, 2, "CDE"},
		{"delete middle", "ABCDE", 1, 3, "ADE"},
		{"delete to end", "ABCDE", 3, 5, "ABC"},
		{"delete all", "ABCDE", 0, 5, ""},
		{"reversed args", "ABCDE", 3, 1, "ADE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gb := newFilledBuffer(tc.initial)
			gb.DeleteRange(tc.s, tc.e)

			eq(t, "String", gb.String(), tc.want)
			eq(t, "Length", gb.Length(), len([]rune(tc.want)))
		})
	}
}

func TestDeleteRangeOutOfBoundsNoPanic(t *testing.T) {
	gb := newFilledBuffer("ABC")
	gb.DeleteRange(-1, 2)
	gb.DeleteRange(0, 100)
}

func TestCapacityInvariantThroughoutMutations(t *testing.T) {
	gb := NewGapBuffer(10)
	checkInvariant(t, gb, "empty")

	for i, r := range "Hello, 世界!" {
		gb.Insert(r)
		checkInvariant(t, gb, "insert "+string(r)+" ("+string(rune('0'+i))+")")
	}
	for range 5 {
		gb.DeleteBackwards()
		checkInvariant(t, gb, "delete")
	}
}

func TestStringRoundTrip(t *testing.T) {
	cases := []string{
		"",
		"a",
		"hello, world",
		"日本語テスト",
		strings.Repeat("x", NEW_GAP_SIZE*3),
	}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			gb := newFilledBuffer(s)
			eq(t, "String", gb.String(), s)
			eq(t, "Length", gb.Length(), len([]rune(s)))
		})
	}
}
