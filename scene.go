package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	time        int
	bg          *sdl.Texture
	bird        *bird
	pipeTexture *sdl.Texture
	pipes       []*Pipe
	speed       int32
	width       int
	height      int
}

func newScene(r *sdl.Renderer, width, height int) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, fmt.Errorf("could not create bird: %v", err)
	}

	pipe, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}

	return &scene{
		bg:          bg,
		bird:        bird,
		pipeTexture: pipe,
		speed:       2,
		width:       width,
		height:      height,
	}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.generatePipes()
				s.moveScene()
				s.deleteHiddenPipes()
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) moveScene() {
	for _, pipe := range s.pipes {
		pipe.x -= 2
	}
}

func (s *scene) deleteHiddenPipes() {
	pipes := []*Pipe{}
	for _, p := range s.pipes {
		if p.x+p.width >= 0 {
			pipes = append(pipes, p)
		}
	}
	s.pipes = pipes

}

func (s *scene) generatePipes() {
	isNewPipesNeeded := false
	if len(s.pipes) == 0 {
		isNewPipesNeeded = true
	} else {
		lastPipe := s.pipes[len(s.pipes)-1]
		d := 300
		if s.width-(lastPipe.x+lastPipe.width) >= d {
			isNewPipesNeeded = true
		}
	}

	if isNewPipesNeeded {
		x := s.width
		spaceBetweenPipes := 100
		width := 52
		upperHeight := random(100, 400)
		bottomHeight := s.height - upperHeight - spaceBetweenPipes
		y := s.height - bottomHeight
		upperPipe := NewPipe(s.pipeTexture, x, 0, width, upperHeight, true)
		bottomPipe := NewPipe(s.pipeTexture, x, y, width, bottomHeight, false)
		s.pipes = append(s.pipes, upperPipe, bottomPipe)
	}
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyDownEvent:
		switch e.Keysym.Scancode {
		case sdl.SCANCODE_UP:
			s.bird.jump()
		}
	}
	return false
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return fmt.Errorf("could paint bird: %v", err)
	}

	for _, p := range s.pipes {
		if err := p.paint(r); err != nil {
			return fmt.Errorf("could not paint pipe: %v", err)
		}
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.pipeTexture.Destroy()
	s.bird.destroy()
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
