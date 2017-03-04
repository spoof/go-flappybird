package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}
func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)

	}
	defer w.Destroy()
	_ = r

	if err := drawTtitle(r); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}
	time.Sleep(5 * time.Second)

	if err := drawBackground(r); err != nil {
		return fmt.Errorf("could not draw background: %v", err)
	}
	time.Sleep(5 * time.Second)

	return nil
}

func drawTtitle(r *sdl.Renderer) error {
	r.Clear()
	f, err := ttf.OpenFont("res/fonts/FlappyBirdy.ttf", 20)
	if err != nil {
		return fmt.Errorf("cound not load font: %v", err)
	}
	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8_Solid("Flappy Gopher", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}
	r.Present()

	return nil
}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()

	t, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return fmt.Errorf("cound not load background image: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	r.Present()
	return nil
}
