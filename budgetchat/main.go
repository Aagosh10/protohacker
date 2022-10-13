package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

const (
	exitOk = iota
	exitError
)

var (
	alphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	defer stop()
	currUsers := make(map[string]net.Conn)
	server, err := net.Listen("tcp", ":10000")
	if err != nil {
		panic(err)
	}
	room := NewRoom()
	defer server.Close()

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				panic(err)
			}

			go handleNewConnection(conn, currUsers, room)
		}
	}()

	<-ctx.Done()

	return exitOk
}

func handleNewConnection(conn net.Conn, currUsers map[string]net.Conn, room *room) {
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	buf := bufio.NewReader(conn)
	name, err := buf.ReadBytes('\n')
	if err != nil {
		conn.Close()
		return
	}
	name = bytes.TrimSpace(name)
	if len(name) == 0 {
		log.Println("invalid username")
		conn.Close()
		return
	}

	if !alphaNumeric.Match(name) {
		log.Println("invalid username")
		conn.Close()
		return
	}
	existingUsers := strings.Join(room.getUsers(), ",")
	err = room.addUser(strings.TrimSpace(string(name)), conn)
	if err != nil {
		conn.Close()
		return
	}
	conn.Write([]byte(fmt.Sprintf("* The room contains: %s\n", existingUsers)))

	go startChat(room, string(name), conn)
}
