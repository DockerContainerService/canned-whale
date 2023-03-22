//go:build linux && amd64

package registry

import (
	"embed"
	"github.com/DockerContainerService/canned-whale/pkg/utils/path"
)

//go:embed bin/linux-amd64/registry
var binFS embed.FS

func InitBinFS() (err error) {
	err = path.CopyFSFile(binFS, "bin/linux-amd64/registry", "bin/registry")
	return err
}
