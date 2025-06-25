package main

import (
	"fmt"

	"github.com/llm-inferno/lpsolve/pkg/config"
	"github.com/llm-inferno/lpsolve/pkg/core"
	"github.com/llm-inferno/lpsolve/pkg/utils"
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
	case config.SINGLE, config.MULTI:
		p, err = core.CreateCplexProblem(numServers, numAccelerators, instanceCost, numInstancesPerReplica,
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
		switch problemType {
		case config.SINGLE:
			SetFileNames(p, "single-limited")
		case config.MULTI:
			SetFileNames(p, "multi-limited")
		}
	} else {
		p.UnSetLimited()
		switch problemType {
		case config.SINGLE:
			SetFileNames(p, "single-unlimited")
		case config.MULTI:
			SetFileNames(p, "multi-unlimited")
		}
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

func SetFileNames(p core.Problem, name string) {
	pc := p.(*core.CplexProblem)
	pc.SetModelFileName(name + ".mod")
	pc.SetDataFileName(name + ".dat")
	pc.SetOutputFileName(name + ".txt")
}
