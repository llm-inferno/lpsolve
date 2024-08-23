package main

import (
	"fmt"

	"github.ibm.com/tantawi/lpsolve/pkg/core"
	"github.ibm.com/tantawi/lpsolve/pkg/utils"
)

var numServers int
var numAccelerators int
var unitCost []float64         // [numAccelerators]
var numUnitsPerReplica [][]int // [numServers][numAccelerators]
var ratePerReplica [][]float64 // [numServers][numAccelerators]
var arrivalRates []float64     // [numServers]

var numAcceleratorTypes int
var unitsAvailByType []int         // [numAcceleratorTypes]
var acceleratorTypesMatrix [][]int // [numAcceleratorTypes][numAccelerators]

// create problem instance, solve it, and print results
func Optimize(problemType core.ProblemType, isLimited bool) {
	var p core.Problem
	var err error
	// create a new problem instance
	switch problemType {
	case core.SINGLE:
		p, err = core.CreateSingleAssignProblem(numServers, numAccelerators, unitCost, numUnitsPerReplica,
			ratePerReplica, arrivalRates)
	case core.MULTI:
		p, err = core.CreateMultiAssignProblem(numServers, numAccelerators, unitCost, numUnitsPerReplica,
			ratePerReplica, arrivalRates)
	default:
		fmt.Printf("Unknown problem type: %s", problemType)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// set accelerator count limited option
	if isLimited {
		if err := p.SetLimited(numAcceleratorTypes, unitsAvailByType, acceleratorTypesMatrix); err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		p.UnSetLimited()
	}

	// solve the problem
	if err = p.Solve(); err != nil {
		fmt.Println(err.Error())
		return
	}

	// print solution details
	fmt.Printf("Solution type: %v\n", p.GetSolutionType())
	fmt.Printf("Solution time: %d msec\n", p.GetSolutionTimeMsec())
	fmt.Printf("Objective value: %v\n", p.GetObjectiveValue())

	numReplicas := p.GetNumReplicas()
	fmt.Println(utils.Pretty2DInt("numReplicas", numReplicas))

	unitsUsed := p.GetUnitsUsed()
	fmt.Println(utils.Pretty1DInt("unitsUsed", unitsUsed))

	if isLimited {
		fmt.Println(utils.Pretty1DInt("unitsAvailByType", unitsAvailByType))
		unitsUsedByType := p.GetUnitsUsedByType()
		fmt.Println(utils.Pretty1DInt("unitsUsedByType", unitsUsedByType))
	}
}
