package utils

import (
	"bytes"
	"fmt"
)

const LargeNumber float64 = 1e+12

func Pretty1DInt(name string, x []int) string {
	n := len(x)
	var b bytes.Buffer
	b.WriteString(name + " = [ ")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d ", x[i])
	}
	b.WriteString("]")
	return b.String()
}

func Pretty1DFloat64(name string, x []float64) string {
	n := len(x)
	var b bytes.Buffer
	b.WriteString(name + " = [ ")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%v ", x[i])
	}
	b.WriteString("]")
	return b.String()
}

func Pretty2DInt(name string, x [][]int) string {
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
			fmt.Fprintf(&b, "%d ", x[i][j])
		}
		b.WriteString("]\n")
	}
	b.WriteString("]")
	return b.String()
}

func Pretty2DFloat64(name string, x [][]float64) string {
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
			fmt.Fprintf(&b, "%v ", x[i][j])
		}
		b.WriteString("]\n")
	}
	b.WriteString("]")
	return b.String()
}
