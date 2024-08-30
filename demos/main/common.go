package main

import (
	"fmt"

	"github.ibm.com/tantawi/lpsolve/pkg/config"
	"github.ibm.com/tantawi/lpsolve/pkg/core"
	"github.ibm.com/tantawi/lpsolve/pkg/utils"
)

var numServers int
var numAccelerators int
var instanceCost []float64         // [numAccelerators]
var numInstancesPerReplica [][]int // [numServers][numAccelerators]
var ratePerReplica [][]float64     // [numServers][numAccelerators]
var arrivalRates []float64         // [numServers]

var numAcceleratorTypes int
var unitsAvail []int               // [numAcceleratorTypes]
var acceleratorTypesMatrix [][]int // [numAcceleratorTypes][numAccelerators]

// create problem instance, solve it, and print results
func Optimize(problemType config.ProblemType, isLimited bool) {
	var p core.Problem
	var err error
	// create a new problem instance
	switch problemType {
	case config.SINGLE:
		p, err = core.CreateSingleAssignProblem(numServers, numAccelerators, instanceCost, numInstancesPerReplica,
			ratePerReplica, arrivalRates)
	case config.MULTI:
		p, err = core.CreateMultiAssignProblem(numServers, numAccelerators, instanceCost, numInstancesPerReplica,
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
		if err := p.SetLimited(numAcceleratorTypes, unitsAvail, acceleratorTypesMatrix); err != nil {
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

	instancesUsed := p.GetInstancesUsed()
	fmt.Println(utils.Pretty1DInt("instancesUsed", instancesUsed))

	if isLimited {
		fmt.Println(utils.Pretty1DInt("unitsAvail", unitsAvail))
		unitsUsed := p.GetUnitsUsed()
		fmt.Println(utils.Pretty1DInt("unitsUsed", unitsUsed))
	}
}
