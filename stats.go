package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
)

type Histogram map[string] [2]int;

type Stats struct {
	Name string
	JumpShots, Attempts int	
	Distance, Defender, DistanceSD, DefenderSD float64
	DistanceMap Histogram
	JumpShots1, Attempts1 int
	Distance1, Defender1, DistanceSD1, DefenderSD1 float64
	DistanceMap1 Histogram
	JumpShots2, Attempts2 int
	Distance2, Defender2, DistanceSD2, DefenderSD2  float64
	DistanceMap2 Histogram
}

func zero(s *Stats) {
	s.JumpShots = 0
	s.Attempts = 0
	s.Distance = 0
	s.DistanceMap = make(Histogram)
	s.Defender = 0
	s.DistanceSD = 0
	s.DefenderSD = 0
	s.JumpShots1 = 0
	s.Attempts1 = 0
	s.Distance1 = 0
	s.DistanceMap1 = make(Histogram)
	s.Defender1 = 0
	s.DistanceSD1 = 0
	s.DefenderSD1 = 0
	s.JumpShots2 = 0
	s.Attempts2 = 0
	s.Distance2 = 0
	s.DistanceMap2 = make(Histogram)
	s.Defender2 = 0
	s.DistanceSD2 = 0
	s.DefenderSD2 = 0
}

func (hist Histogram) AddDistance(d int, made int) {
	distance  := strconv.Itoa(d)	
	old, exists := hist[distance]		
	if !exists {
		old = [2]int { 0, 0}
	}
	hist[distance] = [2]int { old[0] + made, old[1] + 1 }
}

func combineHistograms(h1 Histogram, h2 Histogram) Histogram {
	out := make(Histogram)
	for dist, val := range h1 {
		out[dist] = [2]int { val[0], val[1] }
	}
	for dist, val := range h2 {
		v, exists := out[dist]
		if !exists {
			v = [2]int { 0, 0 }
		} 
		out[dist] = [2]int { val[0] + v[0], val[1] + v[1] }
	}
	return out
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
		DistanceMap: combineHistograms(s1.DistanceMap, s2.DistanceMap),
		JumpShots1: s1.JumpShots1 + s2.JumpShots1,
		Attempts1: s1.Attempts1 + s2.Attempts1,
		Distance1: combineAvg(s1.Distance1, n1, s2.Distance1, m1),
		DistanceSD1: combineSD(s1.DistanceSD1, s1.Distance1, n,
			s2.DistanceSD1, s2.Distance1, m1),
		Defender1: combineAvg(s1.Defender1, n1, s2.Defender1, m1),
		DefenderSD1: combineSD(s1.DefenderSD1, s1.Defender1, n1,
			s2.DefenderSD1, s2.Defender1, m2),
		DistanceMap1: combineHistograms(s1.DistanceMap1, s2.DistanceMap1),
		JumpShots2: s1.JumpShots2 + s2.JumpShots2,
		Attempts2: s1.Attempts2 + s2.Attempts2,
		Distance2: combineAvg(s1.Distance2, n2, s2.Distance2, m2),
		DistanceSD2: combineSD(s1.DistanceSD2, s1.Distance2, n2,
			s2.DistanceSD2, s2.Distance2, m2),
		Defender2: combineAvg(s1.Defender2, n2, s2.Defender2, m2),
		DefenderSD2: combineSD(s1.DefenderSD2, s1.Defender2, n2,
			s2.DefenderSD2, s2.Defender2, m2),
		DistanceMap2: combineHistograms(s1.DistanceMap2, s2.DistanceMap2),
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
	jsonString, err := json.Marshal(stats.DistanceMap)
	if err != nil { fmt.Println(err) }
	jsonString1, err := json.Marshal(stats.DistanceMap1)
	if err != nil { fmt.Println(err) }
	jsonString2, err := json.Marshal(stats.DistanceMap2)
	if err != nil { fmt.Println(err) }
	line := fmt.Sprintf("%s\t%d\t%d\t%f\t%f\t%f\t%f\t%f\t%d\t%d\t%f\t%f\t%f\t%f\t%f\t%d\t%d\t%f\t%f\t%f\t%f\t%f\t%s\t%s\t%s\n", stats.Name,
		stats.JumpShots, stats.Attempts, fg, stats.Distance, stats.DistanceSD, stats.Defender, stats.DefenderSD,
		stats.JumpShots1, stats.Attempts1, fg1, stats.Distance1, Dz1, stats.Defender1, Defz1,
		stats.JumpShots2, stats.Attempts2, fg2, stats.Distance2, Dz2, stats.Defender2, Defz2,
		jsonString, jsonString1, jsonString2)
		out.WriteString(line)
}
