package main

import (
	"context"
	"fmt"
	"homework/internal/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())

	go shutdownListener(cancel)

	http.StartServer(8080)
	// grpc.StartServer(8081)

	<-ctx.Done()
	cancel()

	os.Exit(0)
}

func shutdownListener(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	for {
		sig := <-sigs
		fmt.Fprintf(os.Stderr, "=== received %v ===\n", sig)

		fmt.Fprintln(os.Stderr, "*** goroutine dump...")

		// grpc.StopServer()
		http.StopServer()

		var (
			buf      []byte
			bufsize  int
			stacklen int
		)

		for bufsize = 1e6; bufsize < 100e6; bufsize *= 2 {
			buf = make([]byte, bufsize)
			stacklen = runtime.Stack(buf, true)

			if stacklen < bufsize {
				break
			}
		}

		fmt.Fprintln(os.Stderr, string(buf[:stacklen]))
		fmt.Fprintln(os.Stderr, "*** end of dump")

		cancel()
	}
}
