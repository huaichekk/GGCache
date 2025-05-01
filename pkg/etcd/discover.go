package etcd

import (
	"GGCache/internal/server"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func DisCover() ([]string, int64) {
	res, err := EtcdClient().Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	version := res.Header.Revision
	addrs := make([]string, 0)
	for _, kv := range res.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, version
}
func WatchFromVersion(version int64, sve *server.HTTPPool) {
	// 先获取当前所有节点，建立已存在节点集合
	currentNodes := make(map[string]struct{})
	if res, err := EtcdClient().Get(context.Background(), prefix, clientv3.WithPrefix()); err == nil {
		for _, kv := range res.Kvs {
			currentNodes[string(kv.Value)] = struct{}{}
		}
	}

	watcher := EtcdClient().Watch(context.Background(),
		prefix,
		clientv3.WithPrefix(),
		clientv3.WithRev(version),
		clientv3.WithPrevKV())

	for resp := range watcher {
		for _, event := range resp.Events {
			addr := string(event.Kv.Value)
			if _, exists := currentNodes[addr]; exists {
				continue // 跳过已存在的节点
			}

			switch {
			case event.IsCreate():
				fmt.Println("[NEW]", addr)
				sve.RegisterNode(addr)
				currentNodes[addr] = struct{}{}

			case event.Type == clientv3.EventTypeDelete:

				a := string(event.PrevKv.Value)
				fmt.Println("[DEL]", a)
				delete(currentNodes, a)
				sve.DeleteNode(a)

			}
		}
	}
}
