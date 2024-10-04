package main

import (
	"fmt"
	"log"
	"proj1/cmd/internal/pkg/storage"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	s.Set("key1", "val1")
	fmt.Println(s.Get("key1"))
}
