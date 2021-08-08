package main

import (
	"fmt"
	"golua/api"
	"golua/state"
	"io/ioutil"
)

func main()  {
	//fmt.Println("args", os.Args)
	data, err := ioutil.ReadFile("luac.out")
	if err != nil {
		panic("err binary")
	}
	ls := state.NewLuaState()
	ls.Register("print", print)
	ls.Load(data, "luac.out", "b")

	ls.Call(0, 0)
	fmt.Println("hello world")
}

func print(ls api.LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i < nArgs; i++ {
		fmt.Printf("%v\n", ls.ToString(i))
	}
	fmt.Println("fuck world")
	return 0
}