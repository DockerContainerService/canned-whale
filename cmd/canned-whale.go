package cmd

import (
	"fmt"
	"github.com/DockerContainerService/canned-whale/pkg/client/pack"
	"github.com/DockerContainerService/canned-whale/pkg/client/registry"
	"github.com/DockerContainerService/canned-whale/pkg/utils/path"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

const cachedPath = "bin"

var (
	authFile, imageFile          string
	procNum, retries, listenPort int
	osFilterList, archFilterList []string
)

var rootCmd = &cobra.Command{
	Use:   "canned-whale",
	Short: "Container export tool",
	Long: `A container export tool implement by Go.
    Complete documentation is available at https://github.com/DockerAcCn/canned-whale`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Println("hello")
		return nil
	},
}

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Start a docker registry",
	Long: `Start a docker registry.
    Complete documentation is available at https://github.com/DockerAcCn/canned-whale`,
	RunE: func(cmd *cobra.Command, args []string) error {
		registry.RunRegistryClient(listenPort)
		defer registry.StopRegistryClient()
		return nil
	},
}

var canCmd = &cobra.Command{
	Use:     "can",
	Aliases: []string{"save"},
	Short:   "Image export tool, ",
	Long: `A image export tool implement by Go.
    Complete documentation is available at https://github.com/DockerAcCn/canned-whale`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Println("can")
		return nil
	},
}

var packCmd = &cobra.Command{
	Use:     "pack",
	Aliases: []string{"package"},
	Short:   "Docker registry export tool",
	Long: `A docker registry export tool implement by Go.
    Complete documentation is available at https://github.com/DockerAcCn/canned-whale`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pack.RunPackClient(authFile, imageFile, procNum, retries, osFilterList, archFilterList)
		fmt.Println("Finished")
		return nil
	},
}

func init() {
	if !path.IsPathExist(cachedPath) {
		err := path.MkdirPath(cachedPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	packCmd.PersistentFlags().StringVar(&authFile, "auth", "", "auth file path. This flag need to be pair used with --images")
	packCmd.PersistentFlags().StringVar(&imageFile, "images", "", "images file path. This flag need to be pair used with --auth")
	packCmd.PersistentFlags().IntVarP(&procNum, "proc", "p", 5, "numbers of working goroutines")
	packCmd.PersistentFlags().IntVarP(&retries, "retries", "r", 3, "times to retry failed task")
	packCmd.PersistentFlags().StringArrayVar(&osFilterList, "os", []string{}, "os list to filter source tags, not works for docker v2 schema1 media")
	packCmd.PersistentFlags().StringArrayVar(&archFilterList, "arch", []string{}, "architecture list to filter source tags")

	registryCmd.PersistentFlags().IntVarP(&listenPort, "port", "P", 80, "docker registry listen port")
	rootCmd.AddCommand(packCmd, canCmd, registryCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
