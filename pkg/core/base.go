package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/draffensperger/golp"
	"github.ibm.com/tantawi/lpsolve/pkg/config"
)

// Base optimization problem
type BaseProblem struct {
	numServers             int
	numAccelerators        int
	instanceCost           []float64   // instance cost of accelerators [numAccelerators]
	numInstancesPerReplica [][]int     // number of accelerator instances [numServers][numAccelerators]
	ratePerReplica         [][]float64 // max arrival rate to attain SLOs [numServers][numAccelerators]
	arrivalRates           []float64   // arrival rates to servers [numServers]
	isLimited              bool        // solution limited to the available number of accelerator types

	solutionType     golp.SolutionType
	solutionTimeMsec int64
	objectiveValue   float64 // value of objective function
	numReplicas      [][]int // resulting number of replicas [numServers][numAccelerators]
	instancesUsed    []int   // number of used accelerator instances [numAccelerators]

	numAcceleratorTypes    int
	acceleratorTypesMatrix [][]int // [numAcceleratorTypes][numAccelerators]: number of unit types for an accelerator
	unitsAvail             []int   // available number of accelerator units [numAcceleratorTypes]
	unitsUsed              []int   // number of used units of accelerator [numAcceleratorTypes]

	lp               *golp.LP // lp_solve problem model
	solverTimeoutSec int      // override default timeout

	Setup func() error // pre-solve setup
	Solve func() error // solve problem
}

// create an instance of base problem
func CreateBaseProblem(numServers int, numAccelerators int, instanceCost []float64, numInstancesPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*BaseProblem, error) {
	if numServers <= 0 || numAccelerators <= 0 || len(instanceCost) != numAccelerators ||
		len(numInstancesPerReplica) != numServers || len(numInstancesPerReplica[0]) != numAccelerators ||
		len(ratePerReplica) != numServers || len(ratePerReplica[0]) != numAccelerators ||
		len(arrivalRates) != numServers {
		return nil, errors.New("inconsistent problem size")
	}
	return &BaseProblem{
		numServers:             numServers,
		numAccelerators:        numAccelerators,
		instanceCost:           instanceCost,
		numInstancesPerReplica: numInstancesPerReplica,
		ratePerReplica:         ratePerReplica,
		arrivalRates:           arrivalRates,
		isLimited:              false,
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
	p.unitsAvail = unitsAvail
	p.acceleratorTypesMatrix = acceleratorTypesMatrix
	return nil
}

// unset limited accelerator types option
func (p *BaseProblem) UnSetLimited() {
	p.isLimited = false
}

func (p *BaseProblem) IsLimited() bool {
	return p.isLimited
}

func (p *BaseProblem) SetSolverTimeout(t int) {
	if t > 0 {
		p.solverTimeoutSec = t
	}
}

func (p *BaseProblem) GetSolverTimeout() int {
	return p.solverTimeoutSec
}

// solve MILP problem using a timeout
func (p *BaseProblem) solveWithTimeout() error {
	startTime := time.Now()
	var err error
	var timeoutSec int
	if p.solverTimeoutSec > 0 {
		timeoutSec = p.solverTimeoutSec
	} else {
		timeoutSec = config.DefaultSolverTimeout
	}
	timeout := time.Duration(timeoutSec) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// solution routine
	go func() {
		// p.lp.SetVerboseLevel(golp.DETAILED)
		p.solutionType = p.lp.Solve()
		if p.solutionType != golp.OPTIMAL && p.solutionType != golp.SUBOPTIMAL {
			err = fmt.Errorf("LP solve failed; solutionType=%s", p.solutionType.String())
		}
		cancel()
	}()

	// wait for solve to finish or timeout
	<-ctx.Done()
	endTime := time.Now()
	p.solutionTimeMsec = endTime.Sub(startTime).Milliseconds()
	return err
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

func (p *BaseProblem) GetInstancesUsed() []int {
	return p.instancesUsed
}

func (p *BaseProblem) GetUnitsUsed() []int {
	return p.unitsUsed
}
