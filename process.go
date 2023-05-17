package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ProcessInfo struct {
	PID     int
	OutPath string
	ErrPath string
	Cmd     *exec.Cmd
}

func spawnProcess(binPath string, args []string) (ProcessInfo, error) {
	// Eventually replace this with files
	sOut := os.Stdout
	sErr := os.Stderr

	fmt.Printf("Creating command for %s with args %v\n", binPath, args)
	cmd := exec.Command(binPath, args...)
	cmd.Stdout = sOut
	cmd.Stderr = sErr
	cmd.Dir = filepath.Dir(binPath)

	fmt.Printf("Starting command\n")
	err := cmd.Start()
	if err != nil {
		return ProcessInfo{}, err
	}

	fmt.Printf("Command PID: %d\n", cmd.Process.Pid)
	pid := cmd.Process.Pid

	return ProcessInfo{
		PID:     pid,
		OutPath: "",
		ErrPath: "",
		Cmd:     cmd,
	}, nil
}

func spawnDetachedProcess(binPath string, args []string) (ProcessInfo, error) {
	// Eventually replace this with files
	curDateStamp := time.Now().Unix()
	binName := filepath.Base(binPath)
	outPath := fmt.Sprintf("%s_out_%d.log", binName, curDateStamp)
	errPath := fmt.Sprintf("%s_err_%d.log", binName, curDateStamp)

	outFile, err := os.Create(outPath)
	if err != nil {
		return ProcessInfo{}, err
	}

	errFile, err := os.Create(errPath)
	if err != nil {
		return ProcessInfo{}, err
	}

	sOut := outFile
	sErr := errFile

	fmt.Printf("Creating command for %s with args %v\n", binPath, args)
	cmd := exec.Command(binPath, args...)
	cmd.Stdout = sOut
	cmd.Stderr = sErr
	cmd.Dir = filepath.Dir(binPath)

	fmt.Printf("Starting command\n")
	err = cmd.Start()
	if err != nil {
		return ProcessInfo{}, err
	}

	fmt.Printf("Command PID: %d\n", cmd.Process.Pid)
	pid := cmd.Process.Pid

	fmt.Printf("Releasing command\n")
	err = cmd.Process.Release()
	if err != nil {
		return ProcessInfo{}, err
	}

	return ProcessInfo{
		PID:     pid,
		OutPath: outPath,
		ErrPath: errPath,
		Cmd:     nil,
	}, nil
}
