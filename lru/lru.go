package lru

import (
	"container/list"
	"sync"
	"time"
)

type Value struct {
	key  int64
	data interface{} // 实际存储的值
	t    time.Time   // 存储时的时间戳
}

type Lru struct {
	m      map[int64]*list.Element // map
	list   *list.List              // 双向链表
	mutex  sync.Mutex              // 互斥锁
	max    int                     // 缓存数据的最大数量
	expire time.Duration           // 过期时间
}

// New 创建LRU缓存
// max 最大缓存数量
// expire 缓存过期时间
func New(max int, expire time.Duration) *Lru {
	return &Lru{
		m:      make(map[int64]*list.Element, max),
		list:   list.New(),
		max:    max,
		expire: expire,
	}
}

// Get 从缓存中获取数据，如果为nil,说明缓存不存在或者已经过期
func (l *Lru) Get(k int64) interface{} {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	e, ok := l.m[k]
	if !ok {
		return nil
	}

	v := e.Value.(Value)
	if v.t.Add(l.expire).Before(time.Now()) {
		return nil
	}

	l.list.MoveToFront(e)
	return v.data
}

// Set 设置数据
// 如果缓存中数据超出最大数量，就会删除一个最近最久未使用的数据
// 如果缓存已经存在数据，会删除掉在进行设置
func (l *Lru) Set(key int64, data interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	e, ok := l.m[key]
	if ok {
		delete(l.m, key)
		l.list.Remove(e)
	}

	value := Value{
		key:  key,
		data: data,
		t:    time.Now(),
	}
	e = l.list.PushFront(value)
	l.m[key] = e

	if l.list.Len() > l.max {
		value := l.list.Remove(l.list.Back()).(Value)
		delete(l.m, value.key)
	}
}
