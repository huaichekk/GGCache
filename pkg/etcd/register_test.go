package etcd

import "testing"

func TestRegisterSelfAddr(t *testing.T) {
	RegisterSelfAddr("node1", "127.0.0.1:8888")
}

func TestRegisterSelfAddr2(t *testing.T) {
	RegisterSelfAddr("node2", "127.0.0.1:7777")
}
func TestRegisterSelfAddr3(t *testing.T) {
	RegisterSelfAddr("node3", "127.0.0.1:6666")
}
