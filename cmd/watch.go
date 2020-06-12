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

	fmt.Println("Watching changes")
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
		fmt.Println(end.Sub(start))
		fmt.Println("Deployed changes")

		dc, err := runtimeManager.GetDepChanges()
		if err != nil {
			return err
		}
		runtimeManager.StoreState()

		if dc != nil {
			fmt.Println("Updating dependencies...")
			command := runtime.DepCommands[progInfo.Runtime]
			if len(dc.Added) > 0 {
				installCmd := fmt.Sprintf("%s install", command)
				for _, a := range dc.Added {
					installCmd = fmt.Sprintf("%s %s", installCmd, a)
				}
				o, err := client.UpdateProgDeps(&api.UpdateProgDepsRequest{
					ProgramID: progInfo.ID,
					Command:   installCmd,
				})
				if err != nil {
					return fmt.Errorf("failed to add dependencies: %v", err)
				}
				fmt.Println(o.Output)

				for _, a := range dc.Added {
					progInfo.Deps = append(progInfo.Deps, a)
				}
				runtimeManager.StoreProgInfo(progInfo)
			}
			if len(dc.Removed) > 0 {
				uninstallCmd := fmt.Sprintf("%s uninstall", command)
				for _, d := range dc.Removed {
					uninstallCmd = fmt.Sprintf("%s %s", uninstallCmd, d)
				}
				o, err := client.UpdateProgDeps(&api.UpdateProgDepsRequest{
					ProgramID: progInfo.ID,
					Command:   uninstallCmd,
				})
				if err != nil {
					return fmt.Errorf("failed to remove dependencies: %v", err)
				}
				fmt.Println(o.Output)
				for _, d := range dc.Removed {
					progInfo.Deps = removeFromSlice(progInfo.Deps, d)
				}
				runtimeManager.StoreProgInfo(progInfo)
			}
		}
	}
}
