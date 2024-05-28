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
	string := fmt.Sprintf("Asteroids: %d", len(list_ast))
	rl.DrawText(string, 10, 10, 20, rl.White)
	string = fmt.Sprintf("Bullets: %d", len(player.Bullets))
	rl.DrawText(string, 10, 30, 20, rl.White)

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
		Active:           true,
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

func checkCollisions() {
	playerRadius := float32(PLAYER_WIDTH) / 2

	for i := range list_ast {
		asteroid := list_ast[i]
		asteroidRadius := float32(asteroid.Size * 2) / 2

		distance := rl.Vector2Distance(player.Position, asteroid.Position)
		if distance < playerRadius+asteroidRadius {
			fmt.Println("Collision detected!")
			gameState = false
		}
	}
}

func main() {
	player = pl.Ship{Position: rl.Vector2{X: 400, Y: 225}, Speed: 0, Acceleration: 0.1, Rotation: 0}

	rl.InitWindow(ScreenWidth, ScreenHeight, "Asteroids")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() && gameState {
		getRandomPos()
		checkBoarders()
		getNewAsteroid()
		checkCollisions()
		draw()
	}
}
