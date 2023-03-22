package registry

import (
	_ "embed"
	"github.com/DockerContainerService/canned-whale/pkg/utils/path"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const configFilePath = "bin/config.yml"

//go:embed bin/config.yml
var registryConfig string

var (
	client *exec.Cmd
	killed = false
)

func RunRegistryClient(listenPort int) {
	closeHandler()
	config := strings.ReplaceAll(registryConfig, strconv.Itoa(80), strconv.Itoa(listenPort))
	logrus.Printf("registry run on port: %+v\n", listenPort)
	err := InitBinFS()
	if err != nil {
		logrus.Fatalf("init bin fs failed: %+v", err)
	}

	file, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Printf("create registry config file failed: %+v", err)
	}
	_, err = file.Write([]byte(config))
	if err != nil {
		logrus.Fatalf("write registry config file failed: %+v", err)
	}
	err = file.Close()
	if err != nil {
		logrus.Fatalf("close registry config file failed: %+v", err)
	}
	command := exec.Command("chmod", "a+x", "./bin/registry")
	err = command.Run()
	if err != nil {
		logrus.Fatalf("chmod registry failed: %+v", err)
	}
	client = exec.Command("./bin/registry", "serve", configFilePath)
	err = client.Run()
	if err != nil {
		path.RemovePath("./bin")
		if !killed {
			logrus.Fatalf("run registry failed: %+v", err)
		}
	}
}

func StopRegistryClient() {
	killed = true
	err := client.Process.Kill()
	if err != nil {
		logrus.Fatalf("stop registry client failed: %+v", err)
	}
}

func closeHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logrus.Println("cleaning ...")
		path.RemovePath("./bin")
		os.Exit(0)
	}()
}
