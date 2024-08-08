package particles

import (
	"image/color"
	"log"

	"git.smallzcomputing.com/sand-game/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type Grid struct {
	Width, Height int
	Map           [][]Particle
}

type Vector2 struct {
	X, Y int
}

type Particle struct {
	Active   bool
	Position Vector2
	Pixel    *ebiten.Image
}

func PrepareGrid(width, height, MOUSEX, MOUSEY int) [][]Particle {
	// Prepare Gridmap
	result := make([][]Particle, width)

	for i := 0; i < width; i++ {
		result[i] = make([]Particle, height)

		for j := 0; j < height; j++ {
			result[i][j] = *result[i][j].PrepareParticle(MOUSEX, MOUSEY)
		}

	}
	log.Printf("Grid size = [width: %v, height: %v]\n", len(result), len(result[0]))
	return result
}

func (particle *Particle) PrepareParticle(MOUSEX, MOUSEY int) *Particle {

	result := Particle{Active: false, Position: Vector2{MOUSEX, MOUSEY}, Pixel: ebiten.NewImage(1, 1)}
	return &result
}

func DrawGrid(renderer *ebiten.Image, GRID Grid) {
	// Loop through all grid positions
	for x := 0; x < GRID.Width; x++ {
		for y := 0; y < GRID.Height; y++ {

			// If particle is active, draw it
			if GRID.Map[x][y].Active {

				// Drawing
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x), float64(y))
				GRID.Map[x][y].Pixel.Fill(color.RGBA{255, 255, 255, 255})
				renderer.DrawImage(GRID.Map[x][y].Pixel, op)
			}

		}
	}
}

func SimulateGravity(particle Particle) *Particle {
	//pos := particle.Position
	//floorY := (SCREENHEIGHT / 2) - 1

	//	if pos.Y < floorY {

	return &particle
	//	}

}

func CheckForParticleSpawn(GRID Grid, MOUSEX int, MOUSEY int) {
	// If mouse0 pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton(0)) && MOUSEX >= 0 && MOUSEY >= 0 {
		// If particle pixel is INACTIVE
		particle := GRID.Map[MOUSEX][MOUSEY]

		if !particle.Active {
			// ACTIVATE particle pixel
			//game.PARTICLE_COUNT++
			GRID.Map[MOUSEX][MOUSEY].Active = true
			//log.Printf("Spawning pixel @ [%v, %v] -- #%v", MOUSEX, MOUSEY, game.PARTICLE_COUNT)
		}
	}
}

func SimulateCollision(particle Particle) *Particle {

	belowParticle := game.GRID.Map[particle.Position.X][particle.Position.Y-1]
	aboveParticle := game.GRID.Map[particle.Position.X][particle.Position.Y+1]
	belowLeftParticle := game.GRID.Map[particle.Position.X-1][particle.Position.Y-1]
	belowRightParticle := game.GRID.Map[particle.Position.X+1][particle.Position.Y-1]
	groundParticle := game.GRID.Map[particle.Position.X][particle.Position.Y-game.GRAVITY]
	if !particle.Active {

		// Particle is inactive

	} else if /*PARTICLE ACTIVE*/ belowParticle.Active {

		// am inside particle?
		//if ()

		if groundParticle.Active {
			particle.Position = Vector2{groundParticle.Position.X, groundParticle.Position.Y - 2}
		}

		// has particle above & below
		if aboveParticle.Active && belowParticle.Active {
			return &particle
		}

		// is there a particle below left or right? (sand effect)
		if !belowLeftParticle.Active {
			particle.Position = belowLeftParticle.Position
			//	belowLeftParticle.Active = true
			//particle.Active = false
		} else if !belowRightParticle.Active {
			particle.Position = belowRightParticle.Position
			//	belowRightParticle.Active = true
			//	particle.Active = false
		}

	} else if particle.Active {

		// does not have particle below
		particle.Position = belowParticle.Position
	}
	return &particle
}
