package mat

import (
	"fmt"
)

type Field[T any] interface {
	One() T // ignores parameter
	Zero() T // ignores parameter; ought to be the initialized type anyway.
	Plus(T) T
	Times(T) T
	Minus(T) T // this, instead of arithmetic inverse
	Divide(T) T // this, instead o multiplicative inverse
	Equals(T) bool
}
// Note the possibility of a type that is not T implementing Field[T].

// Vector
type V[T Field[T]] interface {
	// Plus(V[T]) V[T]
	// Times(T) V[T] // Scalar
	// Inner(V[T]) T
	At(i int) T	
	Len() int
}

// Matrix
type M[T Field[T]] interface {
	BinaryForAll(M[T], func(a, b T) T) M[T]
	UnaryForAll(func(a T) T) M[T]
	Transpose() M[T] // = Transpose(M)
	One() M[T] // identity M -- same shape as receiver
	Zero() M[T] // zero M -- same shape as receiver.
	Dims() (rows, cols int)
	Row(i int) V[T]
	Col(i int) V[T]
	At(row, col int) T
	// Times(M[T]) M[T] // linear algebra matrix multiplication
	// ScalarTimes(T) M[T]
	// TimesVector(V[T]) V[T] // yields (V^T M)^T
	// VectorTimes(V[T]) V[T] // M V
}

// SCALAR TYPES

type F float64
type B bool
type C complex128
// Need to do something with GF8, GF16, for sake of the exercise

func (b B) One() B {
	return B(true)
}
func (b F) One() F {
	return 1.0
}
func (b C) One() C {
	return 1.0+0i
}


func (b B) Zero() (r B) {
	return 
}
func (b F) Zero() (r F) {
	return
}
func (b C) Zero() (r C) {
	return 
}


func (b F) Plus(c F) F {
	return F(b+c)
}
func (b F) Times(c F) F {
	return F(b*c)
}

func (b F) Minus(c F) F {
	return F(b-c)
}
func (b F) Divide(c F) F {
	return F(b/c)
}
func (b F) Equals(c F) bool {
	return b==c
}


func (b B) Plus(c B) B {
	return B(b!=c)
}
func (b B) Times(c B) B {
	return B(b && c)
}
func (b B) Minus(c B) B { // x-y == x + (-y) == x + y
	return B(b!=c)
}
func (b B) Divide(c B) B { // x/y == x * (1/y) == x * y
	if !c {
		panic("Boolean divide by zero")
	}
	return B(b && c)
}
func (b B) Equals(c B) bool {
	return b==c
}


func (b C) Plus(c C) C {
	return C(b+c)
}
func (b C) Times(c C) C {
	return C(b*c)
}
func (b C) Minus(c C) C {
	return C(b-c)
}
func (b C) Divide(c C) C {
	return C(b/c)
}
func (b C) Equals(c C) bool {
	return b==c
}


// Vectors

type ContiguousVector[T Field[T]] struct {
	v []T
}

func (v *ContiguousVector[T]) At(i int) T {
	return v.v[i]
}

func (v *ContiguousVector[T]) Len() int {
	return len(v.v)
}

type StridedVector[T Field[T]] struct {
	v []T
	stride, len int
}

func (v *StridedVector[T]) At(i int) T {
	return v.v[i*v.stride]
}

func (v *StridedVector[T]) Len() int {
	return v.len
}

// Mutable Matrix
type MuM[T Field[T]] interface {
	M[T]
	Set(row, col int, v T)
	SetBinaryForAll(A, B M[T], f func(a, b T) T) MuM[T]
	SetUnaryForAll(A M[T], f func(a T) T) MuM[T]
	SetCopy(M[T]) MuM[T]
}

// HELPER FUNCTIONS
func assertSameDims[T Field[T]](a, b M[T]) {
	ai, aj := a.Dims()
	bi, bj := b.Dims()
	if ai != bi ||aj != bj {
		panic("Matrix dimensions must match")
	}
}

func CheckBounds[T Field[T]](a M[T], i, j int) {
	rows, cols := a.Dims()
	if i >= rows || i < 0 {
		panic(fmt.Errorf("Out-of-bounds row index, expected 0 <= %d < %d", i, rows))
	}
	if j >= cols || j < 0 {
		panic(fmt.Errorf("Out-of-bounds column index, expected 0 <= %d < %d", j, cols))
	}	
}

func SetBinaryForAll[T Field[T]](aMu MuM[T], b, c M[T], f func(a, b T) T) MuM[T] {
	rows, cols := aMu.Dims()
	for i:= 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i,j, f(b.At(i, j), c.At(i,j)))
		}
	}
	return aMu
}

