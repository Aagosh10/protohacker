package main

import (
	"context"
	"fmt"
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

	fmt.Println("listening on port 10000...")
	server, err := net.Listen("tcp", ":10000")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				panic(err)
			}
			go func(c net.Conn) {
				defer c.Close()
				x, err := io.ReadAll(conn)
				if err != nil {
					panic(err)
				}
				conn.Write(x)
			}(conn)
		}
	}()

	<-ctx.Done()
	return exitOk
}
