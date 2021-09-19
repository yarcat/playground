// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en_GB": &dictionary{index: en_GBIndex, data: en_GBData},
		"ru_RU": &dictionary{index: ru_RUIndex, data: ru_RUData},
	}
	fallback := language.MustParse("en-GB")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"Resist. max: %.3f":     2,
	"Resistivity":           0,
	"Resistivity min: %.3f": 1,
}

var en_GBIndex = []uint32{ // 4 elements
	0x00000000, 0x0000000c, 0x00000025, 0x0000003a,
} // Size: 40 bytes

const en_GBData string = "" + // Size: 58 bytes
	"\x02Resistivity\x02Resistivity min: %.3[1]f\x02Resist. max: %.3[1]f"

var ru_RUIndex = []uint32{ // 4 elements
	0x00000000, 0x0000001b, 0x00000036, 0x00000053,
} // Size: 40 bytes

const ru_RUData string = "" + // Size: 83 bytes
	"\x02Сопротивление\x02Мин. сопр.: %.3[1]f\x02Макс. сопр.: %.3[1]f"

	// Total table size 221 bytes (0KiB); checksum: C12DD4E9
