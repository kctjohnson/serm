package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type ProcessInfo struct {
	PID     int
	OutPath string
	ErrPath string
}

type Process struct {
	ProcInfo ProcessInfo
	Service  Service
}

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

func (pm *ProcessManager) OnStartup() error {
	for _, serv := range pm.Config.Services {
		// Check to see if we're already running the service
		running, err := pm.ServiceRunning(serv)
		if err != nil {
			return err
		}

		// Begin the start process for this service if it runs on startup
		if serv.OnStartup == true && !running {
			err := pm.StartProcess(serv)
			if err != nil {
				for retries := 0; err != nil && serv.Retry && retries < serv.RetryCount; retries++ {
					err = pm.StartProcess(serv)
				}
			}

			if err != nil {
				fmt.Printf("Failed to start %s\nError: %s\n", serv.Name, err.Error())
			}
		}
	}
	return nil
}

func (pm ProcessManager) ServiceRunning(service Service) (bool, error) {
	// First check to see if the service is in a running process currently
	var foundProc *Process
	for _, proc := range pm.Processes {
		if proc.Service.Name == service.Name {
			foundProc = &proc
			break
		}
	}

	// Then check to see if the PID is running
	if foundProc != nil {
		if foundProc.ProcInfo.PID <= 0 {
			return false, fmt.Errorf("invalid PID %v", foundProc.ProcInfo.PID)
		}

		proc, err := os.FindProcess(int(foundProc.ProcInfo.PID))
		if err != nil {
			return false, err
		}

		err = proc.Signal(syscall.Signal(0))
		if err == nil {
			return true, nil
		}

		if err.Error() == "os: process already finished" {
			return false, nil
		}

		errno, ok := err.(syscall.Errno)
		if !ok {
			return false, err
		}

		switch errno {
		case syscall.ESRCH:
			return false, nil
		case syscall.EPERM:
			return true, nil
		}

		return false, err
	} else {
		// TODO: Reassess if we would ever want to look for detached and missing procs
		// 				Mostly because that shoooouldn't ever really happen if we keep track
		// 				of the processes with a sqlite DB?
		// // Search for the process
		// cmdStr := fmt.Sprintf("ps aux | grep \"%s\" | awk '{print $2}'", "./cmd/printer/printer -type alph")
		// cmd := exec.Command("bash", "-c", cmdStr)
		// outBytes, err := cmd.Output()
		// if err != nil {
		// 	return false, err
		// }
		//
		// outStr := string(outBytes)
		// pids := strings.Split(outStr, "\n")
		// if len(pids) >= 4 {
		// 	programPid := pids[len(pids)-2]
		// 	pm.Processes = append(pm.Processes, Process{
		// 		ProcInfo: ProcessInfo{
		//
		// 		},
		// 		Service:  service,
		// 	})
		// }
	}

	return false, nil
}

func (pm *ProcessManager) StartProcess(service Service) error {
	procInfo, err := spawnDetachedProcess(pm.Config.LogPath, service.Bin, service.Args, service.Name)
	if err != nil {
		return err
	}

	pm.Processes = append(pm.Processes, Process{
		ProcInfo: procInfo,
		Service:  service,
	})

	return nil
}

func (pm *ProcessManager) StopProcess(proc Process) error {
	// Get the actual process
	foundProc, err := os.FindProcess(proc.ProcInfo.PID)
	if err != nil {
		return err
	}

	// Kill the process
	err = foundProc.Kill()
	if err != nil {
		return err
	}

	return nil
}

func (pm *ProcessManager) RestartProcess(proc Process) error {
	err := pm.StopProcess(proc)
	if err != nil {
		return err
	}
	err = pm.StartProcess(proc.Service)
	if err != nil {
		return err
	}
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

	cmd := exec.Command(binPath, args...)
	cmd.Stdout = outFile
	cmd.Stderr = errFile
	cmd.Dir = filepath.Dir(binPath)

	err = cmd.Start()
	if err != nil {
		return ProcessInfo{}, err
	}

	pid := cmd.Process.Pid

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
