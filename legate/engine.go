package legate

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const (
	KindNoAlloc = "noalloc"
	KindTinyGo  = "tinygo"
	KindRust    = "rust"
	KindZig     = "zig"
)

const (
	RCOkay = iota
	RCUnexpectedTerm
)

var ErrInvalidAllocKind = errors.New("unknown 'kind' property given for target")
var ErrUndefinedHandler = errors.New("unable to locate indicated handler function in target wasm")
var ErrClosedModuleOnHandleExec = errors.New("wasm module was closed when call was attempted")

const (
	RTFailureModuleClosed     = 10000
	RTFailureHandlerExecution = 10001
)

type Engine struct {
	wasmCacheDir string
	wasmCache    wazero.CompilationCache
	runtime      wazero.Runtime
	runtimeCfg   wazero.RuntimeConfig
	runtimeCtx   context.Context
	maxTTL       time.Duration

	runtimeModules []RuntimeModule       // Go code callable from wasm
	modules        map[string]wasmModule // Wasm source binaries

	modulesName string
}

type DataReceiver func(data []byte) uint64

type wasmModule struct {
	runtime wazero.Runtime
	apiMod  api.Module
	handler DataReceiver
}

func New(wasmModuleName string, opts *Opts) (*Engine, error) {

	eng := Engine{}
	eng.modulesName = wasmModuleName
	eng.runtimeModules = make([]RuntimeModule, 0)
	eng.prepareEnvironment(opts.Ttl, opts.Pages)

	// Read target configs from disk
	for i, target := range opts.Targets {
		slog.Debug("load target", "idx", i, "path", target)
		t, e := loadTargetConfig(target)
		if e != nil {
			slog.Debug("failed to load target", "idx", i, "path", target, "error", e.Error())
			return nil, e
		}

		// Instantiate from config
		if err := eng.instantiateTarget(target, t); err != nil {
			return nil, err
		}
	}
	return &eng, nil
}

// Add runtime modules to the wasm environment so that they can
// be called by the target at execution time
func (engine *Engine) WithRuntimeModule(m RuntimeModule) *Engine {
	engine.runtimeModules = append(engine.runtimeModules, m)
	return engine
}

func loadTargetConfig(path string) (Target, error) {
	var target Target
	targetFile := filepath.Join(path, defaultTargetName)
	slog.Debug(targetFile)
	t, err := os.ReadFile(targetFile)
	if err != nil {
		slog.Error("failed to load target config", "error", err.Error())
		return target, err
	}
	err = yaml.Unmarshal(t, &target)
	if err != nil {
		slog.Error("failed to load target config", "error", err.Error())
		return target, err
	}
	return target, nil
}

func (engine *Engine) Close() {

	engine.wasmCache.Close(engine.runtimeCtx)
	os.RemoveAll(engine.wasmCacheDir)
}

func (engine *Engine) prepareEnvironment(ttl time.Duration, pages uint32) error {

	engine.modules = make(map[string]wasmModule)

	// Choose the context to use for function calls.
	engine.runtimeCtx = context.Background()

	var err error

	// Prepare a cache directory.
	engine.wasmCacheDir, err = os.MkdirTemp("", "wasm-cache")
	if err != nil {
		slog.Error("failed to create cache directory")
		return err
	}

	// Initializes the new compilation cache with the cache directory.
	// This allows the compilation caches to be shared even across multiple OS processes.
	engine.wasmCache, err = wazero.NewCompilationCacheWithDir(engine.wasmCacheDir)
	if err != nil {
		slog.Error("failed to create compilation cache with dir")
		return err
	}

	// Creates a shared runtime config to share the cache across multiple wazero.Runtime.
	engine.runtimeCfg = wazero.NewRuntimeConfig().
		WithCompilationCache(engine.wasmCache). // Cache
		WithMemoryLimitPages(pages).            // Set max memory in terms of pages
		WithMemoryCapacityFromMax(true)         // Auto allocate memory up-front to avoid reallocs

	engine.runtime = wazero.NewRuntimeWithConfig(engine.runtimeCtx, engine.runtimeCfg)

	engine.maxTTL = ttl
	return nil
}

func loadTargetWasm(path string, target Target) ([]byte, error) {

	expectedFileName := filepath.Join(path, fmt.Sprintf("%s.wasm", target.Name))

	slog.Debug("attempting to load wasm data from disk",
		"path", path,
		"file", expectedFileName)

	b, err := os.ReadFile(filepath.Join(path, defaultTargetName))
	if err != nil {
		slog.Error("failed to load wasm data", "error", err.Error())
	}
	return b, err
}

func (engine *Engine) instantiateTarget(path string, target Target) error {
	slog.Debug("instantiating target", "target", target.Name, "kind", target.Kind)

	// create an instance for the defined target that is scoped
	// by the terms of the prepared environment (above) and has access
	// to emit([]byte), and the list of functions that can work
	// with no-alloc. Some targets might not want to use alloc, so
	// the design should minimize the requirements of direct messaging (for now)

	wasm, err := loadTargetWasm(path, target)
	if err != nil {
		return err
	}

	switch target.Kind {
	case KindNoAlloc:
		break
	case KindTinyGo:
		err = engine.instantiateTargetTinyGo(target, wasm)
		break
	case KindRust:
		break
	case KindZig:
		break
	default:
		return ErrInvalidAllocKind
	}
	return err
}

func (engine *Engine) instantiateTargetNoAlloc(target Target, bin []byte) error {
	slog.Debug("instantiateTargetNoAlloc")

	// Target doesn't need memory allocations (meaning it can't accept binary data )

	return nil
}

func (engine *Engine) instantiateTargetTinyGo(target Target, bin []byte) error {
	slog.Debug("instantiateTargetTinyGo")
	/*

	  NOTE: If we have to reinstantiate later due to close issue we may be
	        able to do the following:


	      instantiate := func() {
	        ..
	        mod.ExportedFunction(target.Handler)
	        ...
	      }

	      instantiate()

	      wasmModule {
	        reinstantiate: instantiate
	      }
	*/

	mod, err := engine.instantiateWithRuntimeModules(bin)
	if err != nil {
		slog.Error("failed to instantiate target with runtime modules", "name", target.Name)
		return err
	}

	targetFn := mod.ExportedFunction(target.Handler)
	if targetFn == nil {
		return ErrUndefinedHandler
	}

	engine.modules[target.Name] = wasmModule{
		runtime: wazero.NewRuntimeWithConfig(engine.runtimeCtx, engine.runtimeCfg),
		handler: func(data []byte) uint64 {

			// Notes from wasmer docs:
			//
			// ExportedFunction returns a function exported from this module or nil if it wasn't.
			//
			// Note: The default wazero.ModuleConfig attempts to invoke `_start`, which
			// in rare cases can close the module. When in doubt, check IsClosed prior
			// to invoking a function export after instantiation.
			if mod.IsClosed() {
				slog.Error(
					"wasm module closed when attempting to execute target's handler",
					"target", target.Name, "kind", target.Kind, "handler", target.Handler)
				return RTFailureModuleClosed
			}

			ctx, cancel := context.WithTimeout(engine.runtimeCtx, engine.maxTTL)
			defer cancel()

			// TODO: Need to alloc space for `data` to be passed to module

			_, callErr := targetFn.Call(ctx) // TODO: Pass the data
			if callErr != nil {
				slog.Error("error attempting to call target handler",
					"target", target.Name,
					"handler", target.Handler)
				return RTFailureHandlerExecution
			}

			// TODO: do something eith result
			return 0
		},
	}
	return nil
}
