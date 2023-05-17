package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	logFile, err := os.Create("print.log")
	if err != nil {
		panic(err)
	}

	c := 'a'
	for {
		_, err = logFile.WriteString(fmt.Sprintf("%c\n", c))
		if err != nil {
			panic(err)
		}
		c = ((c - 97 + 1) % 26) + 97
		time.Sleep(time.Second / 4)
	}
}
