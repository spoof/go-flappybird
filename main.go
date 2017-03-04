package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

var (
	windowWidth  = 800
	windowHeight = 600
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

	w, r, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)

	}
	defer w.Destroy()

	if err := drawTtitle(r); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}
	time.Sleep(2 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case err := <-s.run(ctx, r):
		return err
	case <-time.After(30 * time.Second):
		return nil
	}

}

func drawTtitle(r *sdl.Renderer) error {
	r.Clear()
	f, err := ttf.OpenFont("res/fonts/flappy.ttf", 20)
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
