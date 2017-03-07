package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type bird struct {
	time     int
	textures []*sdl.Texture

	x      int
	y      int
	width  int
	height int
}

func newBird(r *sdl.Renderer, x, y int) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("cound not load bird texture: %v", err)
		}
		textures = append(textures, texture)
	}
	width := 50
	height := 43
	return &bird{textures: textures, x: x - width/2, y: y - height/2, width: width, height: height}, nil
}

func (b *bird) paint(r *sdl.Renderer, drawOutline bool) error {
	b.time++
	b.y++

	rect := &sdl.Rect{X: int32(b.x), Y: int32(b.y), W: int32(b.width), H: int32(b.height)}

	if drawOutline {
		r.SetDrawColor(255, 0, 0, 0)
		r.FillRect(rect)
		r.DrawRect(rect)
	}

	i := b.time / 8 % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	return nil
}

func (b *bird) jump() {
	b.y -= 25
}

func (b *bird) destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}
