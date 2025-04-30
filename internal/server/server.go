package server

import (
	"GGCache/internal/group"
	"fmt"
	"net/http"
	"strings"
)

type HTTPServer struct {
	selfAddr string
}

func NewHTTPServer(addr string) *HTTPServer {
	return &HTTPServer{
		selfAddr: addr,
	}
}

func (s *HTTPServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	//http://127.0.0.1/group/key
	params := strings.Split(req.URL.Path, "/")
	params = params[1:]
	if len(params) != 2 {
		fmt.Println(params, len(params))
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}
	g := group.GetGroup(params[0])
	if g == nil {
		http.Error(resp, "no such group", http.StatusNotFound)
		return
	}
	key := params[1]
	if v, ok := g.Get(key); ok {
		resp.Header().Set("Content-Type", "application/octet-stream")
		_, _ = resp.Write(v.ByteSlice())
		return
	} else {
		http.Error(resp, "key not find by cache and local", http.StatusInternalServerError)
	}
}
