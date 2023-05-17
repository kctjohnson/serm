package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ProcessManager struct {
	Processes []Process
	Config    Config
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		Processes: []Process{},
		Config:    NewConfig(),
	}
}

type Process struct {
	Info   ProcessInfo
	Config Service
}

type ProcessInfo struct {
	PID     int
	OutPath string
	ErrPath string
}

func (pm *ProcessManager) OnStartup() error {
	for _, serv := range pm.Config.Services {
		if serv.OnStartup == true {
			err := pm.startProcess(serv)
			if err != nil {
				for retries := 0; err != nil && serv.Retry && retries < serv.RetryCount; retries++ {
					err = pm.startProcess(serv)
				}
			}

			if err != nil {
				fmt.Printf("Failed to start %s\nError: %s\n", serv.Name, err.Error())
			}
		}
	}
	return nil
}

func (pm *ProcessManager) startProcess(service Service) error {
	procInfo, err := spawnDetachedProcess(pm.Config.LogPath, service.Bin, service.Args, service.Name)
	if err != nil {
		return err
	}

	pm.Processes = append(pm.Processes, Process{
		Info:   procInfo,
		Config: service,
	})

	return nil
}

func spawnDetachedProcess(logPath string, binPath string, args []string, processName string) (ProcessInfo, error) {
	// Get clean paths
	logPathClean := filepath.Clean(logPath)

	// Eventually replace this with files
	curDateStamp := time.Now().Unix()
	outPath := fmt.Sprintf("%s/%s_out_%d.log", logPathClean, processName, curDateStamp)
	errPath := fmt.Sprintf("%s/%s_err_%d.log", logPathClean, processName, curDateStamp)

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
	}, nil
}
