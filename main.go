package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("expect at least argument 'command'")
	}
	cmd := os.Args[1]
	os.Args = os.Args[1:]
	switch cmd {
	case "fmt":
		runFmt()
	}
}
