package main

import (
	"flag"
	"github.com/bosley/emrs/app"
	"github.com/bosley/emrs/badger"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	defaultServerName  = "EMRS Server"
	defaultEnvHome     = "EMRS_HOME"
	defaultBinding     = "127.0.0.1:8080"
	defaultStoragePath = "storage"
	defaultRuntimeFile = ".emrs.pid"
	defaultConfigName  = "server.cfg"
)

type Config struct {
	Binding  string `yaml:binding`
	Key      string `yaml:key`
	Cert     string `yaml:cert`
	Identity string `yaml:identity`
}

func main() {

	/*
	   to launch emrs

	   if not installed:

	    emrs --new [path/to/install]

	    <Set EMRS_HOME to this install to launch later with just "emrs">

	    Or don't and use this:

	    emrs --home [path/to/different/install]


	   --health      if EMRS_HOME is defined in ENV --home can be omitted
	     check to see if the server is alive, and get status


	   --import [import/assets.yaml] | /some_db_info_to_install.yaml


	   --api-key     generates an API key, no duration sets to some default (6 months?)

	     --duration    5h30m40s format for  time.ParseDuration

	     --count    Range: 1-50, to generate many keus

	     --json    Dump in json format

	     [ emrs --api-key --duration 1h --count --json 10 >> keys.json ]

	           "vouchers": [
	             "",
	             "",
	           ]


	   --remove-key <KEY VALUE>

	               invalidates an api key

	*/

	emrsHome := flag.String("home", "", "Home directory")
	createNew := flag.Bool("new", false, "Create a new EMRS instance")
	useForce := flag.Bool("force", false, "Force \"new\" operation, no prompting if item exists")
	coolGuy := flag.Bool("no-prompt", false, "Don't try to be helpful during setup")

	flag.Parse()

	if *emrsHome == "" {
		fromEnv := os.Getenv(defaultEnvHome)
		if fromEnv == "" {
			slog.Error("unable to determine emrs home directory from environment")
			os.Exit(1)
		}
		*emrsHome = fromEnv
	}

	if *createNew {
		writeNewEmrs(*emrsHome, *useForce, *coolGuy)
		return
	}

  cfg := getConfig(*emrsHome)

  println(cfg.Identity)

	emrs := app.New(&app.Opts{
		Binding: "127.0.0.1:8080",
	})

	emrs.Run()
}

func getConfig(home string) Config {
  var config Config 
  target, err := os.ReadFile(filepath.Join(home, defaultConfigName))
  if err != nil {
    slog.Error("failed to load config", "error", err.Error())
    os.Exit(1)
  }
  err = yaml.Unmarshal(target, &config)
  if err != nil {
    slog.Error("failed to load config", "error", err.Error())
    os.Exit(1)
  }
  slog.Debug("loaded config", "binding", config.Binding, "key", config.Key, "cert", config.Cert)
  return config
}

func writeNewEmrs(home string, force bool, noHelp bool) {

	slog.Info("creating new emrs instance", "home", home, "force", force)

	badge, berr := badger.New(defaultServerName)
	if berr != nil {
		slog.Error("badger failed to produce a new identity")
		os.Exit(1)
	}

	info, err := os.Stat(home)

	if err != nil && !os.IsNotExist(err) {
		slog.Error("error retrieving information on home path", "error", err.Error())
		os.Exit(1)
	}

	if err == nil {
		slog.Warn("home directory path already exists", "dir", info.IsDir())
		if force {
			slog.Warn("forcing overwrite..")
			os.RemoveAll(home)
		} else {
			slog.Error("given path exists as directory. Use --force to overwrite")
			os.Exit(1)
		}
	}

	os.MkdirAll(filepath.Join(home, defaultStoragePath), 0755)

	cfg := Config{
		Binding:  defaultBinding,
		Identity: badge.EncodeIdentityString(),
	}

	b, e := yaml.Marshal(&cfg)
	if e != nil {
		slog.Error("Failed to encode new config", "error", e.Error())
		os.Exit(1)
	}

	if err := os.WriteFile(filepath.Join(home, defaultConfigName), b, 0600); err != nil {
		slog.Error("Failed to write configuration file")
		// TODO: Its not required, but we could MOVE the existing, forced, items off
		//        to the side in a temp folder and restore them on a failure via defer
		os.Exit(1)
	}

	slog.Info("Config written to new home directory")

	if os.Getenv(defaultEnvHome) != "" {
		os.Exit(0)
	}

	if !noHelp {
		println(`



   Now that the EMRS server is installed, it is recommended
   that EMRS_HOME is added to the environment. This will enable
   users and scripts to omit the '--home [path]' argument on
   launch if this is a singular, or typical, EMRS instance.

   If you want to enable HTTPS, provide a 'key' and 'cert' value
   in the installed configuration file that can direct the server
   to locate a key and cert file on the host system.

   To change the address/port port of the EMRS services, see the
   relevant settings generated within the configuration file.


          It is now save to start the server


    `)
	}
	os.Exit(0)
}
