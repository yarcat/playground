package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/yarcat/go-persistent/vase"
)

func main() {
	rand.Seed(time.Now().UnixMicro())

	v, err := vase.New("timeseries")
	if err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	stopHammering := vase.Break(v).Every(2 * time.Second)
	defer stopHammering()

	for n := 0; n < 10; n++ {
		startedAt := time.Now()
		for i := 0; i < 30; i++ {
			err = v.Put(func(w io.Writer) error {
				return EncodeTypeValue(w, binary.BigEndian, TimeValue{
					Type: Type1,
					T:    time.Now(),
					V:    rand.Float32(),
				})
			})
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("Finished %d-th batch in %v\n", n, time.Since(startedAt))
		time.Sleep(1 * time.Second)
	}
}
