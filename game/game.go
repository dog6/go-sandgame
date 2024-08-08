package game

import (
	"fmt"
	"log"

	"git.smallzcomputing.com/sand-game/particles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

type Vector2 struct {
	X, Y int
}

type Particle struct {
	Active   bool
	Position Vector2
	Pixel    *ebiten.Image
}

const (
	SCREENWIDTH, SCREENHEIGHT = 1280, 720
	GRAVITY                   = 10
)

var (
	GRID           particles.Grid
	MOUSEX, MOUSEY int
	PARTICLE_COUNT int
)

func (g *Game) Update() error {
	MOUSEX, MOUSEY = ebiten.CursorPosition()              // Capture mouse position
	particles.CheckForParticleSpawn(GRID, MOUSEX, MOUSEY) // Check for particle spawn
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %v\nFPS: %v", ebiten.ActualTPS(), ebiten.ActualFPS()))

	particles.DrawGrid(screen, GRID)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREENWIDTH / 2, SCREENHEIGHT / 2
}

func Start() {

	ebiten.SetWindowSize(SCREENWIDTH, SCREENHEIGHT)
	ebiten.SetWindowTitle("Sand-game")
	ebiten.SetTPS(120) // double max TPS

	GRID = particles.Grid{Width: SCREENHEIGHT, Height: SCREENHEIGHT}
	GRID.Map = particles.PrepareGrid(SCREENWIDTH, SCREENHEIGHT, MOUSEX, MOUSEY)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

func SimulateParticles() {
	// For each particle
	for x := 0; x < GRID.Width-1; x++ {
		for y := 0; y < (GRID.Height/2)-1; y++ {

			if GRID.Map[x][y].Active && !GRID.Map[x][y+1].Active {

				GRID.Map[x][y].Active = false
				if GRID.Map[x][y+1].Active {
					GRID.Map[x][y+1].Active = true
				}

				fmt.Printf("Moved particle to %v, %v\n", x, y+1)

			}
		}
	}
}
