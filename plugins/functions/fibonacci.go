package functions

func FibonacciRecursive(n int) int {
	if n < 0 {
		return -1
	} else if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
	}
}

func FibonacciIterative(n int) int {
	if n < 0 {
		return -1
	}
	prev1, prev2 := 0, 1
	for i := 2; i <= n; i++ {
		prev1, prev2 = prev2, prev1+prev2
	}
	return prev2
}
