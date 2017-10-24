package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

	w, renderer, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	w.SetTitle("Flappy Bird")

	sceneManager, err := NewSceneManager(renderer, windowWidth, windowHeight)
	if err != nil {
		return fmt.Errorf("could not create scene manager: %v", err)
	}
	defer sceneManager.Destroy()

	events := make(chan sdl.Event)
	errc := sceneManager.Run(events, renderer)

	runtime.LockOSThread()
	for {
		event := sdl.PollEvent()
		if event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				close(events)
				<-errc
				return nil
			case *sdl.MouseButtonEvent:
				events <- event
			}
		}

		select {
		case e, ok := <-errc:
			if !ok {
				return nil
			}
			return fmt.Errorf("SceneManager got error %v", e)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
