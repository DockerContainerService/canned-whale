package cmd

import (
	"fmt"
	"github.com/AliyunContainerService/image-syncer/pkg/client"
	"github.com/DockerContainerService/canned-whale/pkg/utils/path"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"os"
	"os/exec"
)

const cachedPath = "bin"

var (
	authFile, imageFile, configFile string
	procNum, retries                int
	osFilterList, archFilterList    []string
)

var RootCmd = &cobra.Command{
	Use:     "canned-whale",
	Aliases: []string{"canned-whale"},
	Short:   "A docker registry export tool",
	Long: `A docker registry export tool implement by Go.
    Complete documentation is available at https://github.com/DockerAcCn/canned-whale`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("start docker registry...")
		err := InitBinFS()
		if err != nil {
			return fmt.Errorf("init registry client error: %+v", err)
		}
		command := exec.Command("chmod", "a+x", "bin/registry")
		err = command.Run()
		if err != nil {
			return fmt.Errorf("start registry client error: %+v", err)
		}
		if configFile == "" {
			err := copyFSFile("bin/config.yml", "bin/config.yml")
			if err != nil {
				return fmt.Errorf("generate registry config file error: %+v", err)
			}
			configFile = "./bin/config.yml"
		}
		runRegistryCmd := exec.Command("./bin/registry", "serve", configFile)
		go func() {
			err := runRegistryCmd.Run()
			if err != nil {
				fmt.Printf("start docker registry error: %+v\n", err)
			}
		}()
		fmt.Println("start sync task...")
		client, err := client.NewSyncClient("", authFile, imageFile, "", procNum, retries, "", "", osFilterList, archFilterList)
		if err != nil {
			return fmt.Errorf("init sync client error: %+v", err)
		}
		client.Run()
		fmt.Println("start package...")
		command = exec.Command("mv", "/tmp/registry", "./registry")
		err = command.Run()
		if err != nil {
			return fmt.Errorf("move package error: %+v", err)
		}
		command = exec.Command("tar", "zcvf", "registry.tgz", "registry")
		err = command.Run()
		if err != nil {
			return fmt.Errorf("package error: %+v", err)
		}
		command = exec.Command("rm", "-rf", "./registry")
		err = command.Run()
		if err != nil {
			return fmt.Errorf("remove cache error: %+v", err)
		}
		runRegistryCmd.Process.Kill()
		path.RemovePath(cachedPath)
		fmt.Println("Finished")
		return nil
	},
}

func copyFSFile(FSFilePath, localPath string) (err error) {
	binFsFileOpen, err := binFS.Open(FSFilePath)
	if err != nil {
		return err
	}
	defer func(binFsFileOpen fs.File) {
		err := binFsFileOpen.Close()
		if err != nil {
			fmt.Printf("binFsFileOpen close error: %+v\n", err)
		}
	}(binFsFileOpen)
	binFsFileOut, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func(binFsFileOut *os.File) {
		err := binFsFileOut.Close()
		if err != nil {
			fmt.Printf("binFsFileOut close error: %+v\n", err)
		}
	}(binFsFileOut)
	_, err = io.Copy(binFsFileOut, binFsFileOpen)
	return err
}

func init() {
	if !path.IsPathExist(cachedPath) {
		err := path.MkdirPath(cachedPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	RootCmd.PersistentFlags().StringVar(&authFile, "auth", "", "auth file path. This flag need to be pair used with --images")
	RootCmd.PersistentFlags().StringVar(&imageFile, "images", "", "images file path. This flag need to be pair used with --auth")
	//RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "docker registry config file")
	RootCmd.PersistentFlags().IntVarP(&procNum, "proc", "p", 5, "numbers of working goroutines")
	RootCmd.PersistentFlags().IntVarP(&retries, "retries", "r", 3, "times to retry failed task")
	RootCmd.PersistentFlags().StringArrayVar(&osFilterList, "os", []string{}, "os list to filter source tags, not works for docker v2 schema1 media")
	RootCmd.PersistentFlags().StringArrayVar(&archFilterList, "arch", []string{}, "architecture list to filter source tags")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
