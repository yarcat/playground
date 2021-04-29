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

	c := redis.NewClient(logging{conn})
	run(c, buf /*buf*/, logResults)

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
	if err := wrap(c.SetString("mykey", "myvalue\x00")).Ignore(); err != nil {
		log.Fatalln("- set:", err)
	} else if withResultLogging {
		log.Println("+ set: OK")
	}

	if n, err := wrap(c.Exists("mykey", "kk", "mykey2")).Int(); err != nil {
		log.Fatalln("- exists:", err)
	} else if withResultLogging {
		log.Println("+ exists:", n)
	}

	for _, k := range []string{"mykey", "nosuchkey"} {
		if n, ok, err := wrap(c.Get(k)).Bulk(buf); err != nil {
			log.Fatalln("- get:", err)
		} else if withResultLogging {
			log.Printf("+ get: %q (found=%v)\n", buf[:n], ok)
		}
	}
}

type logging struct{ io.ReadWriter }

func (l logging) Write(b []byte) (int, error) {
	log.Printf("WRITE: %q", b)
	return l.ReadWriter.Write(b)
}

type wrappedResponse struct {
	resp redis.Response
	err  error
}

func wrap(resp redis.Response, err error) wrappedResponse {
	return wrappedResponse{resp: resp, err: err}
}

func (wr wrappedResponse) Ignore() error {
	if wr.err != nil {
		return wr.err
	}
	return wr.resp.ErrorOrConsume()
}

func (wr wrappedResponse) Int() (n int, err error) {
	if wr.err != nil {
		return 0, wr.err
	}
	return wr.resp.Int()
}

func (wr wrappedResponse) Bulk(buf []byte) (int, bool, error) {
	if wr.err != nil {
		return 0, false, wr.err
	}
	switch n, err := wr.resp.BytesBulk(buf); {
	case err != nil:
		return 0, false, err
	case n < 0:
		return 0, false, nil
	default:
		return n, true, nil
	}
}
