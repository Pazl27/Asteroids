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
	list_size [3]int = [3]int{as.SMALL, as.MEDIUM, as.LARGE}
	list_ast  []as.Asteroid

	target_pos rl.Vector2

	lastSpawnTime float64

	player = pl.Ship{
		Position: rl.Vector2{X: ScreenWidth / 2,
			Y: ScreenHeight / 2},
		Speed:        0,
		Acceleration: 0.1,
		Rotation:     0,
	}

	gameRunning = false

	score float32 = 0.0

	highscore HighScore

	ast_added int = 0

	playerName string = "unknown"

	item it.Item = nil

	lastItemSpawnTime float64 // Track the last item spawn time

	itemEffectStartTimes = make(map[string]float64) // Track when item effects start
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 800
	PLAYER_WIDTH = 48
)

func getRandomPos() {
	centerX := ScreenWidth / 2
	centerY := ScreenHeight / 2

	target_pos = rl.Vector2{
		X: float32(centerX + rand.Intn(101) - 50),
		Y: float32(centerY + rand.Intn(101) - 50),
	}
}

func drawGame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	score_string := fmt.Sprintf("Score: %d", int(score))
	rl.DrawText(score_string, 10, 10, 20, rl.White)

	if item != nil {
		item.Draw()
	}

	for i := range list_ast {
		as.UpdateAsteroid(&list_ast[i])
		as.DrawAsteroid(&list_ast[i])
	}

	player.UpdateShip()
	player.DrawShip("assets/ship.png")

	rl.EndDrawing()
}

func getNewAsteroid() {
	if ast_added >= 10 {
		return
	}

	currentTime := rl.GetTime()
	if currentTime-lastSpawnTime < 1 {
		return
	}
	lastSpawnTime = currentTime

	size_rand := list_size[rand.Intn(len(list_size))]

	// Generate a random position at the edge of the screen
	edge := rand.Intn(4)
	var position rl.Vector2
	switch edge {
	case 0: // top
		position = rl.Vector2{X: float32(rand.Intn(rl.GetScreenWidth())), Y: 0}
	case 1: // right
		position = rl.Vector2{X: float32(rl.GetScreenWidth()), Y: float32(rand.Intn(rl.GetScreenHeight()))}
	case 2: // bottom
		position = rl.Vector2{X: float32(rand.Intn(rl.GetScreenWidth())), Y: float32(rl.GetScreenHeight())}
	case 3: // left
		position = rl.Vector2{X: 0, Y: float32(rand.Intn(rl.GetScreenHeight()))}
	}

	a := as.Asteroid{
		Position:         position,
		Speed:            as.GetSpeed(&target_pos, position),
		Rotation:         0,
		RoatatationSpeed: float32(rand.Intn(5) - 2),
		Size:             size_rand,
	}

	list_ast = append(list_ast, a)
	ast_added++
}

func checkBoarders() {
	for i := len(list_ast) - 1; i >= 0; i-- {
		ast_temp := list_ast[i]

		if ast_temp.Position.X > float32(rl.GetScreenWidth()) || ast_temp.Position.X < 0 || ast_temp.Position.Y > float32(rl.GetScreenHeight()) || ast_temp.Position.Y < 0 {
			list_ast = append(list_ast[:i], list_ast[i+1:]...)
			ast_added--
		}
	}
}

func checkPlayerCollision() {
	if !player.Invincible {
		playerRadius := float32(PLAYER_WIDTH) / 2

		for i := range list_ast {
			asteroid := list_ast[i]
			asteroidRadius := float32(asteroid.Size * 8) // Update radius based on asteroid size

			distance := rl.Vector2Distance(player.Position, asteroid.Position)
			if distance < playerRadius+asteroidRadius {
				gameRunning = false
			}
		}
	}
}

