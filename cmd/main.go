package main

import (
	"flag"
	"fmt"
	"github.com/ch4t5ky/armadillo/internal/helpers"
)

func main() {
	actionPtr := flag.String("--action", "lock", "Start action with files in directory")
	flag.Parse()
	action := *actionPtr
	switch action {
	case "lock":
		helpers.LockFiles()
	case "unlock":
		break
	default:
		fmt.Println("Attempt to launch a non-existent action")
	}
}
