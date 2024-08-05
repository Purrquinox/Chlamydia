package plugins

import (
	"github.com/yuin/gopher-lua"
	"log"
)

func disableLibraries(L *lua.LState) {
	L.PreloadModule("os", func(L *lua.LState) int {
		L.Push(L.NewTable())
		return 1
	})
	L.PreloadModule("debug", func(L *lua.LState) int {
		L.Push(L.NewTable())
		return 1
	})
}

func restrictedPrint(L *lua.LState) int {
	msg := L.ToString(1)
	if len(msg) > 0 {
		log.Println(msg)
	}
	return 0
}

func loadSafeLibraries(L *lua.LState) {
	L.SetGlobal("print", L.NewFunction(restrictedPrint))
	L.PreloadModule("math", lua.OpenMath)
	L.PreloadModule("string", lua.OpenString)
}

func createSandbox(L *lua.LState) {
	sandboxEnv := L.NewTable() // Create sandbox environment

	// Create variables for sandbox
	L.SetField(sandboxEnv, "print", L.NewFunction(restrictedPrint))
	L.SetField(sandboxEnv, "math", L.GetGlobal("math"))
	L.SetField(sandboxEnv, "string", L.GetGlobal("string"))

	// Push variables to sandbox
	L.Push(sandboxEnv)

	// Set the sandbox environment as the new global environment
	L.SetGlobal("_G", sandboxEnv)
}

func LuaTest() {
	L := lua.NewState()
	defer L.Close()

	disableLibraries(L)
	loadSafeLibraries(L)
	createSandbox(L)

	luaScript := `
		local os = require("os")  -- Should fail or be empty
		print("Trying to access os: ", os)
		print("Fibonacci (Recursive) of 10: ", fibonacci_recursive(10))
		print("Fibonacci (Iterative) of 10: ", fibonacci_iterative(10))
	`

	// Adding Fibonacci functions to the Lua state
	L.SetGlobal("fibonacci_recursive", L.NewFunction(func(L *lua.LState) int {
		n := L.ToInt(1)
		L.Push(lua.LNumber(fibonacciRecursive(n)))
		return 1
	}))
	L.SetGlobal("fibonacci_iterative", L.NewFunction(func(L *lua.LState) int {
		n := L.ToInt(1)
		L.Push(lua.LNumber(fibonacciIterative(n)))
		return 1
	}))

	if err := L.DoString(luaScript); err != nil {
		log.Fatal(err)
	}
}

func fibonacciRecursive(n int) int {
	if n < 0 {
		return -1
	} else if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacciRecursive(n-1) + fibonacciRecursive(n-2)
	}
}

func fibonacciIterative(n int) int {
	if n < 0 {
		return -1
	}
	prev1, prev2 := 0, 1
	for i := 2; i <= n; i++ {
		prev1, prev2 = prev2, prev1+prev2
	}
	return prev2
}
