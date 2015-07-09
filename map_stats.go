package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stats struct {
	jump_shots, attempts int
	distance, defender float64
	jump_shots_1, attempts_1 int
	distance_1, defender_1 float64
	jump_shots_2, attempts_2 int
	distance_2, defender_2 float64
}

const map_path = "map"

func mapPlays(game_id int) error {
	fmt.Printf(" %d ", game_id)
	in_name := fmt.Sprintf("%s/%s/plays_%d.json", data_path, play_path, game_id)
	in, err := os.Open(in_name)
	if err != nil { return err }
	defer in.Close()

	out_name := fmt.Sprintf("%s/%s/map_%d.json", data_path, map_path, game_id)
	out, err := os.Create(out_name)
	if err != nil { return err }
	defer out.Close()

	var stats_map map[string]*Stats = make(map[string]*Stats)
	var made_last_1 map[string]bool = make(map[string]bool)
	var made_last_2 map[string]bool = make(map[string]bool)
	
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Jump Shot") || strings.Contains(line, "Jump Bank Shot") {
			index := strings.IndexRune(line, ':')
			name := line[:index-1]			
			stats, exists := stats_map[name]
			if !exists {
				stats = &Stats {}
				stats_map[name] = stats
			}
			stats.attempts += 1

			i_foot_end := strings.LastIndex(line, "' ")
			var distance float64
			if i_foot_end < 0 {
				if strings.Contains(line, "3PT") {
					distance = 22
				} else {
					fmt.Printf("Error getting distance from %s\n", line)
					continue
				}
			} else {
				i_foot_start := strings.LastIndex(line[:i_foot_end], " ") + 1
				distance, err = strconv.ParseFloat(line[i_foot_start:i_foot_end], 64)
				if err != nil {
					fmt.Printf("Error getting distance from %s\n", line)
					continue
				}
			}
			stats.distance = (stats.distance * float64(stats.attempts - 1) + distance) / float64(stats.attempts)
			var adj int
			if strings.Contains(line, "MISS") {
				adj = 0
			} else {
				adj = 1
			}
			stats.jump_shots += adj 
			if made_last_1[name] {
				stats.attempts_1 += 1
				stats.jump_shots_1 += adj
				stats.distance_1 = (stats.distance_1 * float64(stats.attempts_1 - 1) + distance) / float64(stats.attempts_1)
				if made_last_2[name] {
					stats.attempts_2 += 1
					stats.jump_shots_2 += adj
					stats.distance_2 = (stats.distance_2 * float64(stats.attempts_2 - 1) + distance) / float64(stats.attempts_2)
				}
			}

			made_last_2[name] = made_last_1[name]
			made_last_1[name] = adj == 1
		}
	}
	writeStatsMapToFile(stats_map, out)
	return nil
}

func writeStatsMapToFile(stats_map map[string]*Stats, out *os.File) {
	for name, stats := range stats_map {
		line := fmt.Sprintf("%s, %d, %d, %f, %d, %d, %f, %d, %d, %f\n",
			name, stats.jump_shots, stats.attempts, stats.distance,
			stats.jump_shots_1, stats.attempts_1, stats.distance_1,
			stats.jump_shots_2, stats.attempts_2, stats.distance_2)
		out.WriteString(line)
	}
}
