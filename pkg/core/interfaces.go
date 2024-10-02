package core

import (
	"github.com/draffensperger/golp"
)

// interface to an optimization problem
type Problem interface {
	// limiting number of available accelerator types
	SetLimited(numAcceleratorTypes int, unitsAvail []int, acceleratorTypesMatrix [][]int) error
	UnSetLimited()
	IsLimited() bool

	// pre-solve setup
	Setup() error
	// solve problem
	SetSolverTimeout(int)
	GetSolverTimeout() int
	Solve() error

	// problem solution
	GetSolutionType() golp.SolutionType
	GetSolutionTimeMsec() int64
	GetObjectiveValue() float64
	GetNumReplicas() [][]int
	GetInstancesUsed() []int
	GetUnitsUsed() []int
}
