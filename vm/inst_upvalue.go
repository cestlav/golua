package vm

import . "golua/api"

func getTabUp(i Instruction, vm LuaVM)  {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.GetRK(c)
	vm.GetTable(LuaUpValueIndex(b))
	vm.Replace(a)
}

func setTabUp(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpValueIndex(a))
}

func getUpValue(i Instruction, vm LuaVM)  {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.Copy(LuaUpValueIndex(b), a)
}

func setUpValue(i Instruction, vm LuaVM)  {
	a, b, _ := i.ABC()
	a += 1
	b += 1

	vm.Copy(a, LuaUpValueIndex(b))
}