package main

var globalActionBlueprint = `

// All actions need to be defined in "actions" package
// but we set GOPATH to the $EMRS_HOME/actions so
// files placed there can be included by yaegi
// allowing full customizations

package actions 

import (
  "fmt"
  "emrs"
)

func availableEmrsFunctions() {

  fmt.Println("If you are seeing this message, edit your EMRS_HOME/actions/init.go!")

  emrs.Log("Some information")

  emrs.Emit("signal.name", []byte("some data"))

  emrs.Signal("signal.name.no.data")
}

func onData(origin string, route []string, data []byte) error {

  fmt.Println("request from", origin)

  // Now we need to route and parse the data based on what we want to do

  return nil
}

func onInit() error {

  availableEmrsFunctions() 

  return nil
}


`
