package player

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const speed = 8.0

type Bullet struct {
	Position  rl.Vector2
	Speed     float32
	Direction rl.Vector2
}

/*
* function to create a new bullet
* @param rl.Vector2 position
* @param float32 rotation
* @return Bullet
*/
func NewBullet(position rl.Vector2, rotation float32) Bullet {
	// Convert rotation to radians
	rotationRadians := rotation * (math.Pi / 180.0)

	// Calculate direction vector based on rotation
	direction := rl.Vector2{
		X: float32(math.Cos(float64(rotationRadians))),
		Y: float32(math.Sin(float64(rotationRadians))),
	}

	return Bullet{
		Position:  position,
		Speed:     speed,
		Direction: direction,
	}
}

/*
* function to update the bullet position
*/
func (b *Bullet) Update() {
	// Update the position based on speed and direction
	b.Position.X += b.Direction.X * b.Speed
	b.Position.Y += b.Direction.Y * b.Speed
}

/*
* function to draw the bullet
*/
func (b *Bullet) Draw() {
	// Draw the bullet
	rl.DrawCircleV(b.Position, 5, rl.Red)
}

/*
* function to check if the bullet is out of bounds
* or if it is too far from the player
* @param rl.Vector2 playerPosition
* @return bool
*/
func (b *Bullet) DeleteBullet(playerPosition rl.Vector2) bool {
  // Check if the bullet is out of bounds
  if b.Position.X < 0 || b.Position.X > 1000 || b.Position.Y < 0 || b.Position.Y > 800 {
    return true
  }

  distance := rl.Vector2Distance(b.Position, playerPosition)
	if distance > 500 {
		return true
	}
  
  return false
}
