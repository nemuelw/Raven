package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("[*] Initializing Raven ...")
}

func has_internet_access() bool {
	_, err := http.Get("https://google.com")
	return err == nil
}