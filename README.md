# 分布式K/V缓存系统
## 项目背景/能够解决的问题
- 单机缓存受制于内存大小的限制，无法承载过大的数据量
- 使用分布式节点将缓存分布在多个节点的内存中，横向扩展缓存的容量

## 架构图
![img.png](images/img.png)

## 安装方式
```shell
git clone https://github.com/huaichekk/GGCache.git
cd ./GGCache
go run .
```

## 使用流程
1. 新建一个Group,一个Group代表一类数据，并传入回调函数，回调函数一般为从数据库中获取，一般在缓存中无数据时，调用回调函数，查找数据。
2. 新建一个HTTPPool，并启动，需要传入需要向外界暴露服务的地址，节点之间使用rpc通信，本地暴露rpc服务的地址在GGCache/config.json中读取.


```go 使用示例
func main() {
	_ = group.NewGroup("scores", 2<<10, func(key string) ([]byte, bool) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), true
		}
		return nil, false
	})
	s := server.NewHTTPPool(os.Args[1])
	s.Start()
}
```


## 节点初始化流程
1. 注册本地的rpc服务，每个节点即充当rpc的服务端，向外界暴露读取缓存的服务，也作为客户端，在使用一致性哈希算法匹配到非本地节点时，调用对应节点的rpc服务
2. 将自己的的rpc地址存入etcd中
3. 使用get --prefix命令获取etcd中的集群其他节点地址
4. 启动一个协程，使用watch命令监听集群变化，若add/delete，则在本地的哈希环中增删对应地址
5. 启动用户初始化HTTPPool时传入的地址的http服务

## 读操作流程
示例URL:http://example.com/GroupName/key
![img_1.png](images/img_1.png)
1. HTTPPoll使用内部的consistent一致性哈希选择改key对应的节点
2. 若节点时本地ip地址时，直接读取本地缓存，若本地缓存中该key，则从回调函数中获取
3. 若匹配到的时远程节点，则调用改节点的rpc服务


## 后续的优化点：
- [ ] 增加lfu,lru-k,fifo等多种算法
- [x] 增加ttl缓存过期机制（惰性删除，定时删除）
- [x] 缓存雪崩
- [x] 缓存穿透
- [x] 缓存击穿

