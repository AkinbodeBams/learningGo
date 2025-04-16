package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define a flag (--name) with a default value and a description
	name := flag.String("name", "Guest", "Your name")

	// Parse command-line flags
	flag.Parse()

	// Print the greeting
	fmt.Println("Hello,", *name)
}
