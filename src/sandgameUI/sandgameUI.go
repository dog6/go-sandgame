package sandgameUI

import (
	"fmt"

	"git.smallzcomputing.com/sand-game/src/game"
	"git.smallzcomputing.com/sand-game/src/particles"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/shiny/widget"
)

var GameInfoLabel *widget.Text

func UpdateUI(g *game.Game) {
	GameInfoLabel.Label = fmt.Sprintf("TPS: %.1f | FPS: %.1f | PC: %v", ebiten.ActualTPS(), ebiten.ActualFPS(), particles.PARTICLE_COUNT)
	g.ui.Update()
}
