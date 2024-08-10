package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hop-/goi/internal/infra"
	"github.com/hop-/goi/internal/services"
	"github.com/hop-/goi/internal/storages"
	"github.com/hop-/golog"
)

type App struct {
	exitChan    chan os.Signal
	mu          *sync.Mutex
	isRunning   bool
	services    []services.Service
	storageType string
	storageUri  string
	wg          *sync.WaitGroup
}

func newApp(opts Options) (*App, error) {
	srvs := []services.Service{}

	certFile := opts.tls.certFile
	keyFile := opts.tls.keyFile

	if opts.quic != nil {
		s, err := services.NewQuicService(opts.quic.Port, certFile, keyFile)
		if err != nil {
			return nil, err
		}

		srvs = append(srvs, s)
	}

	if opts.tcp != nil {
		s, err := services.NewTcpService(opts.tcp.Port, certFile, keyFile)
		if err != nil {
			return nil, err
		}

		srvs = append(srvs, s)
	}

	app := App{
		make(chan os.Signal, 1),
		&sync.Mutex{},
		false,
		srvs,
		opts.storage.storageType,
		opts.storage.uri,
		&sync.WaitGroup{},
	}

	// Signal handling
	signal.Notify(app.exitChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return &app, nil
}

func New(optionModifiers ...OptionModifier) (*App, error) {
	o := defaultOptions()
	for _, omd := range optionModifiers {
		omd(&o)
	}

	return newApp(o)
}

func (a *App) Start() {
	a.isRunning = true
	// Graceful shutdown handler
	go a.gracefulShutDownHandler()

	// Init storage
	err := storages.InitStorage(a.storageType, a.storageUri)
	if err != nil {
		golog.Fatal("Failed to initialize the storage", err.Error())
	}
	defer storages.GetStorage().Close()

	err = infra.LoadData()
	if err != nil {
		golog.Fatal("Failed to load data", err.Error())
	}

	for _, s := range a.services {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			err := s.Start()
			if err != nil {
				golog.Error("Failed to start service", err.Error())
			}
		}()
	}

	// Wait until all goroutines are done
	a.wg.Wait()
	a.isRunning = false
}

func (a *App) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isRunning {
		return
	}

	a.isRunning = false

	// Iterate and stop all services
	for _, s := range a.services {
		s.Stop()
	}
}

func (a *App) gracefulShutDownHandler() {
	<-a.exitChan

	a.Stop()
}
