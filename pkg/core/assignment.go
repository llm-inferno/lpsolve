package core

import (
	"math"
	"time"

	"github.com/draffensperger/golp"
)

// A special MIP problem with binary variables
//   - assign one accelerator type to a server
type AssignmentProblem struct {
	MIPProblem
}

// create an instance of an assignment problem
func CreateAssignmentProblemInstance(numServers int, numAccelerators int, unitCost []float64, numUnitsPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*AssignmentProblem, error) {
	mipProblem, err := CreateMIPProblemInstance(numServers, numAccelerators, unitCost, numUnitsPerReplica,
		ratePerReplica, arrivalRates)
	if err != nil {
		return nil, err
	}
	return &AssignmentProblem{MIPProblem: *mipProblem}, nil
}

// setup constraints and objective function
func (p *AssignmentProblem) setup() {
	// define LP problem
	numVars := p.numServers * p.numAccelerators
	p.lp = golp.NewLP(0, numVars)
	for k := 0; k < numVars; k++ {
		p.lp.SetBinary(k, true)
	}

	// set objective function: cost coefficients
	costVector := make([]float64, numVars)
	for i := 0; i < p.numServers; i++ {
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			costVector[v0+j] = float64(p.numUnitsPerReplica[i][j]) * p.unitCost[j] *
				math.Ceil(p.arrivalRates[i]/p.ratePerReplica[i][j])
		}
	}
	p.lp.SetObjFn(costVector)
	// fmt.Println(utils.Pretty1DFloat64("costVector", costVector))

	// set binary assignment constraints
	for i := 0; i < p.numServers; i++ {
		assignVector := make([]float64, numVars)
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			assignVector[v0+j] = 1
		}
		p.lp.AddConstraint(assignVector, golp.EQ, 1)
		// fmt.Printf("i=%d; %s; tot=%v\n", i, utils.Pretty1DFloat64("assignVector", assignVector), 1)
	}

	// set count limit constraints
	if p.isLimited {
		for j := 0; j < p.numAccelerators; j++ {
			countVector := make([]float64, numVars)
			for i := 0; i < p.numServers; i++ {
				v0 := i * p.numAccelerators // begin index
				countVector[v0+j] = float64(p.numUnitsPerReplica[i][j]) * math.Ceil(p.arrivalRates[i]/p.ratePerReplica[i][j])
			}
			p.lp.AddConstraint(countVector, golp.LE, float64(p.unitsAvail[j]))
			// fmt.Printf("j=%d; %s; avail=%d\n", j, utils.Pretty1DFloat64("countVector", countVector), p.unitsAvail[j])
		}
	}
}

// solve problem
func (p *AssignmentProblem) Solve() error {
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
			p.numReplicas[i][j] = int(math.Ceil(vars[v0+j]) * (p.arrivalRates[i] / p.ratePerReplica[i][j]))
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
