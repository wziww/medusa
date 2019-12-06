package stream

import (
	"sync"
)

var (
	// Counter 总流量统计
	Counter *counter = &counter{}
	// FlowIn 总流量 (入)
	FlowIn *uint64 = new(uint64)
	// FlowOut 总流量 （出）
	FlowOut *uint64 = new(uint64)
)

type count struct {
	FlowIn  uint64
	FlowOut uint64
}
type counter struct {
	Data map[string]*count
	lock sync.RWMutex
}

func init() {
}

// Set 流量设置
func (c *counter) Set(key string, value uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
}

// Get 流量获取
func (c *counter) Get(key string) uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return 0
}
