package main

import (
	"errors"
	"fmt"
	"sync"
)

// 定义通用接口从不同的数据源获取数据, 只需要实现该接口即可
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (gf GetterFunc) Get(key string) ([]byte, error) {
	return gf(key)
}

// 缓存命名空间
type Group struct {
	// 命名空间名称
	name string
	// 命名空间数据来源
	getter Getter
	// 主缓存对象
	mainCache cache
}

var (
	mu           sync.RWMutex
	groups       = make(map[string]*Group)
	ErrNilGetter = errors.New("nil getter") // 缺少getter方法异常
)

// 创建新的Group实例
func NewGroup(name string, getter Getter, limitedBytes int64) *Group {
	if getter == nil {
		panic(ErrNilGetter)
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			limitedBytes: limitedBytes,
		},
	}
	groups[name] = g
	return g
}

//  获取指定的Group
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.Unlock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("未指定查找的KEY")
	}
	// 从主缓存查找
	v, ok := g.mainCache.get(key)
	if ok {
		return v, nil
	}
	// 主缓存中不存在
	return g.loadCache(key)
}

func (g *Group) loadCache(key string) (ByteView, error) {
	// 从自定义数据源获取
	b, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{
		b: cloneBytes(b),
	}
	g.mainCache.put(key, value)
	return value, nil
}
