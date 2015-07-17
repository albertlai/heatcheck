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

func zero(s *Stats) {
	s.JumpShots = 0
	s.Attempts = 0
	s.Distance = 0
	s.Defender = 0
	s.DistanceSD = 0
	s.DefenderSD = 0
	s.JumpShots1 = 0
	s.Attempts1 = 0
	s.Distance1 = 0
	s.Defender1 = 0
	s.DistanceSD1 = 0
	s.DefenderSD1 = 0
	s.JumpShots2 = 0
	s.Attempts2 = 0
	s.Distance2 = 0
	s.Defender2 = 0
	s.DistanceSD2 = 0
	s.DefenderSD2 = 0
}

func add(s1 Stats, s2 Stats) Stats {
	n := float64(s1.Attempts)
	n1 := float64(s1.Attempts1)
	n2 := float64(s1.Attempts2)
	m := float64(s2.Attempts)
	m1 := float64(s2.Attempts1)
	m2 := float64(s2.Attempts2)	
	out := Stats{
		Name: s1.Name,
		JumpShots: s1.JumpShots + s2.JumpShots,
		Attempts: s1.Attempts + s2.Attempts,
		Distance: combineAvg(s1.Distance, n, s2.Distance, m),
		DistanceSD: combineSD(s1.DistanceSD, s1.Distance, n,
			s2.DistanceSD, s2.Distance, m),
		Defender: combineAvg(s1.Defender, n, s2.Defender, m),
		DefenderSD: combineSD(s1.DefenderSD, s1.Defender, n,
			s2.DefenderSD, s2.Defender, m),
		JumpShots1: s1.JumpShots1 + s2.JumpShots1,
		Attempts1: s1.Attempts1 + s2.Attempts1,
		Distance1: combineAvg(s1.Distance1, n1, s2.Distance1, m1),
		DistanceSD1: combineSD(s1.DistanceSD1, s1.Distance1, n,
			s2.DistanceSD1, s2.Distance1, m1),
		Defender1: combineAvg(s1.Defender1, n1, s2.Defender1, m1),
		DefenderSD1: combineSD(s1.DefenderSD1, s1.Defender1, n1,
			s2.DefenderSD1, s2.Defender1, m2),
		JumpShots2: s1.JumpShots2 + s2.JumpShots2,
		Attempts2: s1.Attempts2 + s2.Attempts2,
		Distance2: combineAvg(s1.Distance2, n2, s2.Distance2, m2),
		DistanceSD2: combineSD(s1.DistanceSD2, s1.Distance2, n2,
			s2.DistanceSD2, s2.Distance2, m2),
		Defender2: combineAvg(s1.Defender2, n2, s2.Defender2, m2),
		DefenderSD2: combineSD(s1.DefenderSD2, s1.Defender2, n2,
			s2.DefenderSD2, s2.Defender2, m2),
	}
	return out
	
}

func combineAvg(Ex1 float64, n1 float64, Ex2 float64, n2 float64) float64 {
	return (Ex1 * n1 + Ex2 * n2) / (n1 + n2)
}

func combineSD(SD1 float64, Ex1 float64, n1 float64, SD2 float64, Ex2 float64, n2 float64) float64{
	Sx_squared := (SD1 * SD1 + Ex1 * Ex1) * n1 + (SD2 * SD2 + Ex2 * Ex2) * n2
	Eout := combineAvg(Ex1, n1, Ex2, n2)
	return sd(Sx_squared, Eout, n1 + n2)
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
		fg2 = float64(stats.JumpShots2) / n2
		Dstderr2 := stats.DistanceSD / math.Sqrt(n2)
		Dz2 = (stats.Distance2 - stats.Distance) / Dstderr2
		Defstderr2 := stats.DefenderSD / math.Sqrt(n2)
		Defz2 = (stats.Defender2 - stats.Defender) / Defstderr2		
	}

	line := fmt.Sprintf("%s, %d, %d, %f, %f, %f, %f, %f, %d, %d, %f, %f, %f, %f, %f, %d, %d, %f, %f, %f, %f, %f\n", stats.Name,
		stats.JumpShots, stats.Attempts, fg, stats.Distance, stats.DistanceSD, stats.Defender, stats.DefenderSD,
		stats.JumpShots1, stats.Attempts1, fg1, stats.Distance1, Dz1, stats.Defender1, Defz1,
		stats.JumpShots2, stats.Attempts2, fg2, stats.Distance2, Dz2, stats.Defender2, Defz2)
		out.WriteString(line)
}
