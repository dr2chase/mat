package mat_test

import (
	"fmt"
	"github.com/dr2chase/mat"
	"testing"
)

func rmAndCm() (rm, cm mat.MuM[mat.F]) {
	rm = mat.RowMajor[mat.F](5, 5)
	cm = mat.ColumnMajor[mat.F](5, 5)
	f := func(i, j int) mat.F {
		diff := j - i
		switch diff {
		case 0:
			return mat.F(1)
		case -1, 1:
			return mat.F(2 * diff)
		}
		return mat.F(0)
	}
	mat.SetByIJ[mat.F](rm, f)
	mat.SetByIJ[mat.F](cm, f)
	return
}

func TestEqual(t *testing.T) {
	rm, cm := rmAndCm()
	fmt.Println("Row major matrix")
	rm.Print()
	fmt.Println("Column major matrix")
	cm.Print()
	if !mat.Equals(mat.M[mat.F](rm), mat.M[mat.F](cm)) {
		t.Fail()
	}
	if !cm.Equals(rm) {
		t.Fail()
	}
}

func TestVM(t *testing.T) {
	// rm, cm := rmAndCm()
	v := mat.Vector[mat.F](5)
	v.SetByI(func(i int) mat.F { return mat.F(i + 1) })
}
