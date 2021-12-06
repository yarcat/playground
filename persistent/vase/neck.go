package vase

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

type Neck struct {
	CRC32     uint32
	Size      uint32
	Version   uint32
	CreatedAt time.Time
	LastAt    time.Time
}

func neckSize() int {
	var n Neck
	return binary.Size(n.CRC32) +
		binary.Size(n.Size) +
		binary.Size(n.Version) +
		binary.Size(n.CreatedAt.UTC().UnixMilli()) +
		binary.Size(n.LastAt.UTC().UnixMilli())
}

var NeckSize = neckSize()

func EncodeNeck(w io.Writer, order binary.ByteOrder, n Neck) error {
	var b bytes.Buffer
	b.Grow(NeckSize)

	n.Size = uint32(NeckSize)
	var err error
	if err == nil {
		err = binary.Write(&b, order, n.CRC32)
	}
	if err == nil {
		err = binary.Write(&b, order, n.Size)
	}
	if err == nil {
		err = binary.Write(&b, order, n.Version)
	}
	if err == nil {
		err = binary.Write(&b, order, n.CreatedAt.UTC().UnixMilli())
	}
	if err == nil {
		err = binary.Write(&b, order, n.LastAt.UTC().UnixMilli())
	}

	if err == nil {
		_, err = io.CopyN(w, &b, int64(NeckSize))
	}
	return err
}

func DecodeNeck(r io.Reader, order binary.ByteOrder, n *Neck) error {
	b := make([]byte, NeckSize)
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}

	r = bytes.NewReader(b) // Now read from the buffer.

	var err error
	if err == nil {
		err = binary.Read(r, order, &n.CRC32)
	}
	if err == nil {
		err = binary.Read(r, order, &n.Size)
	}
	if err == nil {
		err = binary.Read(r, order, &n.Version)
	}
	if err == nil {
		var ms int64
		err = binary.Read(r, order, &ms)
		n.CreatedAt = time.UnixMilli(ms)
	}
	if err == nil {
		var ms int64
		err = binary.Read(r, order, &ms)
		n.LastAt = time.UnixMilli(ms)
	}
	return err
}
