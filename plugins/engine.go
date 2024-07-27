package plugins

import (
	"github.com/yuin/gopher-lua"
)

func Lua() {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(`print("hello")`); err != nil {
		panic(err)
	}
}
