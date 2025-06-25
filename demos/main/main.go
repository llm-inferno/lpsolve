package main

import (
	"fmt"
	"os"

	"github.com/llm-inferno/lpsolve/pkg/config"
)

func main() {
	// get problem type argument, or use default
	problemType := config.SINGLE
	if len(os.Args) > 1 {
		problemType = config.GetProblemType(os.Args[1])
	}

	numServers = 5
	numAccelerators = 8
	numAcceleratorTypes = 8

	// available number of units of accelerator types (numAcceleratorTypes)
	unitsAvail = []int{512, 256, 192, 128, 98, 64, 48, 32}
	acceleratorTypesMatrix = make([][]int, numAcceleratorTypes)
	for i := 0; i < numAcceleratorTypes; i++ {
		acceleratorTypesMatrix[i] = make([]int, numAccelerators)
		acceleratorTypesMatrix[i][i] = 1 // one type per accelerator
	}

	// instance cost of accelerators (numAccelerators)
	instanceCost = []float64{0.5, 1.0, 1.2, 2.3, 2.7, 5.6, 7.0, 10.0}

	// number of accelerator instances for pairs of server and accelerator
	//	(numServers x numAccelerators)
	numInstancesPerReplica = [][]int{
		{3, 2, 2, 2, 1, 1, 1, 1},
		{4, 3, 3, 2, 2, 1, 1, 1},
		{5, 4, 3, 2, 2, 2, 1, 1},
		{5, 4, 3, 3, 2, 2, 2, 2},
		{6, 5, 4, 4, 3, 3, 2, 2},
	}

	// max arrival rate for pairs of server and accelerator
	//	(numServers x numAccelerators)
	ratePerReplica = [][]float64{
		{0.1, 0.2, 0.4, 0.6, 0.9, 1.4, 2.0, 3.2},
		{0.1, 0.2, 0.4, 0.6, 0.9, 1.4, 2.0, 3.2},
		{0.1, 0.2, 0.4, 0.6, 0.9, 1.4, 2.0, 3.2},
		{0.1, 0.2, 0.4, 0.6, 0.9, 1.4, 2.0, 3.2},
		{0.1, 0.2, 0.4, 0.6, 0.9, 1.4, 2.0, 3.2},
	}

	// arrival rates to servers
	arrivalRates = []float64{10, 20, 30, 40, 50}

	fmt.Printf("Problem type: %v\n", problemType)
	fmt.Println()

	// unlimited case
	fmt.Println("Solution of Unlimited case:")
	fmt.Println("---------------------------")
	if p, err := CreateProblem(problemType, false); err != nil || p.Solve() != nil {
		fmt.Println(err)
		return
	} else {
		PrintResults(p)
	}
	fmt.Println()

	//limited case
	fmt.Println("Solution of Limited case:")
	fmt.Println("-------------------------")
	if p, err := CreateProblem(problemType, true); err != nil || p.Solve() != nil {
		fmt.Println(err)
		return
	} else {
		PrintResults(p)
	}
	fmt.Println()
}
