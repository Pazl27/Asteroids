package asteroids

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

const (
	SMALL  = 1
	MEDIUM = 2
	LARGE  = 4

	SPEED = 2
)

var( 
  list_size [3]int = [3]int{SMALL, MEDIUM, LARGE}

  lastSpawnTime float64
)

type Asteroid struct {
	Position         rl.Vector2
	Speed            rl.Vector2
	Rotation         float32
	RoatatationSpeed float32
  texture           rl.Texture2D

	Size  int
	Split bool
}

/* 
* function to draw the asteroid
* @param *Asteroid asteroid
*/
func DrawAsteroid(asteroid *Asteroid) {
  rl.DrawTexturePro(asteroid.texture, 
                    rl.Rectangle{X: 0, Y: 0, Width: float32(asteroid.texture.Width), Height: float32(asteroid.texture.Height)},
                    rl.Rectangle{X: asteroid.Position.X, Y: asteroid.Position.Y, Width: float32(asteroid.texture.Width), Height: float32(asteroid.texture.Height)},
                    rl.Vector2{X: float32(asteroid.texture.Width) / 2, Y: float32(asteroid.texture.Height) / 2},
                    asteroid.Rotation,
                    rl.White,
    )
	// rl.DrawPolyLines(asteroid.Position, 8, float32(16*asteroid.Size), asteroid.Rotation, rl.White)
}

/*
* function to resize the asteroid image
*/
func (asteroid *Asteroid) resizeImage() {
  image := rl.LoadImage("assets/images/asteroid.png")
  rl.ImageResize(image, int32(asteroid.Size * 32), int32(asteroid.Size * 32))
  texture := rl.LoadTextureFromImage(image)
  asteroid.texture = texture
}

/* 
* function that creates new asteroids
* @param *[]Asteroid list_roids
* if the number of asteroids is less than the maximum number of asteroids new ones spawn
* the new asteroids spawn at the edge of the screen
* spawn time is 1 second
* the position of the asteroid is random calls the getRandomTarget function
* the speed of the asteroid is calculated using the GetSpeed function
* the rotation of the asteroid is random
* the size of the asteroid is random
* the resizeImage function is called to resize the asteroid image
*/
func GetNewAsteroid(list_roids *[]Asteroid) {
	const maxAsteroids = 15

	if len(*list_roids) >= maxAsteroids {
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

	target_pos := getRandomTarget()
	a := Asteroid{
		Position:         position,
		Speed:            GetSpeed(&target_pos, position),
		Rotation:         0,
		RoatatationSpeed: float32(rand.Intn(5) - 2),
		Size:             size_rand,
	}
	a.resizeImage()

	*list_roids = append(*list_roids, a)
}

/*
* function to update the asteroid
* @param *Asteroid asteroid
* updates the position of the asteroid based on the speed
* updates the rotation of the asteroid based on the rotation speed
*/
func UpdateAsteroid(asteroid *Asteroid) {
	asteroid.Position.X += asteroid.Speed.X
	asteroid.Position.Y += asteroid.Speed.Y
	asteroid.Rotation += asteroid.RoatatationSpeed
}

/*
* function to calculate the speed of the asteroid
* @param *rl.Vector2 target
* @param rl.Vector2 position
* calculates the speed of the asteroid based on the target and position
* the speed is normalized
* the speed is multiplied by the speed constant
*/
func GetSpeed(target *rl.Vector2, position rl.Vector2) rl.Vector2 {
	speed := rl.Vector2{X: target.X - position.X, Y: target.Y - position.Y}
	speed = rl.Vector2Normalize(speed)
	speed.X *= SPEED
	speed.Y *= SPEED

	return speed
}

/*
* function to split the asteroid
* @param Asteroid ast
* @param *[]Asteroid list_ast
* if the asteroid is small it does not split
* if the asteroid is large it splits into two medium asteroids
* the new asteroids have random directions
* the new asteroids have random rotation speeds
* the new asteroids have the size of the asteroid that split
* the resizeImage function is called to resize the asteroid image
*/
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
			Position:         ast.Position,
			Speed:            direction,
			Rotation:         float32(rand.Intn(360)),   // Random initial rotation
			RoatatationSpeed: float32(rand.Intn(5) - 2), // Random rotation speed
			Size:             newSize,
			Split:            true,
		}
    new_ast.resizeImage()

		*list_ast = append(*list_ast, new_ast)
	}
}

/*
* function to get a random position
* returns a random position within the screen
* the position is the target position for the asteroid
*/
func getRandomTarget() rl.Vector2 {
	centerX := rl.GetScreenWidth() / 2
	centerY := rl.GetScreenHeight() / 2

	return rl.Vector2{
		X: float32(centerX + rand.Intn(101) - 50),
		Y: float32(centerY + rand.Intn(101) - 50),
	}
}
