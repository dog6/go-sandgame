package util

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Vector2 Structure (X, Y int)
type Vector2 struct {
	X int `yaml:"X"`
	Y int `yaml:"Y"`
}

type RGBA struct {
	R int `yaml:"R"`
	G int `yaml:"G"`
	B int `yaml:"B"`
	A int `yaml:"A"`
}

func (c RGBA) ToColor() color.RGBA {
	r := uint8(c.R)
	g := uint8(c.G)
	b := uint8(c.B)
	a := uint8(c.A)

	return color.RGBA{R: r, G: g, B: b, A: a}
}

type Particle struct {
	Active   bool
	Position Vector2
	Pixel    *ebiten.Image
	Color    color.Color
}

type Grid struct {
	Width, Height int
	Map           [][]Particle
}

var VerboseLogging = true

/* [ Logging functions ] */
// Should always log erros
func LogErr(err error) {
	log.Printf("[ERROR] %v\n", err.Error())
}

func LogInfo(msg string) {
	if VerboseLogging {
		log.Println("[INFO] ", msg)
	}
}

func Log(msg string) {
	if VerboseLogging {
		log.Println(msg)
	}
}

func (particle *Particle) PrepareParticle(MOUSEX, MOUSEY int, color color.Color) *Particle {
	return &Particle{Active: false, Position: Vector2{X: MOUSEX, Y: MOUSEY}, Pixel: ebiten.NewImage(1, 1), Color: color}
}
