package stream

import (
	"strings"
	"sync"
	"sync/atomic"
)

/**
 *	流量统计相关数据操作必须通过本模块提供的方法进行
 */
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
	Data map[string]*count `json:"Data"`
	lock sync.RWMutex
}

func init() {
	Counter.Data = make(map[string]*count)
}

// FlowInIncr 入流量增加
func (c *counter) FlowInIncr(key string, value uint64) {
	key = strings.Split(key, ":")[0]
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.Data[key]; ok {
		atomic.AddUint64(&c.Data[key].FlowIn, value)
	} else {
		c.Data[key] = &count{
			FlowIn:  0,
			FlowOut: 0,
		}
	}
}

// FlowOutIncr 出流量增加
func (c *counter) FlowOutIncr(key string, value uint64) {
	key = strings.Split(key, ":")[0]
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.Data[key]; ok {
		atomic.AddUint64(&c.Data[key].FlowOut, value)
	} else {
		c.Data[key] = &count{
			FlowIn:  0,
			FlowOut: 0,
		}
	}
}

// Get 流量获取
func (c *counter) Get(key string) uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return 0
}
func (c *counter) GetAll() map[string]count {
	c.lock.RLock()
	defer c.lock.RUnlock()
	m := make(map[string]count)
	for k, v := range c.Data {
		m[k] = *v
	}
	return m
}
