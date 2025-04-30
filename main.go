package main

import (
	"GGCache/internal/group"
	"GGCache/internal/server"
	"log"
	"net/http"
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
	s := server.NewHTTPServer("127.0.0.1:6666")
	log.Println(http.ListenAndServe(":8888", s))
}
