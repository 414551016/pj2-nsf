package cmd

import (
	"context"
	"fmt"
	"nsf/pkg"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var nsfCmd = &cobra.Command{
	Use:   "nsf",
	Short: "Network Service Function",
	Run:   nsfRun,
}

func init() {
	nsfCmd.Flags().StringP("config", "c", "", "config file")
	if err := nsfCmd.MarkFlagRequired("config"); err != nil {
		panic(err)
	}
}

func nsfRun(cmd *cobra.Command, args []string) {
	configFilePath, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}

	app, err := pkg.NewApp(configFilePath)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx); err != nil {
		return
	}
	defer app.ShutDown()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	fmt.Println()

	cancel()
}

func Execute() {
	if err := nsfCmd.Execute(); err != nil {
		panic(err)
	}
}
