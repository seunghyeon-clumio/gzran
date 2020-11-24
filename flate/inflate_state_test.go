package flate

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

const testDat = "../testdata/Isaac.Newton-Opticks.txt"

type testState struct {
	remaining []byte
	state     []byte
}

func TestRestoreState(t *testing.T) {
	// Read and compress the data into a bytes.Buffer.
	f, err := os.Open(testDat)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	var b bytes.Buffer
	w, err := NewWriter(&b, DefaultCompression)
	if err != nil {
		t.Error(err)
	}
	if _, err := io.Copy(w, f); err != nil {
		t.Error(err)
	}
	if err := w.Close(); err != nil {
		t.Error(err)
	}

	// Decompress the data, saving decompressor state at random points.
	r := NewReader(&b)
	nRead := 0
	nextState := rand.Intn(128 * 1024)
	var states []testState
	for {
		sz := rand.Intn(4096)
		buf := make([]byte, sz)
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
		}

		nRead += n
		if nRead >= nextState {
			remaining := make([]byte, b.Len())
			copy(remaining, b.Bytes())
			state, err := DecompressorState(r)
			if err != nil {
				t.Error(err)
			}
			states = append(states, testState{
				remaining: remaining,
				state:     state,
			})
			nextState += rand.Intn(128 * 1024)
		}
	}
	t.Logf("Saved %d decompressor states", len(states))

	// Verify that we can resume decompression from each of the saved states.
	for _, s := range states {
		t.Logf("Resuming decompression with %d bytes remaining", len(s.remaining))
		compressedR := bytes.NewReader(s.remaining)
		r, err := NewReaderState(compressedR, s.state)
		if err != nil {
			t.Error(err)
		}
		if _, err := io.Copy(ioutil.Discard, r); err != nil {
			t.Error(err)
		}
	}
}
