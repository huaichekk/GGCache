package main

import (
	"GGCache/internal/group"
	"GGCache/internal/server"
	"log"
	"os"
)

// 模拟查询数据库的耗时操作
var db = map[string]string{
	"Tom":                     "630",
	"Jack":                    "589",
	"Sam":                     "567",
	"Alice":                   "712", // 新增正常用户
	"Bob":                     "643",
	"Charlie":                 "598",
	"David":                   "null",   // 测试空值响应
	"Eve":                     "",       // 测试空字符串
	"Frank":                   "0",      // 测试零值
	"Grace":                   "999",    // 测试极值
	"黑客User":                  "xss",    // 测试特殊字符
	"User123":                 "456",    // 测试数字用户名
	"测试用户":                    "测试值",    // 测试中文
	"User.With.Dots":          "dotted", // 测试含点用户名
	"user@mail":               "email",  // 测试含@符号
	"longusernamex1234567890": "long",   // 测试长用户名
	"UPPER":                   "CASE",   // 测试大小写敏感
	"lower":                   "case",
	"Space User":              "space value", // 测试空格
}

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
