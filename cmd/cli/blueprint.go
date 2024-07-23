package main

var globalActionBlueprint = `
// All actions need to be defined in "actions" package
// but we set GOPATH to the $EMRS_HOME/actions so
// files placed there can be included by yaegi
// allowing full customizations

package actions 

import (
  "fmt"
)

type Handler func([]byte)

// All callable handlers that can be triggered by emrs
// must be included in the handlers map, otherwise they
// will be considered internal
var handlers = map[string]Handler = {
  "echo": echoHandler,
}


func echoHandler(data []byte) {
  fmt.Println(data)
}

// Executed iff exists the moment actions are fully loaded
func onInit() {
}

// Executed iff exists right before a valid handler function is
// located
func onBeforeHandle(handler string) {
  // We may want to pass in a ctx to permit us to
  // cancel things or to store a translated variant of
  // the data, then users can define a "Decoder" in the file for 
  // the handler
}
`
