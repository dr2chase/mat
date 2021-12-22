package mat_test

import (
	"fmt"
	"github.com/dr2chase/mat"
	"testing"
)

func TestFoo(t *testing.T) {
	rm := mat.RowMajor[mat.F](5, 5)
	cm := mat.ColumnMajor[mat.F](5, 5)
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
	if !mat.Equals(mat.M[mat.F](rm), mat.M[mat.F](cm)) {
		fmt.Println("Row major matrix")
		mat.Print(mat.M[mat.F](rm))
		fmt.Println("Column major matrix")
		mat.Print(mat.M[mat.F](cm))
		t.Fail()
	}
}
