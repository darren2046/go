package main

import (
	. "github.com/darren2046/go"
)

func main() {
	Lg.Trace("Started")
	for i := range Os.Stdin.Readlines() {
		Os.Stdout.Write(i.S + "--------")
	}
}
