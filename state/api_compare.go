package state

import . "golua/api"

func (s *luaState) Compare(index1, index2 int, op CompareOp) bool {
	if !s.luaStack.isValid(index1) || !s.luaStack.isValid(index2) {
		return false
	}

	a := s.luaStack.get(index1)
	b := s.luaStack.get(index2)
	switch op {
	case LUA_OPEQ:
		return _eq(a, b, s)
	case LUA_OPLT:
		return _lt(a, b, s)
	case LUA_OPLE:
		return _le(a, b, s)
	default:
		panic("invalid compare op!")
	}
}

func _eq(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case *luaTable:
		if y, ok := b.(*luaTable); ok && x != y && ls != nil {
			if result, ok := callMetaMethod(x, y, "__eq", ls); ok {
				return convertToBoolean(result)
			}
		}
		return a == b
	default:
		return a == b
	}
}

func _lt(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}
	if result, ok := callMetaMethod(a, b, "__lt", ls); ok {
		return convertToBoolean(result)
	} else {
		panic("comparison error!")
	}
}

func _le(a, b luaValue, ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	if result, ok := callMetaMethod(a, b, "__le", ls); ok {
		return convertToBoolean(result)
	} else if result, ok := callMetaMethod(b, a, "__lt", ls); ok {
		return !convertToBoolean(result)
	} else {
		panic("comparison error!")
	}
}

func (s *luaState) RawEqual(idx1, idx2 int) bool {
	if !s.luaStack.isValid(idx1) || !s.luaStack.isValid(idx2) {
		return false
	}

	a := s.luaStack.get(idx1)
	b := s.luaStack.get(idx2)
	return _eq(a, b, nil)
}