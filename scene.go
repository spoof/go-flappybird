package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type scene struct {
	time  int
	bg    *sdl.Texture
	bird  *bird
	pipe  *sdl.Texture
	pipes []*pipe
	speed int32
}

func newScene(r *sdl.Renderer) (*scene, error) {
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

	return &scene{bg: bg, bird: bird, pipe: pipe, speed: 2}, nil
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
				s.generatePipe()
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) generatePipe() {
	pipes := make([]*pipe, len(s.pipes))
	copy(pipes, s.pipes)
	for i, p := range pipes {
		if p.isHidden() {
			s.pipes = append(s.pipes[:i], s.pipes[i+1:]...)
		}
	}

	if len(s.pipes) > 0 {
		lastPipe := s.pipes[len(s.pipes)-1]
		d := random(2, 21) * 100
		if windowWidth-(lastPipe.x+lastPipe.width) >= int32(d) {
			pipe := newPipe(s.pipe, windowHeight, windowWidth, s.speed)
			s.pipes = append(s.pipes, pipe)
		}

	} else {
		pipe := newPipe(s.pipe, windowHeight, windowWidth, s.speed)
		s.pipes = append(s.pipes, pipe)
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
		case sdl.SCANCODE_LEFT:
			s.bird.back()
		case sdl.SCANCODE_RIGHT:
			s.bird.forward()
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
	s.pipe.Destroy()
	s.bird.destroy()
}

func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
