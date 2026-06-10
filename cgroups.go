package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const CGROUP_DIR = "/sys/fs/cgroup"

// setupCgroups will initialize the container cgroups. It is
// responsible for limiting the resources that the container
// can use.
func setupCgroups(pid, cpuMax, memMax, pidMax int, contName string) error {
	err, dirName := setupCgroupDir(contName)
	if err != nil {
		return fmt.Errorf("Error in setupCgroups(): %w", err)
	}
	log.Printf("Directory created")

	if err := setupPidCgroup(pid, dirName); err != nil {
		return fmt.Errorf("Error in setupCgroups(): %w", err)
	}
	log.Printf("PID cgroup written")

	if err := setupMemoryMax(memMax, dirName); err != nil {
		return fmt.Errorf("Error in setupCgroups(): %w", err)
	}
	log.Printf("Memory max written")

	if err := setupMaxPid(pidMax, dirName); err != nil {
		return fmt.Errorf("Error in setupCgroups(): %w", err)
	}
	log.Printf("PID max written")

	if err := setupMaxCpu(cpuMax, dirName); err != nil {
		return fmt.Errorf("Error in setupCgroups(): %w", err)
	}
	log.Printf("CPU max written")

	return nil
}

// deleteCgroup receives the container name and deletes
// the directory associated to it in the cgroups.
func deleteCgroup(name string) error {
	dirName := fmt.Sprintf("%s/%s", CGROUP_DIR, name)
	if err := os.RemoveAll(dirName); err != nil {
		return fmt.Errorf("Error in deleteCgroup(): %w", err)
	}
	return nil
}

// setupCgroupDir creates the cgroup directory.
// Once the dir is created, the Linux Kernel automatically populates it
// with virtual text files.
func setupCgroupDir(name string) (error, string) {
	dirName := fmt.Sprintf("%s/%s", CGROUP_DIR, name)
	if err := os.Mkdir(dirName, 0777); err != nil {
		return fmt.Errorf("Error creating dir: %w", err), ""
	}
	return nil, dirName
}

// setupPidCgroup writes the process PID associated to the folder.
func setupPidCgroup(pid int, dirPath string) error {
	filePath := fmt.Sprintf("%s/cgroup.procs", dirPath)
	data := []byte(strconv.Itoa(pid))
	if err := os.WriteFile(filePath, data, 0777); err != nil {
		return fmt.Errorf("Error writing cgroup.procs: %w", err)
	}
	return nil
}

// setupMemoryMax writes to the memory.max virtual file. It is responsible
// for specifying how much RAM the process is allowed to consume.
func setupMemoryMax(max int, dirPath string) error {
	filePath := fmt.Sprintf("%s/memory.max", dirPath)
	data := []byte(strconv.Itoa(max))
	if err := os.WriteFile(filePath, data, 0777); err != nil {
		return fmt.Errorf("Error writing memory.max: %w", err)
	}
	return nil
}

// setupMaxPid writes to the pids.max virtual file. It is responsible
// spcecifying the maximum number of process the container can hold.
func setupMaxPid(max int, dirName string) error {
	filePath := fmt.Sprintf("%s/pids.max", dirName)
	data := []byte(strconv.Itoa(max))
	if err := os.WriteFile(filePath, data, 0777); err != nil {
		return fmt.Errorf("Error writing max PID: %w", err)
	}
	return nil
}

// setupMaxCpu writes to the cpu.max virtual file. It is responsible for
// for specifying the maximum amount of CPU alocated to the container
func setupMaxCpu(max int, dirName string) error {
	filePath := fmt.Sprintf("%s/cpu.max", dirName)
	data := []byte(fmt.Sprintf("%d 100000", max))
	if err := os.WriteFile(filePath, data, 0777); err != nil {
		return fmt.Errorf("Error writing cpu.max: %w", err)
	}
	return nil
}