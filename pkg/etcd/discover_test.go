package etcd

import (
	"GGCache/internal/server"
	"fmt"
	"testing"
	"time"
)

func TestDisCover(t *testing.T) {
	//RegisterSelfAddr("node1", "127.0.0.1:8888")
	//RegisterSelfAddr("node2", "127.0.0.1:7777")
	//RegisterSelfAddr("node3", "127.0.0.1:6666")
	addrs, version := DisCover()
	s := server.NewHTTPPool("127.0.0.1:8888")
	fmt.Println(addrs)
	for _, addr := range addrs {
		s.RegisterNode(addr)
	}
	go WatchFromVersion(version, s)
	time.Sleep(1 * time.Hour)
}
