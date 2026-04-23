//go:build !darwin
// +build !darwin

package build

import _ "embed"

//go:embed appicon.png
var Icon []byte
