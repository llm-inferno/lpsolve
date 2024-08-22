package core

// the type of optimization problem
type ProblemType int

const (
	MULTI ProblemType = iota
	ASSIGN
	UNKNOWN
)

func (pt ProblemType) String() string {
	return [...]string{"MULTI", "ASSIGN", "UNKNOWN"}[pt]
}

func GetProblemType(s string) ProblemType {
	switch s {
	case "MULTI":
		return MULTI
	case "ASSIGN":
		return ASSIGN
	default:
		return UNKNOWN
	}
}
