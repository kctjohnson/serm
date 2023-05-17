package main

import (
	"fmt"
	"os"
	"time"
)

// Service manager handles managing all of the different services, their status, and their logs
//
// Need a config file that stores info about each service
// - Retry On Error
// - Binary location
// - Flags
//
// Keeps track of PID and

func main() {
	procMan := NewProcessManager()
	err := procMan.OnStartup()
	if err != nil {
		return
	}

	for i := 0; i < 20; i++ {
		fmt.Printf("Processes:\n\n")
		for _, proc := range procMan.Processes {
			fmt.Printf("Name: %s\n", proc.Config.Name)
			fmt.Printf("PID: %d\n\n\n", proc.Info.PID)
		}
		time.Sleep(time.Second)
	}

	for _, proc := range procMan.Processes {
		p, err := os.FindProcess(proc.Info.PID)
		if err != nil {
			return
		}
		p.Kill()
	}
}
