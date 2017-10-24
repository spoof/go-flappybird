package scene

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/spoof/flappybird/scene/gameobj"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	distanceBetweenPipes = 300
	birdX                = 200
)

// Game is game scene
type Game struct {
	width  int
	height int

	bg          *sdl.Texture
	bird        *gameobj.Bird
	scoreFont   *ttf.Font
	pipeTexture *sdl.Texture
	pipeWidth   int
	pipePairs   []*gameobj.PipePair

	score      int
	bestScore  int
	isGameOver bool
}

// NewGame creates new Game scene
func NewGame(r *sdl.Renderer, width, height int) (*Game, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	bird, err := gameobj.NewBird(r, birdX, height/2)
	if err != nil {
		return nil, fmt.Errorf("could not create bird: %v", err)
	}

	pipe, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}

	_, _, pipeWidth, _, err := pipe.Query()
	if err != nil {
		return nil, fmt.Errorf("could not get pipe width: %v", err)
	}

	scoreFont, err := ttf.OpenFont("res/fonts/flappy.ttf", 42)
	if err != nil {
		return nil, fmt.Errorf("cound not load font: %v", err)
	}

	return &Game{
		width:  width,
		height: height,

		bg:          bg,
		bird:        bird,
		pipeTexture: pipe,
		pipeWidth:   int(pipeWidth),
		scoreFont:   scoreFont,
	}, nil
}

// Run runs the game scene.
func (g *Game) Run(in <-chan sdl.Event, r *sdl.Renderer) <-chan Event {
	out := make(chan Event)
	go func() {
		defer close(out)

		g.reset()

		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case event, ok := <-in:
				if !ok {
					return
				}
				g.handleEvent(event)
			case <-tick:
				if g.hasCollisions() {
					g.isGameOver = true
				}

				if !g.isGameOver {
					g.generatePipes()
					g.moveScene()
					g.updateScore()
					g.deleteHiddenPipes()
				} else {
					g.bird.Fall()
				}

				g.moveBird()

				if err := g.paint(r); err != nil {
					out <- &ErrorEvent{Err: err}
					return
				}

				if g.doesBirdHitsGround() && g.isGameOver {
					out <- &EndGameEvent{Score: g.score, BestScore: g.bestScore}
					return
				}

			}

		}

	}()
	return out
}

// Destroy frees all resources
func (g *Game) Destroy() {
	g.bg.Destroy()
	g.bird.Destroy()
	g.pipeTexture.Destroy()
}

func (g *Game) reset() {
	g.score = 0
	g.bird.ResetPosition()
	g.pipePairs = nil
	g.isGameOver = false
}

func (g *Game) hasCollisions() bool {
	if g.bird.Y <= 0 {
		return true
	}

	if g.doesBirdHitsGround() {
		return true
	}

	for _, pp := range g.pipePairs {
		if pp.Hits(g.bird) {
			return true
		}

	}

	return false
}

func (g *Game) doesBirdHitsGround() bool {
	if g.bird.Y+g.bird.Height >= g.height {
		return true
	}

	return false
}

func (g *Game) handleEvent(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.MouseButtonEvent:
		if e.Type != sdl.MOUSEBUTTONDOWN {
			return
		}
		if !g.isGameOver {
			g.bird.Jump()
		}
	}
}

func (g *Game) generatePipes() {
	needNewPipe := false
	if len(g.pipePairs) == 0 {
		needNewPipe = true
	} else {
		lastPipe := g.pipePairs[len(g.pipePairs)-1]
		if g.width-(lastPipe.X+lastPipe.Width) >= distanceBetweenPipes {
			needNewPipe = true
		}
	}

	if needNewPipe {
		x := g.width
		pipes := gameobj.NewPipePair(g.pipeTexture, x, int(g.pipeWidth), g.height)
		g.pipePairs = append(g.pipePairs, pipes)
	}
}

func (g *Game) moveBird() {
	if g.bird.Y+g.bird.Height <= g.height {
		g.bird.Move()
	}
}
func (g *Game) moveScene() {
	for _, pp := range g.pipePairs {
		pp.Move(-2)
	}
}

func (g *Game) updateScore() {
	for _, pp := range g.pipePairs {
		if !pp.Counted && pp.X+pp.Width < g.bird.X {
			pp.Counted = true
			g.score++
			g.bestScore = int(math.Max(float64(g.score), float64(g.bestScore)))
		}
	}
}

func (g *Game) deleteHiddenPipes() {
	pipes := []*gameobj.PipePair{}
	for _, pp := range g.pipePairs {
		if pp.X+pp.Width >= 0 {
			pipes = append(pipes, pp)
		}
	}
	g.pipePairs = pipes
}

func (g *Game) paint(renderer *sdl.Renderer) error {
	renderer.Clear()

	if err := renderer.Copy(g.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	drawOutline := false
	if err := g.bird.Paint(renderer, drawOutline); err != nil {
		return fmt.Errorf("could paint bird: %v", err)
	}

	for _, p := range g.pipePairs {
		if err := p.Paint(renderer); err != nil {
			return fmt.Errorf("could paint pipe: %v", err)
		}
	}

	if err := g.paintScore(renderer); err != nil {
		return fmt.Errorf("could not paint score: %v", err)
	}

	renderer.Present()
	return nil
}

func (g *Game) paintScore(renderer *sdl.Renderer) error {
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	text, err := g.scoreFont.RenderUTF8_Solid(strconv.Itoa(g.score), white)
	if err != nil {
		return fmt.Errorf("could not render score: %v", err)
	}
	defer text.Free()

	t, err := renderer.CreateTextureFromSurface(text)
	if err != nil {
		return fmt.Errorf("cound not create texture: %v", err)
	}
	defer t.Destroy()

	var clipRect sdl.Rect
	text.GetClipRect(&clipRect)
	rect := &sdl.Rect{X: int32(g.width)/2 - clipRect.W/2, Y: 60, W: clipRect.W, H: clipRect.H}

	if err := renderer.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("cound not copy texture: %v", err)
	}

	return nil
}
