package player

import (
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_SPEED = 6.0
const IMMUNITY_DURATION = time.Second
const RELOAD_TIME = time.Second * 4

type Ship struct {
	Position     rl.Vector2
	Speed        float32
	Acceleration float32
	Rotation     float32 // in degrees

	Bullets      []Bullet
	Invincible   bool
	InfiniteAmmo bool
	Health       int
	Ammo         int
	Reloading    bool

	textureShip  rl.Texture2D
	textureHeart rl.Texture2D

	resetProtection time.Time
	reloadStart     time.Time
}

/*
* function to create a new ship
* @return Ship
 */
func NewShip() Ship {
	return Ship{
		Position: rl.Vector2{
			X: float32(rl.GetScreenWidth() / 2),
			Y: float32(rl.GetScreenHeight() / 2),
		},
		Speed:        0,
		Acceleration: 0.1,
		Rotation:     0,
		Invincible:   false,
		InfiniteAmmo: false,
		Ammo:         20,

		Health: 3,

		textureShip:  rl.LoadTexture("assets/images/ship.png"),
		textureHeart: resizeTexture("assets/images/heart.png", 32, 32),
	}
}

/*
* function to resize the texture
* @param string image
* @param int width
* @param int height
* @return rl.Texture2D
 */
func resizeTexture(image string, width int, height int) rl.Texture2D {
	img := rl.LoadImage(image)
	rl.ImageResize(img, int32(width), int32(height))
	texture := rl.LoadTextureFromImage(img)
	return texture
}

/*
* function to draw the ship
* TODO: Draw the ship with differnt colors based on the effect
* FIX: Drawing Invincible ship with yellow doesn't work
 */
func (ship *Ship) DrawShip() {
	source := rl.Rectangle{X: 0, Y: 0, Width: 32, Height: 48}
	dest := rl.Rectangle{X: ship.Position.X, Y: ship.Position.Y, Width: 48, Height: 48}
	origin := rl.Vector2{X: dest.Width / 2, Y: dest.Height / 2}
	rl.DrawTexturePro(ship.textureShip, source, dest, origin, ship.Rotation+90, rl.White)

	for _, bullet := range ship.Bullets {
		bullet.Draw()
	}
}

/*
* function that resets the ship to the center of the screen
 */
func (ship *Ship) Reset() {
	ship.Position = rl.Vector2{
		X: float32(rl.GetScreenWidth() / 2),
		Y: float32(rl.GetScreenHeight() / 2),
	}
	ship.Speed = 0
	ship.Rotation = 0

	ship.resetProtection = time.Now()
	ship.Invincible = true
}

/*
* function to draw the health of the ship
 */
func (ship *Ship) DrawHealth() {

	for i := 0; i < ship.Health; i++ {
		rl.DrawTexture(ship.textureHeart, int32(10+32*i), int32(rl.GetScreenHeight()-75), rl.White)
	}
}

/*
* function that draws the ammo of the ship
* if reloading, draw a loading bar
 */
func (ship *Ship) DrawAmmo() {
	if ship.Reloading {
		rl.DrawRectangle(10, int32(rl.GetScreenHeight()-30), 200, 20, rl.Black)
		rl.DrawRectangle(10, int32(rl.GetScreenHeight()-30), int32(200-(time.Since(ship.reloadStart).Seconds()*50)), 20, rl.Red)
	} else {
		rl.DrawRectangle(10, int32(rl.GetScreenHeight()-30), 200, 20, rl.Black)
		rl.DrawRectangle(10, int32(rl.GetScreenHeight()-30), int32(200-(float32(ship.Ammo)*10)), 20, rl.Red)
	}
}

/*
* function to update the ship
* calls the processInput function to handle key inputs
* updates the position of the ship and bullets
 */
func (ship *Ship) UpdateShip() {

	// processes key inputs
	ship.processInput()
	// Delete bullets that are out of bounds
	for i := 0; i < len(ship.Bullets); i++ {
		if ship.Bullets[i].DeleteBullet(ship.Position) {
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

	// Check if the ship is immune
	if time.Since(ship.resetProtection) > IMMUNITY_DURATION {
		ship.Invincible = false
	}

	// reload ammo
	if ship.Ammo == 0 && !ship.Reloading {
		ship.Reloading = true
		ship.reloadStart = time.Now()
	} else if ship.Reloading && time.Since(ship.reloadStart) > RELOAD_TIME {
		ship.Ammo = 20
		ship.Reloading = false
	}
}

/*
* function to handle key inputs
* W - accelerate
* A - rotate left
* D - rotate right
* S - stop
* Space - shoot
 */
func (ship *Ship) processInput() {
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
		ship.Rotation -= 6
	}

	if rl.IsKeyDown(rl.KeyD) {
		ship.Rotation += 6
	}

	// Stop the ship
	if rl.IsKeyDown(rl.KeyS) {
		ship.Speed = 0
	}

	// Shoot a bullet
	if ship.InfiniteAmmo && rl.IsKeyDown(rl.KeySpace) {
		newBullet := NewBullet(ship.Position, ship.Rotation) // Adjust the speed as needed
		ship.Bullets = append(ship.Bullets, newBullet)

	} else if rl.IsKeyPressed(rl.KeySpace) && !ship.Reloading {
		newBullet := NewBullet(ship.Position, ship.Rotation) // Adjust the speed as needed
		ship.Bullets = append(ship.Bullets, newBullet)
		ship.Ammo--
	}
}

/*
* function to check if the ship is in screen bounds
* @param rl.Vector2 pos
* @return bool
 */
func inScreenBounds(pos rl.Vector2) bool {
	return pos.X > 0 && pos.X < float32(rl.GetScreenWidth()) && pos.Y > 0 && pos.Y < float32(rl.GetScreenHeight())
}
