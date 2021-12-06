package main

import (
	"encoding/binary"
	"io"
	"math"
	"time"
)

type Type uint16

const (
	TypeUndefined Type = iota
	Type1
	Type2
	Type3
)

type TimeValue struct {
	Type Type
	T    time.Time
	V    float32
}

func EncodeTypeValue(w io.Writer, order binary.ByteOrder, tv TimeValue) error {
	var err error
	if err == nil {
		err = binary.Write(w, order, uint16(tv.Type))
	}
	if err == nil {
		err = binary.Write(w, order, tv.T.UTC().UnixMilli())
	}
	if err == nil {
		err = binary.Write(w, order, math.Float32bits(tv.V))
	}
	return err
}
