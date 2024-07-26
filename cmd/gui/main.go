package main

import (
	"github.com/bosley/emrs/badger"
	"github.com/bosley/emrs/datastore"
	"log/slog"
	"os"
)

/*

structure:

    core - holds datastrj, cfg, and system info.

    app - has core, has apis from /api, handed to each window so they can
          call methods once user input is completed for whatever task

    windows:

      login   -   show on launch, must login with user password from db


      suite:
              load a file on disk into an editor with go highlighting.
              (no lsp req)

              Next to editor should be an updating terminal log of an intsance of the thing
              when it runs in memory. This can be "local" and not from an actual server.
              Once its loaded, tested, and ready, offer to update the install and restart
              any local server that might exist.

      Dashboard:

      Asset CRUD (text box with button next to:), with list of assets.

      Submission builder for selected asset in list

        drop down list of potential actions

            logger.Log
            logger.Info
            logger.Debug

            alert.Sms
              (

                later, if we add an export to the action file we could load
                 the routes that the function uses internally to display them here:

                    // Action builder could import this and call it to retrieve the info
                    func GetRoutes(exportedFn string) []string{
                      if exportedFn == "Sms":
                        .......


                    // Then, in the dropdown we could show:

                    alert.Sms.info.twilio

               )



        Show EMRSURL in text area, updated as selections are made. This shouldn't be editable.
        Under the URL display there can be a button for "copy"

        Under the Url/ Button:

        Token Generator:
          Days, Hours, Minutes, Seconds selection

          generate button

          text area, not editable, with full token. copy button to the right, or underneath

        Submission builder (greyed-out until all of the above is selected):

        Input text area for "Data" to submit. checkbox above it, pre-selected, labeled "Data"

        If "Data" selected, input text area is enabled and user can type away. If not selected,
        erased, greyed out,


        Text area with curl command generated that will submit data. This will update with
        all changes above so it can be adjusted.


        --- Bonus

        Command saver. Save the commands to a list that can be selected from nav bar, and loaded

        Token saver. Save token from generator, add button to load previously generated token.


*/

type EmrsInfo struct {
	Home  string
	Cfg   Config
	Ds    datastore.DataStore
	Badge badger.Badge
}

func main() {

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	home := mustFindHome("")

	info := EmrsInfo{
		Home: home,
	}

	info.Cfg, info.Badge = mustLoadCfgAndBadge(home)
	info.Ds = mustLoadDefaultDataStore(home)

	engine := MustCreateEngine().
		PushView(NewLoginView(info))

	engine.Run()
}
