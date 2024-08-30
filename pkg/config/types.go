package config

// the type of optimization problem
type ProblemType int

const (
	SINGLE ProblemType = iota // a single kind of accelerator to a server
	MULTI                     // multiple kinds of accelerators to a server
	UNKNOWN
)

func (pt ProblemType) String() string {
	return [...]string{"SINGLE", "MULTI", "UNKNOWN"}[pt]
}

func GetProblemType(s string) ProblemType {
	switch s {
	case "SINGLE":
		return SINGLE
	case "MULTI":
		return MULTI
	default:
		return UNKNOWN
	}
}
