package configs

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig()
	fmt.Println(config)
}

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	fmt.Println(config.CacheConfig.Eviction, config.CacheConfig.ShardNum)
}
