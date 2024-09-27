package utils

import (
	"bytes"
	"fmt"
)

// scalar write: name = value;
func Pretty[T interface{}](name string, x T) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s = %v;", name, x)
	return b.String()
}

// vector write: name = [ v1, v2, ..., vn ];
func Pretty1D[T interface{}](name string, x []T) string {
	n := len(x)
	var b bytes.Buffer
	b.WriteString(name + " = [ ")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%v", x[i])
		if i < n-1 {
			fmt.Fprint(&b, ", ")
		}
	}
	b.WriteString(" ];")
	return b.String()
}

// matrix write:
//
//	name = [
//	[ v11, ..., v1n ],
//	...
//	[ vm1, ..., vmn ]
//	];
func Pretty2D[T interface{}](name string, x [][]T) string {
	n := len(x)
	m := 0
	if n > 0 {
		m = len(x[0])
	}
	var b bytes.Buffer
	b.WriteString(name + " = [\n")
	for i := 0; i < n; i++ {
		b.WriteString("[ ")
		for j := 0; j < m; j++ {
			fmt.Fprintf(&b, "%v", x[i][j])
			if j < m-1 {
				fmt.Fprint(&b, ", ")
			}
		}
		b.WriteString(" ]")
		if i < n-1 {
			fmt.Fprint(&b, ",")
		}
		b.WriteString("\n")
	}
	b.WriteString("];")
	return b.String()
}
