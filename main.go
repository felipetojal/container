package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
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

// ParentNamespace is responsible for creating the namespace,
// basically a sandbox within the OS.
// It then makes a call to a child process to run the program
// we want inside this namespace.
func ParentNamespace() error {
	flags := syscall.CLONE_NEWNS | // Creates a new filesystem
		syscall.CLONE_NEWIPC | // Only process within the namespace will be able to communicate
		syscall.CLONE_NEWNET | // Creates a new network stack for the namespace
		syscall.CLONE_NEWUSER | // New UID and GID for the namespace
		syscall.CLONE_NEWPID | // New process ID tree
		syscall.CLONE_NEWCGROUP | // Creates a new cgroup for the namespace
		syscall.CLONE_NEWUTS // Detaches the host name from the namespace name

	log.Printf("Flag: %v\n", flags)

	callerUID := os.Getuid() // User ID
	log.Printf("callerUID: %v\n", callerUID)

	callerGID := os.Getgid() // Group ID
	log.Printf("callerGID: %v\n", callerGID)

	// Maps the OS user GID to a new GID inside the namespace
	gidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      callerGID,
			Size:        1,
		},
	}
	log.Printf("gidMappings: %v", gidMappings)

	// Maps the OS user UID to a new UID inside the namespace
	uidMappings := []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID:      callerUID,
			Size:        1,
		},
	}
	log.Printf("uidMappings: %v", uidMappings)

	// "/proc/self/exe" -> points to the binary executable file of the program that is currently running
	// "child" -> becomes os.Args[1] in the newly spawned process
	// "/bin/sh" -> becomes os.Args[2] in the newly spawned process
	cmd := exec.Command("/proc/self/exe", "child", "/bin/sh")

	// Sets the namespace attributes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:  uintptr(flags),
		Pdeathsig:   syscall.SIGTERM,
		GidMappings: gidMappings,
		UidMappings: uidMappings,
	}
	log.Printf("cmd.SysProcAttr: %v", cmd.SysProcAttr)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	fmt.Printf("Parent PID: %v\n", os.Getpid())

	return cmd.Run()
}

func ChildProcess() {
	if err := syscall.Sethostname([]byte("binha-host")); err != nil {
		log.Println("Error setting hostname")
		os.Exit(1)
	}

	if hostname, err := os.Hostname(); err != nil {
		fmt.Println("Error retrieving hostname")
		log.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Println(hostname)
	}

	// Changing the root filesystem
	// It tells the kernel: "Whenever this specific process asks to look at /, 
	// do not show it the real hard drive. 
	// Redirect its eyes to /home/lfcor/alpine-rootfs instead."
	syscall.Chroot(ALPINE_ROOT)

	// Forces the process to go to its new root.
	syscall.Chdir("/")

	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2])
	
	// Attaching the OS pipeline to the child process
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	fmt.Printf("Child PID: %v\n", os.Getpid())
	fmt.Println("Hello boyz")

	cmd.Run()

	//Cleaning up the virtual filesystem
	fmt.Println("Unmount()")
	syscall.Unmount("proc", 0)
}
