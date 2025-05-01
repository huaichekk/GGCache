package server

import (
	"GGCache/configs"
	"GGCache/internal/cache/eviction"
	"GGCache/internal/consistent"
	"GGCache/internal/group"
	"GGCache/pkg/etcd"
	pb "GGCache/pkg/rpc"
	"GGCache/pkg/rpc/rpcclient"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const HTTPADDR = "127.0.0.1:9999"

type HTTPPool struct {
	pb.UnimplementedCacheServiceServer // 必须内嵌
	selfAddr                           string
	mu                                 sync.Mutex
	consistent                         *consistent.Consistent
	nodes                              map[string]*rpcclient.GRPCClient
}

func NewHTTPPool(addr string) *HTTPPool {
	return &HTTPPool{
		selfAddr:   addr,
		mu:         sync.Mutex{},
		consistent: consistent.NewConsistent(configs.GetConfig().CacheConfig.Replicas, nil),
		nodes:      make(map[string]*rpcclient.GRPCClient),
	}
}

func (s *HTTPPool) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//http://127.0.0.1/group/key
	params := strings.Split(req.URL.Path, "/")
	params = params[1:]
	if len(params) != 2 {
		fmt.Println(params, len(params))
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}
	if v, err := s.Get(params[0], params[1]); err == nil {
		resp.Header().Set("Content-Type", "application/octet-stream")
		_, _ = resp.Write(v.ByteSlice())
		return
	} else {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPPool) RegisterNode(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consistent.AddNode(addr)
	var err error
	if addr == s.selfAddr {
		return
	}
	s.nodes[addr], err = rpcclient.NewClient(addr, 5*time.Second)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *HTTPPool) Get(groupName, key string) (eviction.ByteView, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if addr := s.consistent.ChooseNode(key); addr == s.selfAddr { //本地找
		g := group.GetGroup(groupName)
		if g == nil {
			return eviction.ByteView{}, fmt.Errorf("no such group")
		}
		if v, ok := g.Get(key); ok {
			return v, nil
		} else {
			return eviction.ByteView{}, fmt.Errorf("[Local]key not find by cache and local")
		}
	} else { //远程节点找
		c := s.nodes[addr]
		if v, ok := c.Get(groupName, key); ok {
			fmt.Println("get from", addr)
			return eviction.ByteView{B: v}, nil
		} else {
			return eviction.ByteView{}, fmt.Errorf("[Peer]key not find by cache and local")
		}
	}
}
func (s *HTTPPool) DeleteNode(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consistent.DeleteNode(addr)
	delete(s.nodes, addr)
}

// 启动服务
func (s *HTTPPool) Start() {
	//注册自己结点的Rpc服务，自己可以调用其他结点的Rpc，其他结点也可以调用我的
	grpcServer := grpc.NewServer()

	pb.RegisterCacheServiceServer(grpcServer, s)
	lis, err := net.Listen("tcp", s.selfAddr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("gRPC 服务已启动，监听 :", s.selfAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	//将自己的地址存入ETCD
	etcd.RegisterSelfAddr(s.selfAddr, s.selfAddr)
	addrs, version := etcd.DisCover()
	for _, addr := range addrs {
		s.RegisterNode(addr)
	}
	fmt.Println(addrs)
	//动态监听ETCD中的结点变化HTTP服务
	go s.WatchFromVersion(version)
	fmt.Println("http server listen at", HTTPADDR)
	//启动对客户端提供服务的HTTPf
	log.Fatalln(http.ListenAndServe(HTTPADDR, s))
}

func (s *HTTPPool) WatchFromVersion(version int64) {
	// 先获取当前所有节点，建立已存在节点集合
	currentNodes := make(map[string]struct{})
	if res, err := etcd.EtcdClient().Get(context.Background(), etcd.Prefix, clientv3.WithPrefix()); err == nil {
		for _, kv := range res.Kvs {
			currentNodes[string(kv.Value)] = struct{}{}
		}
	}

	watcher := etcd.EtcdClient().Watch(context.Background(),
		etcd.Prefix,
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
				s.RegisterNode(addr)
				currentNodes[addr] = struct{}{}

			case event.Type == clientv3.EventTypeDelete:

				a := string(event.PrevKv.Value)
				fmt.Println("[DEL]", a)
				delete(currentNodes, a)
				s.DeleteNode(a)

			}
		}
	}
}

func (s *HTTPPool) GetRpc(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	if v, err := s.Get(req.GetGroup(), req.GetKey()); err == nil {
		return &pb.GetResponse{
			Value: v.ByteSlice(),
			Found: true,
		}, nil
	} else {
		return &pb.GetResponse{
			Value: nil,
			Found: false,
		}, err
	}
}
