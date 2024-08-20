package core

import (
	"errors"
	"math"
	"time"

	"github.com/draffensperger/golp"
)

type MIPProblem struct {
	numServers         int
	numAccelerators    int
	unitCost           []float64   // unit cost of accelerators [numAccelerators]
	numUnitsPerReplica [][]int     // number of accelerator units for pairs of server and accelerator [numServers][numAccelerators]
	ratePerReplica     [][]float64 // max arrival rate for pairs of server and accelerator [numServers][numAccelerators]
	arrivalRates       []float64   // arrival rates to servers [numServers]

	isLimited  bool  // solution limited to the available number of accelerator units
	unitsAvail []int // available number of accelerators [numAccelerators]

	solutionType     golp.SolutionType
	solutionTimeMsec int64
	objectiveValue   float64 // value of objective function
	numReplicas      [][]int // resulting number of replicas for pairs of server and accelerator [numServers][numAccelerators]
	unitsUsed        []int   // number of used accelerator units [numAccelerators]

	lp *golp.LP // problem model
}

// create an instance of MIP problem
func CreateMIPProblemInstance(numServers int, numAccelerators int, unitCost []float64, numUnitsPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*MIPProblem, error) {
	if numServers <= 0 || numAccelerators <= 0 || len(unitCost) != numAccelerators ||
		len(numUnitsPerReplica) != numServers || len(numUnitsPerReplica[0]) != numAccelerators ||
		len(ratePerReplica) != numServers || len(ratePerReplica[0]) != numAccelerators ||
		len(arrivalRates) != numServers {
		return nil, errors.New("inconsistent problem size")
	}
	return &MIPProblem{
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
func (p *MIPProblem) SetLimited(unitsAvail []int) error {
	if len(unitsAvail) != p.numAccelerators {
		return errors.New("inconsistent dimension")
	}
	p.isLimited = true
	p.unitsAvail = unitsAvail
	return nil
}

// unset limited accelerator units option
func (p *MIPProblem) UnSetLimited() {
	p.isLimited = false
}

// setup constraints and objective function
func (p *MIPProblem) setup() {
	// define LP problem
	numVars := p.numServers * p.numAccelerators
	p.lp = golp.NewLP(0, numVars)
	for k := 0; k < numVars; k++ {
		p.lp.SetInt(k, true)
	}

	// set objective function: cost coefficients
	costVector := make([]float64, numVars)
	for i := 0; i < p.numServers; i++ {
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			costVector[v0+j] = float64(p.numUnitsPerReplica[i][j]) * p.unitCost[j]
		}
	}
	p.lp.SetObjFn(costVector)
	// fmt.Println(utils.Pretty1DFloat64("costVector", costVector))

	// set rate constraints: rate coefficients
	for i := 0; i < p.numServers; i++ {
		rateVector := make([]float64, numVars)
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			rateVector[v0+j] = p.ratePerReplica[i][j]
		}
		p.lp.AddConstraint(rateVector, golp.GE, p.arrivalRates[i])
		// fmt.Printf("i=%d; %s; arrv=%v\n", i, utils.Pretty1DFloat64("rateVector", rateVector), p.arrivalRates[i])
	}

	// set count limit constraints
	if p.isLimited {
		for j := 0; j < p.numAccelerators; j++ {
			countVector := make([]float64, numVars)
			for i := 0; i < p.numServers; i++ {
				v0 := i * p.numAccelerators // begin index
				countVector[v0+j] = float64(p.numUnitsPerReplica[i][j])
			}
			p.lp.AddConstraint(countVector, golp.LE, float64(p.unitsAvail[j]))
			// fmt.Printf("j=%d; %s; avail=%d\n", j, utils.Pretty1DFloat64("countVector", countVector), p.unitsAvail[j])
		}
	}
}

// solve problem
func (p *MIPProblem) Solve() error {
	p.setup()

	//lp.SetVerboseLevel(golp.DETAILED)
	startTime := time.Now()
	p.solutionType = p.lp.Solve()
	endTime := time.Now()
	p.solutionTimeMsec = endTime.Sub(startTime).Milliseconds()

	// extract (optimal) solution
	p.objectiveValue = p.lp.Objective()
	vars := p.lp.Variables()

	// obtain number of replicas for pairs of server and accelerator
	p.numReplicas = make([][]int, p.numServers)
	for i := 0; i < p.numServers; i++ {
		p.numReplicas[i] = make([]int, p.numAccelerators)
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			p.numReplicas[i][j] = int(math.Round(vars[v0+j]))
		}
	}

	// calculate number of used accelerator units
	p.unitsUsed = make([]int, p.numAccelerators)
	for i := 0; i < p.numServers; i++ {
		for j := 0; j < p.numAccelerators; j++ {
			p.unitsUsed[j] += p.numReplicas[i][j] * p.numUnitsPerReplica[i][j]
		}
	}
	return nil
}

func (p *MIPProblem) GetSolutionType() golp.SolutionType {
	return p.solutionType
}

func (p *MIPProblem) GetSolutionTimeMsec() int64 {
	return p.solutionTimeMsec
}

func (p *MIPProblem) GetObjectiveValue() float64 {
	return p.objectiveValue
}

func (p *MIPProblem) GetNumReplicas() [][]int {
	return p.numReplicas
}

func (p *MIPProblem) GetUnitsUsed() []int {
	return p.unitsUsed
}
