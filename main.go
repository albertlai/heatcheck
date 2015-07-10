package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"sync"
)

type FetchStats func(int, string) Stats

var messages = make(chan string)
var data_path string = "data"
var season_name = "2014-15"
const num_games = 1230 // There are 1230 NBA regular season games per year
const players_file = "players.gob"

func main() {
	players_file_name := fmt.Sprintf("%s/%s", data_path, players_file)
	var players []Player
	if !exists(players_file_name) {
		players = fetchPlayers()
		err := saveToDisk(players, players_file_name)
		if err != nil { panic(err) }
	} else {
		loadFromDisk(&players, players_file_name)
	}
	
	var num_procs = runtime.NumCPU()
	fmt.Printf("Processing %d players on %d processes\n", len(players), num_procs)
	// Set the max number of processes to the number of CPUs on this machine
	runtime.GOMAXPROCS(num_procs)

	in := gen(players)		
	var channels = make([]<-chan Stats, 0, 4)
	for i := 0; i < num_procs; i++ {
		channels = append(channels, fetchStatsForPlayers(in, fetchShots))
	}

	out_file_name := fmt.Sprintf("%s/shots.csv", data_path)
	out, err := os.Create(out_file_name)
	if err != nil { panic(err) }
	for n := range merge(&channels) {
		writeStatToFile(n, out)			
	}
		
	fmt.Printf("\n")
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
			out <- fn(player.ID, player.Name)
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

func makeDataPath(path string) {
	// Make the output directory if it doesn't exist
	full_path := fmt.Sprintf("%s/%s", data_path, path)
	if !exists(full_path) {
		err := os.Mkdir(full_path, 0666)
		if err != nil {	panic(err) }
	}
}

func dummy(id int) error {
//	fmt.Printf("there are %d dense turds\n", id)
	return nil
}


func saveToDisk(data interface{}, file_name string) error {
	outfile, err := os.Create(file_name)
	if err != nil { return err } else {
		defer outfile.Close()
		dataEncoder := gob.NewEncoder(outfile)
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
