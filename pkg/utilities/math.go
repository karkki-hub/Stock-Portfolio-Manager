package utilities

import "math"

func RoundUp(val float64) float64 { // Round a float to 2 decimal places
	return math.Round(val*100) / 100
}
