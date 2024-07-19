package main

import (
	"fmt"
	"os"

	"github.com/hop-/goconfig"
	"github.com/hop-/goi/internal/app"
	"github.com/hop-/golog"
)

func getTlsOpts() (string, string, error) {
	certFile, err := goconfig.Get[string]("tls.certFile")
	if err != nil {
		return "", "", err
	}
	keyFile, err := goconfig.Get[string]("tls.keyFile")
	if err != nil {
		return "", "", err
	}

	return *certFile, *keyFile, nil
}

func getQuicOpts() (bool, int, error) {
	enabled, err := goconfig.Get[bool]("quic.enabled")
	if err != nil {
		return false, -1, err
	}

	port, err := goconfig.Get[int]("quic.port")
	if err != nil {
		return false, -1, err
	}

	return *enabled, *port, nil
}

func getTcpOtps() (bool, int, error) {
	enabled, err := goconfig.Get[bool]("tcp.enabled")
	if err != nil {
		return false, -1, err
	}

	port, err := goconfig.Get[int]("tcp.port")
	if err != nil {
		return false, -1, err
	}

	return *enabled, *port, nil
}

func getStorageOptions() (string, string, error) {
	storageType, err := goconfig.Get[string]("storage.type")
	if err != nil {
		return "", "", err
	}

	uri, err := goconfig.Get[string]("storage.uri")
	if err != nil {
		return "", "", err
	}

	return *storageType, *uri, nil
}

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

	certFile, keyFile, err := getTlsOpts()
	if err != nil {
		golog.Fatalf("Failed to get TLS configuration %s", err.Error())
	}
	quicEnabled, quicPort, err := getQuicOpts()
	if err != nil {
		golog.Fatalf("Failed to get QUIC configuration %s", err.Error())
	}
	tcpEnabled, tcpPort, err := getTcpOtps()
	if err != nil {
		golog.Fatalf("Failed to get TCP configuration %s", err.Error())
	}
	storageType, storageUri, err := getStorageOptions()
	if err != nil {
		golog.Fatalf("Failed to get storage configuration %s", err.Error())
	}

	opts := []app.OptionModifier{
		app.WithTls(certFile, keyFile),
		app.WithStorage(storageType, storageUri),
	}

	if !quicEnabled && !tcpEnabled {
		golog.Fatalf("At least one of the services should be enabled")
	}

	if quicEnabled {
		opts = append(opts, app.WithQuic(quicPort))
	}

	if tcpEnabled {
		opts = append(opts, app.WithTcp(tcpPort))
	}
	app, err := app.New(opts...)

	if err != nil {
		golog.Fatal("Failed to initialize app", err.Error())
	}

	// Run app
	app.Start()
}
