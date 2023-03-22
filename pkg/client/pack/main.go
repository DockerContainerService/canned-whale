package pack

import (
	"github.com/AliyunContainerService/image-syncer/pkg/client"
	"github.com/DockerContainerService/canned-whale/pkg/client/registry"
	"github.com/DockerContainerService/canned-whale/pkg/utils/path"
	"github.com/DockerContainerService/canned-whale/pkg/utils/tcp"
	"github.com/sirupsen/logrus"
	"time"
)

func RunPackClient(authFile, imageFile string, procNum, retries int, osFilterList, archFilterList []string) {
	listenPort := 10000
	for !tcp.IsPortAvailable(listenPort) {
		listenPort += 1
	}
	go registry.RunRegistryClient(listenPort)

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
