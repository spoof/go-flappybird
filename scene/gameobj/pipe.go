package gameobj

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	minPipeHeight = 100
)

// Pipe is a base model for pipe
type pipe struct {
	texture *sdl.Texture
	x       int
	y       int
	width   int
	height  int
	isUpper bool

	counted bool
}

// NewPipe creates new Pipe
func newPipe(texture *sdl.Texture, x, y, width, height int, isUpper bool) *pipe {
	p := &pipe{
		texture: texture,
		isUpper: isUpper,
		x:       x,
		y:       y,
		width:   width,
		height:  height,
	}
	return p
}

func (p *pipe) hits(b *Bird) bool {
	if p.x < b.X+b.Width &&
		p.x+p.width > b.X &&
		p.y < b.Y+b.Height &&
		p.y+p.height > b.Y {
		return true
	}
	return false
}

func (p *pipe) paint(r *sdl.Renderer) error {
	flip := sdl.FLIP_NONE
	if p.isUpper {
		flip = sdl.FLIP_VERTICAL
	}

	rect := &sdl.Rect{X: int32(p.x), Y: int32(p.y), W: int32(p.width), H: int32(p.height)}
	if err := r.CopyEx(p.texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}
	return nil
}

const (
	spaceBetweenPipes = 160
)

// PipePair is a pair of pipes
type PipePair struct {
	X       int
	Width   int
	Counted bool

	top    *pipe
	bottom *pipe
}

// NewPipePair creates new PipePair with given position x and width
func NewPipePair(texture *sdl.Texture, x, width, windowHeight int) *PipePair {
	topHeight := random(minPipeHeight, windowHeight-minPipeHeight-spaceBetweenPipes)
	bottomHeight := windowHeight - topHeight - spaceBetweenPipes
	bottomY := windowHeight - bottomHeight
	topPipe := newPipe(texture, x, 0, width, topHeight, true)
	bottomPipe := newPipe(texture, x, bottomY, width, bottomHeight, false)

	pp := &PipePair{
		X:     x,
		Width: width,

		top:    topPipe,
		bottom: bottomPipe,
	}
	return pp
}

// Hits checks if bird hits any pipe
func (pp *PipePair) Hits(b *Bird) bool {
	if pp.top.hits(b) || pp.bottom.hits(b) {
		return true
	}
	return false
}

// Move moves pipepair by given x
func (pp *PipePair) Move(x int) {
	pp.X += x
	pp.top.x += x
	pp.bottom.x += x
}

// Paint paints the pair or pipes using r render
func (pp *PipePair) Paint(r *sdl.Renderer) error {
	if err := pp.top.paint(r); err != nil {
		return fmt.Errorf("top pipe: %v", err)
	}

	if err := pp.bottom.paint(r); err != nil {
		return fmt.Errorf("top pipe: %v", err)
	}

	return nil
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
