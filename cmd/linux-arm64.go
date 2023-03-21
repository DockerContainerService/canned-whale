//go:build linux && arm64

package cmd

import (
	"embed"
)

//go:embed bin/linux-arm64/registry
//go:embed bin/config.yml
var binFS embed.FS

func InitBinFS() (err error) {
	err = copyFSFile("bin/linux-arm64/registry", "bin/registry")
	return err
}
