package etcd

import (
	"GGCache/configs"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var (
	client *clientv3.Client
	once   sync.Once
)

func EtcdClient() *clientv3.Client {
	var err error
	if client == nil {
		once.Do(func() {
			client, err = clientv3.New(clientv3.Config{
				Endpoints:   configs.GetConfig().EtcdConfig.Addr,
				DialTimeout: 5 * time.Second,
			})
		})
	}
	if err != nil {
		panic(err)
	}
	return client
}
