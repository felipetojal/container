package main

import (
	"fmt"
	"log"
	"os"
)

const ALPINE_ROOT = "/home/lfcor/alpine-rootfs"

func main() {
	args := os.Args

	switch args[1] {
	case "run":
		if err := ParentNamespace(); err != nil {
			fmt.Println("Error: " + err.Error())
			log.Fatal(err)
		}
		fmt.Println("Exiting ParentNamespace()")
	case "child":
		ChildProcess()
		fmt.Println("Exiting ChildProcess()")
	default:
		fmt.Println("Command not identified")
		os.Exit(1)
	}

}