func SetUnaryForAll[T Field[T]](aMu MuM[T], b M[T], f func(a T) T) MuM[T] {
	rows, cols := aMu.Dims()
	for i:= 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i,j, f(b.At(i, j)))
		}
	}
	return aMu
}

func SetCopy[T Field[T]](aMu MuM[T], b M[T]) MuM[T] {
	rows, cols := aMu.Dims()
	for i:= 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i,j, b.At(i,j))
		}
	}
	return aMu
}

func SetByIJ[T Field[T]](aMu MuM[T], f func(i, j int) T) MuM[T] {
	rows, cols := aMu.Dims()
	for i:= 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i,j, f(i,j))
		}
	}
	return aMu
}

func Equals[T Field[T]](a, b M[T]) bool {
	rows, cols := a.Dims()
	if br, bc := b.Dims(); br != rows || bc != cols {
		return false
	}
	for i:= 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if !a.At(i,j).Equals(b.At(i,j)) {
				return false
			}
		}
	}
	return true
}

func Print[T Field[T]](a M[T]) {
	rows, cols := a.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Printf(" %v", a.At(i,j))
		}
		fmt.Printf("\n")
	}
}

// TRANSPOSE[D]

type Transposed[T Field[T]] struct {
	m M[T]
}

func Transpose[T Field[T]](a M[T]) M[T] {
	switch x := a.(type) {
	case *Transposed[T] : return x.m
	default: return &Transposed[T]{m:a}
	}
}

func (a * Transposed[T]) BinaryForAll(b M[T], f func(T, T) T) M[T] {
	assertSameDims[T](a, b)
	switch x := b.(type) {
	case *Transposed[T] : return Transpose(a.m.BinaryForAll(x.m, f))
	}
	return Transpose(a.m.BinaryForAll(Transpose(b), f))
}

func (a * Transposed[T]) UnaryForAll(f func(T) T) M[T] {
	return Transpose(a.m.UnaryForAll(f))
}

func (a * Transposed[T]) SetBinaryForAll(b,c M[T], f func(T, T) T) MuM[T] {	
	return SetBinaryForAll[T](a, b, c, f)
}

func (a * Transposed[T]) SetUnaryForAll(b M[T], f func(T) T) MuM[T] {
	return SetUnaryForAll[T](a, b, f)
}

func (a * Transposed[T]) SetCopy(b M[T]) MuM[T] {
	return SetCopy[T](a,b)
}

func (a * Transposed[T]) Transpose() M[T] {
	return a.m
}

func (a * Transposed[T]) One() M[T] {
	return Transpose(a.m.One()) // TODO do better
}

func (a * Transposed[T]) Zero() M[T] {
	return Transpose(a.m.Zero()) // TODO do better
}

func (a * Transposed[T]) Dims() (rows int, cols int) {
	cols,rows = a.m.Dims()
	return
}

func (a * Transposed[T]) Row(i int) V[T] {
	return a.m.Col(i)
}

func (a * Transposed[T]) Col(i int) V[T] {
	return a.m.Row(i)
}

func (a * Transposed[T]) At(i,j int) T {
	return a.m.At(j,i)
}

func (a * Transposed[T]) Set(i,j int, v T) {
	a.m.(MuM[T]).Set(j,i, v)
}

// ROW MAJOR / CONTIGUOUS ROW

type ContiguousRowMatrix[T Field[T]] struct {
	x []T
	rows, cols int
}

func RowMajor[T Field[T]](rows, cols int) *ContiguousRowMatrix[T] {
	prod := rows * cols // TODO check for overflow
	return &ContiguousRowMatrix[T]{cols:cols, rows:rows, x:make([]T, prod, prod)}
}

func (a * ContiguousRowMatrix[T]) BinaryForAll(b M[T], f func(T, T) T) M[T] {
	assertSameDims[T](a, b)
	r := RowMajor[T](a.rows, a.cols)
	return r.SetBinaryForAll(a, b, f)
}

func (a * ContiguousRowMatrix[T]) UnaryForAll(f func(T) T) M[T] {
	r := RowMajor[T](a.rows, a.cols)
	return r.SetUnaryForAll(a, f)
}

func (a * ContiguousRowMatrix[T]) Transpose() M[T] {
	return Transpose[T](a)
}

