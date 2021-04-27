package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/yarcat/playground/redis"
)

func main() {
	addr := flag.String("addr", ":6379", "redis endpoint")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer conn.Close()

	c := redis.NewClient(logging{conn})
	run(c, make([]byte, 100) /*buf*/, true /*withResultLogging*/)
}

func run(c redis.Client, buf []byte, withResultLogging bool) {
	resp, err := c.SetString("mykey", "myvalue\x00")
	if err == nil {
		err = resp.ErrorOrConsume()
	}
	if err != nil {
		log.Fatalf("set failed: %v", err)
	}

	resp, err = c.Exists("mykey", "kk")
	if err == nil {
		err = resp.Error()
	}
	if err != nil {
		log.Fatalf("exist failed: %v", err)
	}
	var n int
	if n, err = resp.Int(); err != nil {
		log.Fatalf("int failed: %v", err)
	}
	if withResultLogging {
		log.Printf("exists = %v", n)
	}

	resp, err = c.Get("mykey")
	if err == nil {
		err = resp.Error()
	}
	if err != nil {
		log.Fatalf("get failed: %v", err)
	}
	if n, err = resp.BytesBulk(buf); err != nil {
		log.Fatalf("get failed: %v", err)
	}
	if withResultLogging {
		log.Printf("get = %q", buf[:n])
	}
}

type logging struct{ io.ReadWriter }

func (l logging) Write(b []byte) (int, error) {
	log.Printf("WRITE: %q", b)
	return l.ReadWriter.Write(b)
}
