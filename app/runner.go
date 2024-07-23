package app

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Job struct {
	Ctx         context.Context
	Origin      string
	Destination string
	Data        []byte

	// Tags []string
}

type Runner interface {
	Load(dir string, rootFile string) error
	SubmitJob(job *Job) error
}

type yaegiRunner struct {
	env      *interp.Interpreter
	routeMap map[string]func([]byte)
}

func (r *yaegiRunner) Load(dir string, rootFile string) error {
	slog.Debug("load runner directory", "dir", dir, "rootFile", rootFile)

	r.env = interp.New(interp.Options{
		GoPath: dir,
		/*

		   TODO:

		         Here we can add FS stuff to sandbox the thing in but give it
		         some access

		         We can also setup a buffer for the io so we can log it, etc

		         https://pkg.go.dev/github.com/traefik/yaegi@v0.16.1/interp#Options

		*/
	})

	if err := r.env.Use(stdlib.Symbols); err != nil {
		slog.Error("yaegi failed to import stdlib symbols")
		return err
	}

	// evaluate the root file (init.go)

	_, err := r.env.EvalPath(filepath.Join(dir, rootFile))
	if err != nil {
		slog.Error("yaegi failed to eval root file", "rootFile", rootFile, "error", err.Error())
		return err
	}

	// init actions

	vOnInitFn, err := r.env.Eval("actions.onInit")
	if err != nil {
		slog.Error("failed to retrieve onInit() function", "error", err.Error())
		return err
	}

	onInitFn := vOnInitFn.Interface().(func() error)

	if err := onInitFn(); err != nil {
		slog.Error("error experienced within rootFile onInit function", "error", err.Error())
		return err
	}

	// retrieve handlers map

	vRouteMap, err := r.env.Eval("actions.routeMap")
	if err != nil {
		slog.Error("failed to retrieve route map", "error", err.Error())
		return err
	}

	// This route map contains a key that represents the processing path (see README)
	// of an EMRS URL.

	r.routeMap = vRouteMap.Interface().(map[string]func([]byte))

	return nil
}

func (r *yaegiRunner) SubmitJob(job *Job) error {
	slog.Debug("submit job to runner")

	// TODO: We should add the context here that terminated the thread if it exceeds the
	//        maximum time from the server.cfg file

	go r.processJob(job)

	return nil
}

func (r *yaegiRunner) processJob(job *Job) {

	slog.Info("processing job", "from", job.Origin, "to", job.Destination)

	target, ok := r.routeMap[job.Destination]
	if !ok {
		// TODO: take ctx into account for auto timeout

		slog.Error("failed to locate job destination in route map")
		return
	}

	target(job.Data)
}
