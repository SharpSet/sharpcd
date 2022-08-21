package main

import "fmt"

// use fmt

// ScriptVersion :create struct that stores versions and their unique features
type ScriptVersion struct {
	Version float32
	WorksOn string
}

// list of versions
var scriptVersions = []ScriptVersion{
	{
		Version: 1.1,
		WorksOn: "V3.8 and up",
	},
	{
		Version: 1.0,
		WorksOn: "V3.3 and below. Does not support (depends) arg",
	},
}

// Function that says what version of the script it is
// and what versions it works on
func compareVersions(version float32) {
	var triggered bool = false

	for _, v := range scriptVersions {
		if v.Version == version {
			fmt.Println("SharpCD Version: ", v.Version)
			fmt.Println("Works on: ", v.WorksOn)
			triggered = true
			break
		}
	}

	if !triggered {
		// print float
		fmt.Println("Version not found:", version)
		fmt.Println()
	}
}
