// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
)

type Field[T any] interface {
	One() T  // ignores parameter
	Zero() T // ignores parameter; ought to be the zero value anyway.
	Plus(T) T
	Times(T) T
	Minus(T) T  // this, instead of arithmetic inverse
	Divide(T) T // this, instead of multiplicative inverse
	Equals(T) bool
}

// Note the possibility of a type that is not T implementing Field[T].

// Vector
type V[T Field[T]] interface {
	Plus(V[T]) V[T]
	// Times(T) V[T] // Scalar
	Inner(V[T]) T
	At(i int) T
	Len() int
	Print()
	Equals(V[T]) bool
}

// Matrix
type M[T Field[T]] interface {
	BinaryForAll(M[T], func(a, b T) T) M[T]
	UnaryForAll(func(a T) T) M[T]
	Transpose() M[T] // = Transpose(M)
	One() M[T]       // identity M -- same shape as receiver
	Zero() M[T]      // zero M -- same shape as receiver.
	Dims() (rows, cols int)
	Row(i int) V[T]
	Col(i int) V[T]
	At(row, col int) T
	Equals(b M[T]) bool
	// Times(M[T]) M[T] // linear algebra matrix multiplication
	ScalarTimes(T) M[T]
	LeftTimesVector(V[T]) V[T] // yields (V^T M)^T
	TimesVector(V[T]) V[T]     // M V
	Print()
	// Set(row, col int, v T) // M[row, col] = v
}

// Mutable Matrix
type MuM[T Field[T]] interface {
	M[T]
	MuZero() MuM[T]
	MuOne() MuM[T]

	Set(row, col int, v T)
	SetBinaryForAll(B, C M[T], f func(a, b T) T) MuM[T] // ∀ i,j ∈ range(A), A[i,j] = f(B[i,j], C[i,j])
	SetUnaryForAll(B M[T], f func(a T) T) MuM[T]        // ∀ i,j ∈ range(A), A[i,j] = f(B[i,j])
	SetCopy(B M[T]) MuM[T]                              // ∀ i,j ∈ range(A), A[i,j] = B[i,j]
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
	return 1.0 + 0i
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
	return F(b + c)
}
func (b F) Times(c F) F {
	return F(b * c)
}

func (b F) Minus(c F) F {
	return F(b - c)
}
func (b F) Divide(c F) F {
	return F(b / c)
}
func (b F) Equals(c F) bool {
	return b == c
}

func (b B) Plus(c B) B {
	return B(b != c)
}
func (b B) Times(c B) B {
	return B(b && c)
}
func (b B) Minus(c B) B { // x-y == x + (-y) == x + y
	return B(b != c)
}
func (b B) Divide(c B) B { // x/y == x * (1/y) == x * y
	if !c {
		panic("Boolean divide by zero")
	}
	return B(b && c)
}
func (b B) Equals(c B) bool {
	return b == c
}

func (b C) Plus(c C) C {
	return C(b + c)
}
func (b C) Times(c C) C {
	return C(b * c)
}
func (b C) Minus(c C) C {
	return C(b - c)
}
func (b C) Divide(c C) C {
	return C(b / c)
}
func (b C) Equals(c C) bool {
	return b == c
}

// Vectors

func InnerVV[T Field[T]](v, w V[T]) T {
	var sum T
	for i := 0; i < v.Len(); i++ {
		sum = sum.Plus(v.At(i).Times(w.At(i)))
	}
	return sum
}

func EqualsV[T Field[T]](a, b V[T]) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		if !a.At(i).Equals(b.At(i)) {
			return false
		}
	}
	return true
}

func PrintV[T Field[T]](a V[T]) {
	n := a.Len()
	for i := 0; i < n; i++ {
		fmt.Printf(" %v", a.At(i))
	}
	fmt.Printf("\n")
}

type ContiguousVector[T Field[T]] struct {
	v []T
}

// type MutableContiguousVector[T Field[T]] struct {
// 	ContiguousVector[T]
// }

func (v *ContiguousVector[T]) Set(i int, x T) {
	v.v[i] = x
}

func (v *ContiguousVector[T]) SetByI(f func(i int) T) {
	for i := range v.v {
		v.v[i] = f(i)
	}
}

func Vector[T Field[T]](n int) *ContiguousVector[T] {
	v := &ContiguousVector[T]{v: make([]T, n, n)}
	return v
}

