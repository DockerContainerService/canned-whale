package pack

import (
	"fmt"
	"github.com/AliyunContainerService/image-syncer/pkg/client"
	"github.com/DockerContainerService/canned-whale/pkg/client/registry"
	"github.com/DockerContainerService/canned-whale/pkg/utils/path"
	"github.com/DockerContainerService/canned-whale/pkg/utils/tcp"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func RunPackClient(authFile, imageFile string, procNum, retries int, osFilterList, archFilterList []string) {
	listenPort := 10000
	for !tcp.IsPortAvailable(listenPort) {
		listenPort += 1
	}
	go registry.RunRegistryClient(listenPort)

	authFile = processSed(authFile, "127.0.0.1", fmt.Sprintf("127.0.0.1:%d", listenPort))
	imageFile = processSed(imageFile, "127.0.0.1", fmt.Sprintf("127.0.0.1:%d", listenPort))

	for tcp.IsPortAvailable(listenPort) {
		time.Sleep(1 * time.Second)
		logrus.Println("waiting for registry up...")
	}
	c, err := client.NewSyncClient("", authFile, imageFile, "", procNum, retries, "", "", osFilterList, archFilterList)
	if err != nil {
		logrus.Fatalf("init PackClient failed: %+v", err)
	}
	c.Run()

	registry.StopRegistryClient()

	path.TarPath("/tmp/registry", "registry.tgz")

	path.RemovePath("/tmp/registry")

	path.RemovePath("./bin")
}

func processSed(filePath, src, dst string) (res string) {
	fo, err := os.ReadFile(filePath)
	if err != nil {
		logrus.Fatalf("file %s exchange failed: %+v", filePath, err)
	}
	content := strings.ReplaceAll(string(fo), src, dst)
	res = fmt.Sprintf("./bin/%s", filepath.Base(filePath))
	fw, err := os.OpenFile(res, os.O_CREATE|os.O_WRONLY, 0666)
	fw.Write([]byte(content))
	return res
}