func checkBulletCollision() {
	for i := len(player.Bullets) - 1; i >= 0; i-- {
		bullet := player.Bullets[i]

		for j := len(list_ast) - 1; j >= 0; j-- {
			asteroid := list_ast[j]

			// Check for collision
			distance := rl.Vector2Distance(bullet.Position, asteroid.Position)
			if distance < 5+float32(asteroid.Size*8) { // Assuming bullet radius is 5 and asteroid radius is size*8
				// Remove the bullet
				player.Bullets = append(player.Bullets[:i], player.Bullets[i+1:]...)

				// Split the asteroid
				as.SplitAsteroid(asteroid, &list_ast)

				// Remove the asteroid
				list_ast = append(list_ast[:j], list_ast[j+1:]...)

				// Update the score based on the size of the asteroid
				switch asteroid.Size {
				case as.SMALL:
					score += 10
				case as.MEDIUM:
					score += 20
				case as.LARGE:
					score += 30
				}

				// Break the inner loop as the bullet is already removed
				break
			}
		}
	}
}

func checkItemCollision() {
	if item != nil {
		distance := rl.Vector2Distance(player.Position, item.GetPosition())
		if distance < 20 {
			item.ApplyEffect(&player)
			itemEffectStartTimes[item.GetName()] = rl.GetTime() // Record when the effect starts
			item = nil
		}
	}
}

func disableExpiredEffects() {
	currentTime := rl.GetTime()
	for effect, startTime := range itemEffectStartTimes {
		if currentTime-startTime >= 5 {
			switch effect {
			case "Invincible":
				player.Invincible = false
			case "InfiniteAmmo":
				player.InfiniteAmmo = false
			}
			delete(itemEffectStartTimes, effect) // Remove the effect once it's expired
		}
	}
}

func checkCollisions() {
	checkItemCollision()
	checkPlayerCollision()
	checkBulletCollision()
}

func resetGame() {
	player = pl.Ship{
		Position:     rl.Vector2{X: ScreenWidth / 2, Y: ScreenHeight / 2},
		Speed:        0,
		Acceleration: 0.1,
		Rotation:     0,
	}

	list_ast = nil
	score = 0
	gameRunning = true
	ast_added = 0
	lastItemSpawnTime = rl.GetTime()
	itemEffectStartTimes = make(map[string]float64)
}

func drawMenu() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText("Press 'Enter' to start", ScreenWidth/2-100, ScreenHeight/2, 20, rl.White)
	rl.DrawText("Asteroids", ScreenWidth/2-100, ScreenHeight/2-50, 40, rl.White)
	rl.DrawText("Score: "+fmt.Sprintf("%d", int(score)), ScreenWidth/2-100, ScreenHeight/2+50, 20, rl.White)
	rl.DrawText("Highscore: "+fmt.Sprintf("%d", highscore.Score), ScreenWidth/2-100, ScreenHeight/2+75, 20, rl.White)
	rl.DrawText("Press 'Q' to quit", ScreenWidth/2-100, ScreenHeight/2+100, 20, rl.White)
	// TODO: Add input for player name

	rl.EndDrawing()

	processInput()
}

func processInput() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		resetGame()
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		os.Exit(0)
	}
}

func init() {
	decodeHighScore()
}

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

func checkHighScore() {
	if int(score) > highscore.Score {
		err := saveHighScore()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func saveHighScore() error {
	highscore.Score = int(score)
	highscore.Name = playerName
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

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Asteroids")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	lastItemSpawnTime = rl.GetTime()

	for !rl.WindowShouldClose() {
		if !gameRunning {
			checkHighScore()
			drawMenu()
		} else {
			score += rl.GetFrameTime()
			getRandomPos()
			checkBoarders()
			getNewAsteroid()
			checkCollisions()
			disableExpiredEffects()
			drawGame()

			// Check if 10 seconds have passed since the last item was spawned
			if rl.GetTime()-lastItemSpawnTime >= 10 {
				item = it.SpawnItem()
				lastItemSpawnTime = rl.GetTime()
			}
		}
	}
}

