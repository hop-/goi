package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	exitChan  chan os.Signal
	mu        *sync.Mutex
	isRunning bool
	// TODO
}

func newApp(opts Options) *App {
	// TODO: remove tmp line
	_ = opts

	app := App{
		make(chan os.Signal, 1),
		&sync.Mutex{},
		false,
	}

	// Signal handling
	signal.Notify(app.exitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return &app
}

func New(optionModifiers ...OptionModifier) *App {
	o := defaultOptions()
	for _, omd := range optionModifiers {
		omd(&o)
	}

	return newApp(o)
}

func (a *App) Start() {
	a.isRunning = true
	// Graceful shutdown
	go a.gracefulShutDownTracker()

	// TODO: run

	a.Stop()
}

func (a *App) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isRunning {
		return
	}

	a.isRunning = false
	// TODO: stop
}

func (a *App) gracefulShutDownTracker() {
	<-a.exitChan

	a.Stop()
}
