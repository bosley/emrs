package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bosley/emrs/app"
	"github.com/bosley/emrs/badger"
	"github.com/bosley/emrs/datastore"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultServerName        = "EMRS Server"
	defaultEnvHome           = "EMRS_HOME"
	defaultBinding           = "127.0.0.1:8080"
	defaultStoragePath       = "storage"
	defaultRuntimeFile       = ".emrs.pid"
	defaultConfigName        = "server.cfg"
	defaultUiKeyDuration     = "8760h" // 1 year
	defaultUserGivenDuration = "4320h" // ~6 months
)

type Config struct {
	Binding  string `yaml:binding`
	Key      string `yaml:key`
	Cert     string `yaml:cert`
	Identity string `yaml:identity`
}

func main() {

	/*
		    TODO:
		        CLI Args that would be nice, and potentially required, but not required as-of-yet:

		        --new-ui-key    (uses givenDuration)    Generates a new UI key for the Ring-0 user,
		                                                this will require a sever restart, or a means
		                                                to update the api auth used by the server at runtime.
		                                                This might be a non-issue but since the auth is not
		                                                implemented yet, it may require change later idk

		        --health        Check to see if there is a currently running EMRS instance
		        --down          Kill a currently running emrs instance (could use PID, could locally-bound port
		                        and a "local-machine" api to do the IPC between a cli instance and a running server

		        --disable       Disable all running functionality (submisisons to server, ui, etc)
		                        This could be useful for testing.
	*/

	emrsHome := flag.String("home", "", "Home directory")
	createNew := flag.Bool("new", false, "Create a new EMRS instance")
	useForce := flag.Bool("force", false, "Force \"new\" operation, no prompting if item exists")
	coolGuy := flag.Bool("no-prompt", false, "Don't try to be helpful during setup")
	isRelease := flag.Bool("release", false, "Enable release mode")
	genVouchers := flag.Int("vouchers", 0, "Enter a number >0 to generate a series of vouchers. Use with `duration.`")
	givenDuration := flag.String("duration", defaultUserGivenDuration, "Duration to give to vouchers (ex: 1h15m)")


  createAsset := flag.String("new-asset", "", "Create a new asset")
  listAssets  := flag.Bool("list-assets", false, "List all known assets")
  removeAsset := flag.String("remove-asset", "", "Remove an asset by its UUID")
  updateAsset := flag.String("update-asset", "", "Update an asset's name given its UUID (requires --name)")

  withName:= flag.String("name", "", "Set the name value for a corresponding command")

  // beFancy := flag.Bool("fancy", false, "Perform the action with a fancy tui") // (list assets, using bubbletea)

  // panel  := flag.Bool("panel", false, "Launch the interactive TUI that requires direct access to the datastore (not via web api)

	flag.Parse()

  if !*isRelease {
	  slog.SetDefault(
	  	slog.New(
	  		slog.NewTextHandler(os.Stdout,
	  			&slog.HandlerOptions{
	  				Level: slog.LevelDebug,
	  			})))
      }

	if *emrsHome == "" {
		fromEnv := os.Getenv(defaultEnvHome)
		if fromEnv == "" {
			slog.Error("unable to determine emrs home directory from environment")
			os.Exit(1)
		}
		*emrsHome = fromEnv
	}

  // Create a new EMRS instance on disk, and then exit
	if *createNew {
		writeNewEmrs(*emrsHome, *useForce, *coolGuy)
		return
	}

  // Load the configuration, and then populate the server identity badge

	cfg := getConfig(*emrsHome)
	badge, err := badger.DecodeIdentityString(cfg.Identity)
	if err != nil {
		slog.Error("badger failed to decode server identity", "error", err.Error())
		os.Exit(1)
	}

  // If the user wants to generate vouchers based on the server identity,
  // we do so here and then exist

	if *genVouchers > 0 {
		d, err := time.ParseDuration(*givenDuration)
		if err != nil {
			slog.Error("failed to parse duration", "error", err.Error())
			os.Exit(1)
		}
		generateVouchers(badge, *genVouchers, d)
		return
	}

  // Load the DataStore from the EMRS home directory

	dataStrj, err := datastore.Load(filepath.Join(*emrsHome, defaultStoragePath))
	if err != nil {
		slog.Error("failed to load datastore", "error", err.Error())
		os.Exit(1)
	}

  // Check for asset commands

  if *listAssets {
    assets := dataStrj.GetAssets()
    if len(assets) == 0 {
      fmt.Println("There are no assets contained in the EMRS data storage system")
      return
    }
    for i, a := range assets {
      fmt.Printf("%6d | %s | %s\n", i, a.Id, a.DisplayName)
    }
    return
  }
  if strings.Trim(*createAsset, " ") != "" {
    id, err := badger.GenerateId()
    if err != nil {
      slog.Error("badger failed to create a unique id for asset", "error", err.Error())
      os.Exit(1)
    }
    if !dataStrj.AddAsset(datastore.Asset{
      Id: id,
      DisplayName: *createAsset,
    }) {
      slog.Error("failed to add asset", "id", id, "name", *createAsset)
      os.Exit(1)
    }
    return
  }
  if strings.Trim(*removeAsset, " ") != "" {
    if !dataStrj.RemoveAsset(*removeAsset) {
      slog.Error("failed to remove asset", "id", *removeAsset)
      os.Exit(1)
    }
    return
  }
  if strings.Trim(*updateAsset, " ") != "" {
    if !dataStrj.UpdateAsset(datastore.Asset{
      Id: *updateAsset,
      DisplayName: *withName,
    }) {
      slog.Error("failed to add asset", "id", *updateAsset, "name", *withName)
      os.Exit(1)
    }
    return
  }

  // Create the EMRS server application 

	emrs := app.New(&app.Opts{
		Badge:     badge,
		Binding:   cfg.Binding,
		DataStore: dataStrj,
	})

  // Check if we can use HTTPS

	if strings.Trim(cfg.Key, " ") != "" && strings.Trim(cfg.Cert, " ") != "" {
		emrs.UseHttps(cfg.Key, cfg.Cert)
	}

  // RUN

	emrs.Run(*isRelease)
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

	newUser := RunUserInfoTui()

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

	strj := filepath.Join(home, defaultStoragePath)
	os.MkdirAll(strj, 0755)

	oneYear, err := time.ParseDuration(defaultUiKeyDuration)
	if err != nil {
		slog.Error("failed to setup voucher duration")
		os.Exit(1)
	}

	voucher, err := badger.NewVoucher(badge, oneYear)
	if err != nil {
		slog.Error("failed to generate ui voucher")
		os.Exit(1)
	}

	datastore.SetupDisk(strj, datastore.User{
		DisplayName: newUser.Name,
		Hash:        newUser.Hash,
		UiKey:       voucher, // 1 year
	})

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

func generateVouchers(badge badger.Badge, n int, durr time.Duration) {
	vouchers := make([]string, n)
	for i := range n {
		voucher, err := badger.NewVoucher(badge, durr)
		if err != nil {
			slog.Error("failed to generate vouchers")
			os.Exit(1)
		}
		vouchers[i] = voucher
	}

	b, _ := json.Marshal(vouchers)
	fmt.Println(string(b))
}
