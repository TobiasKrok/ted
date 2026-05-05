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
		{"one char", "ABCDE", 0, 0, "BCDE"},
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

func TestDeleteRangeOutOfBoundsPanic(t *testing.T) {
	gb := newFilledBuffer("ABC")
	err1 := gb.DeleteRange(-1, 2)
	err2 := gb.DeleteRange(0, 100)
	if err1 == nil {
		t.Errorf("deleted from range -1 with no error returned")
	}
	if err2 == nil {
		t.Errorf("deleted to outside the buffer with no error returned")
	}
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

func TestMoveGap(t *testing.T) {
	cases := []struct {
		name      string
		initial   string
		moveTo    int
		wantStr   string
		wantStart int
	}{
		{"to beginning", "ABCDE", 0, "ABCDE", 0},
		{"to middle", "ABCDE", 2, "ABCDE", 2},
		{"to end (no-op)", "ABCDE", 5, "ABCDE", 5},
		{"single char to start", "A", 0, "A", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gb := newFilledBuffer(tc.initial)
			err := gb.MoveGap(tc.moveTo)

			eq(t, "err", err, nil)
			eq(t, "String", gb.String(), tc.wantStr)
			eq(t, "start", gb.start, tc.wantStart)
			eq(t, "cursor", gb.cursor, tc.wantStart)
			checkInvariant(t, gb, "after move")
		})
	}
}

func TestMoveGapLeftThenRight(t *testing.T) {
	gb := newFilledBuffer("ABCDE")
	// gap is at 5, move left to 1
	gb.MoveGap(1)
	eq(t, "after left", gb.String(), "ABCDE")
	eq(t, "start after left", gb.start, 1)
	checkInvariant(t, gb, "after left")

	// now move right to 4
	gb.MoveGap(4)
	eq(t, "after right", gb.String(), "ABCDE")
	eq(t, "start after right", gb.start, 4)
	checkInvariant(t, gb, "after right")
}

func TestMoveGapThenInsert(t *testing.T) {
	cases := []struct {
		name    string
		initial string
		moveTo  int
		insert  rune
		want    string
	}{
		{"insert at beginning", "BCD", 0, 'A', "ABCD"},
		{"insert in middle", "ACD", 1, 'B', "ABCD"},
		{"insert at end", "ABC", 3, 'D', "ABCD"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gb := newFilledBuffer(tc.initial)
			gb.MoveGap(tc.moveTo)
			gb.Insert(tc.insert)

			eq(t, "String", gb.String(), tc.want)
			eq(t, "Length", gb.Length(), len([]rune(tc.want)))
			checkInvariant(t, gb, "after move+insert")
		})
	}
}

func TestMoveGapPreservesGapSize(t *testing.T) {
	gb := newFilledBuffer("ABCDE")
	gapBefore := gb.gapLen()

	gb.MoveGap(2)
	eq(t, "gap after left", gb.gapLen(), gapBefore)

	gb.MoveGap(4)
	eq(t, "gap after right", gb.gapLen(), gapBefore)

	gb.MoveGap(0)
	eq(t, "gap after far left", gb.gapLen(), gapBefore)
}

func TestMoveGapRepeated(t *testing.T) {
	gb := newFilledBuffer("ABCDE")
	positions := []int{2, 0, 4, 1, 5, 3}
	for _, pos := range positions {
		gb.MoveGap(pos)
		eq(t, "String", gb.String(), "ABCDE")
		checkInvariant(t, gb, "repeated move")
	}
}

func TestMoveGapOutOfBounds(t *testing.T) {
	gb := newFilledBuffer("ABC")

	err := gb.MoveGap(-1)
	if err == nil {
		t.Error("expected error for negative position")
	}

	err = gb.MoveGap(100)
	if err == nil {
		t.Error("expected error for position beyond length")
	}
}

func TestDeleteRangeWithGapInMiddle(t *testing.T) {
	cases := []struct {
		name    string
		initial string
		gapAt   int
		s, e    int
		want    string
	}{
		{"after gap", "ABCDE", 1, 2, 4, "ABE"},
		{"across gap", "ABCDE", 2, 1, 4, "AE"},
		{"across gap all", "ABCDE", 2, 0, 5, ""},
		{"before gap", "ABCDE", 4, 0, 2, "CDE"},
		{"last char after gap", "ABCDE", 0, 4, 5, "ABCD"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gb := newFilledBuffer(tc.initial)
			err := gb.MoveGap(tc.gapAt)
			eq(t, "move err", err, nil)

			err = gb.DeleteRange(tc.s, tc.e)
			eq(t, "delete err", err, nil)
			eq(t, "String", gb.String(), tc.want)
			eq(t, "Length", gb.Length(), len([]rune(tc.want)))
			checkInvariant(t, gb, "after delete")
		})
	}
}

func TestCharAt(t *testing.T) {
	cases := []struct {
		name    string
		initial string
		pos     int
		want    rune
	}{
		{"first character", "ABCD", 0, 'A'},
		{"last character", "ABCD", 3, 'D'},
		{"middle character", "ABCDE", 2, 'C'},
		{"out of bounds left", "ABC", -2, -1},
		{"out of bounds right", "ABC", 10, -1},
		{"empty buffer", "", 0, -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gb := newFilledBuffer(tc.initial)
			ch := gb.CharAt(tc.pos)
			eq(t, "Character", ch, tc.want)
			checkInvariant(t, gb, "after char at")
		})
	}
}
