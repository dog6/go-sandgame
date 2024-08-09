package game

import (
	"fmt"
	"image/color"
	"log"

	particles "git.smallzcomputing.com/sand-game/particles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Grid struct {
	Width, Height int
	Map           [][]Particle
}

var (
	GRID           Grid
	MOUSEX, MOUSEY int
	PARTICLE_COUNT int
)

type Game struct{}

/*
	type Vector2 struct {
		X, Y int
	}
*/
type Particle struct {
	Active   bool
	Position particles.Vector2
	Pixel    *ebiten.Image
}

const (
	SCREENWIDTH, SCREENHEIGHT = 1280, 720
	GRAVITY                   = 10
)

func (g *Game) Update() error {
	MOUSEX, MOUSEY = ebiten.CursorPosition()    // Capture mouse position
	CheckForParticleSpawn(GRID, MOUSEX, MOUSEY) // Check for particle spawn
	//CheckForSortColumns()
	SimulateParticles()
	return nil
}

/*


func CountActiveInColumn(c int) int {
	var result int
	fmt.Printf("Sorting column #%v with %v active particles\n", c, len(GRID.Map[c]))
	for i := 0; i < len(GRID.Map[c])-1; i++ {
		if GRID.Map[c][i].Active {
			result++
		}
	}
	return result
}

func SwapColumns(col1, col2 int) {

	var newCol1 []Particle
	var newCol2 []Particle

	for i := 0; i < (GRID.Height/2)-1; i++ {
		if GRID.Map[col1][i].Active {
			newCol2 = append(newCol2, Particle{Active: true, Position: GRID.Map[col2][i].Position, Pixel: GRID.Map[col1][i].Pixel})
		}
		if GRID.Map[col2][i].Active {
			newCol1 = append(newCol1, Particle{Active: true, Position: GRID.Map[col1][i].Position, Pixel: GRID.Map[col2][i].Pixel})
		}
	}
	GRID.Map[col1] = newCol1
	GRID.Map[col2] = newCol2

}

/*func CheckForSortColumns() {

	if ebiten.IsKeyPressed(ebiten.KeySpace) {

		for x := 0; x < GRID.Width-2; x++ {
			fmt.Printf("Column %v has %v active\n", x, CountActiveInColumn(x))
			if CountActiveInColumn(x) < CountActiveInColumn(x+1) { /*CountActiveInColumn(x) > CountActiveInColumn(x+1) {
				SwapColumns(x, x+1)
				log.Printf("SWAPPING COLUMNS, %v with %v", x, x+1)
			}
		}
	}
}*/

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %v\nFPS: %v\nPC: %v", ebiten.ActualTPS(), ebiten.ActualFPS(), PARTICLE_COUNT))

	DrawGrid(screen, GRID)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREENWIDTH / 2, SCREENHEIGHT / 2
}

func Start() {

	ebiten.SetWindowSize(SCREENWIDTH, SCREENHEIGHT)
	ebiten.SetWindowTitle("Sand-game")
	ebiten.SetTPS(120) // double max TPS

	GRID = Grid{Width: SCREENHEIGHT, Height: SCREENHEIGHT}
	GRID.Map = PrepareGrid(SCREENWIDTH, SCREENHEIGHT, MOUSEX, MOUSEY)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

func (particle *Particle) PrepareParticle(MOUSEX, MOUSEY int) *Particle {

	result := Particle{Active: false, Position: particles.Vector2{X: MOUSEX, Y: MOUSEY}, Pixel: ebiten.NewImage(1, 1)}
	return &result
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

func SimulateParticles() {
	// For each particle
	for x := 0; x < GRID.Width; x++ {
		for y := 0; y < GRID.Height/2; y++ {

			p := GRID.Map[x][y]

			// if particle below is inactive
			if p.Active && y+1 < SCREENHEIGHT && !GRID.Map[x][y+1].Active {
				GRID.Map[x][y+1].Active = true
				GRID.Map[x][y].Active = false
				fmt.Printf("Moved particleY from %v to  %v\n", y, y-1)
			}

			/*if !GRID.Map[x][y+1].Active && !GRID.Map[x-1][y+1].Active && !GRID.Map[x+1][y+1].Active {
				SimulateCollision(GRID.Map[x][y])
			}*/
			// Sand effect
			/*if !GRID.Map[x-1][y-1].Active && GRID.Map[x][y-1].Active {
				p.Active = false
				GRID.Map[x-1][y-1].Active = true
			}*/

		}
	}
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

				//SimulateParticles()
			}

		}
	}
}

/*
func SimulateGravity(particle Particle) *Particle {

		pos := particle.Position
		floorY := (SCREENHEIGHT / 2) - 1

		if pos.Y < floorY {
			//fmt.Printf("Moving pixel @ %v to %v", particle.Position.Y, particle.Position.Y-1)
			pos.Y -= 1
			return &particle
		}
		return nil
	}
*/
func CheckForParticleSpawn(GRID Grid, MOUSEX int, MOUSEY int) {
	// If mouse0 pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton(0)) && MOUSEX >= 0 && MOUSEY >= 0 {
		// If particle pixel is INACTIVE
		particle := GRID.Map[MOUSEX][MOUSEY]

		if !particle.Active {
			// ACTIVATE particle pixel
			PARTICLE_COUNT++
			GRID.Map[MOUSEX][MOUSEY].Active = true
			log.Printf("Activating pixel @ [%v, %v] -- #%v", MOUSEX, MOUSEY, PARTICLE_COUNT)
		}
	}
}

/*
func SimulateCollision(particle Particle) *Particle {

	belowParticle := GRID.Map[particle.Position.X][particle.Position.Y-1]
	aboveParticle := GRID.Map[particle.Position.X][particle.Position.Y+1]
	belowLeftParticle := GRID.Map[particle.Position.X-1][particle.Position.Y-1]
	belowRightParticle := GRID.Map[particle.Position.X+1][particle.Position.Y-1]
	groundParticle := GRID.Map[particle.Position.X][particle.Position.Y-GRAVITY]

	if particle.Active && belowParticle.Active {

		if groundParticle.Active {
			particle.Position = particles.Vector2{X: groundParticle.Position.X, Y: groundParticle.Position.Y - 2}
		}

		// has particle above & below
		if aboveParticle.Active && belowParticle.Active {
			return &particle
		}

		// is there a particle below left or right? (sand effect)
		if !belowLeftParticle.Active {
			//particle.Position = belowLeftParticle.Position
			belowLeftParticle.Active = true
			particle.Active = false
		} else if !belowRightParticle.Active {
			//particle.Position = belowRightParticle.Position
			belowRightParticle.Active = true
			particle.Active = false
		}

	} else if particle.Active && !belowParticle.Active {

		// does not have particle below
		//particle.Position = belowParticle.Position
		particle.Active = false
		belowParticle.Active = true
	}
	return &particle
}
*/
