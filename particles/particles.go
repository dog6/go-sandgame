package particles

import (
	"fmt"
	"image/color"

	"git.smallzcomputing.com/sand-game/util"
	"github.com/hajimehoshi/ebiten/v2"
)

func IsParticleStable(GRID util.Grid, showSkipped bool, x, y int) bool {
	// Check if bottom at screen
	if y == GRID.Height/2-1 {
		if showSkipped {
			GetParticle(GRID, x, y).Color = color.RGBA{255, 120, 120, 255}
		}
		return true
	}

	// Check if has 3 particles below ( cannot fall ) ( needs RNG, favors falling left atm )
	if GetParticle(GRID, x-1, y+1).Active && GetParticle(GRID, x+1, y+1).Active && GetParticle(GRID, x, y+1).Active {
		if showSkipped {
			GetParticle(GRID, x, y).Color = color.RGBA{255, 120, 120, 255}
		}
		return true
	}

	return false
}

func DrawParticle(GRID util.Grid, renderer *ebiten.Image, x, y int) {
	renderer.Set(x, y, GetParticle(GRID, x, y).Color)
}

func SetParticle(GRID util.Grid, x, y int, isActive bool) {
	GetParticle(GRID, x, y).Active = isActive
}

func GetParticle(GRID util.Grid, x, y int) *util.Particle {
	return &GRID.Map[x][y]
}

func SpawnParticle(GRID util.Grid, maxP, pCount *int, x, y int) {
	if *pCount+1 <= *maxP {
		// ACTIVATE particle pixel
		*pCount++
		SetParticle(GRID, x, y, true)
		util.Log(fmt.Sprintf("Activating pixel @ [%v, %v] -- #%v", x, y, &pCount))
	}
}

func DisableParticle(pCount *int, GRID util.Grid, x, y int) {
	if GetParticle(GRID, x, y).Active {
		SetParticle(GRID, x, y, false)
		*pCount--
		util.Log(fmt.Sprintf("Deactivating pixel @ [%v, %v] -- #%v", x, y, &pCount))
	}
}
