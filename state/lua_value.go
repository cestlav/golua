package state

import (
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