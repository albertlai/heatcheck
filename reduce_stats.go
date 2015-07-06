package main

import (
	"bufio"
	"io/ioutil"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func reducePlays() {
	path := "map"
	files, err := ioutil.ReadDir(path)
	if err != nil { panic(err) }
	var stats_map map[string]*Stats = make(map[string]*Stats)
	for i := 0; i < len(files); i++ {
		file_info := files[i]
		in_name := fmt.Sprintf("%s/%s", path, file_info.Name())
		in, err := os.Open(in_name)
		fmt.Printf("opening %s\n", in_name)
		if err != nil {
			fmt.Printf("Failed to open %s\n", in_name)
			panic(err)
			continue
		}

		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			line := scanner.Text()
			tokens := strings.Split(line, ", ")
			name := tokens[0]
			
			js, err := strconv.Atoi(tokens[1])
			jsa, err := strconv.Atoi(tokens[2])
			d, err := strconv.ParseFloat(tokens[3], 64)
			js1, err := strconv.Atoi(tokens[4])
			jsa1, err := strconv.Atoi(tokens[5])
			d1, err := strconv.ParseFloat(tokens[6], 64)
			js2, err := strconv.Atoi(tokens[7])
			jsa2, err := strconv.Atoi(tokens[8])
			d2, err := strconv.ParseFloat(tokens[9], 64)

			if err != nil { fmt.Printf("uh oh") }
			stats, exists := stats_map[name]
			if !exists {
				stats = &Stats {}
				stats_map[name] = stats
			}

			stats.distance = (stats.distance * float64(stats.attempts) + d * float64(jsa)) / float64(stats.attempts + int(jsa))
			stats.jump_shots += int(js)
			stats.attempts += int(jsa)
			if stats.attempts_1 + jsa1 > 0 {
				stats.distance_1 = (stats.distance_1 * float64(stats.attempts_1) + d1 * float64(jsa1)) / float64(stats.attempts_1 + int(jsa1))
			}
			stats.jump_shots_1 += int(js1)
			stats.attempts_1 += int(jsa1)
			if stats.attempts_2 + jsa2 > 0 {
				stats.distance_2 = (stats.distance_2 * float64(stats.attempts_2) + d2 * float64(jsa2)) / float64(stats.attempts_2 + int(jsa2))
			}
			stats.jump_shots_2 += int(js2)
			stats.attempts_2 += int(jsa2)
		}
	}

	out_name := "stats.csv"
	out, err := os.Create(out_name)
	if err != nil { panic(err) }
	defer out.Close()

	writeStatsMapToFile(stats_map, out)
}
