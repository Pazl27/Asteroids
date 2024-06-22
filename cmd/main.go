package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	as "example.com/asteroids/asteroids"
	it "example.com/asteroids/items"
	pl "example.com/asteroids/player"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type HighScore struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

var (
	list_roids              []as.Asteroid
	asteroid_added          int     = 0
	player                  pl.Ship
	game_running                    = false
	score                   float32 = 0.0
	highscore               HighScore
	player_name                                        = make([]rune, 0)
	item                    it.Item                    = nil
	last_item_spawn_time    float64                    // Track the last item spawn time
	item_effect_start_times = make(map[string]float64) // Track when item effects start
	framesCounter           int
)

const (
	ScreenWidth   = 1000
	ScreenHeight  = 800
	PLAYER_WIDTH  = 48
	maxInputChars = 9
)

/*
* function to get a random position on the screen
* the position is within 101 pixels of the center of the screen
*/
func getRandomPos() rl.Vector2 {
	centerX := ScreenWidth / 2
	centerY := ScreenHeight / 2

	return rl.Vector2{
		X: float32(centerX + rand.Intn(101) - 50),
		Y: float32(centerY + rand.Intn(101) - 50),
	}
}

/*
* function to draw the game
* draws the background, player, asteroids and items
*/ 
func drawGame(background_texture rl.Texture2D) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	// Background image
	rl.DrawTexture(background_texture, 0, 0, rl.White)

	score_string := fmt.Sprintf("Score: %d", int(score))
	rl.DrawText(score_string, 10, 10, 20, rl.White)
	// Powerup timer display
	for effect, startTime := range item_effect_start_times {
		timeLeft := 5 - (rl.GetTime() - startTime)
		if timeLeft > 0 {
			rl.DrawText(fmt.Sprintf("%s: %.2f", effect, timeLeft), 10, 30, 20, rl.White)
		}
	}

	if item != nil {
		item.Draw()
	}

	for i := range list_roids {
		as.UpdateAsteroid(&list_roids[i])
		as.DrawAsteroid(&list_roids[i])
	}

	player.UpdateShip()
	player.DrawShip()
  player.DrawHealth()
  player.DrawAmmo()

	rl.EndDrawing()
}

/*
* function to delete asteroids that go off the screen
*/
func checkBoarders() {
	for i := len(list_roids) - 1; i >= 0; i-- {
		ast_temp := list_roids[i]

		if ast_temp.Position.X > float32(rl.GetScreenWidth()) || ast_temp.Position.X < 0 || ast_temp.Position.Y > float32(rl.GetScreenHeight()) || ast_temp.Position.Y < 0 {
			list_roids = append(list_roids[:i], list_roids[i+1:]...)
		}
	}
}

/*
* function to check if the player has collided with an asteroid
* if that happens the player loses health
* player gets reset to the center of the screen
* if the player has no health the game ends
*/
func checkPlayerCollision() {
	if !player.Invincible {
		playerRadius := float32(PLAYER_WIDTH) / 2

		for i := range list_roids {
			asteroid := list_roids[i]
			asteroidRadius := float32(asteroid.Size * 8) // Update radius based on asteroid size

			distance := rl.Vector2Distance(player.Position, asteroid.Position)
			if distance < playerRadius+asteroidRadius {
        player.Health --
        player.Reset()
			}
		}

    if player.Health <= 0 {
      game_running = false
    }
	}
}

/*
* function to check if a bullet has collided with an asteroid
* if that happens the bullet is removed
* if that happens the asteroid is split and the score is updated
*/
func checkBulletCollision() {
	for i := len(player.Bullets) - 1; i >= 0; i-- {
		bullet := player.Bullets[i]

		for j := len(list_roids) - 1; j >= 0; j-- {
			asteroid := list_roids[j]

			// Check for collision
			distance := rl.Vector2Distance(bullet.Position, asteroid.Position)
			if distance < 5+float32(asteroid.Size*16) { // Assuming bullet radius is 5 and asteroid radius is size*8
				// Remove the bullet
				player.Bullets = append(player.Bullets[:i], player.Bullets[i+1:]...)

				// Split the asteroid
				as.SplitAsteroid(asteroid, &list_roids)

				// Remove the asteroid
				list_roids = append(list_roids[:j], list_roids[j+1:]...)

				// Update the score based on the size of the asteroid
				switch asteroid.Size {
				case as.SMALL:
					score += 10
				case as.MEDIUM:
					score += 20
				case as.LARGE:
					score += 30
				}

				break
			}
		}
	}
}

/* 
* function to check if the player has collided with an item
*/
func checkItemCollision() {
	if item != nil {
		distance := rl.Vector2Distance(player.Position, item.GetPosition())
		if distance < 20 {
			item.ApplyEffect(&player)
			item_effect_start_times[item.GetName()] = rl.GetTime() // Record when the effect starts
			item = nil
		}
	}
}

