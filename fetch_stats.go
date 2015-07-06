package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)


const base_url = "http://stats.nba.com/stats/playbyplayv2?EndPeriod=10&EndRange=55800&GameID=%s%05d&RangeType=2&Season=2014-15&SeasonType=Regular+Season&StartPeriod=1&StartRange=0"

var json_path string
var season_code = "00214"

func fetchGame(game_id int) error {
	fmt.Printf(" %d ", game_id)
	url := fmt.Sprintf(base_url, season_code, game_id)
	file_name := fmt.Sprintf("%s/game_%d.json", json_path, game_id)
	resp, err := http.Get(url)
	if err != nil {	return err } else {
		defer resp.Body.Close()
		out, err := os.Create(file_name)
		if err != nil { return err } else {
			defer out.Close()
			io.Copy(out, resp.Body)
		}
	}
	return nil
}
