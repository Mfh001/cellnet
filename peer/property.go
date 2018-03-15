package peer

import (
	"reflect"
	"sync"
)

type ctx struct {
	key   interface{}
	value interface{}
}

type CoreContextSet struct {
	ctxes      []ctx
	ctxesGuard sync.RWMutex
}

func (self *CoreContextSet) GetContext(key, valuePtr interface{}) bool {

	pv, ok := self.RawGetContext(key)
	if !ok {
		return false
	}

	switch rawValue := valuePtr.(type) {
	case *string:
		*rawValue = pv.(string)
	case *int:
		*rawValue = pv.(int)
	case *int32:
		*rawValue = pv.(int32)
	case *int64:
		*rawValue = pv.(int64)
	case *uint:
		*rawValue = pv.(uint)
	case *uint32:
		*rawValue = pv.(uint32)
	case *uint64:
		*rawValue = pv.(uint64)
	case *bool:
		*rawValue = pv.(bool)
	case *float32:
		*rawValue = pv.(float32)
	case *float64:
		*rawValue = pv.(float64)
	case *[]byte:
		*rawValue = pv.([]byte)
	default:
		v := reflect.Indirect(reflect.ValueOf(valuePtr))

		v.Set(reflect.ValueOf(pv))
	}

	return true
}

func (self *CoreContextSet) RawGetContext(key interface{}) (interface{}, bool) {

	self.ctxesGuard.RLock()
	defer self.ctxesGuard.RUnlock()

	for _, t := range self.ctxes {
		if t.key == key {
			return t.value, true
		}
	}

	return nil, false
}

func (self *CoreContextSet) SetContext(key, v interface{}) {

	self.ctxesGuard.Lock()
	defer self.ctxesGuard.Unlock()

	for i, t := range self.ctxes {
		if t.key == key {
			self.ctxes[i] = ctx{key, v}
			return
		}
	}

	self.ctxes = append(self.ctxes, ctx{key, v})
}
