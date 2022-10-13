package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

func startChat(room *room, userName string, conn net.Conn) {
	defer conn.Close()

	bufReader := bufio.NewReader(conn)

	for {
		msg, err := bufReader.ReadBytes('\n')
		if err != nil && err == io.EOF {
			room.removeUser(userName, conn)
			return
		} else if err != nil {
			return
		}
		msg = bytes.TrimSpace(msg)
		if len(msg) == 0 {
			continue
		}
		finalMsgText := fmt.Sprintf("[%s] %s\n", userName, string(msg))

		err = room.publish([]byte(finalMsgText), userName)
		if err != nil {
			return
		}
	}
}
