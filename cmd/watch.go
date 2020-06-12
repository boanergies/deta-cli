package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deta/deta-cli/api"
	"github.com/deta/deta-cli/runtime"
	"github.com/rjeczalik/notify"
	"github.com/spf13/cobra"
)

var (
	watchCmd = &cobra.Command{
		Use:   "watch [path]",
		Short: "Deploy changes in real time",
		RunE:  watch,
		Args:  cobra.MaximumNArgs(1),
	}
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

func watch(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if len(args) != 0 {
		wd = args[0]
	}

	runtimeManager, err := runtime.NewManager(&wd)
	if err != nil {
		return err
	}

	isInitialized, err := runtimeManager.IsInitialized()
	if err != nil {
		return err
	}

	if !isInitialized {
		return fmt.Errorf("deta program not initilialized. see `deta new --help` to create a program")
	}

	progInfo, err := runtimeManager.GetProgInfo()
	if err != nil {
		return err
	}

	lc := api.NewLambdaClient()

	c := make(chan notify.EventInfo, 1)

	if err := notify.Watch(filepath.Join(wd, "main.py"), c, notify.InCloseWrite); err != nil {
		return err
	}

	for {
		<-c
		archive, err := runtimeManager.Zipp()
		if err != nil {
			return err
		}
		start := time.Now()
		err = lc.DeployLambda(progInfo.ID, archive)
		if err != nil {
			return err
		}
		end := time.Now()
		fmt.Println("Deploy lambda took", end.Sub(start))
		fmt.Println("Deployed changes")
		runtimeManager.StoreState()
	}
}
