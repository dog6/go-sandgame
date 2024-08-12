package game

import (
	"C"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu"
)
import (
	"image/color"
	"os"

	"git.smallzcomputing.com/sand-game/config"
	"git.smallzcomputing.com/sand-game/particles"
	"git.smallzcomputing.com/sand-game/util"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
)

var (
	settingsUI                ebitenui.UI
	MOUSEX, MOUSEY            int
	PARTICLE_COUNT            int
	MAX_PARTICLES             int = SCREENWIDTH * SCREENHEIGHT // max particles allowed on screen at once (in a perfect world this is SCREENWIDTH*SCREENHEIGHT)
	Config                    config.Configuration
	SCREENWIDTH, SCREENHEIGHT = 1600, 900
	VERSION                   string // version of game
	Conf                      config.Configuration
	ShowSkippedParticles      = false // renders particles not being simulated as red
)

type Game struct {
	ui *ebitenui.UI
}

var GRID util.Grid

var GameInfoLabel *widget.Text

// CONST GAME VARIABLES
const GRAVITY = 1

func (g *Game) Update() error {
	GameInfoLabel.Label = fmt.Sprintf("TPS: %.2f\nFPS: %.2f\nPC: %v", ebiten.ActualTPS(), ebiten.ActualFPS(), PARTICLE_COUNT)
	g.ui.Update()

	MOUSEX, MOUSEY = ebiten.CursorPosition() // Capture mouse position

	CheckForParticleSpawn(GRID, MOUSEX, MOUSEY /*&wg*/) // Check for particle spawn

	if Conf.RainRate > 0 {
		SpawnRain(Conf.RainRate)
	}

	SimulateParticles()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(Conf.BackgroundColor.ToColor())
	//ebitenutil.DebugPrint(screen, )

	g.ui.Draw(screen)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go DrawGrid(screen, GRID, &wg)
	wg.Wait()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREENWIDTH / 2, SCREENHEIGHT / 2
}

func SetupUI() ebitenui.UI {
	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer()

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	data, err := os.ReadFile(Conf.FontFilePath)

	if err != nil {
		log.Fatalf("%v\n", data)
	}

	// This loads a font and creates a font face.
	ttfFont, err := truetype.Parse(data)
	if err != nil {
		log.Fatal("Error Parsing Font", err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size: 24,
	})

	// This creates a text widget that says "Hello World!"
	GameInfoLabel = widget.NewText(
		widget.TextOpts.Text("", fontFace, Conf.UITextColor.ToColor()),
	)

	// To display the text widget, we have to add it to the root container.
	rootContainer.AddChild(GameInfoLabel)

	return *eui
}

