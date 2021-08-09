package state

import (
	"golua/number"
	"math"
)

type luaTable struct {
	luaArray []luaValue
	luaMap map[luaValue]luaValue
	metatable *luaTable
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}

	if nArr > 0 {
		t.luaArray = make([]luaValue, 0, nArr)
	}
	
	if nRec > 0 {
		t.luaMap = make(map[luaValue]luaValue, nRec)
	}

	return t
}

func (self *luaTable) hasMetaField(fieldName string) bool {
	return self.metatable != nil &&
		self.metatable.get(fieldName) != nil
}

func (t *luaTable) len() int {
	return len(t.luaArray)
}

func (t *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(t.luaArray)) {
			return t.luaArray[idx-1]
		}
	}
	return t.luaMap[key]
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (t *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(t.luaArray))
		if idx <= arrLen {
			t.luaArray[idx-1] = val
			if idx == arrLen && val == nil {
				t._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(t.luaMap, key)
			if val != nil {
				t.luaArray = append(t.luaArray, val)
				t._expandArray()
			}
			return
		}
	}
	if val != nil {
		if t.luaMap == nil {
			t.luaMap = make(map[luaValue]luaValue, 8)
		}
		t.luaMap[key] = val
	} else {
		delete(t.luaMap, key)
	}
}

func (t *luaTable) _shrinkArray() {
	for i := len(t.luaArray) - 1; i >= 0; i-- {
		if t.luaArray[i] == nil {
			t.luaArray = t.luaArray[0:i]
		} else {
			break
		}
	}
}

func (t *luaTable) _expandArray() {
	for idx := int64(len(t.luaArray)) + 1; true; idx++ {
		if val, found := t.luaMap[idx]; found {
			delete(t.luaMap, idx)
			t.luaArray = append(t.luaArray, val)
		} else {
			break
		}
	}
}