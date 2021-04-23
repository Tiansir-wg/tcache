package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

var db = map[string]string{
	"key1": "value1",
	"key2": "value2",
	"key3": "value3",
}

func TestHttpPool_ServeHTTP(t *testing.T) {
	NewGroup("scores", GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)

	addr := "localhost:9999"
	peers := NewHttpPool(addr)
	log.Println("tcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
