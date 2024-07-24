package main

var globalActionBlueprint = `

// TODO: We need to find a way to make the interpretr recognize the install directory
//        of actions as the GO_PATH and have it such that this file can 
//        import "whatever" defined in EMRS_HOME/actions/whatever
//      This will allow for a good deal of modularity and customization

package actions 

import (
  "fmt"
  "emrs"
)

// TODO: Need to figure out the best interface back to EMRS instance that is running this
//        perhaps a message system or something? 
//       
func availableEmrsFunctions() {

  fmt.Println("If you are seeing this message, edit your EMRS_HOME/actions/init.go!")

  emrs.Log("Some information")

  emrs.Emit("signal.name", []byte("some data"))

  emrs.Signal("signal.name.no.data")
}

// NOTE:    At one point I would like to have it such that the user can do something
//          like the following:
//
//                emrs actions --fetch github.com/someone/emrs-something 
//
//          Then have that action installed in EMRS_HOME or somewhere so that the user
//          immediatly has the functionality offered.
//
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
