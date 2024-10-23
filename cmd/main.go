package main

import (
	"fmt"
	"log"
	"time"

	"proj1/internal/pkg/server"
	"proj1/internal/pkg/storage"
)

const (
	file = "slice_storage.json"
)

func main() {
	// s, err := storage.NewStorage()
	// if err != nil {
	// 	panic(err)
	// }
	//

	stor2, err := storage.NewSliceStorage()
	if err != nil {
		log.Fatal(err)
	}
	stor2.LoadFromFile(file)
	closeChan := make(chan struct{})
	go stor2.IWantToSleepFor(closeChan, 10*time.Second, file)
	srv := server.New(":8090", &stor2)
	srv.Start()
	// stor2.Set("intval", "ingfdt")
	// fmt.Println(stor2.Get("intval"))
	// fmt.Println(stor2.HSet("mkey", []map[string]string{{"hs1": "v1"}, {"hs3": "v3"}}))
	// fmt.Println(stor2.HSet("mkey2", []map[string]string{{"hs2": "77"}}))
	//res, _ := stor2.Get("intval")
	//
	// stor2.LPush("key1", []string{"hhh", "gggg"})
	// stor2.LPush("key1", []string{"llll", "aaaa"})
	// stor2.RPush("key2", []string{"kkkk", "pppp"})
	// stor2.RAddToSet("key1", []string{"jjjj", "hhh", "gggg", "mmmm"})
	// fmt.Println(stor2)
	// fmt.Println(stor2.LGet("key1", 3))
	// fmt.Println(stor2.RPop("key1", 2, 78))
	// res, _ := stor2.HGet("mkey", "hs1")
	// fmt.Println(*res)
	// res, _ = stor2.HGet("mkey2", "hs2")
	// if res != nil {
	// 	fmt.Println(*res)
	// }
	// fmt.Println(stor2.LPop("key1", -6, 4))
	// fmt.Println(stor2.LPop("key2", 1))
	// fmt.Println(stor2.GetKind("intvallll"))
	// s.Set("key1", "val1")
	// s.Set("key2", "754")
	// res1, ok := s.Get("key1")
	// if ok {
	// 	fmt.Println(res1, s.GetKind("key1"))
	// }
	// res2, ok := s.Get("key2")
	// if ok {
	// 	fmt.Println(res2, s.GetKind("key2"))
	// }
	err = stor2.SaveToFile("slice_storage.json")
	if err != nil {
		fmt.Println("Error saving storage:", err)
	}
	//s.SaveToFile("storage.json")

}
