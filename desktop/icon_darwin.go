//go:build darwin
// +build darwin

package main

import _ "embed"

//go:embed build/appicon.darwin.png
var appIcon []byte
