/*

  This "script" is for development purposes only.

  Once we have subcommands built into the emrs application then we can
  have something like:


        emrs mod --install /path/to/mod


  When that is created, it will update the configurations with the requisite legate opts


        emrs mod --uninstall


  Potential CNC functionality:


        emrs cnc --reload-modules



*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	defaultEnvHome    = "EMRS_HOME"
	defaultInstallDir = "modules"
	defaultLegateYaml = "legate.yaml"
)

var modules = []string{
	"echo",
}

type target struct {
	path     string
	name     string
	wasm     string
	yamlFile string
}

func main() {
	emrsHome := flag.String("home", "", "Home directory")
	flag.Parse()
	if *emrsHome == "" {
		fromEnv := os.Getenv(defaultEnvHome)
		if fromEnv == "" {
			fmt.Println("unable to locate EMRS home directory")
			os.Exit(1)
		}
		*emrsHome = fromEnv
	}

	modulesDir := filepath.Join(*emrsHome, defaultInstallDir)
	os.RemoveAll(modulesDir)

	os.MkdirAll(modulesDir, 0755)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("working directory:", cwd)

	for _, module := range modules {
		path := filepath.Join(cwd, module)

		installDir := filepath.Join(modulesDir, module)
		wasmFile := fmt.Sprintf("%s.wasm", module)

		os.MkdirAll(installDir, 0755)

		install(installDir, wasmFile, target{
			path:     path,
			name:     module,
			wasm:     filepath.Join(path, wasmFile),
			yamlFile: filepath.Join(path, defaultLegateYaml),
		})
	}
}

func install(destDir string, wasmFileName string, module target) {

	fmt.Println("building module:", module.name, ">>", module.wasm)

	{
		cmd := exec.Command("make")
		cmd.Dir = module.path
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to build module", err.Error())
			os.Exit(1)
		}
	}

	fmt.Println("installing to:", destDir)

	{
		size, err := Copy(module.wasm, filepath.Join(destDir, wasmFileName))
		if err != nil {
			fmt.Println("failed to copy wasm to destination", err.Error())
			os.Exit(1)
		}
		fmt.Println("wasm file installed. size:", fmt.Sprintf("%dK", size))
	}

	{
		size, err := Copy(module.yamlFile, filepath.Join(destDir, defaultLegateYaml))
		if err != nil {
			fmt.Println("failed to copy yaml to destination", err.Error())
			os.Exit(1)
		}
		fmt.Println("yaml file installed. size:", fmt.Sprintf("%dK", size))
	}

}

func Copy(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	srcFileState, err := srcFile.Stat()
	if err != nil {
		return 0, err
	}

	if !srcFileState.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}
