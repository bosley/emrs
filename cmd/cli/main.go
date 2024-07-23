package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bosley/emrs/api"
	"github.com/bosley/emrs/app"
	"github.com/bosley/emrs/badger"
	"github.com/bosley/emrs/datastore"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
  "io"
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

  defaultActionsDir        = "actions"
  defaultActionsBaseFile   = "_actions/blueprint.go"
  defaultActionsInstalled  = "init.go"
)

const (
	ttlEphemeralVoucher = "30s"
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

	*/

	emrsHome := flag.String("home", "", "Home directory")
	createNew := flag.Bool("new", false, "Create a new EMRS instance")
	useForce := flag.Bool("force", false, "Force \"new\" operation, no prompting if item exists")
	coolGuy := flag.Bool("no-prompt", false, "Don't try to be helpful during setup")
	isRelease := flag.Bool("release", false, "Enable release mode")
	genVouchers := flag.Int("vouchers", 0, "Enter a number >0 to generate a series of vouchers. Use with `duration.`")
	givenDuration := flag.String("duration", defaultUserGivenDuration, "Duration to give to vouchers (ex: 1h15m)")

	createAsset := flag.String("new-asset", "", "Create a new asset")
	listAssets := flag.Bool("list-assets", false, "List all known assets")
	removeAsset := flag.String("remove-asset", "", "Remove an asset by its UUID")
	updateAsset := flag.String("update-asset", "", "Update an asset's name given its UUID (requires --name)")

	withName := flag.String("name", "", "Set the name value for a corresponding command")

	emit := flag.String("submit", "", "Submit event to a server. Format>  Asset-UUID:deceoder.proc0.proc1.proc2@https://127.0.0.1:8080")
	withData := flag.String("data", "", "Add data to a submission")

	getStatus := flag.String("stat", "", "Submit getStatus reaquest to server Format> https://127.0.0.1:8080")

	down := flag.Bool("down", false, "Stop a server instance using the loaded server configuration")

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

	if *getStatus != "" {
		executeGetStatus(*getStatus, cfg)
		return
	}

	// Check to see if we are just emitting an event

	if *emit != "" {
		executeSubmission(badge, cfg, *emit, *withData)
		return
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

	// Check for "DOWN"

	if *down {
		executeDown(cfg, badge, dataStrj)
		return
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
			Id:          id,
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
			Id:          *updateAsset,
			DisplayName: *withName,
		}) {
			slog.Error("failed to add asset", "id", *updateAsset, "name", *withName)
			os.Exit(1)
		}
		return
	}

	// Create the EMRS server application

	emrs, launchErr := app.New(&app.Opts{
		Badge:     badge,
		Binding:   cfg.Binding,
		DataStore: dataStrj,
    ActionsPath: filepath.Join(*emrsHome, defaultActionsDir),
    ActionRootFile: defaultActionsInstalled,
	})

  if launchErr != nil { 
    slog.Error("failed to create emrs application", "error", launchErr.Error())
    os.Exit(1)
  }

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

  actions := filepath.Join(home, defaultActionsDir)
  os.MkdirAll(actions, 0755)

  actionInitDest := filepath.Join(actions, defaultActionsInstalled)

	if err := os.WriteFile(actionInitDest, []byte(globalActionBlueprint), 0600); err != nil {
		slog.Error("Failed to write actions init file")
		os.Exit(1)
	}

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

func copyFile(src, dst string) (int64, error) {
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

// Takes in the server's badge, user-supplied url and data (optional)
// and submits an event to the targeted EMRS server.
// The badge is utilized to generate a very short-lived voucher
// (30 sec) for each request. Whats important to realize is that we
// are using the local server's identity, meaning that this will
// only be valid for the local EMRS instance, and not any others
// unless they share the same identity
func executeSubmission(badge badger.Badge, cfg Config, url string, data string) {

	slog.Debug("submission execution request", "url", url, "data", data)

	emrsUrl := new(api.EmrsAddress)

	if err := emrsUrl.From(url); err != nil {
		slog.Error("failed to parse url", "error", err.Error())
		fmt.Println("Malformed URL given. Proper format>   ASSET_ID:PATH@SERVER:PORT")
		os.Exit(1)
	}

	dur, err := time.ParseDuration(ttlEphemeralVoucher)
	if err != nil {
		slog.Error("failed to generate duration", "error", err.Error())
		os.Exit(1)
	}

	voucher, err := badger.NewVoucher(badge, dur)
	if err != nil {
		slog.Error("failed to generate ui voucher")
		os.Exit(1)
	}

	var info *api.HttpsInfo

	if strings.Trim(cfg.Key, " ") != "" && strings.Trim(cfg.Cert, " ") != "" {
		info = new(api.HttpsInfo)
		info.Cert = cfg.Cert
		info.Key = cfg.Key
	}

	client := api.HttpSubmissions(api.Options{
		Binding:     emrsUrl.Server,
		AssetId:     emrsUrl.Asset,
		AccessToken: voucher,
	},
		info,
	)

	composed, _ := api.ComposeRoute(emrsUrl.Route)

	if e := client.Submit(composed, []byte(data)); e != nil {
		fmt.Println("Error from HTTP Client:", e.Error())
		os.Exit(1)
	}
	return
}

func executeGetStatus(binding string, cfg Config) {

	var info *api.HttpsInfo

	if strings.Trim(cfg.Key, " ") != "" && strings.Trim(cfg.Cert, " ") != "" {
		info = new(api.HttpsInfo)
		info.Cert = cfg.Cert
		info.Key = cfg.Key
	}

	client := api.HttpStats(api.Options{
		Binding: binding,
	},
		info,
	)

	ut, err := client.GetUptime()

	if err != nil {
		slog.Error("error fetching server getStatus", "binding", binding, "error", err.Error)
		os.Exit(1)
	}

	fmt.Println("server is up. uptime:", ut.String())
}

func executeDown(cfg Config, badge badger.Badge, db datastore.DataStore) {

	// TODO: Each of these commands build their own api which is intended, but once the different
	//        apis expand we should restructure this main application to route the commands
	//        to a specific api handler that is constructed once for all of the different commands

	var info *api.HttpsInfo

	if strings.Trim(cfg.Key, " ") != "" && strings.Trim(cfg.Cert, " ") != "" {
		info = new(api.HttpsInfo)
		info.Cert = cfg.Cert
		info.Key = cfg.Key
	}

	o, e := db.GetOwner()
	if e != nil {
		slog.Error("failed to obtain user's access code for CNC", "error", e.Error())
		os.Exit(1)
	}

	if !badger.ValidateVoucher(badge.PublicKey(), o.UiKey) {
		slog.Error("user's current Ui Key is no longer valid. Please replace the key with a new voucher")
		os.Exit(2)
	}

	client := api.HttpCNC(cfg.Binding, o.UiKey, info)

	if err := client.Shutdown(); err != nil {
		slog.Info("failed to request shutdown on server", "error", err.Error())
		os.Exit(1)
	}

	fmt.Println("shutdown request sent")
}
