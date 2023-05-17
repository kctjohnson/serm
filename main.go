package main

import "fmt"

// Service manager handles managing all of the different services, their status, and their logs
//
// Need a config file that stores info about each service
// - Retry On Error
// - Binary location
// - Flags
//
// Keeps track of PID and

func main() {
	cfg := NewConfig()
	fmt.Printf("%#v\n", cfg)

	// pInfo, err := spawnProcess("/usr/bin/tail", []string{"-f", "/USers/killean/Documents/Projects/go/service-manager/cmd/printer/print.log"})
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("%#v\n", pInfo)
	//
	// for {
	// 	// out, err := pInfo.Cmd.CombinedOutput()
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	// fmt.Printf("Out: %s\n", out)
	// }
}
