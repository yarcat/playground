package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"runtime/trace"
	"sort"
	"strings"
	"sync"

	"net/http"
	_ "net/http/pprof"

	"github.com/yarcat/playground/redis"
	"github.com/yarcat/playground/redis/protocol"
)

func main() {
	runners := versionRunners{
		"v1": (runner).RunV1,
		"v2": (runner).RunV2,
		"v3": (runner).RunV3,
	}

	redisAddr := flag.String("redis_addr", ":6379", "redis backend `host:port`")

	iterations := flag.Int("iterations", 1, "how many `times` to loop")
	extraLogging := flag.Bool("extra_logging", true, "whether results and writes should be logged")
	versions := flag.String("versions", strings.Join(runners.All(), ","),
		"comma-separated string `list` of implementation versions to run")

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

	r := runner{
		redisAddr:    *redisAddr,
		extraLogging: *extraLogging,
		iterations:   *iterations,
		runners:      runners,
	}
	var wg sync.WaitGroup
	for _, v := range strings.Split(*versions, ",") {
		wg.Add(1)
		v := v
		go func() {
			defer wg.Done()
			r.Run(v)
		}()
	}
	wg.Wait()
}

type versionRunners map[string]func(runner, io.ReadWriter)

func (vr versionRunners) All() (versions []string) {
	for v := range vr {
		versions = append(versions, v)
	}
	sort.Strings(versions)
	return versions
}

type runner struct {
	redisAddr    string
	extraLogging bool
	iterations   int
	runners      versionRunners
}

func (r runner) Run(version string) {
	version = strings.ToLower(strings.TrimSpace(version))
	f, ok := r.runners[version]
	if !ok {
		return
	}
	conn := mustConnect(r.redisAddr, r.extraLogging)
	defer conn.Close()
	f(r, conn)
}

func (r runner) RunV1(rw io.ReadWriter) {
	client := redis.NewClient(rw)
	run(r.iterations, "v1", func(buf []byte) {
		runV1(client, buf, r.extraLogging)
	})
}

func (r runner) RunV2(rw io.ReadWriter) {
	stream := redis.NewStream(rw)
	run(r.iterations, "v2", func(buf []byte) {
		runV2(stream, buf, r.extraLogging)
	})
}

func (r runner) RunV3(rw io.ReadWriter) {
	p := protocol.New(redis.NewConn(rw))
	run(r.iterations, "v3", func(buf []byte) {
		runV3(p, buf, r.extraLogging)
	})
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

type loggingExecutor struct{ p *protocol.Protocol }

func (le loggingExecutor) Exec(cmd string, args protocol.ArgFunc, res protocol.ResFunc) {
	if err := protocol.Exec(le.p, cmd, args, res); err != nil {
		log.Fatalf("!%v: %v", strings.ToLower(cmd), err)
	}
}

type logger struct {
	str string
	log bool
}

func (sl logger) Status(data []byte, ok bool) {
	if !sl.log {
		return
	}
	if ok {
		log.Printf("+%s: %s", sl.str, data)
	} else {
		log.Printf("-%s: %s", sl.str, data)
	}
}

func (sl logger) Int(n int, status []byte) {
	if status != nil {
		log.Printf("-%s: %v", sl.str, status)
	} else if sl.log {
		log.Printf("+%s: %v", sl.str, n)
	}
}

type logFactory struct{ log bool }

func (lf logFactory) Status(cmd string) protocol.SimpleStrFunc {
	return logger{cmd, lf.log}.Status
}

func (lf logFactory) Int(cmd string) protocol.IntFunc {
	return logger{cmd, lf.log}.Int
}

func runV3(p *protocol.Protocol, b []byte, withResultLogging bool) {
	log := logFactory{log: withResultLogging}
	must := loggingExecutor{p}
	must.Exec("SET",
		protocol.WriteStrings("mykey", "my\x00value"),
		protocol.IgnoreOutput())
	must.Exec("SET",
		protocol.WriteStrings("mykey", "my\x00value"),
		protocol.AcceptStatus(log.Status("set")))
	must.Exec("EXISTS",
		protocol.WriteStrings("mykey", "kk", "does not exist"),
		protocol.AcceptInt(log.Int("exists")))
	for _, k := range []string{"mykey", "nosuchkey"} {
		must.Exec("GET",
			protocol.WriteStrings(k),
			protocol.AcceptBulk(b, log.Status("get")))
	}
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
