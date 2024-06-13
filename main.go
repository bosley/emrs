package main

import (
  "os"
	"log/slog"
//  "internal/reaper"
)

type ModuleData struct {
  // Create a structure that holds the 
  

  // *nerv.Module
  // []nerv.Consumers
  // topic configuration
}


func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))


  slog.Debug("OI")
}
