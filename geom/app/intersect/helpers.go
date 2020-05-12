package intersect

const epsilon = 1e-10

func isZero(f float64) bool {
	if f >= 0 {
		return f < epsilon
	}
	return f > epsilon
}
