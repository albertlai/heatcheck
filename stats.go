package main

import (
	"fmt"
	"os"
)
type Stats struct {
	Name string
	JumpShots, Attempts int
	Distance, Defender float64
	JumpShots1, Attempts1 int
	Distance1, Defender1 float64
	JumpShots2, Attempts2 int
	Distance2, Defender2 float64
}

func writeStatToFile(stats Stats, out *os.File) {
	line := fmt.Sprintf("%s, %d, %d, %f, %f, %d, %d, %f, %f, %d, %d, %f, %f\n",
		stats.Name, stats.JumpShots, stats.Attempts, stats.Distance, stats.Defender,
		stats.JumpShots1, stats.Attempts1, stats.Distance1, stats.Defender1,
		stats.JumpShots2, stats.Attempts2, stats.Distance2, stats.Defender2)
		out.WriteString(line)
}
