package main

import (
	. "github.com/ChaunceyShannon/golanglibs"
)

func main() {
	Lg.Trace("Started")
	for i := range Os.Stdin.Readlines() {
		Os.Stdout.Write(i.S + "--------")
	}
}
