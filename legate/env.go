package legate

import (
	"context"
	"errors"
	"fmt"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"log/slog"
)

type Export struct {
	Name string
	Fn   interface{}
}

type RuntimeModule interface {
	GetExports() []Export
}

func (engine *Engine) instantiateWithRuntimeModules(target []byte) (api.Module, error) {

	// Build the runtime modules
	exports := make(map[string]bool)
	for _, m := range engine.runtimeModules {
		for _, export := range m.GetExports() {

			_, ok := exports[export.Name]
			if ok {
				slog.Error("duplicate name in exported function", "name", export.Name)
				return nil, errors.New(fmt.Sprintf("duplicate exported function name [%s]", export.Name))
			}

			_, err := engine.runtime.NewHostModuleBuilder(engine.modulesName).
				NewFunctionBuilder().WithFunc(export.Fn).Export(export.Name).
				Instantiate(engine.runtimeCtx)
			if err != nil {
				slog.Error("failed to instantiate module")
				return nil, err
			}

			exports[export.Name] = true
		}
	}

	// Instantiate target
	modCtx, modCancel := context.WithTimeout(
		engine.runtimeCtx,
		engine.maxTTL)
	defer modCancel()

	// TinyGo needs
	wasi_snapshot_preview1.MustInstantiate(modCtx, engine.runtime)

	mod, err := engine.runtime.Instantiate(
		modCtx,
		target)

	return mod, err
}