func (a * ContiguousRowMatrix[T]) One() M[T] {
	b := RowMajor[T](a.rows, a.cols)
	var z T
	one := z.One()
	for i := 0; i < len(a.x); i += a.cols+1 {
		a.x[i] = one
	}
	return b
}

func (a * ContiguousRowMatrix[T]) Zero() M[T] {
	return RowMajor[T](a.rows, a.cols)
}

func (a * ContiguousRowMatrix[T]) Dims() (rows int, cols int) {
	return a.rows, a.cols
}

func (a * ContiguousRowMatrix[T]) Row(i int) V[T] {
	start := i*a.cols
	return &ContiguousVector[T]{a.x[start:start+a.cols]}
}

func (a * ContiguousRowMatrix[T]) Col(i int) V[T] {
	return &StridedVector[T]{v:a.x[i:], stride:a.cols, len:a.rows}
}

func (a * ContiguousRowMatrix[T]) At(i,j int) T {
	CheckBounds[T](a, i, j)
	return a.x[i*a.cols + j]
}

// Mutable matrix methods

func (a * ContiguousRowMatrix[T]) Set(i,j int, v T) {
	CheckBounds[T](a, i, j)
	a.x[i*a.cols + j] = v
}

func (a * ContiguousRowMatrix[T]) SetBinaryForAll(b,c M[T], f func(T, T) T) MuM[T] {	
	return SetBinaryForAll[T](a, b, c, f)
}

func (a * ContiguousRowMatrix[T]) SetUnaryForAll(b M[T], f func(T) T) MuM[T] {
	return SetUnaryForAll[T](a, b, f)
}

func (a * ContiguousRowMatrix[T]) SetCopy(b M[T]) MuM[T] {
	return SetCopy[T](a,b)
}


// COLUMN MAJOR / CONTIGUOUS COLUMN

type ContiguousColumnMatrix[T Field[T]] struct {
	x []T
	rows, cols int
}

func ColumnMajor[T Field[T]](rows, cols int) *ContiguousColumnMatrix[T] {
	prod := rows * cols // TODO check for overflow
	return &ContiguousColumnMatrix[T]{cols:cols, rows:rows, x:make([]T, prod, prod)}
}

func (a * ContiguousColumnMatrix[T]) BinaryForAll(b M[T], f func(T, T) T) M[T] {
	assertSameDims[T](a, b)
	r := ColumnMajor[T](a.rows, a.cols)
	return r.SetBinaryForAll(a, b, f)
}

func (a * ContiguousColumnMatrix[T]) UnaryForAll(f func(T) T) M[T] {
	r := ColumnMajor[T](a.rows, a.cols)
	return r.SetUnaryForAll(a, f)
}

func (a * ContiguousColumnMatrix[T]) Transpose() M[T] {
	return Transpose[T](a)
}

func (a * ContiguousColumnMatrix[T]) One() M[T] {
	b := ColumnMajor[T](a.rows, a.cols)
	var z T
	one := z.One()
	for i := 0; i < len(a.x); i += a.rows+1 {
		a.x[i] = one
	}
	return b
}

func (a * ContiguousColumnMatrix[T]) Zero() M[T] {
	return ColumnMajor[T](a.rows, a.cols)
}

func (a * ContiguousColumnMatrix[T]) Dims() (rows int, cols int) {
	return a.rows, a.cols
}

func (a * ContiguousColumnMatrix[T]) Row(i int) V[T] {  // fixme
	return &StridedVector[T]{v:a.x[i:], stride:a.rows, len:a.cols}
}

func (a * ContiguousColumnMatrix[T]) Col(i int) V[T] { // fixme
	start := i*a.rows
	return &ContiguousVector[T]{a.x[start:start+a.rows]}
}

func (a * ContiguousColumnMatrix[T]) At(i,j int) T {
	CheckBounds[T](a, i, j)
	return a.x[i + j*a.rows]
}

// Mutable matrix methods

func (a * ContiguousColumnMatrix[T]) Set(i,j int, v T) {
	CheckBounds[T](a, i, j)
	a.x[i + j*a.rows] = v
}

func (a * ContiguousColumnMatrix[T]) SetBinaryForAll(b,c M[T], f func(T, T) T) MuM[T] {	
	return SetBinaryForAll[T](a, b, c, f)
}

func (a * ContiguousColumnMatrix[T]) SetUnaryForAll(b M[T], f func(T) T) MuM[T] {
	return SetUnaryForAll[T](a, b, f)
}

func (a * ContiguousColumnMatrix[T]) SetCopy(b M[T]) MuM[T] {
	return SetCopy[T](a,b)
}





