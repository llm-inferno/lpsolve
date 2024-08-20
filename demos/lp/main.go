package main

import (
	"fmt"

	"github.com/draffensperger/golp"
)

func main() {
	lp := golp.NewLP(0, 2)
	lp.AddConstraint([]float64{110.0, 30.0}, golp.LE, 4000.0)
	lp.AddConstraint([]float64{1.0, 1.0}, golp.LE, 75.0)
	lp.SetObjFn([]float64{143.0, 60.0})
	lp.SetMaximize()

	lp.Solve()
	vars := lp.Variables()
	fmt.Printf("Plant %.3f acres of barley\n", vars[0])
	fmt.Printf("And  %.3f acres of wheat\n", vars[1])
	fmt.Printf("For optimal profit of $%.2f\n", lp.Objective())

	// No need to explicitly free underlying C structure as golp.LP finalizer will
}
