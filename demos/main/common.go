package main

import (
	"fmt"

	"github.ibm.com/modeling-analysis/lpsolve/pkg/config"
	"github.ibm.com/modeling-analysis/lpsolve/pkg/core"
	"github.ibm.com/modeling-analysis/lpsolve/pkg/utils"
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

// create problem instance
func CreateProblem(problemType config.ProblemType, isLimited bool) (core.Problem, error) {
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
		return nil, fmt.Errorf("unknown problem type: %s", problemType)
	}
	if err != nil {
		return nil, err
	}

	// set accelerator count limited option
	if isLimited {
		if err := p.SetLimited(numAcceleratorTypes, unitsAvail, acceleratorTypesMatrix); err != nil {
			return nil, err
		}
	} else {
		p.UnSetLimited()
	}

	return p, nil
}

// print solution details
func PrintResults(p core.Problem) {
	fmt.Printf("Solution type: %v\n", p.GetSolutionType())
	fmt.Printf("Solution time: %d msec\n", p.GetSolutionTimeMsec())
	fmt.Printf("Objective value: %v\n", p.GetObjectiveValue())

	numReplicas := p.GetNumReplicas()
	fmt.Println(utils.Pretty2D("numReplicas", numReplicas))

	instancesUsed := p.GetInstancesUsed()
	fmt.Println(utils.Pretty1D("instancesUsed", instancesUsed))

	if p.IsLimited() {
		fmt.Println(utils.Pretty1D("unitsAvail", unitsAvail))
		unitsUsed := p.GetUnitsUsed()
		fmt.Println(utils.Pretty1D("unitsUsed", unitsUsed))
	}
}
