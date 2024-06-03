package main

import (
	"fmt"
	"math/rand"

	as "example.com/asteroids/asteroids"
	pl "example.com/asteroids/player"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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

	gameState = true

	score float32 = 0.0
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

func draw() {

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	score_string := fmt.Sprintf("Score: %d", int(score))
	rl.DrawText(score_string, 10, 10, 20, rl.White)
	// string := fmt.Sprintf("Asteroids: %d", len(list_ast))
	// rl.DrawText(string, 10, 10, 20, rl.White)
	// string = fmt.Sprintf("Bullets: %d", len(player.Bullets))
	// rl.DrawText(string, 10, 30, 20, rl.White)

	for i := range list_ast {
		as.UpdateAsteroid(&list_ast[i])
		as.DrawAsteroid(&list_ast[i])
	}

	player.UpdateShip()
	player.DrawShip("assets/ship.png")

	rl.EndDrawing()
}

func getNewAsteroid() {
	if len(list_ast) >= 10 {
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
}

func checkBoarders() {
	for i := len(list_ast) - 1; i >= 0; i-- {
		ast_temp := list_ast[i]

		if ast_temp.Position.X > float32(rl.GetScreenWidth()) || ast_temp.Position.X < 0 || ast_temp.Position.Y > float32(rl.GetScreenHeight()) || ast_temp.Position.Y < 0 {
			list_ast = append(list_ast[:i], list_ast[i+1:]...)
		}
	}
}

func checkPlayerCollision() {
	playerRadius := float32(PLAYER_WIDTH) / 2

	for i := range list_ast {
		asteroid := list_ast[i]
		asteroidRadius := float32(asteroid.Size*2) / 2

		distance := rl.Vector2Distance(player.Position, asteroid.Position)
		if distance < playerRadius+asteroidRadius {
			fmt.Println("Collision detected!")
			gameState = false
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

func checkCollisions() {
	checkPlayerCollision()
	checkBulletCollision()
}

func main() {
	player = pl.Ship{
		Position:     rl.Vector2{X: 400, Y: 225},
		Speed:        0,
		Acceleration: 0.1,
		Rotation:     0,
	}

	rl.InitWindow(ScreenWidth, ScreenHeight, "Asteroids")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() && gameState {

		score += rl.GetFrameTime()
		getRandomPos()
		checkBoarders()
		getNewAsteroid()
		checkCollisions()
		draw()
	}
}
