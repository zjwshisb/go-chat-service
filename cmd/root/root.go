package root

import (
	"github.com/spf13/cobra"
	"ws/cmd/conns"
	"ws/cmd/fake"
	"ws/cmd/migrate"
	"ws/cmd/serve"
	"ws/cmd/stop"
	"ws/config"
)

func NewRootCommand(name string) *cobra.Command {

	var configFile string

	var rootCmd = &cobra.Command{
		Use:                        name,
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Setup(configFile)
		},
	}

	flag := rootCmd.PersistentFlags()
	flag.StringVarP(&configFile, "config", "c","config.yaml", "config file")


	rootCmd.AddCommand(serve.NewServeCommand(),
		migrate.NewMigrateCommand(),
		fake.NewFakeCommand(),
		stop.NewStopCommand(),
		conns.NewConnsCommand(),
	)

	return rootCmd
}
