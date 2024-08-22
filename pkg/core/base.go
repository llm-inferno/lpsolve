package core

import (
	"errors"

	"github.com/draffensperger/golp"
)

type BaseProblem struct {
	numServers         int
	numAccelerators    int
	unitCost           []float64   // unit cost of accelerators [numAccelerators]
	numUnitsPerReplica [][]int     // number of accelerator units for pairs of server and accelerator [numServers][numAccelerators]
	ratePerReplica     [][]float64 // max arrival rate for pairs of server and accelerator [numServers][numAccelerators]
	arrivalRates       []float64   // arrival rates to servers [numServers]

	solutionType     golp.SolutionType
	solutionTimeMsec int64
	objectiveValue   float64 // value of objective function
	numReplicas      [][]int // resulting number of replicas for pairs of server and accelerator [numServers][numAccelerators]
	unitsUsed        []int   // number of used accelerator units [numAccelerators]

	isLimited              bool // solution limited to the available number of accelerator units
	numAcceleratorTypes    int
	acceleratorTypesMatrix [][]int // 0-1 matrix [numAcceleratorTypes][numAccelerators]
	unitsAvailByType       []int   // available number of accelerators [numAcceleratorTypes]
	unitsUsedByType        []int   // number of used units of accelerator types [numAcceleratorTypes]

	lp *golp.LP // problem model

	Setup func()
	Solve func() error
}

// create an instance of base problem
func CreateBaseProblem(numServers int, numAccelerators int, unitCost []float64, numUnitsPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*BaseProblem, error) {
	if numServers <= 0 || numAccelerators <= 0 || len(unitCost) != numAccelerators ||
		len(numUnitsPerReplica) != numServers || len(numUnitsPerReplica[0]) != numAccelerators ||
		len(ratePerReplica) != numServers || len(ratePerReplica[0]) != numAccelerators ||
		len(arrivalRates) != numServers {
		return nil, errors.New("inconsistent problem size")
	}
	return &BaseProblem{
		numServers:         numServers,
		numAccelerators:    numAccelerators,
		unitCost:           unitCost,
		numUnitsPerReplica: numUnitsPerReplica,
		ratePerReplica:     ratePerReplica,
		arrivalRates:       arrivalRates,
		isLimited:          false,
	}, nil
}

// set limited accelerator units option
func (p *BaseProblem) SetLimited(numAcceleratorTypes int, unitsAvail []int, acceleratorTypesMatrix [][]int) error {
	if len(unitsAvail) != numAcceleratorTypes || len(acceleratorTypesMatrix) != numAcceleratorTypes ||
		len(acceleratorTypesMatrix[0]) != p.numAccelerators {
		return errors.New("inconsistent dimension")
	}
	p.isLimited = true
	p.numAcceleratorTypes = numAcceleratorTypes
	p.unitsAvailByType = unitsAvail
	p.acceleratorTypesMatrix = acceleratorTypesMatrix
	return nil
}

// unset limited accelerator units option
func (p *BaseProblem) UnSetLimited() {
	p.isLimited = false
}

func (p *BaseProblem) GetSolutionType() golp.SolutionType {
	return p.solutionType
}

func (p *BaseProblem) GetSolutionTimeMsec() int64 {
	return p.solutionTimeMsec
}

func (p *BaseProblem) GetObjectiveValue() float64 {
	return p.objectiveValue
}

func (p *BaseProblem) GetNumReplicas() [][]int {
	return p.numReplicas
}

func (p *BaseProblem) GetUnitsUsed() []int {
	return p.unitsUsed
}

func (p *BaseProblem) GetUnitsUsedByType() []int {
	return p.unitsUsedByType
}
