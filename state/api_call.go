package state

import (
	"golua/api"
	"golua/binary"
	"golua/vm"
)

func (s *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binary.Undump(chunk)
	c := newLuaClosure(proto)
	s.luaStack.push(c)

	if len(proto.UpValues) > 0 {
		env := s.registry.get(api.LUA_RIDX_GLOBALS)
		c.upValues[0] = &upValue{&env}
	}
	return 0
}

func (s *luaState) Call(nArgs, nResults int)  {
	f := s.luaStack.get(-(nArgs + 1))
	c, ok := f.(*closure)

	if !ok {
		if mf := getMetaField(f, "__call", s); mf != nil {
			if c, ok = mf.(*closure); ok {
				s.luaStack.push(f)
				s.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			s.callLuaClosure(nArgs, nResults, c)
		} else {
			s.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("error call instruction")
	}
}

func (s *luaState) callLuaClosure(nArgs, nResults int, f *closure)  {
	nRegs := int(f.proto.MaxStackSize)
	nParams := int(f.proto.NumParams)
	isVararg := f.proto.IsVararg == 1

	newStack := newLuaStack(nRegs + 20, s)
	newStack.closure = f

	funcAndArgs := s.luaStack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams + 1:]
	}

	s.pushLuaStack(newStack)
	s.runLuaColosure()
	s.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		s.luaStack.check(len(results))
		s.luaStack.pushN(results, nResults)
	}
}

func (s *luaState) runLuaColosure()  {
	for {
		inst := vm.Instruction(s.Fetch())
		inst.Execute(s)
		if inst.OpCode() == vm.OP_RETURN {
			break
		}
	}
}

func (s *luaState) callGoClosure(nArgs, nResults int, c *closure)  {
	newStack := newLuaStack(nArgs + 20, s)
	newStack.closure = c
	args := s.luaStack.popN(nArgs)
	newStack.pushN(args, nArgs)
	s.luaStack.pop()
	s.pushLuaStack(newStack)
	r := c.goFunc(s)
	s.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(r)
		s.luaStack.check(len(results))
		s.luaStack.pushN(results, nResults)
	}
}

func (s *luaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := s.luaStack
	status = api.LUA_ERRRUN
	defer func() {
		if err := recover(); err != nil {
			for s.luaStack != caller {
				s.popLuaStack()
			}
			s.luaStack.push(err)
		}
	}()

	s.Call(nArgs, nResults)
	status = api.LUA_OK
	return
}