package main

import (
	"log"
	"net"
	"testing"

	"github.com/yarcat/playground/redis"
	"github.com/yarcat/playground/redis/protocol"
)

func init() {
	// We care about memory allocations here, so testing this against real Redis
	// sucks, but is acceptable.
	conn, err := net.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	client = redis.NewClient(conn)
	stream = redis.NewStream(conn)
	proto = protocol.New(redis.NewConn(conn))
}

const noResultLogging = false

var (
	client redis.Client
	stream *redis.Stream
	buf    = make([]byte, 100)
	proto  *protocol.Protocol
)

func BenchmarkRun(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runV1(client, buf, noResultLogging)
	}
}

func BenchmarkRun2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runV2(stream, buf, noResultLogging)
	}
}

func BenchmarkRun3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runV3(proto, buf, noResultLogging)
	}
}
