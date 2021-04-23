package lru

import "container/list"

// 缓存值类型
type Value interface {
	// 获取占用的字节数
	Len() int
}

// 注意在缓存中最近访问的元素在链表头,待清除的元素在链表尾
type Cache struct {
	// 缓存最大能使用的内存, 为0表示无限制
	limitedBytes int64
	// 已经使用的内存
	usedBytes int64
	// 存放缓存元素的双向链表
	ll *list.List
	// 哈希表
	cache map[string]*list.Element
	// 移除元素时的回调函数,可以为空
	OnEvicted func(key string, value Value)
}

// 双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

func New(limitedBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		limitedBytes: limitedBytes,
		OnEvicted:    onEvicted,
		ll:           list.New(),
		cache:        make(map[string]*list.Element),
	}
}

// 查询
func (c *Cache) Get(key string) (value Value, ok bool) {
	ele, ok := c.cache[key]
	// 缓存中存在该值
	if ok {
		// 移动到链表头
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 添加或修改
// 添加时需要添加到链表头，并在哈希表中指向对应链表的元素，然后更新已使用的字节数
// 修改时需要将元素移动到链表头，然后修改对应节点的值，最后更新使用的字节数
func (c *Cache) Put(key string, value Value) {
	// 添加
	if ele, ok := c.cache[key]; !ok {
		ele = c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.usedBytes += int64(len(key)) + int64(value.Len())
	} else {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	}

	// 如果限制了缓存能使用的最大内存并且当前内存已经超过了设置的最大内存则进行清理
	for c.limitedBytes != 0 && c.usedBytes > c.limitedBytes {
		c.RemoveOldest()
	}
}

// 清理链表头部的缓存节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		c.ll.Remove(ele)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int64 {
	return int64(c.ll.Len())
}

//
