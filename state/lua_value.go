package state

import . "golua/api"

type LuaValue interface {}

func typeOf(luaValue LuaValue) LuaType {
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

func convertToBoolean(v LuaValue) bool {
	switch x := v.(type) {
	case nil: return false
	case bool: return v
	default:
		return true
	}
}