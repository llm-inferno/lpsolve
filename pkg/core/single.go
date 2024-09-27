package core

import (
	"math"

	"github.com/draffensperger/golp"
)

// A special MILP problem with binary variables
//   - assign one accelerator kind to a server
type SingleAssignProblem struct {
	BaseProblem

	// calculated maximum number of replicas [numServers][numAccelerators]
	maxNumReplicas [][]int
}

// create an instance of an assignment problem
func CreateSingleAssignProblem(numServers int, numAccelerators int, instanceCost []float64, numInstancesPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*SingleAssignProblem, error) {
	bp, err := CreateBaseProblem(numServers, numAccelerators, instanceCost, numInstancesPerReplica,
		ratePerReplica, arrivalRates)
	if err != nil {
		return nil, err
	}
	p := &SingleAssignProblem{
		BaseProblem: *bp}
	p.BaseProblem.Setup = p.Setup
	p.BaseProblem.Solve = p.Solve
	return p, nil
}

// setup constraints and objective function
func (p *SingleAssignProblem) Setup() {
	// define LP problem
	numVars := p.numServers * p.numAccelerators
	p.lp = golp.NewLP(0, numVars)
	for k := 0; k < numVars; k++ {
		p.lp.SetBinary(k, true)
	}

	// excluded infeasible variables (for a given server accelerator pair)
	excluded := make([]float64, numVars)

	// calculate max number of replicas
	p.maxNumReplicas = make([][]int, p.numServers)
	for i := 0; i < p.numServers; i++ {
		p.maxNumReplicas[i] = make([]int, p.numAccelerators)
		for j := 0; j < p.numAccelerators; j++ {
			if p.ratePerReplica[i][j] > 0 {
				p.maxNumReplicas[i][j] = int(math.Ceil(p.arrivalRates[i] / p.ratePerReplica[i][j]))
			} else {
				excluded[i*p.numAccelerators+j] = 1
			}
		}
	}
	// fmt.Println(utils.Pretty2D("maxNumReplicas", p.maxNumReplicas))

	// set objective function: cost coefficients
	costVector := make([]float64, numVars)
	for i := 0; i < p.numServers; i++ {
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			costVector[v0+j] = float64(p.numInstancesPerReplica[i][j]*p.maxNumReplicas[i][j]) * p.instanceCost[j]
		}
	}
	p.lp.SetObjFn(costVector)
	// fmt.Println(utils.Pretty1D("costVector", costVector))

	// set binary assignment constraints - only one variable set to one per server
	for i := 0; i < p.numServers; i++ {
		assignVector := make([]float64, numVars)
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			assignVector[v0+j] = 1
		}
		p.lp.AddConstraint(assignVector, golp.EQ, 1)
		// fmt.Printf("i=%d; %s tot=%v\n", i, utils.Pretty1D("assignVector", assignVector), 1)
	}

	// set count limit constraints
	if p.isLimited {
		for k := 0; k < p.numAcceleratorTypes; k++ {
			countVector := make([]float64, numVars)
			for i := 0; i < p.numServers; i++ {
				for j := 0; j < p.numAccelerators; j++ {
					if p.acceleratorTypesMatrix[k][j] > 0 {
						idx := i*p.numAccelerators + j
						countVector[idx] = float64(p.numInstancesPerReplica[i][j] * p.maxNumReplicas[i][j] * p.acceleratorTypesMatrix[k][j])
					}
				}
			}
			p.lp.AddConstraint(countVector, golp.LE, float64(p.unitsAvail[k]))
			// fmt.Printf("k=%d; %s avail=%d\n", k, utils.Pretty1D("countVector", countVector), p.unitsAvail[k])
		}
	}

	p.lp.AddConstraint(excluded, golp.EQ, 0)
	// fmt.Println(utils.Pretty1D("excluded", excluded))
}

// solve problem
func (p *SingleAssignProblem) Solve() error {
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
			p.numReplicas[i][j] = int(math.Round(vars[v0+j])) * p.maxNumReplicas[i][j]
			p.instancesUsed[j] += p.numReplicas[i][j] * p.numInstancesPerReplica[i][j]
		}
	}

	// calculate number of used accelerator types
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
