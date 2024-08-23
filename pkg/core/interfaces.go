package core

import "github.com/draffensperger/golp"

// interface to an optimization problem
type Problem interface {
	// accelerator types are limited
	SetLimited(numAcceleratorTypes int, unitsAvail []int, acceleratorTypesMatrix [][]int) error
	UnSetLimited()

	GetSolutionType() golp.SolutionType
	GetSolutionTimeMsec() int64
	GetObjectiveValue() float64
	GetNumReplicas() [][]int
	GetUnitsUsed() []int
	GetUnitsUsedByType() []int

	// pre-solve setup
	Setup()
	// solve problem
	Solve() error
}
