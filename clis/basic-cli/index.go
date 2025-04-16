package main

import (
	"fmt"
	"os"
)

func main() {
	// Check if an argument is provided
	println(len(os.Args))
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <your-name>")  // Expected to print if no args are provided
		return  // Exit the program early
	}

	// Print the name argument
	name := os.Args[1]
	fmt.Println("Hello,", name)
}