package main

import (
	"fmt"
	"golua/state"
	"io/ioutil"
)

func main()  {
	//fmt.Println("args", os.Args)
	data, err := ioutil.ReadFile("luac.out")
	if err != nil {
		panic("err binary")
	}
	ls := state.NewLuaState(20)
	ls.Load(data, "luac.out", "b")
	fmt.Printf("%v\n", )
	//ls.Call(0, 0)
	fmt.Println("hello world")
}