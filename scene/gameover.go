package scene

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

// GameOver is game over scene
type GameOver struct {
	width  int
	height int

	captionFont *ttf.Font

	bestScore int
}

// NewGameOver creates new GameOver scene
func NewGameOver(r *sdl.Renderer, width, height int) (*GameOver, error) {
	captionFont, err := ttf.OpenFont("res/fonts/flappy.ttf", 42)
	if err != nil {
		return nil, fmt.Errorf("cound not load font: %v", err)
	}

	return &GameOver{
		width:       width,
		height:      height,
		captionFont: captionFont,
	}, nil
}

// Run run the scene's loop
func (gos *GameOver) Run(in <-chan sdl.Event, r *sdl.Renderer) <-chan Event {
	out := make(chan Event, 1)
	go func() {
		defer close(out)

		if err := gos.paint(r); err != nil {
			out <- &ErrorEvent{Err: err}
			return
		}

		for {
			select {
			case event, ok := <-in:
				if !ok {
					return
				}
				if gos.handleEvent(event) {
					out <- &StartGameEvent{}
					return
				}
			}
		}

	}()

	return out
}

// SetBestScore sets new best score of game
func (gos *GameOver) SetBestScore(bestScore int) {
	gos.bestScore = bestScore
}

// Destroy frees all resources used by GameOver scene
func (gos *GameOver) Destroy() {
	gos.captionFont.Close()
}

func (gos *GameOver) handleEvent(event sdl.Event) (playAgain bool) {
	switch e := event.(type) {
	case *sdl.MouseButtonEvent:
		if e.Type == sdl.MOUSEBUTTONDOWN {
			return true
		}
	}

	return false
}

func (gos *GameOver) paint(renderer *sdl.Renderer) error {
	rect := &sdl.Rect{X: 0, Y: 0, W: int32(gos.width), H: int32(gos.height)}

	renderer.SetDrawColor(0, 0, 0, 128)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.FillRect(rect)
	renderer.DrawRect(rect)

	if err := gos.paintCaption(renderer); err != nil {
		return fmt.Errorf("could not render caption: %v", err)
	}

	if err := gos.paintBestScoreCaption(renderer); err != nil {
		return fmt.Errorf("could not render best score caption: %v", err)
	}

	renderer.Present()
	return nil
}

func (gos *GameOver) paintCaption(renderer *sdl.Renderer) error {
	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	captionSurface, err := gos.captionFont.RenderUTF8_Solid("Game Over", c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer captionSurface.Free()

	t, err := renderer.CreateTextureFromSurface(captionSurface)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	var clipRect sdl.Rect
	captionSurface.GetClipRect(&clipRect)
	rect := &sdl.Rect{X: int32(gos.width)/2 - clipRect.W/2, Y: 200, W: clipRect.W, H: clipRect.H}

	if err := renderer.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}

	return nil
}

func (gos *GameOver) paintBestScoreCaption(renderer *sdl.Renderer) error {
	c := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	text := "Best Score: " + strconv.Itoa(gos.bestScore)
	captionSurface, err := gos.captionFont.RenderUTF8_Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer captionSurface.Free()

	t, err := renderer.CreateTextureFromSurface(captionSurface)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	var clipRect sdl.Rect
	captionSurface.GetClipRect(&clipRect)
	rect := &sdl.Rect{X: int32(gos.width)/2 - clipRect.W/2, Y: 300, W: clipRect.W, H: clipRect.H}

	if err := renderer.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}

	return nil
}
