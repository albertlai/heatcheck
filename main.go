package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
)

type process func(int) error

var messages = make(chan string)
var data_path string
const num_games = 1230 // There are 1230 NBA regular season games per year
const players_file = "players.gob"

func main() {	
	args := os.Args[1:]
	data_path = "data"

	players_file_name := fmt.Sprintf("%s/%s", data_path, players_file)
	if !exists(players_file_name) {
		err := saveToDisk(fetchPlayers(), players_file_name)
		if err != nil { panic(err) }
	}
	var fn process
	if len(args) > 0 {
		if len(args) > 1 {
			data_path = args[1]
		}
		switch args[0] {
		case "fetch":
			makeDataPath(json_path)
			fn = fetchGame
		case "process":
			makeDataPath(play_path)
			fn = processGameJSON
		case "map":
			makeDataPath(map_path)
			fn = mapPlays
		case "reduce":
			reducePlays()
			return
		default: fn = dummy
		}
	} else {
		fmt.Printf("Usage is go run *.go [action] [outfile]\n")
		return
	}
	var num_blocks = runtime.NumCPU()
	// Set the max number of processes to the number of CPUs on this machine
	runtime.GOMAXPROCS(num_blocks)
	var block = num_games / num_blocks
	var start int
	var end int
	for i := 0; i < num_blocks; i++ {
		start = i * block + 1
		if i == num_blocks - 1 {
			end = num_games + 1
		} else {
			end = start + block
		}
		go runInRange(start, end, fn)
	}
	for i := 0; i < num_blocks; i++ {
		<- messages
	}
	fmt.Printf("\n")
}

func runInRange(start int, finish int, fn process) {
	for i := start; i < finish; i++ {
		err := fn(i)
		if err != nil {
			fmt.Printf("Failed to run function on game %d\n", i)
		}
	}
	messages <- fmt.Sprintf("\nFetched games %d to %d\n", start, finish-1)
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
	fmt.Printf("there are %d dense turds\n", id)
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
