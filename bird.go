package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type bird struct {
	time     int
	textures []*sdl.Texture

	y float64
	x float64
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("cound not load bird texture: %v", err)
		}
		textures = append(textures, texture)
	}
	return &bird{textures: textures, y: 300}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y--

	y := 600 - int32(b.y) - 43/2
	if y >= 600-43 {
		y = 600 - 43
	}
	if y < 0 {
		y = 0
	}
	rect := &sdl.Rect{X: int32(b.x), Y: y, W: 50, H: 43}

	i := b.time / 8 % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	return nil
}

func (b *bird) jump() {
	b.y += 25
}

func (b *bird) forward() {
	b.x += 10
}

func (b *bird) back() {
	b.x -= 10
}

func (b *bird) destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}
