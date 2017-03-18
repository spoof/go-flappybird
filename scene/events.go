package scene

type Event interface{}

type QuitEvent struct{}

type ErrorEvent struct {
	Err error
}

type StartGameEvent struct{}

type EndGameEvent struct {
	Score     int
	BestScore int
}
