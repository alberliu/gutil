package lru

import (
	"container/list"
	"fmt"
	"testing"
	"time"
)

type User struct {
	Id   int64
	Name int64
}

// TestLru 缓存失效测试
func TestLru(t *testing.T) {
	lru := New(5, 5*time.Second)
	user := &User{
		Id:   1,
		Name: 1,
	}
	fmt.Println("get 1:", lru.Get(1))
	lru.Set(1, user)
	fmt.Println("get 1:", lru.Get(1))
	time.Sleep(5 * time.Second)
	fmt.Println("get 1:", lru.Get(1))
}

// TestLru 获取缓存设置
func TestLruGet(t *testing.T) {
	lru := New(5, 5*time.Second)

	var i int64
	for i = 0; i < 5; i++ {
		user := &User{
			Id:   i,
			Name: i,
		}
		lru.Set(i, user)
		printfList(lru.list)
	}

	for i = 0; i < 5; i++ {
		fmt.Println(lru.Get(i))
		printfList(lru.list)
	}
}

// TestLruSet 设置缓存设置
func TestLruSet(t *testing.T) {
	lru := New(5, 5*time.Second)

	var i int64
	for i = 0; i < 10; i++ {
		user := &User{
			Id:   i,
			Name: i,
		}
		lru.Set(i, user)
		printfList(lru.list)
	}

}

// TestLruSet 设置缓存设置
func TestLruReset(t *testing.T) {
	lru := New(5, 5*time.Second)

	var i int64
	for i = 0; i < 5; i++ {
		user := &User{
			Id:   i,
			Name: i,
		}
		lru.Set(i, user)
		printfList(lru.list)
	}

	user := &User{
		Id:   2,
		Name: 2,
	}
	lru.Set(2, user)
	printfList(lru.list)
}

func printfList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value.(Value).key, "  ")
	}
	println()
}
func printfMap(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e)
	}
}
