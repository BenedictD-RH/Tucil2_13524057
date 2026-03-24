package matrix

import "fmt"

type Matrix struct {
	Buffer [][]float32
	rows   int
	cols   int
}

func NewMatrix(row, col int) *Matrix {
	mat := &Matrix{[][]float32{}, row, col}

	for i := 0; i < row; i++ {
		rowbuff := []float32{}
		for j := 0; j < col; j++ {
			rowbuff = append(rowbuff, 0)
		}
		mat.Buffer = append(mat.Buffer, rowbuff)
	}
	mat.rows = row
	mat.cols = col
	return mat
}

func AddMatrix(m1, m2 *Matrix) *Matrix {
	if m1.rows != m2.rows || m1.cols != m2.cols {
		return NewMatrix(0, 0)
	} else {
		m3 := NewMatrix(m1.rows, m2.cols)
		for i := 0; i < m1.rows; i++ {
			for j := 0; j < m2.cols; j++ {
				m3.Buffer[i][j] = m1.Buffer[i][j] + m2.Buffer[i][j]
			}
		}
		return m3
	}
}

func SubtractMatrix(m1, m2 *Matrix) *Matrix {
	if m1.rows != m2.rows || m1.cols != m2.cols {
		return NewMatrix(0, 0)
	} else {
		m3 := NewMatrix(m1.rows, m2.cols)
		for i := 0; i < m1.rows; i++ {
			for j := 0; j < m2.cols; j++ {
				m3.Buffer[i][j] = m1.Buffer[i][j] - m2.Buffer[i][j]
			}
		}
		return m3
	}
}

func MultiplyMatrix(m1, m2 *Matrix) *Matrix {
	if m1.cols != m2.rows {
		fmt.Println(m1.cols, "!=", m2.rows)
		return NewMatrix(0, 0)
	} else {
		m3 := NewMatrix(m1.rows, m2.cols)
		for i := 0; i < m1.rows; i++ {
			for j := 0; j < m2.cols; j++ {
				for k := 0; k < m1.cols; k++ {
					m3.Buffer[i][j] += m1.Buffer[i][k] * m2.Buffer[k][j]
				}
			}
		}
		return m3
	}
}

func (m *Matrix) Transpose() (*Matrix){
	mNew := NewMatrix(m.cols, m.rows)
	for i := 0; i < mNew.rows; i++ {
		for j := 0; j < mNew.cols; j++ {
			mNew.Buffer[i][j] = m.Buffer[j][i]
		}
	}
	return mNew
}

func Identity(n int) (*Matrix) {
	m := NewMatrix(n,n)
	for i:= 0; i < n; i++ {
		m.Buffer[i][i] = 1
	}
	return m
}

func (m Matrix) String() string {
	s := ""
	for i := 0; i < m.rows; i++ {
		s += "| "
		for j := 0; j < m.cols; j++ {
			s += fmt.Sprintf("%f ", m.Buffer[i][j])
		}
		s += "|\n"
	}
	return s
}
