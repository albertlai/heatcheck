package main

import (
	"fmt"
	"os"
	"encoding/json"
)

// { resultSets: { rowSet: [ x[7] ] } }
func processGameJSON(game_id int) error {
	fmt.Printf("Processing game %d\n", game_id)
	in_name := fmt.Sprintf("%s/game_%d.json", json_path, game_id)
	in, err := os.Open(in_name)
	if err != nil { return err }
	defer in.Close()
	
	out_name := fmt.Sprintf("plays/plays_%d.json", game_id)	
	out, err := os.Create(out_name)
	if err != nil { return err }
	defer out.Close()

	dec := json.NewDecoder(in)
	var data map[string]interface{}
	if err := dec.Decode(&data); err != nil { return err }		
	results, exists := data["resultSets"].([]interface{})
	if exists && len(results) > 0 {
		rowSetContainer := results[0].(map[string]interface {})
		rows, exists := rowSetContainer["rowSet"].([]interface {})
		if exists {
			for i := 0; i < len(rows); i++ {
				row := rows[i].([]interface {})
				play1, ok1 := row[7].(string)
				play2, ok2 := row[9].(string)
				player, ok_player := row[14].(string)
				if (ok1 || ok2)  && ok_player {	
					out.WriteString(fmt.Sprintf("%s : %s - %s\n", player, play1, play2))
				}
			}
		}
	}
	return nil
}
