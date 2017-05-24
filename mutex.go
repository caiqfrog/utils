package utils

import (
	"sync/atomic"
	"time"
)

// chan实现的trylock
type Mutex struct {
	c chan int32
	s int32
}

// 创建
func NewMutex() *Mutex {
	return &Mutex{
		c: make(chan int32, 1),
		s: int32(0),
	}
}

// 获取每次锁的序号
func (obj *Mutex) serial() int32 {
	prev, next := int32(0), int32(0)
	for {
		prev = atomic.LoadInt32(&obj.s)
		next = prev + 1
		if atomic.CompareAndSwapInt32(&obj.s, prev, next) {
			break
		}
	}

	return next
}

// 锁定
func (obj *Mutex) Lock() {
	s := obj.serial()

	// log.Println(`lock|serial: `, s)

	obj.c <- s
}

// 解锁
func (obj *Mutex) Unlock() {
	// s :=
	<-obj.c

	// log.Println(`unlock|id: `, s)
}

// 尝试锁
// 返回true则锁定
func (obj *Mutex) Trylock(timeout time.Duration) bool {
	id := obj.serial()

	select {
	case obj.c <- id:
		// log.Println(`trylock|serial|true: `, id)
		return true
	case <-time.After(timeout):
		// log.Println(`trylock|serial|false: `, id)
		return false
	}
}
