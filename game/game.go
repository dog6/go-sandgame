package game

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"sync"

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
var ShowSkippedParticles = true // colors any particles that aren't being simulated red

type Game struct{}

type Particle struct {
	Active   bool
	Position particles.Vector2
	Pixel    *ebiten.Image
	Color    color.Color
}

// CONST GAME VARIABLES
const (
	SCREENWIDTH, SCREENHEIGHT     = 540, 360
	GRAVITY                       = 1
	MAX_PARTICLES             int = 20000 // max particles allowed on screen at once (in a perfect world this is SCREENWIDTH*SCREENHEIGHT)
)

func (g *Game) Update() error {

	MOUSEX, MOUSEY = ebiten.CursorPosition() // Capture mouse position

	CheckForParticleSpawn(GRID, MOUSEX, MOUSEY /*&wg*/) // Check for particle spawn

	SpawnRain(5) // laggy atm
	SimulateParticles()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %v\nFPS: %v\nPC: %v", ebiten.ActualTPS(), ebiten.ActualFPS(), PARTICLE_COUNT))

	wg := sync.WaitGroup{}

	wg.Add(1)
	go DrawGrid(screen, GRID, &wg)
	wg.Wait()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREENWIDTH / 2, SCREENHEIGHT / 2
}

func Start() {

	ebiten.SetWindowSize(SCREENWIDTH, SCREENHEIGHT)
	ebiten.SetWindowTitle("Sand-game")
	ebiten.SetTPS(60) // double max TPS

	GRID = Grid{Width: SCREENHEIGHT, Height: SCREENHEIGHT}
	GRID.Map = PrepareGrid(SCREENWIDTH, SCREENHEIGHT, MOUSEX, MOUSEY)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func (particle *Particle) PrepareParticle(MOUSEX, MOUSEY int) *Particle {
	result := Particle{Active: false, Position: particles.Vector2{X: MOUSEX, Y: MOUSEY}, Pixel: ebiten.NewImage(1, 1), Color: color.RGBA{255, 255, 255, 255}}
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

func IsParticleStable(x, y int) bool {
	// Check if bottom at screen
	if y == GRID.Height/2-1 {
		if ShowSkippedParticles {
			GetParticle(x, y).Color = color.RGBA{255, 120, 120, 255}
		}
		return true
	}

	// Check if has 3 particles below ( cannot fall )
	if GetParticle(x-1, y+1).Active && GetParticle(x+1, y+1).Active && GetParticle(x, y+1).Active {
		if ShowSkippedParticles {
			GetParticle(x, y).Color = color.RGBA{255, 120, 120, 255}
		}
		return true
	}

	// Check if particle above & below
	if GetParticle(x, y-1).Active && GetParticle(x, y+1).Active /*|| GetParticle(x, y+1).Active && (GetParticle(x+1, y).Active) && (GetParticle(x-1, y).Active)*/ {
		if ShowSkippedParticles {
			GetParticle(x, y).Color = color.RGBA{255, 120, 120, 255}
		}
		return true
	}
	return false
}

func SimulateParticles() {
	// For each particle
	for x := GRID.Width - 1; x > 1; x-- {
		for y := GRID.Height / 2; y > 0; y-- {

			if GetParticle(x, y).Active && y > 0 {

				if IsParticleStable(x, y) {
					continue
				}

				// Check if can fall
				if !GRID.Map[x][y+GRAVITY].Active {
					SetParticle(x, y, false)
					SetParticle(x, y+GRAVITY, true)
				} else {

					// Sand effect
					if !GRID.Map[x-1][y+GRAVITY].Active {
						SetParticle(x, y, false)
						SetParticle(x-1, y+GRAVITY, true)
					} else if !GRID.Map[x+1][y+GRAVITY].Active {
						SetParticle(x, y, false)          // disable this particle
						SetParticle(x+1, y+GRAVITY, true) // set particle below to active
					} else {
						continue
					}
				}
			}
		}
	}
}

/*
	func IndexToPos(idx int) (int, int) {
		x := idx % SCREENWIDTH
		var y int
		if idx > SCREENWIDTH {
			y = idx / SCREENWIDTH
		} else {
			y = 1
		}
		return x, y
	}
*/
func DrawGrid(renderer *ebiten.Image, GRID Grid, wg *sync.WaitGroup) {
	defer wg.Done()
	// Loop through all grid positions
	for x := GRID.Width; x > 0; x-- {
		for y := GRID.Height - 1; y > 0; y-- {
			//for i := 0; i < SCREENHEIGHT*SCREENWIDTH; i++ {

			//	x, y := IndexToPos(i)

			if GetParticle(x, y).Active {
				// Drawing
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x), float64(y))
				GetParticle(x, y).Pixel.Fill(GetParticle(x, y).Color)
				renderer.Set(x, y, GetParticle(x, y).Color)
				//renderer.DrawImage(GetParticle(x, y).Pixel, op)

			}
		}

	}
}

func SpawnRain(spawnRate int) {
	if PARTICLE_COUNT+spawnRate <= MAX_PARTICLES {
		for drops := 0; drops < spawnRate; drops++ {
			PARTICLE_COUNT++
			xPos := rand.Intn(SCREENWIDTH)
			yPos := rand.Intn(GRID.Height/2-2) + 100

			if !GRID.Map[xPos][yPos].Active {
				GRID.Map[xPos][yPos].Active = true
			}
		}
	}
}

// func SetParticle(particle *Particle, isActive bool) {
func SetParticle(x, y int, isActive bool) {
	GetParticle(x, y).Active = isActive
}

func GetParticle(x, y int) *Particle {
	return &GRID.Map[x][y]
}

func SpawnParticle(x, y int) {
	if PARTICLE_COUNT+1 <= MAX_PARTICLES {
		// ACTIVATE particle pixel
		PARTICLE_COUNT++
		SetParticle(x, y, true)
		log.Printf("Activating pixel @ [%v, %v] -- #%v", x, y, PARTICLE_COUNT)
	}
}

func DisableParticle(x, y int) {
	if GetParticle(x, y).Active {
		SetParticle(x, y, false)
		PARTICLE_COUNT--
		log.Printf("Deactivating pixel @ [%v, %v] -- #%v", MOUSEX, MOUSEY, PARTICLE_COUNT)
	}
}

func CheckForParticleSpawn(GRID Grid, MOUSEX int, MOUSEY int) {
	// If mouse0 pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton(0)) && MOUSEX >= 0 && MOUSEY >= 0 {
		// If particle pixel is INACTIVE
		if !GetParticle(MOUSEX, MOUSEY).Active {
			SpawnParticle(MOUSEX, MOUSEY)
		} else {
			DisableParticle(MOUSEX, MOUSEY)
		}
	}
}
