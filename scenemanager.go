package main

import (
	"fmt"

	"github.com/spoof/flappybird/scene"
	"github.com/veandco/go-sdl2/sdl"
)

// Scene is an interface for all scenes
type Scene interface {
	Run(<-chan sdl.Event, *sdl.Renderer) <-chan scene.Event
	Destroy()
}

// SceneManager represents main object for managing scenes
type SceneManager struct {
	splash   *scene.Splash
	game     *scene.Game
	gameOver *scene.GameOver

	currentScene Scene
	sceneEvents  chan sdl.Event
	sceneOutc    <-chan scene.Event
}

// NewSceneManager creates new SceneManager
func NewSceneManager(r *sdl.Renderer, w, h int) (*SceneManager, error) {
	splashScene, err := scene.NewSplash(r, w, h)
	if err != nil {
		return nil, fmt.Errorf("could not create Splash scene %v", err)
	}

	gameScene, err := scene.NewGame(r, w, h)
	if err != nil {
		return nil, fmt.Errorf("could not create Game scene %v", err)
	}

	gameOverScene, err := scene.NewGameOver(r, w, h)
	if err != nil {
		return nil, fmt.Errorf("could not create Gamve Over scene%v", err)
	}

	return &SceneManager{
		splash:   splashScene,
		game:     gameScene,
		gameOver: gameOverScene,
	}, nil
}

// Run starts the loop
func (sm *SceneManager) Run(events <-chan sdl.Event, renderer *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)

		sm.currentScene = sm.splash
		sm.sceneEvents = make(chan sdl.Event)
		sceneOutc := sm.currentScene.Run(sm.sceneEvents, renderer)

		for {
			select {
			case e, ok := <-events:
				if !ok {
					close(sm.sceneEvents)
					<-sceneOutc
					return
				}
				sm.sceneEvents <- e

			case e := <-sceneOutc:
				switch event := e.(type) {
				case *scene.ErrorEvent:
					errc <- fmt.Errorf("Error from scene %v", event.Err)

				case *scene.StartGameEvent:
					<-sceneOutc
					sceneOutc = sm.game.Run(sm.sceneEvents, renderer)

				case *scene.EndGameEvent:
					<-sceneOutc
					sm.gameOver.SetBestScore(event.BestScore)
					sceneOutc = sm.gameOver.Run(sm.sceneEvents, renderer)
				}
			}
		}
	}()

	return errc
}

// Destroy frees all resources of SceneManager
func (sm *SceneManager) Destroy() {
	sm.splash.Destroy()
	sm.game.Destroy()
}
