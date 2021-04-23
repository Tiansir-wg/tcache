package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpPool struct {
	addr    string // 当前节点的地址
	baseUrl string // 请求路径的前缀
}

const defaultBaseUrl = "/_tcache/"

func NewHttpPool(add string) *HttpPool {
	return &HttpPool{
		addr:    add,
		baseUrl: defaultBaseUrl,
	}
}

func (hp *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[server %s] %s", hp.addr, fmt.Sprintf(format, v...))
}

func (hp *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, hp.baseUrl) {
		http.Error(w, "错误的请求路径:"+path, http.StatusBadRequest)
		return
	}
	hp.Log("%s %s", r.Method, path)
	parts := strings.SplitN(path[len(hp.baseUrl):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "错误的请求路径", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "请求查询的组不存在", http.StatusBadRequest)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, "内部错误", http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
