package shared

import (
	"fmt"
	"sync"
	"time"
)

type MapConv struct {
	m  map[string]interface{}
	mu sync.RWMutex
}

func NewMapConv(m map[string]interface{}) *MapConv {
	return &MapConv{
		m:  m,
		mu: sync.RWMutex{},
	}
}

func (c *MapConv) String() string {
	return fmt.Sprintf("%v", c.m)
}

// Map return original map
func (c *MapConv) Map() map[string]interface{} {
	return c.m
}

func (c *MapConv) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[key]
	return v, ok
}

// MustGet return value of key, panic if not found
func (c *MapConv) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func (c *MapConv) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *MapConv) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *MapConv) GetFloat64(key string) (f float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f, _ = val.(float64)
	}
	return
}

func (c *MapConv) GetFloat32(key string) (f float32) {
	if val, ok := c.Get(key); ok && val != nil {
		f, _ = val.(float32)
	}
	return
}

func (c *MapConv) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *MapConv) GetInt64(key string) (i int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int64)
	}
	return
}

func (c *MapConv) GetInt32(key string) (i int32) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int32)
	}
	return
}

func (c *MapConv) GetUint(key string) (u uint) {
	if val, ok := c.Get(key); ok && val != nil {
		u, _ = val.(uint)
	}
	return
}

func (c *MapConv) GetUint32(key string) (u uint32) {
	return uint32(c.GetFloat64(key))
}

func (c *MapConv) GetUint64(key string) (u uint64) {
	if val, ok := c.Get(key); ok && val != nil {
		u, _ = val.(uint64)
	}
	return
}

func (c *MapConv) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *MapConv) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *MapConv) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *MapConv) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

func (c *MapConv) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *MapConv) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
