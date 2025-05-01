package client

import (
	"testing"
)

func TestClient_Get(t *testing.T) {
	client := NewClient("127.0.0.1:8888")
	if v, ok := client.Get("scores", "Tom"); !ok || v.String() != "630" {
		t.Error(v.String())
	}
}
