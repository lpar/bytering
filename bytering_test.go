package bytering

import (
	"bytes"
	"testing"
)

func runCapTest(t *testing.T, c int) {
	b := NewByteRing(c)
	for i := 0; i < 2*c; i++ {
		// Make sure we go through non-ASCII byte values so we will detect
		// any UTF-8 encoding issues
		b.WriteByte(byte(255 - i%256))
	}
	s := b.Bytes()
	if len(s) > c {
		t.Errorf("capacity was %d but string was %d long", c, len(s))
	}
}

func TestCapacity(t *testing.T) {
	runCapTest(t, 8)
	runCapTest(t, 17)
	runCapTest(t, 267)
}

func runFindTest(t *testing.T, sdata string, sfind string) {
	search := []byte(sfind)
	data := []byte(sdata)
	// Should find it after writing this many bytes
	cori := bytes.Index(data, search) + len(search)

	ring := NewByteRing(len(search))
	for i, b := range data {
		ring.WriteByte(b)
		if ring.Compare(search) {
			s := ring.String()
			if s != sfind {
				t.Errorf("failed search test, '%s' != '%s'", s, sfind)
			}
			// Bytes written = i+1
			if cori != i+1 {
				t.Errorf("failed search test, %d != %d", i+1, cori)
			}
			return
		}
	}
	t.Errorf("failed search test - didn't find it")
}

func TestCompare(t *testing.T) {
	runFindTest(t, "foofoo", "foo")
	runFindTest(t, "eedle eedle needl haysneedletack", "needle")
	runFindTest(t, " test 🍀 stringLucky", "Lucky")
}

func TestNotFound(t *testing.T) {
	ring := NewByteRing(200)
	for i := 0; i < 255; i++ {
		ring.WriteByte(byte(i))
	}
	if ring.Compare([]byte("bogus")) {
		t.Errorf("failed on search for nonexistent bytes - found them")
	}
}

func TestStringing(t *testing.T) {
	ring := NewByteRing(20)
	ring.WriteByte(65)
	ring.WriteByte(66)
	ring.WriteByte(67)
	s := ring.String()
	if len(s) != 3 {
		t.Errorf("failed to return correct length of string, expected 3 got %d", len(s))
	}
}

func TestStringWrap(t *testing.T) {
	ring := NewByteRing(4)
	for i := 65; i <= 69; i++ {
		ring.WriteByte(byte(i))
	}
	s := ring.String()
	if s != "BCDE" {
		t.Errorf("failed to return correct wraparound string, got '%s'", s)
	}
}

func TestTooSmall(t *testing.T) {
	ring := NewByteRing(8)
	data := []byte("ABCDEFGHI")
	find := []byte("ABCDEFGH")
	for _, b := range data {
		ring.WriteByte(b)
	}
	if ring.Compare(find) {
		t.Errorf("failed on search for bytes which won't fit in buffer - found them")
	}
}
