package core

import (
	"math"

	"github.com/draffensperger/golp"
)

// MILP problem with potential multiple kinds of accelerators assigned to a server
type MultiAssignProblem struct {
	BaseProblem
}

// create an instance of the problem
func CreateMultiAssignProblem(numServers int, numAccelerators int, instanceCost []float64, numInstancesPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*MultiAssignProblem, error) {
	bp, err := CreateBaseProblem(numServers, numAccelerators, instanceCost, numInstancesPerReplica,
		ratePerReplica, arrivalRates)
	if err != nil {
		return nil, err
	}
	p := &MultiAssignProblem{
		BaseProblem: *bp}
	p.BaseProblem.Setup = p.Setup
	p.BaseProblem.Solve = p.Solve
	return p, nil
}

// setup constraints and objective function
func (p *MultiAssignProblem) Setup() {
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
			costVector[v0+j] = float64(p.numInstancesPerReplica[i][j]) * p.instanceCost[j]
		}
	}
	p.lp.SetObjFn(costVector)
	// fmt.Println(utils.Pretty1DFloat64("costVector", costVector))

	// excluded infeasible variables (for a given server accelerator pair)
	excluded := make([]float64, numVars)

	// set rate constraints: rate coefficients
	for i := 0; i < p.numServers; i++ {
		rateVector := make([]float64, numVars)
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			rateVector[v0+j] = p.ratePerReplica[i][j]
			if p.ratePerReplica[i][j] == 0 {
				excluded[v0+j] = 1
			}
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
					if p.acceleratorTypesMatrix[k][j] > 0 {
						idx := i*p.numAccelerators + j
						countVector[idx] = float64(p.numInstancesPerReplica[i][j] * p.acceleratorTypesMatrix[k][j])
					}
				}
			}
			p.lp.AddConstraint(countVector, golp.LE, float64(p.unitsAvail[k]))
			// fmt.Printf("k=%d; %s; avail=%d\n", k, utils.Pretty1DFloat64("countVector", countVector), p.unitsAvailByType[k])
		}
	}

	p.lp.AddConstraint(excluded, golp.EQ, 0)
	// fmt.Println(utils.Pretty1DFloat64("excluded", excluded))
}

// solve problem
func (p *MultiAssignProblem) Solve() error {
	// setup up problem
	p.Setup()

	// solve problem with timeout
	if err := p.solveWithTimeout(); err != nil {
		return err
	}

	// extract (optimal) solution
	p.objectiveValue = p.lp.Objective()
	vars := p.lp.Variables()

	// obtain number of replicas and number of used accelerator units
	p.numReplicas = make([][]int, p.numServers)
	p.instancesUsed = make([]int, p.numAccelerators)
	for i := 0; i < p.numServers; i++ {
		p.numReplicas[i] = make([]int, p.numAccelerators)
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			p.numReplicas[i][j] = int(math.Round(vars[v0+j]))
			p.instancesUsed[j] += p.numReplicas[i][j] * p.numInstancesPerReplica[i][j]
		}
	}

	// calculate number of used accelerator units
	p.unitsUsed = make([]int, p.numAcceleratorTypes)
	for k := 0; k < p.numAcceleratorTypes; k++ {
		for j := 0; j < p.numAccelerators; j++ {
			if p.acceleratorTypesMatrix[k][j] > 0 {
				p.unitsUsed[k] += p.instancesUsed[j] * p.acceleratorTypesMatrix[k][j]
			}
		}
	}
	return nil
}
