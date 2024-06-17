package player

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SPEED = 6.0

type Ship struct {
	Position     rl.Vector2
	Speed        float32
	Acceleration float32
	Rotation     float32 // in degrees

	Bullets      []Bullet

  Invincible bool
  InfiniteAmmo bool
}

func NewShip() Ship {
  return Ship{
    Position: rl.Vector2{
      X: float32(rl.GetScreenWidth() / 2),
      Y: float32(rl.GetScreenHeight() / 2),
    },
    Speed:        0,
    Acceleration: 0.1,
    Rotation:     0,
    Invincible: false,
    InfiniteAmmo: false,
  }
}

func (ship *Ship) DrawShip(image string) {
  texture := rl.LoadTexture(image)
	// Draw the ship using its position and rotation
	// rl.DrawPolyLines(ship.Position, 3, 16, ship.Rotation, rl.White)
	source := rl.Rectangle{X: 0, Y: 0, Width: 32, Height: 48}
	dest := rl.Rectangle{X: ship.Position.X, Y: ship.Position.Y, Width: 48, Height: 48}
	origin := rl.Vector2{X: dest.Width / 2, Y: dest.Height / 2}
	rl.DrawTexturePro(texture, source, dest, origin, ship.Rotation+90, rl.White)

	for _, bullet := range ship.Bullets {
		bullet.Draw()
	}
}

func (ship *Ship) UpdateShip() {
  
  // processes key inputs
  ship.processInput()
	// Delete bullets that are out of bounds
	for i := 0; i < len(ship.Bullets); i++ {
		if ship.Bullets[i].DeleteBullet() {
			ship.Bullets = append(ship.Bullets[:i], ship.Bullets[i+1:]...)
		}
	}

	// Update all bullets
	for i := range ship.Bullets {
		ship.Bullets[i].Update()
	}

	rotationRadians := ship.Rotation * (math.Pi / 180.0)

	// Calculate direction vector based on rotation
	direction := rl.Vector2{
		X: float32(math.Cos(float64(rotationRadians))),
		Y: float32(math.Sin(float64(rotationRadians))),
	}

	// Calculate new position based on speed and direction
	newPosition := rl.Vector2{
		X: ship.Position.X + direction.X*ship.Speed,
		Y: ship.Position.Y + direction.Y*ship.Speed,
	}

	// Update position if it's within screen bounds
	if inScreenBounds(newPosition) {
		ship.Position = newPosition
	}
}

func (ship *Ship)processInput() {
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

	// Shoot a bullet
  if ship.InfiniteAmmo && rl.IsKeyDown(rl.KeySpace) {
    newBullet := NewBullet(ship.Position, ship.Rotation) // Adjust the speed as needed
    ship.Bullets = append(ship.Bullets, newBullet)

  } else if rl.IsKeyPressed(rl.KeySpace) {
		newBullet := NewBullet(ship.Position, ship.Rotation) // Adjust the speed as needed
		ship.Bullets = append(ship.Bullets, newBullet)
	}
}

func inScreenBounds(pos rl.Vector2) bool {
	return pos.X > 0 && pos.X < float32(rl.GetScreenWidth()) && pos.Y > 0 && pos.Y < float32(rl.GetScreenHeight())
}
