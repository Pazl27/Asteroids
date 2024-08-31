package game

import(
  "encoding/json"
  "fmt"
  "os"
)

type HighScore struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

/* 
* function to decode the highscore
*/
func decodeHighScore() {
	file, err := os.Open("highscore.json")
	if err != nil {
		fmt.Println("Error opening highscore file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	highscore = HighScore{}
	err = decoder.Decode(&highscore)
	if err != nil {
		fmt.Println("Error decoding highscore:", err)
	}
}

/*
* function to check if the current score is higher than the highscore
* if yes the function saveHighScore is called
*/
func checkHighScore() {
	if int(score) > highscore.Score {
		err := saveHighScore()
		if err != nil {
			fmt.Println(err)
		}
	}
}

/*
* function to save the highscore
* @return error
*/
func saveHighScore() error {
	highscore.Score = int(score)
	if player_name == nil || len(player_name) == 0{
    highscore.Name = "Unknown"
	} else {
    highscore.Name = string(player_name)
	}
	file, err := os.Create("highscore.json")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(highscore); err != nil {
		return fmt.Errorf("failed to encode highscore: %w", err)
	}

	return nil
}
