package player

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MAX_SPEED = 8.0
)

type Ship struct {
	Position     rl.Vector2
	Speed        float32
	Acceleration float32
	Rotation     float32
}

func (ship *Ship) DrawShip() {
	// Draw the ship using its position and rotation
	rl.DrawPolyLines(ship.Position, 3, 16, ship.Rotation, rl.White)
}

func (ship *Ship) UpdateShip() {
	// Handle acceleration and deceleration
	if rl.IsKeyDown(rl.KeyW) {
		ship.Speed += ship.Acceleration
		if ship.Speed > MAX_SPEED {
			ship.Speed = MAX_SPEED
		}
	} else {
		ship.Speed -= ship.Acceleration
		if ship.Speed < 0 {
			ship.Speed = 0
		}
	}

	// Handle rotation
	if rl.IsKeyDown(rl.KeyA) {
		ship.Rotation -= 5
	}

	if rl.IsKeyDown(rl.KeyD) {
		ship.Rotation += 5
	}

	// Stop the ship
	if rl.IsKeyDown(rl.KeyS) {
		ship.Speed = 0
	}

  if rl.IsKeyDown(rl.KeySpace) {
    fmt.Println("Pew pew!")
  }

	// Convert rotation to radians
	rotationRadians := ship.Rotation * (math.Pi / 180.0)

	// Calculate direction vector based on rotation
	direction := rl.Vector2{
		X: float32(math.Cos(float64(rotationRadians))),
		Y: float32(math.Sin(float64(rotationRadians))),
	}

	// Update position based on speed and direction
	ship.Position.X += direction.X * ship.Speed
	ship.Position.Y += direction.Y * ship.Speed
}