/*
* function that check if the item timer has expired
* if it has the effect is removed
*/
func disableExpiredEffects() {
	currentTime := rl.GetTime()
	for effect, startTime := range item_effect_start_times {
		if currentTime-startTime >= 5 {
			switch effect {
			case "Invincible":
				player.Invincible = false
			case "InfiniteAmmo":
				player.InfiniteAmmo = false
			}
			delete(item_effect_start_times, effect) // Remove the effect once it's expired
		}
	}
}

/* 
* function that checks for collisions
* checks if the player has picked up an item
* checks if the player has collided with an asteroid
* checks if a bullet has collided with an asteroid
* also spawns new asteroids if needed
*/
func checkCollisions() {
	checkItemCollision()
	checkPlayerCollision()
	checkBulletCollision()

	as.GetNewAsteroid(&list_roids)
}

/* 
* function to reset the game
* resets the player, asteroids, score and game state
* remvoes all items
*/
func resetGame() {
	player = pl.NewShip()

	list_roids = nil
	score = 0
	game_running = true
	asteroid_added = 0
	last_item_spawn_time = rl.GetTime()
	item_effect_start_times = make(map[string]float64)

	item = nil
}

/* 
* function that draws the menu
* draws the title, score, highscore and instructions
* also draws the text box for the player to enter their name
*/
func drawMenu() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText("Press 'Enter' to start", ScreenWidth/2-100, ScreenHeight/2, 20, rl.White)
	rl.DrawText("Asteroids", ScreenWidth/2-100, ScreenHeight/2-50, 40, rl.White)
	rl.DrawText("Score: "+fmt.Sprintf("%d", int(score)), ScreenWidth/2-100, ScreenHeight/2+50, 20, rl.White)
	rl.DrawText("Highscore: "+fmt.Sprintf("%d", highscore.Score), ScreenWidth/2-100, ScreenHeight/2+75, 20, rl.White)
	rl.DrawText("Press 'Q' to quit", ScreenWidth/2-100, ScreenHeight/2+100, 20, rl.White)
	rl.DrawText("Enter Name:", ScreenWidth-200, ScreenHeight-70, 10, rl.White)

	textBox := rl.Rectangle{X: ScreenWidth - 200, Y: ScreenHeight - 50, Width: 180, Height: 30}
	mouseOnText := rl.CheckCollisionPointRec(rl.GetMousePosition(), textBox)

	if mouseOnText {
		framesCounter++
	} else {
		framesCounter = 0
	}

	// Draw the text box
	rl.DrawRectangleRec(textBox, rl.LightGray)
	if mouseOnText {
		rl.DrawRectangleLines(int32(textBox.X), int32(textBox.Y), int32(textBox.Width), int32(textBox.Height), rl.Red)
	} else {
		rl.DrawRectangleLines(int32(textBox.X), int32(textBox.Y), int32(textBox.Width), int32(textBox.Height), rl.DarkGray)
	}

	rl.DrawText(string(player_name), int32(textBox.X)+5, int32(textBox.Y)+5, 20, rl.Maroon)

	if mouseOnText {
		rl.SetMouseCursor(rl.MouseCursorIBeam)

		key := rl.GetCharPressed()
		for key > 0 {
			if key >= 32 && key <= 125 && len(player_name) < maxInputChars {
				player_name = append(player_name, key)
			}
			key = rl.GetCharPressed()
		}
		if (framesCounter/20)%2 == 0 {
			rl.DrawText("_", int32(textBox.X)+8+rl.MeasureText(string(player_name), 20), int32(textBox.Y)+10, 20, rl.Maroon)
		}

		if rl.IsKeyPressed(rl.KeyBackspace) {
			if len(player_name) > 0 {
				player_name = player_name[:len(player_name)-1]
			}
		}
	} else {
		rl.SetMouseCursor(rl.MouseCursorDefault)
	}

	rl.EndDrawing()

	processInput()
}

/* 
* function to process input in the menu
*/
func processInput() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		resetGame()
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		os.Exit(0)
	}
}

/* 
* init function to load the highscore
*/
func init() {
	decodeHighScore()
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

/*
* function that loads the background texture
* @return rl.Texture2D
*/
func loadTextures() rl.Texture2D {
	background := rl.LoadImage("assets/background_3.jpg")
	rl.ImageResize(background, ScreenWidth, ScreenHeight)
	background_texture := rl.LoadTextureFromImage(background)

	return background_texture
}

/* 
* main function
* initializes the window, sets the target fps
* loads the background texture
* runs the game loop
*/
func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Asteroids")
  rl.SetExitKey(0)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	background_texture := loadTextures()

  player = pl.NewShip()

	last_item_spawn_time = rl.GetTime()

	for !rl.WindowShouldClose() {
		if !game_running {
			checkHighScore()
			drawMenu()
		} else {
			score += rl.GetFrameTime()
			checkBoarders()
			as.GetNewAsteroid(&list_roids) // Continuously check and add new asteroids if needed
			checkCollisions()
			disableExpiredEffects()
			drawGame(background_texture)

			// Check if 10 seconds have passed since the last item was spawned
			if rl.GetTime()-last_item_spawn_time >= 10 {
				item = it.SpawnItem()
				last_item_spawn_time = rl.GetTime()
			}
		}
	}
}

