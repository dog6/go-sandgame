package game

import (
	"C"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"sync"

	particles "git.smallzcomputing.com/sand-game/particles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "github.com/silbinarywolf/preferdiscretegpu"
)
import (
	"git.smallzcomputing.com/sand-game/config"
	"git.smallzcomputing.com/sand-game/util"
)

type Grid struct {
	Width, Height int
	Map           [][]Particle
}

var (
	GRID                      Grid
	MOUSEX, MOUSEY            int
	PARTICLE_COUNT            int
	Config                    config.Configuration
	SCREENWIDTH, SCREENHEIGHT        = 1600, 900
	MAX_PARTICLES             int    = SCREENWIDTH * SCREENHEIGHT // max particles allowed on screen at once (in a perfect world this is SCREENWIDTH*SCREENHEIGHT)
	VERSION                   string                              // version of game
	Conf                      config.Configuration
	ShowSkippedParticles      = false // renders particles not being simulated as red
)

type Game struct{}

type Particle struct {
	Active   bool
	Position particles.Vector2
	Pixel    *ebiten.Image
	Color    color.Color
}

// CONST GAME VARIABLES
const GRAVITY = 1

func (g *Game) Update() error {

	MOUSEX, MOUSEY = ebiten.CursorPosition() // Capture mouse position

	CheckForParticleSpawn(GRID, MOUSEX, MOUSEY /*&wg*/) // Check for particle spawn

	if Conf.RainAmount > 0 {
		SpawnRain(Conf.RainAmount)
	}

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

func Start(Config *config.Configuration) {

	// Log config
	Conf = *Config
	VERSION = Config.VersionNumber
	Config.LogConfig()
	ebiten.SetWindowSize(SCREENWIDTH, SCREENHEIGHT)
	ebiten.SetWindowTitle(fmt.Sprintf("Sandgame %v", VERSION))
	ebiten.SetTPS(Config.MaxTPS) // double max TPS
	MAX_PARTICLES = Config.MaxParticles
	SCREENWIDTH, SCREENHEIGHT = Config.ScreenWidth, Config.ScreenHeight
	ShowSkippedParticles = Config.ShowSkippedParticles
	// Log about rain
	if Config.RainAmount != 0 {
		util.Log(fmt.Sprintf("Raining ENABLED -> %v drops/frame", Config.RainAmount))
	} else {
		util.Log("Rain DISABLED")
	}

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
	util.Log(fmt.Sprintf("Grid size = [width: %v, height: %v]\n", len(result), len(result[0])))
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
	/*if GetParticle(x, y-1).Active && GetParticle(x, y+1).Active {
		if ShowSkippedParticles {
			GetParticle(x, y).Color = color.RGBA{255, 120, 120, 255}
		}
		return true
	}
	return false*/
	return false
}

func SimulateParticles() {
	// For each particle
	for x := GRID.Width; x > 0; x-- {
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
Cool, but inefficient

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

TODO: Chunking
func IsSurrounded(particle *Particle) bool {
	pos := particle.Position
	p1 := GetParticle(pos.X-1, pos.Y+1).Active
	p2 := GetParticle(pos.X, pos.Y+1).Active
	p3 := GetParticle(pos.X+1, pos.Y+1).Active
	p4 := GetParticle(pos.X-1, pos.Y).Active
	p6 := GetParticle(pos.X+1, pos.Y).Active
	p7 := GetParticle(pos.X-1, pos.Y-1).Active
	p8 := GetParticle(pos.X, pos.Y-1).Active
	p9 := GetParticle(pos.X+1, pos.Y-1).Active

	if p1 && p2 && p3 && p4 && p6 && p7 && p8 && p9 {
		return true
	} else {
		return false
	}
}

func Draw3x3Chunk(renderer *ebiten.Image, GRID Grid, x, y int) {
	for i := 0; i < 9; i++ {
		switch i {
		case 1:
			DrawParticle(renderer, x-1, y+1)
			break
		case 2:
			DrawParticle(renderer, x, y+1)
			break
		case 3:
			DrawParticle(renderer, x+1, y+1)
			break
		case 4:
			DrawParticle(renderer, x-1, y)
			break
		case 5:
			DrawParticle(renderer, x, y)
			break
		case 6:
			DrawParticle(renderer, x+1, y)
			break
		case 7:
			DrawParticle(renderer, x-1, y-1)
			break
		case 8:
			DrawParticle(renderer, x, y-1)
			break
		case 9:
			DrawParticle(renderer, x+1, y-1)
			break
		}
	}
}
*/

func DrawParticle(renderer *ebiten.Image, x, y int) {
	renderer.Set(x, y, GetParticle(x, y).Color)
}

func DrawGrid(renderer *ebiten.Image, GRID Grid, wg *sync.WaitGroup) {
	defer wg.Done()
	// draw from bottom right to top left
	// Loop through all grid positions
	for x := GRID.Width; x > 0; x-- {
		for y := GRID.Height - 1; y > 0; y-- {

			if GetParticle(x, y).Active {
				DrawParticle(renderer, x, y)
			}

		}

	}
}

func SpawnRain(spawnRate int) {
	if PARTICLE_COUNT+spawnRate <= MAX_PARTICLES {
		for drops := 0; drops < spawnRate; drops++ {
			PARTICLE_COUNT++
			xPos := rand.Intn(SCREENWIDTH)
			yPos := rand.Intn(50) //rand.Intn(GRID.Height/2-2) + 100

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
		util.Log(fmt.Sprintf("Activating pixel @ [%v, %v] -- #%v", x, y, PARTICLE_COUNT))
	}
}

func DisableParticle(x, y int) {
	if GetParticle(x, y).Active {
		SetParticle(x, y, false)
		PARTICLE_COUNT--
		util.Log(fmt.Sprintf("Deactivating pixel @ [%v, %v] -- #%v", MOUSEX, MOUSEY, PARTICLE_COUNT))
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
