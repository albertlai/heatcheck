package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

const players_url = "http://stats.nba.com/stats/commonallplayers?IsOnlyCurrentSeason=1&LeagueID=00&Season=%s"// %s = 2014-15
const shots_url = "http://stats.nba.com/stats/playerdashptshotlog?LastNGames=0&LeagueID=00&Location=&Month=0&OpponentTeamID=0&Outcome=&Period=0&PlayerID=%d&Season=%s&SeasonType=Regular+Season&TeamID=0"// %d = 201935, %s = 2014-15

type RowProcessor func([]interface{})

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
				fmt.Printf("rowset exists with %d rows\n", len(rows))
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

type Player struct {
	ID int
	Name string
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
func fetchShots(player_id int) Stats {
	fmt.Println("Fetching shots")	
	stats := Stats{}
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
		shot, ok2 := row[15].(float64)
		defender, ok3 := row[16].(float64)
		
		if ok1 && ok2 && ok3 {
			made := int(shot)	
			stats.attempts += 1
			stats.jump_shots += made
			stats.distance += distance
			stats.defender += defender
			if made_1 {
				stats.attempts_1 += 1
				stats.jump_shots_1 += made
				stats.distance_1 += distance
				stats.defender_1 += defender
				if made_2 {
					stats.attempts_2 += 1
					stats.jump_shots_2 += made
					stats.distance_2 += distance
					stats.defender_2 += defender
				}
			}
			made_2 = made_1
			made_1 = made == 1
		}
	}
	url := fmt.Sprintf(shots_url, season_name, player_id)
	err := processNBAResponse(url, process_shots)
	if err != nil { panic(err) } else {		
		return stats
	}
}
