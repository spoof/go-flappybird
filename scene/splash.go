package scene

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Splash is the first game scene
type Splash struct {
	bg         *sdl.Texture
	logoFont   *ttf.Font
	buttonFont *ttf.Font

	width  int
	height int
}

// NewSplash creates new TitleScreen
func NewSplash(r *sdl.Renderer, width, height int) (*Splash, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	logoFont, err := ttf.OpenFont("res/fonts/flappy.ttf", 26)
	if err != nil {
		return nil, fmt.Errorf("cound not load font: %v", err)
	}

	buttonFont, err := ttf.OpenFont("res/fonts/VanillaExtractRegular.ttf", 16)
	if err != nil {
		return nil, fmt.Errorf("cound not load font: %v", err)
	}

	return &Splash{
		bg:         bg,
		logoFont:   logoFont,
		buttonFont: buttonFont,
		width:      width,
		height:     height,
	}, nil
}

// Run runs the scene loop by listening events from in channel, renders all staff to r. Returns
// channel to listen for events from the scene.
func (s *Splash) Run(in <-chan sdl.Event, r *sdl.Renderer) <-chan Event {
	out := make(chan Event, 1)

	go func() {
		defer close(out)

		if err := s.paint(r); err != nil {
			out <- &ErrorEvent{Err: err}
			return
		}

		for {
			select {
			case e, ok := <-in:
				if !ok {
					return
				}

				switch event := e.(type) {
				case *sdl.MouseButtonEvent:
					if event.Type == sdl.MOUSEBUTTONDOWN {
						out <- &StartGameEvent{}
						return
					}
				}
			}
		}
	}()

	return out
}

// Destroy frees all resources
func (s *Splash) Destroy() {
	s.bg.Destroy()
	s.logoFont.Close()
	s.buttonFont.Close()
}

func (s *Splash) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := s.paintBg(r); err != nil {
		return fmt.Errorf("could not paint background: %v", err)
	}

	if err := s.paintLogo(r); err != nil {
		return fmt.Errorf("could not paint logo: %v", err)
	}
	if err := s.paintButton(r); err != nil {
		return fmt.Errorf("could not paint logo: %v", err)
	}

	r.Present()
	return nil
}

func (s *Splash) paintBg(r *sdl.Renderer) error {
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	return nil
}

func (s *Splash) paintLogo(r *sdl.Renderer) error {
	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	logoSurface, err := s.logoFont.RenderUTF8_Solid("Flappy Bird", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer logoSurface.Free()

	t, err := r.CreateTextureFromSurface(logoSurface)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	var clipRect sdl.Rect
	logoSurface.GetClipRect(&clipRect)
	rect := &sdl.Rect{X: 100 / 2, Y: 40, W: int32(s.width - 100), H: int32(s.height / 2)}
	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}

	return nil
}

func (s *Splash) paintButton(r *sdl.Renderer) error {
	c := sdl.Color{R: 150, G: 155, B: 45, A: 0}
	buttonSurface, err := s.buttonFont.RenderUTF8_Solid("Press any key to start", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer buttonSurface.Free()

	t, err := r.CreateTextureFromSurface(buttonSurface)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	var clipRect sdl.Rect
	buttonSurface.GetClipRect(&clipRect)
	rect := &sdl.Rect{X: 200 / 2, Y: 400, W: int32(s.width - 200), H: int32(80)}
	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}

	return nil
}
