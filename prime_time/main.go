package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type ReqData struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type RespData struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	defer stop()
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
			buf := bufio.NewReader(conn)
			go func(conn net.Conn) {
				defer conn.Close()
				for {
					req, err := buf.ReadBytes('\n')
					if err == io.EOF {
						return
					}
					if err != nil {
						panic(err)
					}

					var reqData ReqData
					err = json.Unmarshal(req, &reqData)
					if err != nil {
						conn.Write([]byte(err.Error()))
						return
					}
					if err = validateReq(reqData); err != nil {
						conn.Write([]byte(err.Error()))
						return
					}
					resp, err := prepareResponse(reqData)
					if err != nil {
						conn.Write([]byte(err.Error()))
						return
					}
					conn.Write(append(resp, []byte("\n")...))
					fmt.Println(string(resp))
				}

			}(conn)
		}
	}()

	<-ctx.Done()

	return 0
}

func validateReq(reqData ReqData) error {
	if reqData.Method == nil || reqData.Number == nil {
		return errors.New("invalid request")
	}
	if *reqData.Method != "isPrime" {
		return errors.New("invalid request method")
	}

	return nil
}

func prepareResponse(data ReqData) ([]byte, error) {
	var isPrime bool
	if *data.Number == float64(int(*data.Number)) {
		isPrime = true
	}

	if isPrime && checkPrime(*data.Number) {
		isPrime = true
	} else {
		isPrime = false
	}
	resp := RespData{
		Method: "isPrime",
		Prime:  isPrime,
	}

	return json.Marshal(resp)
}

func checkPrime(data float64) bool {
	for i := 2.0; i <= math.Sqrt(data); i++ {
		if math.Mod(data, i) == 0.0 {
			return false
		}
	}

	return data > 1
}
