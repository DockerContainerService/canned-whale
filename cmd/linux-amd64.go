//go:build linux && amd64

package cmd

import (
	"embed"
)

//go:embed bin/linux-amd64/registry
//go:embed bin/config.yml
var binFS embed.FS

func InitBinFS() (err error) {
	err = copyFSFile("bin/linux-amd64/registry", "bin/registry")
	return err
}
