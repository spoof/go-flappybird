package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	minPipeHeight = 100
)

type pipe struct {
	texture          *sdl.Texture
	x                int32
	speed            int32
	width            int32
	topPipeHeight    int32
	bottomPipeHeight int32
	bottomPipeY      int32
}

func newPipe(texture *sdl.Texture, windowHeight int, x, speed int32) *pipe {
	space := random(120, 300)
	topPipeHeight := random(minPipeHeight, windowHeight-space-minPipeHeight)
	bottomPipeHeight := windowHeight - space - topPipeHeight
	bottomPipeY := topPipeHeight + space
	p := &pipe{
		texture:          texture,
		x:                x,
		speed:            speed,
		width:            52,
		topPipeHeight:    int32(topPipeHeight),
		bottomPipeHeight: int32(bottomPipeHeight),
		bottomPipeY:      int32(bottomPipeY),
	}
	return p
}

func (p *pipe) paint(r *sdl.Renderer) error {
	p.x -= p.speed

	rect := &sdl.Rect{X: p.x, Y: 0, W: p.width, H: p.topPipeHeight}
	if err := r.CopyEx(p.texture, nil, rect, 0, nil, sdl.FLIP_VERTICAL); err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}

	rect = &sdl.Rect{X: p.x, Y: p.bottomPipeY, W: p.width, H: p.bottomPipeHeight}
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
