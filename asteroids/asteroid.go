package asteroids

import (
	rl "github.com/gen2brain/raylib-go/raylib"
  "math/rand"
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
  Split bool
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

func SplitAsteroid(ast Asteroid, list_ast *[]Asteroid) {
    if ast.Size == SMALL {
        return
    }

    newSize := SMALL
    if ast.Size == LARGE {
        newSize = MEDIUM
    }

    for i := 0; i < 2; i++ {
        // Generate a random direction vector
        direction := rl.Vector2{
            X: rand.Float32()*2 - 1,
            Y: rand.Float32()*2 - 1,
        }
        direction = rl.Vector2Normalize(direction)
        direction.X *= SPEED
        direction.Y *= SPEED

        new_ast := Asteroid{
            Position: ast.Position,
            Speed:    direction,
            Rotation: float32(rand.Intn(360)), // Random initial rotation
            RoatatationSpeed: float32(rand.Intn(5) - 2), // Random rotation speed
            Size:     newSize,
            Split:    true,
        }

        *list_ast = append(*list_ast, new_ast)
    }
}

