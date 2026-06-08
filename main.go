package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	hostname := "binha-space"

	must(ParentNamespace())
	must(ChildNamespace(hostname))
	fmt.Println("hi")
}

func ParentNamespace() (error) {

	flags := syscall.CLONE_NEWNS | // Creates a new filesystem
		syscall.CLONE_NEWIPC | // Only process within the namespace will be able to communicate
		syscall.CLONE_NEWNET | // Creates a new network stack for the namespace
		syscall.CLONE_NEWUSER | // New UID and GID for the namespace
		syscall.CLONE_NEWPID | // New process ID tree
		syscall.CLONE_NEWCGROUP | // Creates a new cgroup for the namespace
		syscall.CLONE_NEWUTS // Detaches the host name from the namespace name

	callerUID := os.Getuid() // User ID
	callerGID := os.Getgid() // Group ID

	// Maps the OS user GID to a new GID inside the namespace
	gidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID: callerGID,
			Size: 1,
		},
	}

	
	// Maps the OS user UID to a new UID inside the namespace
	uidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID: callerUID,
			Size: 1,
		},
	}

	// "/proc/self/exe" -> points to the binary executable file of the program that is currently running
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// Sets the namespace attributes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: uintptr(flags),
		Pdeathsig: syscall.SIGTERM,
		GidMappings: gidMappings,
		UidMappings: uidMappings,	
	}

	// Matching the pipelines		
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func ChildNamespace(hostname string) (error) {
	
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		return err
	}

	return nil
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}