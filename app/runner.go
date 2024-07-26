package app

import (
	"fmt"
	"log/slog"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Job struct {
	Origin      string
	Destination []string
	Data        []byte

	// Tags []string
}

type Runner interface {
	Load(actionsPath string, actionMap map[string]string, exports interp.Exports) error
	SubmitJob(job *Job) error
}

type yaegiRunner struct {
	actions map[string]actionModule
}

type actionModule struct {
	env *interp.Interpreter

	// output buffer
	// input mechanism, etc
}

func (r *yaegiRunner) Load(actionsPath string, actionMap map[string]string, exports interp.Exports) error {
	slog.Debug("load runner directory", "actions", len(actionMap))

	r.actions = make(map[string]actionModule)

	for name, path := range actionMap {

		slog.Debug("loading module", "name", name, "path", path)

		m := actionModule{
			env: interp.New(interp.Options{GoPath: actionsPath}),
		}

		if err := m.env.Use(stdlib.Symbols); err != nil {
			slog.Error("yaegi failed to import stdlib symbols")
			return err
		}

		if err := m.env.Use(exports); err != nil {
			slog.Error("yaegi failed to import EMRS functionality")
			return err
		}

		_, err := m.env.EvalPath(path)
		if err != nil {
			slog.Error("yaegi failed to eval module file", "file", path, "error", err.Error())
			return err
		}

		vOnInitFn, getErr := m.env.Eval(fmt.Sprintf("%s.OnInit", name))
		if getErr == nil {
			onInitFn := vOnInitFn.Interface().(func() error)
			if err := onInitFn(); err != nil {
				slog.Error("error experienced within rootFile onInit function", "error", err.Error())
				return err
			}
		} else {
			slog.Info("module does not contain an `OnInit` function", "name", name)
		}

		// ---

		r.actions[name] = m
	}

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

	if len(job.Destination) < 2 {
		slog.Error("invalid destination; expected at least 'action.function'", "route-len", len(job.Destination))
		return
	}

	actionName := job.Destination[0]

	target, ok := r.actions[actionName]
	if !ok {
		slog.Error("unknown action for chunk in request", "name", actionName)
		return
	}

	symbols := target.env.Symbols(actionName)

	packageMap, pmOk := symbols[actionName]
	if !pmOk {
		slog.Error("failed to retrieve expected package from module", "expected-package", actionName)
		return
	}

	targetFnName := job.Destination[1]

	vTargetFn, ok := packageMap[targetFnName]

	if !ok {
		slog.Error("failed to locate function in action module", "module", actionName, "fn", targetFnName)
		return
	}

	targetFn := vTargetFn.Interface().(func(string, []string, []byte) error)
	if err := targetFn(job.Origin, job.Destination[2:], job.Data); err != nil {
		slog.Error("error experienced while processing function call", "module", actionName, "fn", targetFnName, "error", err.Error())
		return
	}

	slog.Debug("processing complete")
}
