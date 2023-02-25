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

func c2_comms(conn net.Conn) {

}

func connect_to_c2() {
	if has_internet_access() {
		addr := fmt.Sprintf("%s:%d", HOST, PORT)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(time.Second * 30)
			connect_to_c2()
		} else {
			// communication with C2
			c2_comms(conn)
		}
	} else {
		time.Sleep(time.Second * 30)
		connect_to_c2()
	}
}