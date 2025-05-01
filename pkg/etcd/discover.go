package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func DisCover() ([]string, int64) {
	res, err := EtcdClient().Get(context.Background(), Prefix, clientv3.WithPrefix())
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
