package main

import (
	"flag"
	"fmt"
	"time"
)

func logNumbers() {
	nums := "0123456789"

	c := 0
	for {
		fmt.Printf("%c\n", nums[c])
		c = (c + 1) % len(nums)
		time.Sleep(time.Second / 4)
	}
}

func logAlphabet() {
	alph := "abcdefghijklnopqrstuvwxyz"

	c := 0
	for {
		fmt.Printf("%c\n", alph[c])
		c = (c + 1) % len(alph)
		time.Sleep(time.Second / 4)
	}
}

func logSymbols() {
	symb := "!@#$%^&*()_+"

	c := 0
	for {
		fmt.Printf("%c\n", symb[c])
		c = (c + 1) % len(symb)
		time.Sleep(time.Second / 4)
	}
}

func main() {
	typePtr := flag.String("type", "num", "Log Type")
	flag.Parse()

	switch *typePtr {
	case "num":
		logNumbers()
	case "alph":
		logAlphabet()
	case "symb":
		logSymbols()
	default:
		logNumbers()
	}
}
