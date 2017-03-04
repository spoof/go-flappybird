package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	texture *sdl.Texture
	x       int32
	speed   int32
	width   int32
}

func newPipe(texture *sdl.Texture, x, speed int32) *pipe {
	p := &pipe{texture: texture, x: x, speed: speed, width: 52}
	return p
}

func (p *pipe) paint(r *sdl.Renderer) error {
	p.x -= p.speed
	rect := &sdl.Rect{X: p.x, Y: windowHeight - 100, W: p.width, H: 100}
	if err := r.Copy(p.texture, nil, rect); err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}
	return nil
}

func (p *pipe) isHidden() bool {
	if p.x+p.width < 0 {
		return true
	}
	return false
}
