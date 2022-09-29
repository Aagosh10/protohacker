package main

import (
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	exitOk = iota
	exitError
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	defer stop()
	server, err := net.Listen("tcp", ":9090")
	if err != nil {
		return exitError
	}

	defer server.Close()

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				x, err := io.ReadAll(conn)
				_, err = conn.Write(x)
				if err != nil {
					return
				}
			}(conn)
		}
	}()

	<-ctx.Done()

	return exitOk
}
