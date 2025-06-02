package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/draffensperger/golp"
	"github.com/llm-inferno/lpsolve/pkg/utils"
)

const (
	DefaultModelFileName  = "inferno.mod"
	DefaultDataFileName   = "data.dat"
	DefaultOutputFileName = "out.txt"

	DefaultOPLCommand = "oplrun"
)

// path includes '/' at the end
var (
	ModelPath = os.Getenv("CPLEX_MODEL_PATH")
	DataPath  = os.Getenv("CPLEX_DATA_PATH")
)

// Optimization problem solved by CPLEX
type CplexProblem struct {
	BaseProblem

	dataFileName   string
	modelFileName  string
	outputFileName string
}

func CreateCplexProblem(numServers int, numAccelerators int, instanceCost []float64, numInstancesPerReplica [][]int,
	ratePerReplica [][]float64, arrivalRates []float64) (*CplexProblem, error) {
	bp, err := CreateBaseProblem(numServers, numAccelerators, instanceCost, numInstancesPerReplica,
		ratePerReplica, arrivalRates)
	if err != nil {
		return nil, err
	}
	p := &CplexProblem{BaseProblem: *bp,
		dataFileName:   DefaultDataFileName,
		modelFileName:  DefaultModelFileName,
		outputFileName: DefaultOutputFileName}

	p.BaseProblem.Setup = p.Setup
	p.BaseProblem.Solve = p.Solve
	return p, nil
}

// generate data file for CPLEX
func (p *CplexProblem) Setup() error {
	dataString := p.generateDataFile()
	dataFile := DataPath + p.dataFileName
	return os.WriteFile(dataFile, []byte(dataString), 0644)
}

// solve problem using CPLEX by running OPL on data
func (p *CplexProblem) Solve() error {
	startTime := time.Now()

	// generate data file
	if err := p.Setup(); err != nil {
		fmt.Println(err)
		return err
	}

	// invoke OPL on model
	modelFile := ModelPath + p.modelFileName
	dataFile := DataPath + p.dataFileName
	cmdRun := exec.Command(DefaultOPLCommand, modelFile, dataFile)
	stdoutRun, err := cmdRun.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// fmt.Println(string(stdoutRun))
	outFile := DataPath + p.outputFileName
	os.WriteFile(outFile, stdoutRun, 0644)

	// extract objective value
	obj := "grep OBJECTIVE: " + outFile + " | awk '{print $2}'"
	cmdObj := exec.Command("bash", "-c", obj)
	stdoutObj, err := cmdObj.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	stdoutObj = stdoutObj[:len(stdoutObj)-1] // trim end of line
	objective, _ := strconv.ParseFloat(string(stdoutObj), 64)
	// fmt.Printf("objective=%v \n", objective)
	p.objectiveValue = objective

	// obtain number of replicas and number of used accelerator units
	sed := "sed -n '/numReplicas/,/^$/p' " + outFile
	cmdSed := exec.Command("bash", "-c", sed)
	stdoutSed, err := cmdSed.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	nrString := string(stdoutSed)

	nrString = strings.ReplaceAll(nrString, "\n", "")
	nrString = strings.ReplaceAll(nrString, "    ", "")
	nrString = strings.Replace(nrString, "numReplicas = ", "", 1)
	nrString = strings.Replace(nrString, "[", "", -1)
	nrString = strings.Replace(nrString, "]", "", -1)
	nrStringArray := strings.Split(nrString, " ")

	p.numReplicas = make([][]int, p.numServers)
	p.instancesUsed = make([]int, p.numAccelerators)
	for i := 0; i < p.numServers; i++ {
		p.numReplicas[i] = make([]int, p.numAccelerators)
		v0 := i * p.numAccelerators // begin index
		for j := 0; j < p.numAccelerators; j++ {
			p.numReplicas[i][j], _ = strconv.Atoi(nrStringArray[v0+j])
			p.instancesUsed[j] += p.numReplicas[i][j] * p.numInstancesPerReplica[i][j]
		}
	}

	p.solutionType = golp.OPTIMAL

	// calculate number of used accelerator units
	p.unitsUsed = make([]int, p.numAcceleratorTypes)
	for k := 0; k < p.numAcceleratorTypes; k++ {
		for j := 0; j < p.numAccelerators; j++ {
			if p.acceleratorTypesMatrix[k][j] > 0 {
				p.unitsUsed[k] += p.instancesUsed[j] * p.acceleratorTypesMatrix[k][j]
			}
		}
	}

	endTime := time.Now()
	p.solutionTimeMsec = endTime.Sub(startTime).Milliseconds()
	return nil
}

func (p *BaseProblem) generateDataFile() string {
	var b bytes.Buffer

	b.WriteString(utils.Pretty("numServers", p.numServers) + "\n")
	b.WriteString(utils.Pretty("numAccelerators", p.numAccelerators) + "\n")
	if p.isLimited {
		b.WriteString(utils.Pretty("numAcceleratorTypes", p.numAcceleratorTypes) + "\n")
	}
	b.WriteString("\n")

	b.WriteString(utils.Pretty1D("arrivalRates", p.arrivalRates) + "\n")
	b.WriteString(utils.Pretty1D("instanceCost", p.instanceCost) + "\n")
	if p.isLimited {
		b.WriteString(utils.Pretty1D("unitsAvail", p.unitsAvail) + "\n")
	}
	b.WriteString("\n")

	b.WriteString(utils.Pretty2D("numInstancesPerReplica", p.numInstancesPerReplica) + "\n")
	b.WriteString(utils.Pretty2D("ratePerReplica", p.ratePerReplica) + "\n")
	if p.isLimited {
		b.WriteString(utils.Pretty2D("acceleratorTypesMatrix", p.acceleratorTypesMatrix) + "\n")
	}
	b.WriteString("\n")

	return b.String()
}

func (p *CplexProblem) SetDataFileName(dataFileName string) {
	p.dataFileName = dataFileName
}

func (p *CplexProblem) GetDataFileName() string {
	return p.dataFileName
}

func (p *CplexProblem) SetModelFileName(modelFileName string) {
	p.modelFileName = modelFileName
}

func (p *CplexProblem) GetModelFileName() string {
	return p.modelFileName
}

func (p *CplexProblem) SetOutputFileName(outputFileName string) {
	p.outputFileName = outputFileName
}

func (p *CplexProblem) GetOutputFileName() string {
	return p.outputFileName
}
