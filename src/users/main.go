package main

import (
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "saltgram-users", log.LstdFlags)
}