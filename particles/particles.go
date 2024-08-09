package particles

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Vector2 struct {
	X, Y int
}

type Particle struct {
	Active   bool
	Position Vector2
	Pixel    *ebiten.Image
}
