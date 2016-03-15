// Package bytering implements a high speed circular byte buffer for
// scanning files and streams for a given byte sequence.
//
// Internally, a single static byte buffer is used to hold the data.
// WriteByte and Compare require no memory allocations.
//
// General usage:
//
//  1. Create a ByteRing with enough capacity to hold the slice of data you
//     are looking for.
//  2. WriteByte bytes (from a buffered IO source!) until Compare finds
//     the data you are looking for.
//  3. Rewind your source by the length of the search data, if appropriate.
//
package bytering

import "bytes"

type ByteRing struct {
	data  []byte
	start int
	end   int
}

// NewByteRing creates and returns a new ByteRing with the specified capacity.
func NewByteRing(capacity int) *ByteRing {
	var b ByteRing
	b.data = make([]byte, capacity, capacity)
	return &b
}

// WriteByte writes a byte to the ByteRing
func (b *ByteRing) WriteByte(x byte) {
	b.data[b.end] = x
	if b.start == b.end {
		b.start = (b.start + 1) % cap(b.data)
	}
	b.end = (b.end + 1) % cap(b.data)
}

// Compare checks to see if the ByteRing's current contents match the specified
// byte slice. Comparison begins at the current start of the ByteRing, which
// will be the oldest byte written to the ByteRing and still in its buffer.
//
// If the ByteRing's capacity is smaller than the length of the byte slice
// you are searching for, you will never find it.
func (b *ByteRing) Compare(find []byte) bool {
	j := b.start
	for _, c := range find {
		if b.data[j] != c {
			return false
		}
		j = (j + 1) % cap(b.data)
	}
	return true
}

// Return the current contents of the ByteRing as an array of bytes.
// Uses bytes.Buffer to make a copy and hence requires a memory allocation,
// so don't use this in a loop.
func (b *ByteRing) Bytes() []byte {
	var s bytes.Buffer
	if b.end > b.start {
		s.Write(b.data[b.start:b.end])
	} else {
		s.Write(b.data[b.start:])
		s.Write(b.data[:b.end])
	}
	return s.Bytes()
}

// Return the current contents of the ByteRing as a string. See Bytes() for
// caveats.
func (b *ByteRing) String() string {
	bytes := b.Bytes()
	return string(bytes)
}
