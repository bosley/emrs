package main

import (
	"fmt"
)

// --

func main() {

	CreateConfigTemplate().WriteTo("emrs.cfg")

	fmt.Println("yo")
}
