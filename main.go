package main

import (
	"GGCache/internal/group"
	"GGCache/internal/server"
	"fmt"
	"log"
	"net/http"
	"os"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	_ = group.NewGroup("scores", 2<<10, func(key string) ([]byte, bool) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), true
		}
		return nil, false
	})
	selfAddr := os.Args[1]
	fmt.Println("selfAddr:", selfAddr)
	s := server.NewHTTPPool(selfAddr)
	s.RegisterNode(selfAddr)
	s.RegisterNode(os.Args[2])
	s.RegisterNode(os.Args[3])

	log.Println("server listen at ", selfAddr)
	log.Println(http.ListenAndServe(selfAddr, s))
}
