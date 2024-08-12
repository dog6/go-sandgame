package particles

import (
	"fmt"
	"image/color"
	"sync"

	"git.smallzcomputing.com/sand-game/src/util"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ShowSkippedParticles bool = false
)

func Init(screenSize util.Vector2, showSkipped bool) {
	ShowSkippedParticles = showSkipped
}

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
		util.Log(fmt.Sprintf("Activating pixel @ [%v, %v] -- #%v", x, y, *pCount))
	}
}

func DisableParticle(pCount *int, GRID util.Grid, x, y int) {
	if GetParticle(GRID, x, y).Active {
		SetParticle(GRID, x, y, false)
		*pCount--
		util.Log(fmt.Sprintf("Deactivating pixel @ [%v, %v] -- #%v", x, y, *pCount))
	}
}

func DrawColLength(renderer *ebiten.Image, GRID util.Grid, x, y *int) {
	var tmpY int
	particleEnabled := true
	for tmpY = *y; particleEnabled; tmpY-- {
		if !GetParticle(GRID, *x, tmpY).Active {
			particleEnabled = false
			break
		}
		//particles.GetParticle(GRID, x, tmpY).Color = color.RGBA{0, 0, 255, 255} // Colors particles drawn as group blue
		DrawParticle(GRID, renderer, *x, tmpY)

	}
	//util.Log(fmt.Sprintf("Skipped %v particle draw cycles in column %v", y-tmpY, x))
	*y -= *y - tmpY // += crashes

}

func DrawGrid(renderer *ebiten.Image, GRID util.Grid, wg *sync.WaitGroup) {
	defer wg.Done()
	// draw from bottom right to top left
	// Loop through all grid positions
	for x := GRID.Width; x > 0; x-- {
		for y := GRID.Height - 1; y > 0; y-- {

			if GetParticle(GRID, x, y).Active {

				if GetParticle(GRID, x, y+1).Active {
					DrawColLength(renderer, GRID, &x, &y)
				} else {
					DrawParticle(GRID, renderer, x, y)
				}
			}

		}

	}
}

// func Setutil.Particle(particle *util.Particle, isActive bool) {
func CheckForParticleSpawn(GRID util.Grid, MOUSEX, MOUSEY int, maxParticles, particleCount *int) {
	// If mouse0 pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton(0)) && MOUSEX >= 0 && MOUSEY >= 0 {
		// If particle pixel is INACTIVE
		if !GetParticle(GRID, MOUSEX, MOUSEY).Active {
			SpawnParticle(GRID, maxParticles, particleCount, MOUSEX, MOUSEY)
		} else {
			DisableParticle(particleCount, GRID, MOUSEX, MOUSEY)
		}
	}
}

func SimulateParticles(GRID util.Grid, GRAVITY int) {
	// For each particle
	for x := GRID.Width; x > 0; x-- {
		for y := GRID.Height / 2; y > 0; y-- {

			if GetParticle(GRID, x, y).Active && y > 0 {

				if IsParticleStable(GRID, ShowSkippedParticles, x, y) {
					continue
				}

				// Check if can fall
				if !GRID.Map[x][y+GRAVITY].Active {
					SetParticle(GRID, x, y, false)
					SetParticle(GRID, x, y+GRAVITY, true)
				} else {

					// Sand effect
					if !GRID.Map[x-1][y+GRAVITY].Active {
						SetParticle(GRID, x, y, false)
						SetParticle(GRID, x-1, y+GRAVITY, true)
					} else if !GRID.Map[x+1][y+GRAVITY].Active {
						SetParticle(GRID, x, y, false)          // disable this particle
						SetParticle(GRID, x+1, y+GRAVITY, true) // set particle below to active
					} else {
						continue
					}
				}
			}
		}
	}
}
