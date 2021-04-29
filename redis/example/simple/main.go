package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"runtime/trace"
	"sync"

	"net/http"
	_ "net/http/pprof"

	"github.com/yarcat/playground/redis"
)

func main() {
	addr := flag.String("addr", ":6379", "redis `endpoint`")

	iterations := flag.Int("iterations", 1, "how many `times` to loop")
	extraLogging := flag.Bool("extra_logging", true, "whether results and writes should be logged")

	pprof := flag.String("pprof", "", "serve `addr` for pprof")
	traceFile := flag.String("trace_profile", "", "trace output file")

	flag.Parse()

	if *pprof != "" {
		go func() { log.Fatal(http.ListenAndServe(*pprof, nil)) }()
	}
	if *traceFile != "" {
		f := mustCreateFile(*traceFile)
		defer f.Close()
		stop := mustTraceStart(f)
		defer stop()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	conn1 := mustConnect(*addr, *extraLogging)
	defer conn1.Close()

	go func() {
		defer wg.Done()
		client := redis.NewClient(conn1)
		run(*iterations, "V1", func(buf []byte) {
			runV1(client, buf, *extraLogging)
		})
	}()

	conn2 := mustConnect(*addr, *extraLogging)
	defer conn2.Close()

	go func() {
		defer wg.Done()
		stream := redis.NewStream(conn2)
		run(*iterations, "V2", func(buf []byte) {
			runV2(stream, buf, *extraLogging)
		})
	}()

	wg.Wait()
}

func mustCreateFile(name string) *os.File {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func mustTraceStart(w io.Writer) (stop func()) {
	if err := trace.Start(w); err != nil {
		log.Fatal(err)
	}
	return trace.Stop
}

func run(times int, version string, f func([]byte)) {
	buf := make([]byte, 100)

	for i := 1; i <= times; i++ {
		if i%500 == 1 {
			log.Printf("%s: %d of %d (%.2f%%)",
				version, i, times, 100.0*float64(i)/float64(times))
		}
		f(buf)
	}
	log.Printf("%s: done", version)
}

func mustConnect(endpoint string, extraLogging bool) io.ReadWriteCloser {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		panic(err)
	}
	if extraLogging {
		return logging{conn}
	}
	return conn
}

func runV2(stream *redis.Stream, b []byte, withResultLogging bool) {
	switch status, err := redis.Execute(stream, "SET",
		redis.String("mykey"), redis.String("my\x00value")).Void(); {
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

	switch status, err := redis.Execute(stream, "MGET",
		redis.String("mykey"), redis.String("nosuchkey"), redis.String("mykey"),
	).Array(&n); {
	case err != nil:
		log.Fatalf("mget: %v", err)
	case !status.OK():
		log.Printf("-mget: %v", status.Bytes())
	default:
		for i := 0; i < n; i++ {
			var (
				n  int
				ok bool
			)
			switch status, err := redis.ArrayItem(stream).Bytes(b, &n, &ok); {
			case err != nil:
				log.Fatalf("mget[%d]: %v", i, err)
			case !status.OK():
				log.Printf("-mget[%d]: %v", i, status.Bytes())
			case withResultLogging:
				log.Printf("+mget[%d]: %q (found=%v)", i, b[:n], ok)
			}
		}
	}
}

func runV1(c redis.Client, buf []byte, withResultLogging bool) {
	if err := wrap(c.SetString("mykey", "my\x00value")).Ignore(); err != nil {
		log.Fatalln("-set:", err)
	} else if withResultLogging {
		log.Println("+set: OK")
	}

	if n, err := wrap(c.Exists("mykey", "kk", "mykey2")).Int(); err != nil {
		log.Fatalln("-exists:", err)
	} else if withResultLogging {
		log.Println("+exists:", n)
	}

	for _, k := range []string{"mykey", "nosuchkey"} {
		if n, ok, err := wrap(c.Get(k)).Bulk(buf); err != nil {
			log.Fatalln("-get:", err)
		} else if withResultLogging {
			log.Printf("+get: %q (found=%v)\n", buf[:n], ok)
		}
	}
}

type logging struct{ io.ReadWriteCloser }

func (l logging) Write(b []byte) (int, error) {
	log.Printf("WRITE: %q", b)
	return l.ReadWriteCloser.Write(b)
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
	return wr.resp.Void()
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
