package rpcclient

import (
	"context"
	"log"
	"time"

	pb "GGCache/pkg/rpc" // 替换为您的protobuf包路径
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.CacheServiceClient
}

// NewClient 创建gRPC客户端
func NewClient(addr string, timeout time.Duration) (*GRPCClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024), // 10MB
		),
	)
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		conn:   conn,
		client: pb.NewCacheServiceClient(conn),
	}, nil
}

// Get 实现缓存获取
func (c *GRPCClient) Get(group, key string) ([]byte, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := c.client.GetRpc(ctx, &pb.GetRequest{
		Group: group,
		Key:   key,
	})
	if err != nil {
		log.Printf("gRPC Get failed: %v", err)
		return nil, false
	}

	return resp.Value, resp.Found
}

// Close 关闭连接
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}
