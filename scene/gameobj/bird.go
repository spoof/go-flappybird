package gameobj

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity = 0.1
)

// Bird is a main character of this game
type Bird struct {
	time     int
	textures []*sdl.Texture

	X         int
	Y         int
	Width     int
	Height    int
	speedY    float32
	angle     float64
	isJumping bool

	startX int
	startY int
}

// NewBird creates new bird object
func NewBird(r *sdl.Renderer, x, y int) (*Bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("cound not load bird texture: %v", err)
		}
		textures = append(textures, texture)
	}

	_, _, birdWidth, birdHeight, err := textures[0].Query()
	if err != nil {
		return nil, fmt.Errorf("could not get bird texure info: %v", err)
	}
	width := int(birdWidth)
	height := int(birdHeight)

	bird := &Bird{textures: textures, startX: x, startY: y, Width: width, Height: height}
	bird.ResetPosition()

	return bird, nil
}

// ResetPosition resets position of bird to the start one
func (b *Bird) ResetPosition() {
	b.X = b.startX - b.Width/2
	b.Y = b.startY - b.Height/2
	b.speedY = 0
	b.angle = 0
}

// Jump makes bird jump
func (b *Bird) Jump() {
	if b.isJumping {
		b.speedY--
		return
	}

	b.isJumping = true
	b.angle = 0
	b.speedY = -4
}

// Fall makes bird fall
func (b *Bird) Fall() {
	b.speedY = 10
}

// Move moves bird
func (b *Bird) Move() {
	b.speedY += gravity
	b.Y += int(b.speedY)

	if b.isJumping && b.speedY >= 0 {
		b.isJumping = false
		b.speedY = 0
	}
}

// Paint paints the bird.
func (b *Bird) Paint(r *sdl.Renderer, drawOutline bool) error {
	b.time++

	rect := &sdl.Rect{X: int32(b.X), Y: int32(b.Y), W: int32(b.Width), H: int32(b.Height)}
	if drawOutline {
		r.SetDrawColor(255, 0, 0, 0)
		r.FillRect(rect)
		r.DrawRect(rect)
	}

	i := b.time / 8 % len(b.textures)
	if b.speedY >= 5 && b.angle < 90 {
		b.angle += 3.0
	}
	if err := r.CopyEx(b.textures[i], nil, rect, b.angle, nil, sdl.FLIP_NONE); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	return nil
}

// Destroy frees all resources of Bird
func (b *Bird) Destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}
