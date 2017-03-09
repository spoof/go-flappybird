package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

const (
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

	w.SetTitle("Flappy Bird")

	titleScreen, err := NewTitleScreen(r, windowWidth, windowHeight)
	if err != nil {
		return fmt.Errorf("could not create title screen: %v", err)
	}

	events := make(chan sdl.Event)
	titleExitc := titleScreen.run(events, r)

	s, err := newScene(r, windowWidth, windowHeight)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case <-titleExitc:
			return nil
			// errc := s.run(events, r)
			// case err := <-errc:
			// return err
		}
	}
}
