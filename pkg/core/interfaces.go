package core

import "github.com/draffensperger/golp"

// interface to optimization problem
type Problem interface {
	SetLimited(numAcceleratorTypes int, unitsAvail []int, acceleratorTypesMatrix [][]int) error
	UnSetLimited()

	GetSolutionType() golp.SolutionType
	GetSolutionTimeMsec() int64
	GetObjectiveValue() float64
	GetNumReplicas() [][]int
	GetUnitsUsed() []int
	GetUnitsUsedByType() []int

	Setup()
	Solve() error
}
