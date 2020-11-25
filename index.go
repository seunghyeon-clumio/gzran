package gzseek

import (
	"sort"
)

// Index collects decompressor state at offset Points.
// gzseek.Reader adds points to the index on the fly as decompression proceeds.
type Index []Point

func (idx Index) lastUncompressedOffset() int64 {
	if len(idx) == 0 {
		return 0
	}

	return idx[len(idx)-1].UncompressedOffset
}

func (idx Index) closestPointBefore(offset int64) Point {
	j := sort.Search(len(idx), func(j int) bool {
		return idx[j].UncompressedOffset <= offset
	})

	if j == len(idx) {
		return Point{}
	}

	return idx[j]
}

// Point holds the decompressor state at a given offset within the uncompressed data.
type Point struct {
	CompressedOffset   int64
	UncompressedOffset int64
	DecompressorState  []byte
}
