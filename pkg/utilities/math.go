package utilities

import "math"

func RoundUp(val float64) float64 {
	return math.Round(val*100) / 100
}
