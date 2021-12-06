package main

import (
	"encoding/binary"
	"io"
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

type timeValueRepr struct {
	Type uint16
	MSec int64
	Val  float32
}

func EncodeTypeValue(w io.Writer, order binary.ByteOrder, tv TimeValue) error {
	return binary.Write(w, order, timeValueRepr{
		uint16(tv.Type), tv.T.UTC().UnixMilli(), tv.V,
	})
}

func DecodeTypeValue(r io.Reader, order binary.ByteOrder, tv *TimeValue) error {
	var tvr timeValueRepr
	if err := binary.Read(r, order, &tvr); err != nil {
		return err
	}
	*tv = TimeValue{Type(tvr.Type), time.UnixMilli(tvr.MSec), tv.V}
	return nil
}
