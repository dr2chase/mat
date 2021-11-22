package main

import (
	"github.com/dr2chase/mat"
)

func main() {
	a := mat.RowMajor[mat.F](3, 3)
	rows, cols := a.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := mat.F(0.0)
			if i == j {
				v = 1
			} else {
				v = mat.F(i - j)
			}
			a.Set(i, j, v)
		}
	}
	b := a.BinaryForAll(a, func(x, y mat.F) mat.F { return mat.F(x + y) })
	mat.Print(b)
}
