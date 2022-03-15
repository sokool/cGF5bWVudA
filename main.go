package main

import (
	"log"
)

func main() {
	if err := NewService().Run(); err != nil {
		log.Fatal(err)
	}
}
