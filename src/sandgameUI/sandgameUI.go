package sandgameUI

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"git.smallzcomputing.com/sand-game/src/config"
	"git.smallzcomputing.com/sand-game/src/util"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var GameInfoLabel *widget.Text

var ReloadConfigBtn *widget.Button

var (
	ColWidth, RowHeight = [3]int{}, [3]int{}
)

// Updates information displayed on the GameInfoLabel
func UpdateGameInfoLabel(tps, fps float64, particleCount int) {
	GameInfoLabel.Label = fmt.Sprintf("TPS: %.1f\n\nFPS: %.1f\n\nPC: %v", tps, fps, particleCount)
}

// Runs Inital setup of UI using given configuration, returns ebitenui.UI
func SetupUI(Conf *config.Configuration) ebitenui.UI {
	// This creates the root container for this UI.
	// All other UI elements must be added to this container.
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{0, 0, 0, 0})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// This adds the root container to the UI, so that it will be rendered.
	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	ReloadConfigBtn = Load_ReloadConfigBtn(Conf)
	GameInfoLabel = Load_GameInfoLabel(Conf)

	// To display the text widget, we have to add it to the root container.
	rootContainer.AddChild(GameInfoLabel)
	rootContainer.AddChild(ReloadConfigBtn)

	return *eui
}

// Loads fontFace and sets inital text for GameInfoLabel
func Load_GameInfoLabel(Conf *config.Configuration) *widget.Text {
	fontFace, err := loadFont(Conf.FontFilePath)
	if err != nil {
		util.LogErr(err)
	}
	// This creates a text widget that says "Hello World!"
	GameInfoLabel = widget.NewText(
		widget.TextOpts.Text("", fontFace, Conf.UITextColor.ToColor()),
	)
	return GameInfoLabel
}

func Load_ReloadConfigBtn(Conf *config.Configuration) *widget.Button {
	btnImg, _ := loadButtonImage()

	fontFace, _ := loadFont(Conf.FontFilePath)

	// construct a button
	button := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(btnImg),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Reload config", fontFace, &widget.ButtonTextColor{
			Idle:    color.RGBA{60, 60, 60, 255},
			Hover:   color.RGBA{80, 80, 80, 255},
			Pressed: color.RGBA{40, 40, 40, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   5,
			Right:  5,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			// Needs some sort of handler to call so values can be updated after reading config again
			Conf.ReadConfig()
			util.Log("Reloaded config")
		}),
	)
	return button
}

func loadFont(fontPath string) (font.Face, error) {

	data, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("%v\n", data)
	}

	ttfFont, err := truetype.Parse(data)
	if err != nil {
		util.LogErr(err)
	}
	fontFace := truetype.NewFace(ttfFont, &truetype.Options{
		Size:       12,
		DPI:        36,
		SubPixelsX: 12,
		SubPixelsY: 12,
	})
	return fontFace, err
}

func loadButtonImage() (*widget.ButtonImage, error) {

	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})
	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})
	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
