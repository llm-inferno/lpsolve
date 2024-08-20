package main

import (
	"fmt"

	"github.ibm.com/tantawi/lpsolve/pkg/core"
	"github.ibm.com/tantawi/lpsolve/pkg/utils"
)

var numServers int
var numAccelerators int
var unitsAvail []int
var unitCost []float64
var numUnitsPerReplica [][]int
var ratePerReplica [][]float64
var arrivalRates []float64

func main() {
	numServers = 5
	numAccelerators = 8

	// available number of accelerators (numAccelerators)
	unitsAvail = []int{512, 256, 192, 128, 98, 64, 48, 32}

	// unit cost of accelerators (numAccelerators)
	unitCost = []float64{0.5, 1.0, 1.2, 2.3, 2.7, 5.6, 7.0, 10.0}

	// number of accelerator units for pairs of server and accelerator
	//	(numServers x numAccelerators)
	numUnitsPerReplica = [][]int{
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

	// unlimited case
	fmt.Println("Solution of Unlimited case:")
	fmt.Println("---------------------------")
	optimize(false)
	fmt.Println()

	//limited case
	fmt.Println("Solution of Limited case:")
	fmt.Println("-------------------------")
	optimize(true)
	fmt.Println()
}

func optimize(isLimited bool) {
	// create a new MIP problem instance
	mip, err := core.CreateAssignmentProblemInstance(numServers, numAccelerators, unitCost, numUnitsPerReplica,
		ratePerReplica, arrivalRates)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// set acccelerator count limited option
	if isLimited {
		mip.SetLimited(unitsAvail)
		fmt.Println(utils.Pretty1DInt("unitsAvail", unitsAvail))
	} else {
		mip.UnSetLimited()
	}

	// solve the problem
	if err = mip.Solve(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// print solution details
	fmt.Printf("Solution type: %v\n", mip.GetSolutionType())
	fmt.Printf("Solution time: %d msec\n", mip.GetSolutionTimeMsec())
	fmt.Printf("Objective value: %v\n", mip.GetObjectiveValue())

	numReplicas := mip.GetNumReplicas()
	fmt.Println(utils.Pretty2DInt("numReplicas", numReplicas))

	unitsUsed := mip.GetUnitsUsed()
	fmt.Println(utils.Pretty1DInt("unitsUsed", unitsUsed))
}
