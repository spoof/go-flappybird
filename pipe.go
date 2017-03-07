package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	minPipeHeight = 100
)

// Pipe is a base model for pipe
type Pipe struct {
	texture *sdl.Texture
	x       int
	y       int
	width   int
	height  int
	isUpper bool
}

// NewPipe creates new Pipe
func NewPipe(texture *sdl.Texture, x, y, width, height int, isUpper bool) *Pipe {
	p := &Pipe{
		texture: texture,
		isUpper: isUpper,
		x:       x,
		y:       y,
		width:   width,
		height:  height,
	}
	return p
}

func (p *Pipe) paint(r *sdl.Renderer) error {
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
