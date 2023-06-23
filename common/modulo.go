package common

import "math"

func Modulo(x, y float64) float64 {
	// returns the modulo of x,y
	// golang's math.Mod is not suitable
	return x - y*math.Floor(x/y)
}
