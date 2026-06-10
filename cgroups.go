package main

import (
	"fmt"
	"os"
	"strconv"
)

func setupCgroups(pid int, name string) error {
	if err := setupCgroupDir(name); err != nil {
		return fmt.Errorf("Error in setupCgroups(): %w", err)
	}

	if err := setupPidCgroup(pid, name); err != nil {
		return fmt.Errorf("Error in setupCgroups(): %w", err)
	}

	return nil
}

// setupCgroupDir creates the cgroup directory.
// Once the dir is created, the Linux Kernel automatically populates it
// with virtual text files.
func setupCgroupDir(name string) error {
	dirName := fmt.Sprintf("/sys/fs/cgroups/%s", name)
	if err := os.Mkdir(dirName, 0755); err != nil {
		return fmt.Errorf("Error creating dir: %w", err)
	}
	return nil
}

// setupPidCgroup writes the process PID associated to the folder.
func setupPidCgroup(pid int, dirPath string) error {
	filePath := fmt.Sprintf("%s/cgroups.procs", dirPath)
	data := []byte(strconv.Itoa(pid))
	if err := os.WriteFile(filePath, data, 0755); err != nil {
		return fmt.Errorf("Error writing cgroups.procs: %w", err)
	}
	return nil
}