func (v *ContiguousVector[T]) At(i int) T {
	return v.v[i]
}

func (v *ContiguousVector[T]) Len() int {
	return len(v.v)
}

func (v *ContiguousVector[T]) Inner(w V[T]) T {
	return InnerVV[T](v, w)
}

func (v *ContiguousVector[T]) Plus(w V[T]) V[T] {
	u := Vector[T](v.Len())
	for i := 0; i < v.Len(); i++ {
		u.v[i] = v.At(i).Plus(w.At(i))
	}
	return u
}

func (v *ContiguousVector[T]) Equals(w V[T]) bool {
	return EqualsV[T](v, w)
}

func (v *ContiguousVector[T]) Print() {
	PrintV[T](v)
}

type StridedVector[T Field[T]] struct {
	v           []T
	stride, len int
}

func (v *StridedVector[T]) At(i int) T {
	return v.v[i*v.stride]
}

func (v *StridedVector[T]) Len() int {
	return v.len
}

func (v *StridedVector[T]) Inner(w V[T]) T {
	return InnerVV[T](v, w)
}

func (v *StridedVector[T]) Plus(w V[T]) V[T] {
	u := Vector[T](v.Len())
	for i := 0; i < v.Len(); i++ {
		u.v[i] = v.At(i).Plus(w.At(i))
	}
	return u
}

func (v *StridedVector[T]) Equals(w V[T]) bool {
	return EqualsV[T](v, w)
}

func (v *StridedVector[T]) Print() {
	PrintV[T](v)
}

// MATRIX HELPER FUNCTIONS
func assertSameDims[T Field[T]](a, b M[T]) {
	ai, aj := a.Dims()
	bi, bj := b.Dims()
	if ai != bi || aj != bj {
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
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i, j, f(b.At(i, j), c.At(i, j)))
		}
	}
	return aMu
}

func SetUnaryForAll[T Field[T]](aMu MuM[T], b M[T], f func(a T) T) MuM[T] {
	rows, cols := aMu.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i, j, f(b.At(i, j)))
		}
	}
	return aMu
}

func SetCopy[T Field[T]](aMu MuM[T], b M[T]) MuM[T] {
	rows, cols := aMu.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i, j, b.At(i, j))
		}
	}
	return aMu
}

func SetByIJ[T Field[T]](aMu MuM[T], f func(i, j int) T) MuM[T] {
	rows, cols := aMu.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aMu.Set(i, j, f(i, j))
		}
	}
	return aMu
}

func EqualsM[T Field[T]](a, b M[T]) bool {
	rows, cols := a.Dims()
	if br, bc := b.Dims(); br != rows || bc != cols {
		return false
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if !a.At(i, j).Equals(b.At(i, j)) {
				return false
			}
		}
	}
	return true
}

func PrintM[T Field[T]](a M[T]) {
	rows, cols := a.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Printf(" %v", a.At(i, j))
		}
		fmt.Printf("\n")
	}
}

func Transpose[T Field[T]](a M[T]) M[T] {
	switch t := a.(type) {
	case *Transposed[T]:
		return t.m
	default:
		r := &Transposed[T]{m: a}
		r.self = r
		return r
	}
}

func MuTranspose[T Field[T]](a MuM[T]) MuM[T] {
	switch t := a.(type) {
	case *Transposed[T]:
		return t.m.(MuM[T])
	default:
		r := &Transposed[T]{m: a}
		r.self = r
		return r
	}
}

// DEFAULT METHODS

type Default[T Field[T]] struct {
	self MuM[T]
}

func (a *Default[T]) SetBinaryForAll(b, c M[T], f func(T, T) T) MuM[T] {
	return SetBinaryForAll[T](a.self, b, c, f)
}

func (a *Default[T]) SetUnaryForAll(b M[T], f func(T) T) MuM[T] {
	return SetUnaryForAll[T](a.self, b, f)
}

func (a *Default[T]) SetCopy(b M[T]) MuM[T] {
	return SetCopy[T](a.self, b)
}

func (a *Default[T]) ScalarTimes(x T) M[T] {
	s := a.self
	b := s.MuZero()
	r, c := s.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			b.Set(i, j, x.Times(s.At(i, j)))
		}
	}
	return b
}

