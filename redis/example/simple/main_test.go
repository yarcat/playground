package main

import (
	"log"
	"net"
	"testing"

	"github.com/yarcat/playground/redis"
)

func init() {
	// We care about memory allocations here, so testing this against real Redis
	// sucks, but is acceptable.
	conn, err := net.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	client = redis.NewClient(conn)
}

var (
	client redis.Client
	buf    = make([]byte, 100)
)

func BenchmarkRun(b *testing.B) {
	for n := 0; n < b.N; n++ {
		run(client, buf, false /*withResultLogging*/)
	}
}
