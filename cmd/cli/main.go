package main

import (
  "os"
  "log/slog"
  "flag"
	"github.com/bosley/emrs/app"
	"github.com/bosley/emrs/badger"
)

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

  flag.Parse()

  if *emrsHome == "" {
    fromEnv := os.Getenv("EMRS_HOME")
    if fromEnv == "" {
      slog.Error("unable to determine EMRS_HOME from environment, with no --home arg given")
      os.Exit(1)
    }
    *emrsHome = fromEnv
  }

  if *createNew {
    
  }

	emrs := app.New(&app.Opts{
    Binding: "127.0.0.1:8080",
  })

	badge, _ := badger.New("EMRS Server")


  slog.Info("initialization complete", "home", *emrsHome, "server", badge.Id())


	emrs.Run()
}

func writeNewEmrs(home string, force bool) {

  slog.Info("creating new emrs instance", "home", home, "force", foce)

  info, err := os.Stat(home)

  if err != nil {

    BREAK THE BUILD !
    // handle os stat errors

  }

}
