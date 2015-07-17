package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

type FetchStats func(int, string) Stats

var messages = make(chan string)
var season_path string
var season_name = "2014-15"
const last_season_name = "2014-15"
const num_games = 1230 // There are 1230 NBA regular season games per year
const players_file = "players.gob"

func main() {
	data_path := "data"
	
	clArgs := os.Args[1:]
	if len(clArgs) > 0 {
		if clArgs[0] == "collate" {
			combineStats(data_path)
			return
		} 
		season_name = clArgs[0]
	}
	
	mkdirIfNotExists(data_path)
	results_path := fmt.Sprintf("%s/results", data_path)
	mkdirIfNotExists(results_path)
	season_path = fmt.Sprintf("%s/%s", data_path, season_name)
	mkdirIfNotExists(season_path)
	shots_path := fmt.Sprintf("%s/stats", season_path)
	mkdirIfNotExists(shots_path)
	
	var players []Player
	players_file_name := fmt.Sprintf("%s/%s", season_path, players_file)
	if !exists(players_file_name) {
		players = fetchPlayers()
		err := saveToDisk(players, players_file_name)
		if err != nil { panic(err) }
	} else {
		loadFromDisk(&players, players_file_name)
	}

	if season_name != last_season_name {
		var recent_players []Player
		last_season_path := fmt.Sprintf("%s/%s", data_path, last_season_name)
		mkdirIfNotExists(last_season_path)
		players_file_name := fmt.Sprintf("%s/%s", last_season_path, players_file)
		if !exists(players_file_name) {
			recent_players = fetchPlayers()
			err := saveToDisk(recent_players, players_file_name)
			if err != nil { panic(err) }
		} else {
			loadFromDisk(&recent_players, players_file_name)
		}
		players = append(players, recent_players...)
	}
	fmt.Printf("Fetching shot statistics for %d players\n", len(players))
	
	var num_procs = runtime.NumCPU()
	fmt.Printf("Processing %d players on %d processes\n", len(players), num_procs)
	// Set the max number of processes to the number of CPUs on this machine
	runtime.GOMAXPROCS(num_procs)

	in := gen(players)		
	var channels = make([]<-chan Stats, 0, 4)
	for i := 0; i < num_procs; i++ {
		channels = append(channels, fetchStatsForPlayers(in, fetchShots))
	}

	out_file_name := fmt.Sprintf("%s/%s.csv", results_path, season_name)
	out, err := os.Create(out_file_name)
	if err != nil { panic(err) }
	defer out.Close()
	for n := range merge(&channels) {
		writeStatToFile(n, out)
	}
		
	fmt.Printf("\n")
}

func combineStats(data_path string) {
	seasons, _ := ioutil.ReadDir(data_path)
	player_map := make(map[string]Stats)

	out_file_name := fmt.Sprintf("%s/results/total.csv", data_path)
	out, err := os.Create(out_file_name)
	if err != nil { panic(err) }
	var stats Stats
	for _, f := range seasons {
		if f.IsDir() && f.Name() != "results" {
			dir := fmt.Sprintf("%s/%s/stats", data_path, f.Name())
			files, _ := ioutil.ReadDir(dir)
			for _, p_file := range files {
				file_name := fmt.Sprintf("%s/%s", dir, p_file.Name())
				zero(&stats)
				err := loadFromDisk(&stats, file_name)
				if err != nil { panic(err) }
				name := stats.Name
				player, ok := player_map[name]
				if ok {
					stats = add(stats, player)
				}
				player_map[name] = stats
			}
		}
	}
	var i int = 0
	for _, stats := range player_map {
		i++
		writeStatToFile(stats, out)
	}
	fmt.Printf("Processed %d players\n", i)
}

func gen(players []Player) <-chan Player {
	out := make(chan Player)
	go func() {
		for _, player := range players {
			out <- player
		}
		close(out)
	}()
	return out
}

func merge(channels *[]<-chan Stats) <-chan Stats {
	var wg sync.WaitGroup
	out := make(chan Stats)
	output := func(c <- chan Stats) {
		for s := range c {
			out <- s
		}
		wg.Done()
	}
	wg.Add(len(*channels))
	for _, c := range *channels {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func printProcessId(id int, name string) Stats {
	fmt.Printf("Process ID %d for %s \n", os.Getpid(), name)
	return Stats {Name: name}
}

func fetchStatsForPlayers(in <-chan Player, fn FetchStats) <-chan Stats {
	out := make(chan Stats)
	go func() {
		defer close(out)
		for player := range in {
			stats := fn(player.ID, player.Name)
			if stats.Attempts > 0 {
				out <- stats
			}
		}
	}()
	return out
}

// exists returns whether the given file or directory exists or not
func exists(f string) bool {
	_, err := os.Stat(f)
	if err == nil { return true }
	if os.IsNotExist(err) { return false }
	return true
}

func mkdirIfNotExists(dir string) {
	if !exists(dir) {

		err := os.Mkdir(dir, 0755)
		if err != nil { panic(err) }
	}
}

func saveToDisk(data interface{}, file_name string) error {
	out_file, err := os.Create(file_name)
	if err != nil { return err } else {
		defer out_file.Close()
		dataEncoder := gob.NewEncoder(out_file)
		return dataEncoder.Encode(data)
	}
	return nil
}

func loadFromDisk(target interface{}, file_name string) error {
	in_file, err := os.Open(file_name)
	if err != nil { return nil } else {
		defer in_file.Close()
		dataDecoder := gob.NewDecoder(in_file)
		return dataDecoder.Decode(target)
	}
	return nil
}
