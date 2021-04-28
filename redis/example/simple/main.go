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

	const logResults = true
	buf := make([]byte, 100)

	// c := redis.NewClient(logging{conn})
	// run(c, buf /*buf*/, logResults)

	stream := redis.NewStream(logging{conn})
	run2(stream, buf, logResults)
}

func run2(stream *redis.Stream, b []byte, withResultLogging bool) {
	switch status, err := redis.Execute(stream, "SET",
		redis.String("mykey"), redis.String("myvalue")).Void(); {
	case err != nil:
		log.Fatalf("set: %v", err)
	case !status.OK():
		log.Printf("-set: %s", status.Bytes())
	case withResultLogging:
		log.Printf("+set: %s", status.Bytes())
	}

	var n int
	switch status, err := redis.Execute(stream, "EXISTS",
		redis.String("mykey"), redis.String("kk"), redis.String("does not exist")).Int(&n); {
	case err != nil:
		log.Fatalf("exists: %v", err)
	case !status.OK():
		log.Printf("-exists: %s", status.Bytes())
	case withResultLogging:
		log.Printf("+exists: %d (expected=%v)", n, n == 2)
	}

	for _, k := range []string{"mykey", "nosuchkey"} {
		var ok bool
		switch status, err := redis.Execute(stream, "GET", redis.String(k)).Bytes(b[:], &n, &ok); {
		case err != nil:
			log.Fatalf("get: %v", err)
		case !status.OK():
			log.Printf("-get: %v", status.Bytes())
		case withResultLogging:
			log.Printf("+get: %q (found=%v)", b[:n], ok)
		}
	}
}

func run(c redis.Client, buf []byte, withResultLogging bool) {
	resp, err := c.SetString("mykey", "myvalue\x00")
	if err == nil {
		err = resp.ErrorOrConsume()
	}
	if err != nil {
		log.Fatalf("set failed: %v", err)
	}

	resp, err = c.Exists("mykey", "kk", "mykey2")
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

	resp, err = c.Get("nosuchkey")
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
		if n >= 0 {
			log.Printf("get = %q", buf[:n])
		} else {
			log.Println("get = nil")
		}
	}
}

type logging struct{ io.ReadWriter }

func (l logging) Write(b []byte) (int, error) {
	log.Printf("WRITE: %q", b)
	return l.ReadWriter.Write(b)
}
