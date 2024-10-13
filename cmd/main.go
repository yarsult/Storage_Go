package main

import (
	"fmt"
	"log"
	"proj1/internal/pkg/storage"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	s.Set("key1", "val1")
	s.Set("key2", "754")
	fmt.Println(*s.Get("key1"), s.GetKind("key1"))
	fmt.Println(*s.Get("key2"), s.GetKind("key2"))
}
