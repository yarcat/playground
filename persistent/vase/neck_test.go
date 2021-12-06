package vase

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestEncodeNeck(t *testing.T) {
	var (
		n Neck
		b bytes.Buffer
	)
	if err := EncodeNeck(&b, binary.BigEndian, n); err != nil {
		t.Errorf("EncodeNeck(...) = %v, want = nil", err)
	}
	if b.Len() != NeckSize {
		t.Errorf("EncodeNeck(...) wrote %d bytes, want = %d", b.Len(), NeckSize)
	}
}
