package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"flag"
)

func main(){
	lines := flag.Bool("l",false,"Count lines")
	flag.Parse()
fmt.Println(count(os.Stdin, *lines))
}

func count(text io.Reader, lines bool) int{

 scanner := bufio.NewScanner(text)
 
 if !lines{
 scanner.Split(bufio.ScanWords)
	
 }
 wc := 0

 for scanner.Scan(){wc++}


 return wc
}