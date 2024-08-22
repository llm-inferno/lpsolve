package core

import (
	"math"
	"time"

	"github.com/draffensperger/golp"
)

// MIP problem with potential multiple accelerator types assigned to a server
type MIPProblem struct {
	BaseProblem
}

// create an instance of MIP problem
func CreateMIPProblem(numServers int, numAccelerators int, unitCost []float64, numUnitsPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*MIPProblem, error) {
	bp, err := CreateBaseProblem(numServers, numAccelerators, unitCost, numUnitsPerReplica,
		ratePerReplica, arrivalRates)
	if err != nil {
		return nil, err
	}
	mip := &MIPProblem{
		BaseProblem: *bp}
	mip.BaseProblem.Setup = mip.Setup
	mip.BaseProblem.Solve = mip.Solve
	return mip, nil
}

// setup constraints and objective function
func (p *MIPProblem) Setup() {
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
		for k := 0; k < p.numAcceleratorTypes; k++ {
			countVector := make([]float64, numVars)
			for i := 0; i < p.numServers; i++ {
				for j := 0; j < p.numAccelerators; j++ {
					if p.acceleratorTypesMatrix[k][j] == 1 {
						idx := i*p.numAccelerators + j
						countVector[idx] = float64(p.numUnitsPerReplica[i][k])
					}
				}
			}
			p.lp.AddConstraint(countVector, golp.LE, float64(p.unitsAvailByType[k]))
			// fmt.Printf("k=%d; %s; avail=%d\n", k, utils.Pretty1DFloat64("countVector", countVector), p.unitsAvailByType[k])
		}
	}
}

// solve problem
func (p *MIPProblem) Solve() error {
	p.Setup()

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

	// calculate number of used accelerator units
	p.unitsUsedByType = make([]int, p.numAcceleratorTypes)
	for k := 0; k < p.numAcceleratorTypes; k++ {
		for j := 0; j < p.numAccelerators; j++ {
			if p.acceleratorTypesMatrix[k][j] == 1 {
				p.unitsUsedByType[k] += p.unitsUsed[j]
			}
		}
	}
	return nil
}
