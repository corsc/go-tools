package examples

import (
	"fmt"
	"math/rand"
)

// Example1 ...
func Example1() {
	fmt.Printf("Roll: %d", rand.Intn(6))
	fmt.Printf("Roll: %d", rand.Intn(10))
	fmt.Printf("Roll: %d", rand.Intn(12))
}
