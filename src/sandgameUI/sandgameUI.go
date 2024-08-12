package sandgameUI

import (
	"fmt"
	"log"
	"os"

	"git.smallzcomputing.com/sand-game/src/config"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
)

var GameInfoLabel *widget.Text

func UpdateUI(tps, fps float64, particleCount int) {
	GameInfoLabel.Label = fmt.Sprintf("TPS: %.1f | FPS: %.1f | PC: %v", tps, fps, particleCount)
}

func SetupUI(Conf *config.Configuration) ebitenui.UI {
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
		DPI:        36,
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
