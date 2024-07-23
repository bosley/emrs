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
// must be included in the route map, otherwise they
// will be considered internal
//
// Note: routeMap is extracted from this file AFTER
//       onInit() has been executed, meaning that onInit
//       can be used to populate routeMap.
var routeMap = map[string]Handler {

  // Each handler placed here can be called at any time so any
  // data shared beteween methods should be properly guarded
  "echo": echoHandler,
}

// Handle the "echo" route
func echoHandler(data []byte) {
  fmt.Println("ECHO\t>>---> ", string(data))
}

// Executed iff exists the moment actions are fully loaded
func onInit() error {

  // This is executed BEFORE handlers map is extracted, so this FN
  // can dynamically populate the handlers map!

  fmt.Println("ACTIONS INIT!")

  // An error here will indicate to the server that it has 
  // an internal error and can not provice the EMRS service
  return nil
}

`
