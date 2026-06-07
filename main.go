package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Println("hi")
}

func parent() {
	flags := syscall.CLONE_NEWNS | // Creates a new filesystem
		syscall.CLONE_NEWIPC | // Only process within the namespace will be able to communicate
		syscall.CLONE_NEWNET | // Creates a new network stack for the namespace
		syscall.CLONE_NEWUSER | // New UID and GID for the namespace
		syscall.CLONE_NEWPID | // New process ID tree
		syscall.CLONE_NEWCGROUP | // Creates a new cgroup for the namespace
		syscall.CLONE_NEWUTS // Detaches the host name from the namespace name

	callerUID := os.Getuid() // User ID
	callerGID := os.Getgid() // Group ID

	

}
