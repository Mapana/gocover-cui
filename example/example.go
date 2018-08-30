package example

import "fmt"

func SUN(a ...int) (int, error) {
	if len(a) < 2 {
		return 0, fmt.Errorf("at least two numbers")
	}

	fmt.Printf("\033[0;33mcalculate sun(%v)\n", a)
	return sun(a...), nil
}

func sun(a ...int) int {
	var total = 0

	for _, v := range a {
		total += v
	}

	return total
}
