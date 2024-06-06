package items

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	pl "example.com/asteroids/player"
)

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

func (i *Invincible) NewItem() {
	i.Position = getRandomPos()
	i.Invincible = true
	i.InvincibleTime = 10
	i.Color = rl.Yellow
}

func (i *Invincible) Draw() {
	rl.DrawCircleV(i.Position, 10, i.Color)
}

func (i *Invincible) GetPosition() rl.Vector2 {
	return i.Position
}

func (i *Invincible) ApplyEffect(player *pl.Ship) {
	player.Invincible = true
}

func (i *Invincible) GetTime() float32 {
	return i.InvincibleTime
}

func (i *Invincible) GetName() string {
	return "Invincible"
}

type InfiniteAmmo struct {
	Position         rl.Vector2
	InfiniteAmmoTime float32
	InfiniteAmmo     bool
	Color            rl.Color
}

func (i *InfiniteAmmo) NewItem() {
	i.Position = getRandomPos()
	i.InfiniteAmmo = true
	i.InfiniteAmmoTime = 10
	i.Color = rl.Blue
}

func (i *InfiniteAmmo) Draw() {
	rl.DrawCircleV(i.Position, 10, i.Color)
}

func (i *InfiniteAmmo) GetPosition() rl.Vector2 {
	return i.Position
}

func (i *InfiniteAmmo) ApplyEffect(player *pl.Ship) {
	player.InfiniteAmmo = true
}

func (i *InfiniteAmmo) GetTime() float32 {
	return i.InfiniteAmmoTime
}

func (i *InfiniteAmmo) GetName() string {
	return "InfiniteAmmo"
}

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

func getRandomPos() rl.Vector2 {
	return rl.Vector2{
		X: float32(rand.Intn(rl.GetScreenWidth())),
		Y: float32(rand.Intn(rl.GetScreenHeight())),
	}
}

