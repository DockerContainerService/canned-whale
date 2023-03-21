package cmd

import (
	"fmt"
	"github.com/AliyunContainerService/image-syncer/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

var (
	authFile, imageFile          string
	procNum, retries             int
	osFilterList, archFilterList []string
)

var RootCmd = &cobra.Command{
	Use:     "canned-whale",
	Aliases: []string{"canned-whale"},
	Short:   "A docker registry export tool",
	Long: `A docker registry export tool implement by Go.
    Complete documentation is available at https://github.com/DockerAcCn/canned-whale`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("start docker registry...")

		fmt.Println("start sync task...")
		client, err := client.NewSyncClient("", authFile, imageFile, "", procNum, retries, "", "", osFilterList, archFilterList)
		if err != nil {
			return fmt.Errorf("init sync client error: %+v", err)
		}
		client.Run()
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&authFile, "auth", "", "auth file path. This flag need to be pair used with --images")
	RootCmd.PersistentFlags().StringVar(&imageFile, "images", "", "images file path. This flag need to be pair used with --auth")
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
