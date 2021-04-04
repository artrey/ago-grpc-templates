package main

import (
	"log"
	"os"
)

func main() {
	cli := CLI{}

	if err := cli.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
