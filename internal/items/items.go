package items

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	pl "github.com/Pazl27/Asteroids/internal/player"
)

/* 
* interface for items
* used to create new items, draw them, apply effects to the player
* useful for a list with Items with different structs
*/
type Item interface {
	NewItem()
	Draw()
	GetPosition() rl.Vector2
	ApplyEffect(player *pl.Ship)
	GetTime() float32
	GetName() string
}

type Invincible struct {
	Position       rl.Vector2
	InvincibleTime float32
	Invincible     bool
	Color          rl.Color
}

/* 
* function to create a new invincible item
*/
func (i *Invincible) NewItem() {
	i.Position = getRandomPos()
	i.Invincible = true
	i.InvincibleTime = 10
	i.Color = rl.Yellow
}

/*
* function to draw the invincible item
*/
func (i *Invincible) Draw() {
	rl.DrawCircleV(i.Position, 10, i.Color)
}

/*
* function to get the position of the invincible item
*/
func (i *Invincible) GetPosition() rl.Vector2 {
	return i.Position
}

/*
* function to apply the invincible effect to the player
*/
func (i *Invincible) ApplyEffect(player *pl.Ship) {
	player.Invincible = true
}

/*
* function to get the time of the invincible item
*/
func (i *Invincible) GetTime() float32 {
	return i.InvincibleTime
}

/*
* function to get the name of the item
*/
func (i *Invincible) GetName() string {
	return "Invincible"
}

type InfiniteAmmo struct {
	Position         rl.Vector2
	InfiniteAmmoTime float32
	InfiniteAmmo     bool
	Color            rl.Color
}

/*
* function to create a new infinite ammo item
*/
func (i *InfiniteAmmo) NewItem() {
	i.Position = getRandomPos()
	i.InfiniteAmmo = true
	i.InfiniteAmmoTime = 10
	i.Color = rl.Blue
}

/*
* function to draw the infinite ammo item
*/
func (i *InfiniteAmmo) Draw() {
	rl.DrawCircleV(i.Position, 10, i.Color)
}

/*
* function to get the position of the infinite ammo item
*/
func (i *InfiniteAmmo) GetPosition() rl.Vector2 {
	return i.Position
}

/*
* function to apply the infinite ammo effect to the player
*/
func (i *InfiniteAmmo) ApplyEffect(player *pl.Ship) {
	player.InfiniteAmmo = true
}

/*
* function to get the time of the infinite ammo item
*/
func (i *InfiniteAmmo) GetTime() float32 {
	return i.InfiniteAmmoTime
}

/*
* function to get the name of the item
*/
func (i *InfiniteAmmo) GetName() string {
	return "InfiniteAmmo"
}

/*
* function to spawn a random item
* @return Item
*/
func SpawnItem() Item {
	random := rand.Intn(2)
	var item Item

	switch random {
	case 0:
		item = &Invincible{}
	case 1:
		item = &InfiniteAmmo{}
	}
	item.NewItem()

	return item
}

/*
* function to get a random position in the window
* @return rl.Vector2
*/
func getRandomPos() rl.Vector2 {
	return rl.Vector2{
		X: float32(rand.Intn(rl.GetScreenWidth())),
		Y: float32(rand.Intn(rl.GetScreenHeight())),
	}
}

