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

	for n := 0; n < 100; n++ {
		startedAt := time.Now()
		v.Put(func(w io.Writer) error {
			for i := 0; i < 30; i++ {
				err := EncodeTypeValue(w, binary.BigEndian, TimeValue{
					Type: Type1,
					T:    time.Now(),
					V:    rand.Float32(),
				})
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		})
		fmt.Printf("Finished %d-th batch in %v\n", n, time.Since(startedAt))
		v.Close()
		// time.Sleep(1 * time.Second)
	}

	started := time.Now()
	shard := vase.Shards(v)
	defer shard.Close()
	for shard.Next() {
		cnt := 0
		err := shard.Read(func(r io.Reader) error {
			for {
				var tv TimeValue
				if err := DecodeTypeValue(r, binary.BigEndian, &tv); err != nil {
					return err
				}
				cnt++
			}
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Shard %v contains %d records\n", shard.Time(), cnt)
	}
	fmt.Println("read all in", time.Since(started))
}
