package game

import (
	"math"
)

func levelFromXp(xp int) int {
	const base float64 = 0.9
	level := base * math.Cbrt(float64(xp))
	return int(level)
}
