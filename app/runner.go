package app

import (
	"log/slog"
	"path/filepath"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const (
  expectedModuleEntry = "_emrs.go"
)

type Job struct {
	Origin      string
	Destination []string
	Data        []byte

	// Tags []string
}

type Runner interface {
  Import(imports []string) error
	Load(dir string, rootFile string, exports interp.Exports) error
	SubmitJob(job *Job) error
}

type yaegiRunner struct {
	env    *interp.Interpreter
	onData func(string, []string, []byte) error

  actionsPath string
}

func (r *yaegiRunner) Import(imports []string) error {
  for idx, path := range imports {
    if err := r.importEmrsModule(path); err != nil {
      slog.Error("failed to import emrs module", "idx", idx, "name", path, "error", err.Error())
      return err
    }
  }
  return nil
}

func (r *yaegiRunner) Load(dir string, rootFile string, exports interp.Exports) error {
	slog.Debug("load runner directory", "dir", dir, "rootFile", rootFile)

  r.actionsPath = dir
	r.env = interp.New(interp.Options{
		GoPath: r.actionsPath,
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

	if err := r.env.Use(exports); err != nil {
		slog.Error("yaegi failed to import EMRS functionality")
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

	vOnData, err := r.env.Eval("actions.onData")
	if err != nil {
		slog.Error("failed to retrieve onInit() function", "error", err.Error())
		return err
	}

	r.onData = vOnData.Interface().(func(string, []string, []byte) error)

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
	evalErr := r.onData(job.Origin, job.Destination, job.Data)
	if evalErr != nil {
		slog.Error("evaluateion error reported", "error", evalErr.Error())
	}
}

func (r *yaegiRunner) importEmrsModule(name string) error {

  path := filepath.Join(r.actionsPath, name)

  slog.Debug("load emrs module", "name", name, "suspected-path", path)

  entry := filepath.Join(path, expectedModuleEntry)

  slog.Debug("attempting to load entry file", "name", entry)

	_, err := r.env.EvalPath(entry)
	if err != nil {
		slog.Error("yaegi failed to eval import", "file", entry, "error", err.Error())
		return err
	}

	vList, err := r.env.Eval("emrs.Exports")
	if err != nil {
		slog.Error("failed to retrieve Exports list", "error", err.Error())
		return err
	}

	list := vList.Interface().([]string)

  // Not technically an error, but we should let it be known
  if len(list) == 0 {
    slog.Warn("no imports listed from import", "file", entry)
    return nil
  }

  for idx, path := range list {

    if err := r.importEmrsFile(path); err != nil {
      slog.Error("error importing file for emrs module", "idx", idx, "path", path, "error", err.Error())
      return err
    }
  }
  return nil
}

func (r *yaegiRunner) importEmrsFile(file string) error {
  slog.Debug("import emrs module file", "file", file)
  _, err := r.env.EvalPath(file)
	if err != nil {
		slog.Error("yaegi failed to eval emrs import", "file", file, "error", err.Error())
		return err
	}
  return nil
}
