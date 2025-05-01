package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

const prefix = "/service/GGCache"

func RegisterSelfAddr(key, selfAddr string) {
	k := fmt.Sprintf("%s/%s", prefix, key)
	c := EtcdClient()
	grant, err := c.Grant(context.Background(), 3)
	if err != nil {
		panic(err)
	}
	_, err = c.Put(context.Background(), k, selfAddr, clientv3.WithLease(grant.ID))
	if err != nil {
		panic(err)
	}
	go KeepAlive(grant.ID)
	time.Sleep(1 * time.Hour)
}

func KeepAlive(id clientv3.LeaseID) {
	ch, err := client.KeepAlive(context.Background(), id)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				return
			}
		}
	}
}
