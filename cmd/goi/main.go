package main

import (
	"fmt"
	"os"

	"github.com/hop-/goconfig"
	"github.com/hop-/goi/internal/app"
	"github.com/hop-/golog"
)

func main() {
	// Load config
	if err := goconfig.Load(); err != nil {
		fmt.Printf("Failed to load configs %s\n", err.Error())
		os.Exit(1)
	}

	logMode, err := goconfig.Get[string]("logs.mode")
	if err != nil {
		mode := "INFO"
		fmt.Printf("Failed to get log mode default is %s\n", mode)
		logMode = &mode
	}
	// Init Logging
	golog.Init(*logMode)

	opts := []app.OptionModifier{
		// TODO: add
	}

	app := app.New(opts...)

	// Run app
	app.Start()
}
