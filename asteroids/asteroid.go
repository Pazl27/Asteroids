package asteroids

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SMALL  = 1
	MEDIUM = 2
	LARGE  = 4

	SPEED  = 2
)

type Asteroid struct {
	Position         rl.Vector2
	Speed            rl.Vector2
	Rotation         float32
	RoatatationSpeed float32

	Size int
}

func DrawAsteroid(asteroid *Asteroid) {
	rl.DrawPolyLines(asteroid.Position, 3, float32(16*asteroid.Size), asteroid.Rotation, rl.White)
}

func UpdateAsteroid(asteroid *Asteroid) {
	asteroid.Position.X += asteroid.Speed.X
	asteroid.Position.Y += asteroid.Speed.Y
	asteroid.Rotation += asteroid.RoatatationSpeed
}

func GetSpeed(target *rl.Vector2, position rl.Vector2) rl.Vector2 {
	speed := rl.Vector2{X: target.X - position.X, Y: target.Y - position.Y}
	speed = rl.Vector2Normalize(speed)
	speed.X *= SPEED
	speed.Y *= SPEED

	return speed
}
