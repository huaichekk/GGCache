package etcd

import (
	"fmt"
	"testing"
)

func TestEtcdClient(t *testing.T) {
	etcdClient := EtcdClient()
	fmt.Println(etcdClient)
}
