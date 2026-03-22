package matrix

import "math"

func DegreeToRad(deg float32) (float32) {
	return math.Pi*(deg/180)
}

func RotationMatrix_X(rad float32) (*Matrix) {
	m := NewMatrix(3,3)
	m.Buffer[0][0] = 1
	m.Buffer[1][1] = float32(math.Cos(float64(rad)))
	m.Buffer[1][2] = float32(math.Sin(float64(rad)))*-1
	m.Buffer[2][1] = float32(math.Sin(float64(rad)))
	m.Buffer[2][2] = float32(math.Cos(float64(rad)))
	return m
}

func RotationMatrix_Y(rad float32) (*Matrix) {
	m := NewMatrix(3,3)
	m.Buffer[0][0] = float32(math.Cos(float64(rad)))
	m.Buffer[0][2] = float32(math.Sin(float64(rad)))
	m.Buffer[1][1] = 1
	m.Buffer[2][0] = float32(math.Sin(float64(rad)))*-1
	m.Buffer[2][2] = float32(math.Cos(float64(rad)))
	return m
}

func RotationMatrix_Z(rad float32) (*Matrix) {
	m := NewMatrix(3,3)
	m.Buffer[0][0] = float32(math.Cos(float64(rad)))
	m.Buffer[0][1] = float32(math.Sin(float64(rad)))*-1
	m.Buffer[1][0] = float32(math.Sin(float64(rad)))
	m.Buffer[1][1] = float32(math.Cos(float64(rad)))
	m.Buffer[2][2] = 1
	return m
}

func RotationMatrixOnAxis(ux, uy float64, rad float32) (*Matrix) {
	m := NewMatrix(3,3)
	m.Buffer[0][0] = float32(ux*ux*(1 - math.Cos(float64(rad))) + math.Cos(float64(rad)))
	m.Buffer[0][1] = float32(ux*uy*(1 - math.Cos(float64(rad))))
	m.Buffer[0][2] = float32(uy*math.Sin(float64(rad)))
	m.Buffer[1][0] = float32(ux*uy*(1 - math.Cos(float64(rad))))
	m.Buffer[1][1] = float32(uy*uy*(1 - math.Cos(float64(rad))) + math.Cos(float64(rad)))
	m.Buffer[1][2] = float32(ux*math.Sin(float64(rad)))*-1
	m.Buffer[2][0] = float32(uy*math.Sin(float64(rad)))*-1
	m.Buffer[2][1] = float32(ux*math.Sin(float64(rad)))
	m.Buffer[2][2] = float32(math.Cos(float64(rad)))
	return m
}

func RotateMatrixByDrag2D(dx, dy float64, deg float32, r_Mat *Matrix) (*Matrix) {
	vec_mag := math.Sqrt(dx*dx + dy*dy)
	ux := dx/vec_mag
	uy := dy/vec_mag
	rad := DegreeToRad(deg)
	return MultiplyMatrix(RotationMatrixOnAxis(ux,uy,rad), r_Mat)
}