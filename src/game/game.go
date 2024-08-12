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

	"git.smallzcomputing.com/sand-game/src/config"
	"git.smallzcomputing.com/sand-game/src/particles"
	"git.smallzcomputing.com/sand-game/src/sandgameUI"
	"git.smallzcomputing.com/sand-game/src/util"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
)

var (
	settingsUI                ebitenui.UI
	MOUSEX, MOUSEY            int
	MAX_PARTICLES             int = SCREENWIDTH * SCREENHEIGHT // max particles allowed on screen at once (in a perfect world this is SCREENWIDTH*SCREENHEIGHT)
	Config                    config.Configuration
	SCREENWIDTH, SCREENHEIGHT = 1600, 900
	VERSION                   string // version of game
	Conf                      config.Configuration
)

type Game struct {
	ui *ebitenui.UI
}

var GRID util.Grid

// CONST GAME VARIABLES
const GRAVITY = 1

func (g *Game) Update() error {
	sandgameUI.UpdateUI(g)

	MOUSEX, MOUSEY = ebiten.CursorPosition() // Capture mouse position

	particles.CheckForParticleSpawn(GRID, MOUSEX, MOUSEY, &MAX_PARTICLES, &particles.PARTICLE_COUNT /*&wg*/) // Check for particle spawn

	if Conf.RainRate > 0 {
		SpawnRain(Conf.RainRate)
	}

	particles.SimulateParticles(GRID, GRAVITY)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(Conf.BackgroundColor.ToColor())
	//ebitenutil.DebugPrint(screen, )

	g.ui.Draw(screen)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go particles.DrawGrid(screen, GRID, &wg)
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
		Size:       24,
		DPI:        72,
		SubPixelsX: 8,
		SubPixelsY: 8,
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

	particles.Init(MAX_PARTICLES, Config.ScreenSize, Conf.ShowSkippedParticles)

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

func SpawnRain(spawnRate int) {
	if particles.PARTICLE_COUNT+spawnRate <= MAX_PARTICLES {
		for drops := 0; drops < spawnRate; drops++ {
			particles.PARTICLE_COUNT++
			xPos := rand.Intn(SCREENWIDTH)
			yPos := rand.Intn(50) //rand.Intn(GRID.Height/2-2) + 100

			if !GRID.Map[xPos][yPos].Active {
				GRID.Map[xPos][yPos].Active = true
			}
		}
	}
}
