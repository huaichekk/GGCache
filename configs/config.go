package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	CacheConfig CacheConfig `json:"cache"`
	EtcdConfig  EtcdConfig  `json:"etcd"`
	RpcConfig   RpcConfig   `json:"rpc"`
}

type CacheConfig struct {
	Eviction string `json:"eviction"`
	ShardNum int    `json:"shardNum"`
	Replicas int    `json:"replicas"`
}
type EtcdConfig struct {
	Addr []string `json:"addr"`
}

type RpcConfig struct {
	Addr string
}

func LoadConfig() *Config {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	paths := []string{
		"./config.json",
		"../config.json",
		"../../config.json",
		"../../../config.json",
		"../../../../config.json",
		"../../../../../config.json",
		filepath.Dir(pwd) + "/../config.json",
	}

	var configPath string
	for i := range paths {
		if ok, _ := PathExists(paths[i]); ok {
			configPath = paths[i]
			break
		}
	}
	fmt.Printf("load config template path : %s \n", configPath)
	jsondata, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	// 解析JSON
	var res Config
	if err := json.Unmarshal(jsondata, &res); err != nil {
		panic(err)
	}
	fmt.Println(res)
	return &res
}
func GetConfig() *Config {
	if config == nil {
		once.Do(func() {
			config = LoadConfig()
		})
	}
	return config
}
func PathExists(path string) (bool, error) {
	if path == "" {
		return false, errors.New("路径为空,请检查")
	}
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
