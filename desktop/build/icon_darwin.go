//go:build darwin
// +build darwin

package build

import _ "embed"

//go:embed appicon.darwin.png
var Icon []byte