func (a *Default[T]) Equals(b M[T]) bool {
	return EqualsM[T](a.self, b)
}

func (a *Default[T]) Print() {
	PrintM[T](a.self)
}

func (a *Default[T]) Transpose() M[T] {
	return Transpose[T](a.self)
}

// func (a *Default[T]) Times(b M[T]) M[T] { // linear algebra matrix multiplication
// }

// func (a *Default[T]) ScalarTimes(x T) M[T] {
// }

func (a *Default[T]) LeftTimesVector(v V[T]) V[T] { // yields (V^T M)^T
	// result vector has as many elements as columns
	s := a.self
	_, cols := s.Dims()
	w := Vector[T](cols)
	for i := 0; i < cols; i++ {
		w.v[i] = w.v[i].Plus(v.Inner(s.Row(i)))
	}
	return w
}

func (a *Default[T]) TimesVector(v V[T]) V[T] { // M V
	// result vector has as many elements as rows
	s := a.self
	rows, _ := s.Dims()
	w := Vector[T](rows)
	for i := 0; i < rows; i++ {
		w.v[i] = w.v[i].Plus(v.Inner(s.Col(i)))
	}
	return w
}

// TRANSPOSE[D]

type Transposed[T Field[T]] struct {
	Default[T]
	m M[T]
}

func (a *Transposed[T]) BinaryForAll(b M[T], f func(T, T) T) M[T] {
	assertSameDims[T](a, b)
	switch x := b.(type) {
	case *Transposed[T]:
		return Transpose(a.m.BinaryForAll(x.m, f))
	}
	return Transpose(a.m.BinaryForAll(Transpose(b), f))
}

func (a *Transposed[T]) UnaryForAll(f func(T) T) M[T] {
	return Transpose(a.m.UnaryForAll(f))
}

func (a *Transposed[T]) SetBinaryForAll(b, c M[T], f func(T, T) T) MuM[T] {
	return SetBinaryForAll[T](a, b, c, f)
}

func (a *Transposed[T]) SetUnaryForAll(b M[T], f func(T) T) MuM[T] {
	return SetUnaryForAll[T](a, b, f)
}

func (a *Transposed[T]) SetCopy(b M[T]) MuM[T] {
	return SetCopy[T](a, b)
}

func (a *Transposed[T]) Equals(b M[T]) bool {
	return EqualsM[T](a.m, Transpose[T](b))
}

func (a *Transposed[T]) Transpose() M[T] {
	return a.m
}

func (a *Transposed[T]) One() M[T] {
	return Transpose(a.m.One()) // TODO do better
}

func (a *Transposed[T]) Zero() M[T] {
	return Transpose(a.m.Zero()) // TODO do better
}

func (a *Transposed[T]) MuZero() MuM[T] {
	return MuTranspose[T](a.m.(MuM[T]).MuZero()) // Note the dynamic check, this can fail.
}

func (a *Transposed[T]) MuOne() MuM[T] {
	return MuTranspose[T](a.m.(MuM[T]).MuOne()) // Note the dynamic check, this can fail.
}

func (a *Transposed[T]) Dims() (rows int, cols int) {
	cols, rows = a.m.Dims()
	return
}

func (a *Transposed[T]) Row(i int) V[T] {
	return a.m.Col(i)
}

func (a *Transposed[T]) Col(i int) V[T] {
	return a.m.Row(i)
}

func (a *Transposed[T]) At(i, j int) T {
	return a.m.At(j, i)
}

func (a *Transposed[T]) Set(i, j int, v T) {
	a.m.(MuM[T]).Set(j, i, v)
}

// ROW MAJOR / CONTIGUOUS ROW

type ContiguousRowMatrix[T Field[T]] struct {
	Default[T]
	x          []T
	rows, cols int
}

func RowMajor[T Field[T]](rows, cols int) *ContiguousRowMatrix[T] {
	prod := rows * cols // TODO check for overflow
	x := &ContiguousRowMatrix[T]{cols: cols, rows: rows, x: make([]T, prod, prod)}
	x.self = x
	return x
}

func (a *ContiguousRowMatrix[T]) BinaryForAll(b M[T], f func(T, T) T) M[T] {
	assertSameDims[T](a, b)
	r := RowMajor[T](a.rows, a.cols)
	return r.SetBinaryForAll(a, b, f)
}

