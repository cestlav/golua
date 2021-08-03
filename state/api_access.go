package state

import "fmt"
import . "golua/api"

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_typename
func (s *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_type
func (s *luaState) Type(idx int) LuaType {
	if s.luaStack.isValid(idx) {
		val := s.luaStack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnone
func (s *luaState) IsNone(idx int) bool {
	return s.Type(idx) == LUA_TNONE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnil
func (s *luaState) IsNil(idx int) bool {
	return s.Type(idx) == LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnoneornil
func (s *luaState) IsNoneOrNil(idx int) bool {
	return s.Type(idx) <= LUA_TNIL
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isboolean
func (s *luaState) IsBoolean(idx int) bool {
	return s.Type(idx) == LUA_TBOOLEAN
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_istable
func (s *luaState) IsTable(idx int) bool {
	return s.Type(idx) == LUA_TTABLE
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isfunction
func (s *luaState) IsFunction(idx int) bool {
	return s.Type(idx) == LUA_TFUNCTION
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isthread
func (s *luaState) IsThread(idx int) bool {
	return s.Type(idx) == LUA_TTHREAD
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isstring
func (s *luaState) IsString(idx int) bool {
	t := s.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isnumber
func (s *luaState) IsNumber(idx int) bool {
	_, ok := s.ToNumberX(idx)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isinteger
func (s *luaState) IsInteger(idx int) bool {
	val := s.luaStack.get(idx)
	_, ok := val.(int64)
	return ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_toboolean
func (s *luaState) ToBoolean(idx int) bool {
	val := s.luaStack.get(idx)
	return convertToBoolean(val)
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointeger
func (s *luaState) ToInteger(idx int) int64 {
	i, _ := s.ToIntegerX(idx)
	return i
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tointegerx
func (s *luaState) ToIntegerX(idx int) (int64, bool) {
	val := s.luaStack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumber
func (s *luaState) ToNumber(idx int) float64 {
	n, _ := s.ToNumberX(idx)
	return n
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_tonumberx
func (s *luaState) ToNumberX(idx int) (float64, bool) {
	val := s.luaStack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

// [-0, +0, m]
// http://www.lua.org/manual/5.3/manual.html#lua_tostring
func (s *luaState) ToString(idx int) string {
	str, _ := s.ToStringX(idx)
	return str
}

func (s *luaState) ToStringX(idx int) (string, bool) {
	val := s.luaStack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		str := fmt.Sprintf("%v", x) // todo
		s.luaStack.set(idx, s)
		return str, true
	default:
		return "", false
	}
}
