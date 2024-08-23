package core

import (
	"errors"

	"github.com/draffensperger/golp"
)

type BaseProblem struct {
	numServers         int
	numAccelerators    int
	unitCost           []float64   // unit cost of accelerators [numAccelerators]
	numUnitsPerReplica [][]int     // number of accelerator units [numServers][numAccelerators]
	ratePerReplica     [][]float64 // max arrival rate to attain SLOs [numServers][numAccelerators]
	arrivalRates       []float64   // arrival rates to servers [numServers]
	isLimited          bool        // solution limited to the available number of accelerator types

	solutionType     golp.SolutionType
	solutionTimeMsec int64
	objectiveValue   float64 // value of objective function
	numReplicas      [][]int // resulting number of replicas [numServers][numAccelerators]
	unitsUsed        []int   // number of used accelerator units [numAccelerators]

	numAcceleratorTypes    int
	acceleratorTypesMatrix [][]int // 0-1 matrix [numAcceleratorTypes][numAccelerators]
	unitsAvailByType       []int   // available number of accelerators [numAcceleratorTypes]
	unitsUsedByType        []int   // number of used units of accelerator types [numAcceleratorTypes]

	lp *golp.LP // lp_solve problem model

	Setup func()       // pre-solve setup
	Solve func() error // solve problem
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
func (p *BaseProblem) SetLimited(numAcceleratorTypes int, unitsAvailByType []int, acceleratorTypesMatrix [][]int) error {
	if len(unitsAvailByType) != numAcceleratorTypes || len(acceleratorTypesMatrix) != numAcceleratorTypes ||
		len(acceleratorTypesMatrix[0]) != p.numAccelerators {
		return errors.New("inconsistent dimension")
	}
	p.isLimited = true
	p.numAcceleratorTypes = numAcceleratorTypes
	p.unitsAvailByType = unitsAvailByType
	p.acceleratorTypesMatrix = acceleratorTypesMatrix
	return nil
}

// unset limited accelerator types option
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