func (a *ContiguousRowMatrix[T]) UnaryForAll(f func(T) T) M[T] {
	r := RowMajor[T](a.rows, a.cols)
	return r.SetUnaryForAll(a, f)
}

func (a *ContiguousRowMatrix[T]) One() M[T] {
	return a.MuOne()
}

func (a *ContiguousRowMatrix[T]) MuOne() MuM[T] {
	b := RowMajor[T](a.rows, a.cols)
	var z T
	one := z.One()
	for i := 0; i < len(a.x); i += a.cols + 1 {
		a.x[i] = one
	}
	return b
}

func (a *ContiguousRowMatrix[T]) MuZero() MuM[T] {
	return RowMajor[T](a.rows, a.cols)
}

func (a *ContiguousRowMatrix[T]) Zero() M[T] {
	return RowMajor[T](a.rows, a.cols)
}

func (a *ContiguousRowMatrix[T]) Dims() (rows int, cols int) {
	return a.rows, a.cols
}

func (a *ContiguousRowMatrix[T]) Row(i int) V[T] {
	start := i * a.cols
	return &ContiguousVector[T]{a.x[start : start+a.cols]}
}

func (a *ContiguousRowMatrix[T]) Col(i int) V[T] {
	return &StridedVector[T]{v: a.x[i:], stride: a.cols, len: a.rows}
}

func (a *ContiguousRowMatrix[T]) At(i, j int) T {
	CheckBounds[T](a, i, j)
	return a.x[i*a.cols+j]
}

// Mutable matrix methods

func (a *ContiguousRowMatrix[T]) Set(i, j int, v T) {
	CheckBounds[T](a, i, j)
	a.x[i*a.cols+j] = v
}

// COLUMN MAJOR / CONTIGUOUS COLUMN

type ContiguousColumnMatrix[T Field[T]] struct {
	Default[T]
	x          []T
	rows, cols int
}

func ColumnMajor[T Field[T]](rows, cols int) *ContiguousColumnMatrix[T] {
	prod := rows * cols // TODO check for overflow
	x := &ContiguousColumnMatrix[T]{cols: cols, rows: rows, x: make([]T, prod, prod)}
	x.self = x
	return x
}

func (a *ContiguousColumnMatrix[T]) BinaryForAll(b M[T], f func(T, T) T) M[T] {
	assertSameDims[T](a, b)
	r := ColumnMajor[T](a.rows, a.cols)
	return r.SetBinaryForAll(a, b, f)
}

func (a *ContiguousColumnMatrix[T]) UnaryForAll(f func(T) T) M[T] {
	r := ColumnMajor[T](a.rows, a.cols)
	return r.SetUnaryForAll(a, f)
}

func (a *ContiguousColumnMatrix[T]) One() M[T] {
	return a.MuOne()
}
func (a *ContiguousColumnMatrix[T]) MuOne() MuM[T] {
	b := ColumnMajor[T](a.rows, a.cols)
	var z T
	one := z.One()
	for i := 0; i < len(a.x); i += a.rows + 1 {
		a.x[i] = one
	}
	return b
}

func (a *ContiguousColumnMatrix[T]) MuZero() MuM[T] {
	return ColumnMajor[T](a.rows, a.cols)
}

func (a *ContiguousColumnMatrix[T]) Zero() M[T] {
	return ColumnMajor[T](a.rows, a.cols)
}

func (a *ContiguousColumnMatrix[T]) Dims() (rows int, cols int) {
	return a.rows, a.cols
}

func (a *ContiguousColumnMatrix[T]) Row(i int) V[T] { // fixme
	r, c := a.Dims()
	return &StridedVector[T]{v: a.x[i:], stride: r, len: c}
}

func (a *ContiguousColumnMatrix[T]) Col(i int) V[T] { // fixme
	start := i * a.rows
	return &ContiguousVector[T]{a.x[start : start+a.rows]}
}

func (a *ContiguousColumnMatrix[T]) At(i, j int) T {
	CheckBounds[T](a, i, j)
	return a.x[i+j*a.rows]
}

// Mutable matrix methods

func (a *ContiguousColumnMatrix[T]) Set(i, j int, v T) {
	CheckBounds[T](a, i, j)
	a.x[i+j*a.rows] = v
}
