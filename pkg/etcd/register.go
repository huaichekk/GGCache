package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

const Prefix = "/service/GGCache"

func RegisterSelfAddr(key, selfAddr string) {
	k := fmt.Sprintf("%s/%s", Prefix, key)
	c := EtcdClient()
	grant, err := c.Grant(context.Background(), 3)
	if err != nil {
		panic(err)
	}
	_, err = c.Put(context.Background(), k, selfAddr, clientv3.WithLease(grant.ID))
	if err != nil {
		panic(err)
	}
	go keepAlive(grant.ID)
}

func keepAlive(id clientv3.LeaseID) {
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
