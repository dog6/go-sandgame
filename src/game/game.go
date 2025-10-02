package game

import (
	"C"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu"

	"image/color"
	"git.smallzcomputing.com/sand-game/src/config"
	"git.smallzcomputing.com/sand-game/src/particles"
	"git.smallzcomputing.com/sand-game/src/sandgameUI"
	"git.smallzcomputing.com/sand-game/src/util"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

var (
	PARTICLE_COUNT            int = 0
	settingsUI                ebitenui.UI
	MOUSEX, MOUSEY            int
	MAX_PARTICLES             int = SCREENWIDTH * SCREENHEIGHT // max particles allowed on screen at once (in a perfect world this is SCREENWIDTH*SCREENHEIGHT)
	Config                    config.Configuration
	SCREENWIDTH, SCREENHEIGHT = 1600, 900
	VERSION                   string // version of game
	Conf                      config.Configuration
)

type Game struct {
	ui  *ebitenui.UI
	btn *widget.Button
}

var GRID util.Grid

// Updates game, called every frame
func (g *Game) Update() error {
	sandgameUI.UpdateGameInfoLabel(ebiten.ActualTPS(), ebiten.ActualFPS(), PARTICLE_COUNT)
	g.ui.Update()

	MOUSEX, MOUSEY = ebiten.CursorPosition() // Capture mouse position

	particles.CheckForParticleSpawn(GRID, MOUSEX, MOUSEY, &MAX_PARTICLES, &PARTICLE_COUNT /*&wg*/) // Check for particle spawn

	if Conf.RainRate > 0 {
		SpawnRain(Conf.RainRate)
	}

	particles.SimulateParticles(GRID, Conf.GRAVITY)
	return nil
}

// Renders game to screen, called every frame
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(Conf.BackgroundColor.ToColor())
	//ebitenutil.DebugPrint(screen, )
	g.ui.Draw(screen)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go particles.DrawGrid(screen, GRID, &wg)
	wg.Wait()
}

// Returns center of screen
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREENWIDTH / 2, SCREENHEIGHT / 2
}

// Called at start of game
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

	particles.Init(Config.ScreenSize, Conf.ShowSkippedParticles)

	// Log about rain
	if Config.RainRate != 0 {
		util.Log(fmt.Sprintf("Raining ENABLED -> %v drops/frame", Config.RainRate))
	} else {
		util.Log("Rain DISABLED")
	}

	// Setup UI
	eui := sandgameUI.SetupUI(Config)

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

// Prepare the grid before rendering
func PrepareGrid(width, height, MOUSEX, MOUSEY int, col color.RGBA) [][]util.Particle {
	util.Log(fmt.Sprintf("Particle color: R: %v, G: %v, B: %v, A: %v", col.R, col.G, col.B, col.A))
	result := make([][]util.Particle, width)
	for i := 0; i < width; i++ {
		result[i] = make([]util.Particle, height)

		for j := 0; j < height; j++ {
			result[i][j] = *result[i][j].PrepareParticle(MOUSEX, MOUSEY, col)
		}

	}
	util.Log(fmt.Sprintf("Grid size = [width: %v, height: %v]\n", len(result), len(result[0])))
	return result
}

// Spawn rain to help benchmark max particle count
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
