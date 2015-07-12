package main

import (
	"fmt"
	"math"
	"os"
)
type Stats struct {
	Name string
	JumpShots, Attempts int
	Distance, Defender, DistanceSD, DefenderSD float64
	JumpShots1, Attempts1 int
	Distance1, Defender1, DistanceSD1, DefenderSD1 float64
	JumpShots2, Attempts2 int
	Distance2, Defender2, DistanceSD2, DefenderSD2  float64
}

func writeStatToFile(stats Stats, out *os.File) {
	n := float64(stats.Attempts)
	n1 := float64(stats.Attempts1)
	n2 := float64(stats.Attempts2)

	var fg, fg1, fg2 float64
	var Dz1, Dz2, Defz1, Defz2 float64

	if n == 0 { fg = 0 } else {
		fg = float64(stats.JumpShots) / n
	}
	if n1 == 0 { fg1 = 0; Dz1 = 0; Defz1 = 0 } else {
		fg1 = float64(stats.JumpShots1) / n1
		Dstderr1 := stats.DistanceSD / math.Sqrt(n1)
		Dz1 = (stats.Distance1 - stats.Distance) / Dstderr1
		Defstderr1 := stats.DefenderSD / math.Sqrt(n1)
		Defz1 = (stats.Defender1 - stats.Defender) / Defstderr1
	}
	if n2 == 0 { fg2 = 0; Dz2 = 0; Defz2 = 0 } else {	
		fg2 := float64(stats.JumpShots2) / n2
		Defstderr2 := stats.DefenderSD / math.Sqrt(n2)
		Dz2 = (stats.Distance2 - stats.Distance) / Dstderr2
		Dstderr2 := stats.DistanceSD / math.Sqrt(n2)
		Defz2 = (stats.Defender2 - stats.Defender) / Defstderr2		
	}

	line := fmt.Sprintf("%s, %d, %d, %f, %f, %f, %f, %f, %d, %d, %f, %f, %f, %f, %f, %d, %d, %f, %f, %f, %f, %f\n", stats.Name,
		stats.JumpShots, stats.Attempts, fg, stats.Distance, stats.DistanceSD, stats.Defender, stats.DefenderSD,
		stats.JumpShots1, stats.Attempts1, fg1, stats.Distance1, Dz1, stats.Defender1, Defz1,
		stats.JumpShots2, stats.Attempts2, fg2, stats.Distance2, Dz2, stats.Defender2, Defz2)
		out.WriteString(line)
}
