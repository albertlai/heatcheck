package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const players_url = "http://stats.nba.com/stats/commonallplayers?IsOnlyCurrentSeason=1&LeagueID=00&Season=%s"// %s = 2014-15
const shots_url = "http://stats.nba.com/stats/playerdashptshotlog?LastNGames=0&LeagueID=00&Location=&Month=0&OpponentTeamID=0&Outcome=&Period=0&PlayerID=%d&Season=%s&SeasonType=Regular+Season&TeamID=0"// %d = 201935, %s = 2014-15

type RowProcessor func([]interface{})

type Player struct {
	ID int
	Name string
}

func processNBAResponse(url string, row_processor RowProcessor) error {
	resp, err := http.Get(url)
	fmt.Printf("Fetching %s\n", url)
	if err != nil { return err } else {
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		var data map[string]interface{}
		if err := dec.Decode(&data); err != nil { return err }
		results, exists := data["resultSets"].([]interface{})
		if exists && len(results) > 0 {
			rowSetContainer := results[0].(map[string]interface{})
			if rows, exists := rowSetContainer["rowSet"].([]interface{}); exists {
				fmt.Printf("rowset exists with %d rows\n\n", len(rows))
				for i := 0; i < len(rows); i++ {
					if row, exists := rows[i].([]interface{}); exists {
						row_processor(row)
					}
				}				
			}
		}
	}
	return nil
}

func fetchPlayers() []Player {
	fmt.Println("Fetching players")
	var players = make([]Player, 0, 400)
	process_player := func (row []interface{}) {
		id_raw, ok1 := row[0].(float64)
		if ok1 {
			id := int(id_raw)
			name, ok2 := row[1].(string)
			if ok2 {
				players = append(players, Player{ID:id, Name: name})
			}
		}
	}	
	url := fmt.Sprintf(players_url, season_name)
	err := processNBAResponse(url, process_player)
	if err != nil {	panic(err) }
	return players
}

// row =["GAME_ID","MATCHUP","LOCATION","W","FINAL_MARGIN","SHOT_NUMBER",
//  "PERIOD","GAME_CLOCK","SHOT_CLOCK","DRIBBLES","TOUCH_TIME","SHOT_DIST",
//  "PTS_TYPE","SHOT_RESULT","CLOSEST_DEFENDER","CLOSEST_DEFENDER_PLAYER_ID",
//  "CLOSE_DEF_DIST","FGM","PTS"]
func fetchShots(player_id int, name string) Stats {
	fmt.Printf("Fetching shots for %s on pid %d\n", name, os.Getpid())
	var stats Stats
	shots_file_name := fmt.Sprintf("%s/shots/%d.gob", data_path, player_id)
	if exists(shots_file_name) {
		loadFromDisk(&stats, shots_file_name)
		return stats
	}
	stats = Stats{Name: name}
	
	var current_game string
	var made_1, made_2 bool	
	process_shots := func (row []interface{}) {
		game_id := row[0].(string)
		if game_id != current_game {
			current_game = game_id
			made_1 = false
			made_2 = false
		}
		distance, ok1 := row[11].(float64)
		defender, ok2 := row[16].(float64)
		shot, ok3 := row[17].(float64)
		if ok1 && ok2 && ok3 && distance >= 5.5 {
			made := int(shot)	
			stats.Attempts += 1
			stats.JumpShots += made
			stats.Distance += distance
			stats.Defender += defender
			if made_1 {
				stats.Attempts1 += 1
				stats.JumpShots1 += made
				stats.Distance1 += distance
				stats.Defender1 += defender
				if made_2 {
					stats.Attempts2 += 1
					stats.JumpShots2 += made
					stats.Distance2 += distance
					stats.Defender2 += defender
				}
			}
			made_2 = made_1
			made_1 = made == 1
		}
	}
	url := fmt.Sprintf(shots_url, player_id, season_name)
	err := processNBAResponse(url, process_shots)
	if err != nil { panic(err) } else {
		stats.Defender = stats.Defender / float64(stats.Attempts)
		stats.Distance = stats.Distance / float64(stats.Attempts)
		stats.Defender1 = stats.Defender1 / float64(stats.Attempts1)
		stats.Distance1 = stats.Distance1 / float64(stats.Attempts1)
		stats.Defender2 = stats.Defender2 / float64(stats.Attempts2)
		stats.Distance2 = stats.Distance2 / float64(stats.Attempts2)
		saveToDisk(stats, shots_file_name)
		return stats
	}
}