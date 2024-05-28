package main

import (
	"fmt"
	"math/rand"

	as "example.com/asteroids/asteroids"
	rl "github.com/gen2brain/raylib-go/raylib"
  pl "example.com/asteroids/player"
)

var (
  list_size [3]int = [3]int{as.SMALL, as.MEDIUM, as.LARGE}
  list_ast []as.Asteroid

  target_pos rl.Vector2

  player pl.Ship
)

const (
  ScreenWidth = 800
  ScreenHeight = 450
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

    for i := range list_ast {
      as.UpdateAsteroid(&list_ast[i])
      as.DrawAsteroid(&list_ast[i])
    }

    player.UpdateShip()
    player.DrawShip()

		rl.EndDrawing()
}

func getNewAsteroid() {
    if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
      size_rand := list_size[rand.Intn(len(list_size))]

      a := as.Asteroid{
        Active: true,
        Position: rl.GetMousePosition(),
        Speed: as.GetSpeed(&target_pos, rl.GetMousePosition()),
        Rotation: 0,
        RoatatationSpeed: float32(rand.Intn(5) - 2),
        Size: size_rand,
      }

      list_ast = append(list_ast, a)
    }
}

func checkBoarders() {
    for i := len(list_ast) - 1; i >= 0; i-- {
      ast_temp := list_ast[i]

      if ast_temp.Position.X > float32(rl.GetScreenWidth()) || ast_temp.Position.X < 0 || ast_temp.Position.Y > float32(rl.GetScreenHeight()) || ast_temp.Position.Y < 0 {
        list_ast = append(list_ast[:i], list_ast[i+1:]...)
      }
    }
}

func main() {
  player = pl.Ship{Position: rl.Vector2{X: 400, Y: 225}, Speed: 0, Acceleration: 0.1, Rotation: 0}

  rl.InitWindow(ScreenWidth, ScreenHeight, "Asteroids")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
    getRandomPos()

    checkBoarders()

    getNewAsteroid()

    draw()

	}
}