func Start(Config *config.Configuration) {

	// Log config
	Conf = *Config
	VERSION = Config.VersionNumber // set version number
	Config.LogConfig()

	// Set window size
	SCREENWIDTH, SCREENHEIGHT = Config.ScreenSize.X, Config.ScreenSize.Y
	util.Log(fmt.Sprintf("Setting window size to X: %v, Y: %v", SCREENWIDTH, SCREENHEIGHT))
	ebiten.SetWindowSize(SCREENWIDTH, SCREENHEIGHT)
	ebiten.SetWindowTitle(fmt.Sprintf("Sandgame %v", VERSION))
	ebiten.SetTPS(Config.MaxTPS) // double max TPS

	// Set max particles ( check for 0, which = SCREENWIDTH*SCREENHEIGHT)
	if Config.MaxParticles != 0 {
		MAX_PARTICLES = Config.MaxParticles
	} else {
		MAX_PARTICLES = SCREENWIDTH * SCREENHEIGHT
	}
	ShowSkippedParticles = Config.ShowSkippedParticles

	// Log about rain
	if Config.RainRate != 0 {
		util.Log(fmt.Sprintf("Raining ENABLED -> %v drops/frame", Config.RainRate))
	} else {
		util.Log("Rain DISABLED")
	}

	// Setup UI
	eui := SetupUI()

	// Prepare grid
	GRID = util.Grid{Width: SCREENHEIGHT, Height: SCREENHEIGHT}
	GRID.Map = PrepareGrid(SCREENWIDTH, SCREENHEIGHT, MOUSEX, MOUSEY, Config.ParticleColor.ToColor())

	game := Game{
		ui: &eui,
	}

	// Run game
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func PrepareGrid(width, height, MOUSEX, MOUSEY int, col color.RGBA) [][]util.Particle {
	util.Log(fmt.Sprintf("Particle color: R: %v, G: %v, B: %v, A: %v", col.R, col.G, col.B, col.A))
	result := make([][]util.Particle, width)
	for i := 0; i < width; i++ {
		result[i] = make([]util.Particle, height)

		for j := 0; j < height; j++ {
			result[i][j] = *result[i][j].PrepareParticle(MOUSEX, MOUSEY, col /*color.RGBA{R:255, G:255, B: 255, A: 255}*/)
		}

	}
	util.Log(fmt.Sprintf("Grid size = [width: %v, height: %v]\n", len(result), len(result[0])))
	return result
}

func SimulateParticles() {
	// For each particle
	for x := GRID.Width; x > 0; x-- {
		for y := GRID.Height / 2; y > 0; y-- {

			if particles.GetParticle(GRID, x, y).Active && y > 0 {

				if particles.IsParticleStable(GRID, ShowSkippedParticles, x, y) {
					continue
				}

				// Check if can fall
				if !GRID.Map[x][y+GRAVITY].Active {
					particles.SetParticle(GRID, x, y, false)
					particles.SetParticle(GRID, x, y+GRAVITY, true)
				} else {

					// Sand effect
					if !GRID.Map[x-1][y+GRAVITY].Active {
						particles.SetParticle(GRID, x, y, false)
						particles.SetParticle(GRID, x-1, y+GRAVITY, true)
					} else if !GRID.Map[x+1][y+GRAVITY].Active {
						particles.SetParticle(GRID, x, y, false)          // disable this particle
						particles.SetParticle(GRID, x+1, y+GRAVITY, true) // set particle below to active
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
func IsSurrounded(particle *util.Particle) bool {
	pos := particle.Position
	p1 := Getutil.Particle(pos.X-1, pos.Y+1).Active
	p2 := Getutil.Particle(pos.X, pos.Y+1).Active
	p3 := Getutil.Particle(pos.X+1, pos.Y+1).Active
	p4 := Getutil.Particle(pos.X-1, pos.Y).Active
	p6 := Getutil.Particle(pos.X+1, pos.Y).Active
	p7 := Getutil.Particle(pos.X-1, pos.Y-1).Active
	p8 := Getutil.Particle(pos.X, pos.Y-1).Active
	p9 := Getutil.Particle(pos.X+1, pos.Y-1).Active

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
			Drawutil.Particle(renderer, x-1, y+1)
			break
		case 2:
			Drawutil.Particle(renderer, x, y+1)
			break
		case 3:
			Drawutil.Particle(renderer, x+1, y+1)
			break
		case 4:
			Drawutil.Particle(renderer, x-1, y)
			break
		case 5:
			Drawutil.Particle(renderer, x, y)
			break
		case 6:
			Drawutil.Particle(renderer, x+1, y)
			break
		case 7:
			Drawutil.Particle(renderer, x-1, y-1)
			break
		case 8:
			Drawutil.Particle(renderer, x, y-1)
			break
		case 9:
			Drawutil.Particle(renderer, x+1, y-1)
			break
		}
	}
}
*/

func DrawGrid(renderer *ebiten.Image, GRID util.Grid, wg *sync.WaitGroup) {
	defer wg.Done()
	// draw from bottom right to top left
	// Loop through all grid positions
	for x := GRID.Width; x > 0; x-- {
		for y := GRID.Height - 1; y > 0; y-- {

			if particles.GetParticle(GRID, x, y).Active {
				particles.DrawParticle(GRID, renderer, x, y)
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

// func Setutil.Particle(particle *util.Particle, isActive bool) {
func CheckForParticleSpawn(GRID util.Grid, MOUSEX int, MOUSEY int) {
	// If mouse0 pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton(0)) && MOUSEX >= 0 && MOUSEY >= 0 {
		// If particle pixel is INACTIVE
		if !particles.GetParticle(GRID, MOUSEX, MOUSEY).Active {
			particles.SpawnParticle(GRID, &MAX_PARTICLES, &PARTICLE_COUNT, MOUSEX, MOUSEY)
		} else {
			particles.DisableParticle(&PARTICLE_COUNT, GRID, MOUSEX, MOUSEY)
		}
	}
}
