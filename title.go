package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// TitleScreen is the first screen of game
type TitleScreen struct {
	bg *sdl.Texture

	width  int
	height int
}

// NewTitleScreen creates new TitleScreen
func NewTitleScreen(r *sdl.Renderer, width, height int) (*TitleScreen, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	return &TitleScreen{
		bg:     bg,
		width:  width,
		height: height,
	}, nil
}

func (ts *TitleScreen) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)

		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := ts.handleEvent(e); done {
					return
				}
			case <-tick:
				ts.paint(r)

			}
		}
	}()

	return errc
}

func (ts *TitleScreen) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		fmt.Printf("key pressed %+v", e)
		return true
	}

	return false
}

func (ts *TitleScreen) paint(r *sdl.Renderer) {
	r.Clear()
	ts.paintBg(r)
	ts.paintLogo(r)
	ts.paintButton(r)
	r.Present()
}

func (ts *TitleScreen) paintBg(r *sdl.Renderer) error {
	if err := r.Copy(ts.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	return nil
}

func (ts *TitleScreen) paintLogo(r *sdl.Renderer) error {
	f, err := ttf.OpenFont("res/fonts/flappy.ttf", 26)
	if err != nil {
		return fmt.Errorf("cound not load font: %v", err)
	}
	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8_Solid("Flappy Bird", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	var clipRect sdl.Rect
	s.GetClipRect(&clipRect)
	rect := &sdl.Rect{X: 100 / 2, Y: 40, W: int32(ts.width - 100), H: int32(ts.height / 2)}
	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}

	return nil
}

func (ts *TitleScreen) paintButton(r *sdl.Renderer) error {
	f, err := ttf.OpenFont("res/fonts/VanillaExtractRegular.ttf", 16)
	if err != nil {
		return fmt.Errorf("cound not load font: %v", err)
	}
	c := sdl.Color{R: 150, G: 155, B: 45, A: 0}
	s, err := f.RenderUTF8_Solid("Press any key to start", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	var clipRect sdl.Rect
	s.GetClipRect(&clipRect)
	rect := &sdl.Rect{X: 200 / 2, Y: 400, W: int32(ts.width - 200), H: int32(80)}
	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}

	return nil
}
