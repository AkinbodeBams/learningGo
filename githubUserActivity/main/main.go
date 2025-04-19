package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	var gitHubUsername string
	println("Enter github Username")
	fmt.Scan(&gitHubUsername)
	url := fmt.Sprintf("https://api.github.com/users/%s/events", gitHubUsername)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("Response body:")
	fmt.Println(string(body))
}