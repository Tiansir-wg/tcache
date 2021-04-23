package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

func TestGet(t *testing.T) {
	loadCnt := make(map[string]int, len(db))
	g := NewGroup("test", GetterFunc(func(key string) ([]byte, error) {
		log.Printf("指定数据源查找键%s\n", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCnt[key]; !ok {
				loadCnt[key] = 0
			}
			loadCnt[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exists", key)
	}), 2<<10)

	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		}
		if view, err := g.Get(k); err != nil || view.String() != v || loadCnt[k] > 1 {
			t.Fatalf("cache of %s miss", k)
		}
	}
	if view, err := g.Get("key4"); err == nil {
		t.Fatalf("the value of %s should be empty,but got %s", "key4", view)
	}
}
