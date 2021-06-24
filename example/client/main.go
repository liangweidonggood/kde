package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	//addr:="47.114.51.90:18888"
	addr := "127.0.0.1:8999"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	defer conn.Close()
	go func() {
		input := make([]byte, 1024)
		s := []string{`P`, `1`, `55`, `2`, `9032`, `1`, `13`, `16`, `5`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `0`, `B`, `K`}
		for {
			s[9] = strconv.Itoa(rand.Intn(100))
			s[38] = strconv.Itoa(128)
			str := strings.Join(s, `*`)
			log.Println(str)
			input = []byte(str + "\n")
			conn.Write(input)
			time.Sleep(5 * time.Second)
			break
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("read err:", err)
			continue
		}
		fmt.Println(string(buf[:n]))

	}
}
