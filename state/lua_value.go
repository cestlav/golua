package state

import (
	"fmt"
	. "golua/api"
	"golua/number"
)

type luaValue interface {}

func typeOf(luaValue luaValue) LuaType {
	switch luaValue.(type) {
	case nil: return LUA_TNIL
	case bool: return LUA_TBOOLEAN
	case int64: return LUA_TNUMBER
	case float64: return LUA_TNUMBER
	case string: return LUA_TSTRING
	case *luaTable: return LUA_TTABLE
	case *closure: return LUA_TFUNCTION

	default:
		panic("invalid type")
	}
}

func convertToBoolean(v luaValue) bool {
	switch x := v.(type) {
	case nil: return false
	case bool: return x
	default:
		return true
	}
}

func convertToFloat(v luaValue) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

func convertToInteger(v luaValue) (int64, bool) {
	switch x := v.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		return _stringToInteger(x)
	default:
		return 0, false
	}
}

func _stringToInteger(s string) (int64, bool) {
	if i, ok := number.ParseInteger(s); ok {
		return i, ok
	}
	if f, ok := number.ParseFloat(s); ok {
		return number.FloatToInteger(f)
	}

	return 0, false
}

func setMetatable(v luaValue, mt *luaTable, ls *luaState)  {
	if t, ok := v.(*luaTable); ok {
		t.metatable = mt
		return
	}
	key := fmt.Sprintf("_MT%d", typeOf(v))
	ls.registry.put(key, mt)
}

func getMetatable(v luaValue, ls *luaState) *luaTable {
	if t, ok := v.(*luaTable); ok {
		return t.metatable
	}

	key := fmt.Sprintf("_MT%d", typeOf(v))
	if mt := ls.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}
	return nil
}

func callMetaMethod(a, b luaValue, mmName string, ls *luaState) (luaValue, bool)  {
	var mm luaValue
	if mm = getMetaField(a, mmName, ls); mm == nil {
		if mm = getMetaField(b, mmName, ls); mm == nil {
			return nil,false
		}
	}
	ls.luaStack.check(4)
	ls.luaStack.push(mm)
	ls.luaStack.push(a)
	ls.luaStack.push(b)
	ls.Call(2, 1)
	return ls.luaStack.pop(), true
}

func getMetaField(v luaValue, fieldName string, ls *luaState) luaValue {
	if mt := getMetatable(v, ls); mt != nil {
		return mt.get(fieldName)
	}
	return nil
}