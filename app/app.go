package app

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type Opts struct {
	Binding string
}

type httpsInfo struct {
	keyPath  string
	certPath string
}

type App struct {
	binding string

	httpsSettings *httpsInfo // nil if not using https
}

func New(options *Opts) *App {
	return &App{
		binding: options.Binding,
	}
}

func (a *App) WithHttps(keyPath string, certPath string) *App {
	a.httpsSettings = &httpsInfo{
		keyPath:  keyPath,
		certPath: certPath,
	}
	return a
}

func (a *App) Run() {

	slog.Info("Application starting", "binding", a.binding)

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

	var err error
	if a.httpsSettings != nil {
		slog.Info("Using TLS")
		err = http.ListenAndServeTLS(
			a.binding,
			a.httpsSettings.certPath,
			a.httpsSettings.keyPath,
			nil)
	} else {
		slog.Warn("Not using TLS")
		err = http.ListenAndServe(a.binding, nil)
	}

	if errors.Is(err, http.ErrServerClosed) {
	} else if err != nil {
		slog.Error("error starting the server", "error", err.Error())
		os.Exit(1)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}
