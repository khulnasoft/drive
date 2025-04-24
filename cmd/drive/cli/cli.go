package cli

import (
	"github.com/anchore/clio"
	"github.com/khulnasoft/drive/cmd/drive/cli/internal/command"
	"github.com/khulnasoft/drive/cmd/drive/cli/internal/ui"
	"github.com/khulnasoft/drive/internal/bus"
	"github.com/khulnasoft/drive/internal/log"
	"github.com/spf13/cobra"
)

func Application(id clio.Identification) clio.Application {
	app, _ := create(id)
	return app
}

func Command(id clio.Identification) *cobra.Command {
	_, cmd := create(id)
	return cmd
}

func create(id clio.Identification) (clio.Application, *cobra.Command) {
	clioCfg := clio.NewSetupConfig(id).
		WithGlobalConfigFlag().   // add persistent -c <path> for reading an application config from
		WithGlobalLoggingFlags(). // add persistent -v and -q flags tied to the logging config
		WithConfigInRootHelp().   // --help on the root command renders the full application config in the help text
		WithUI(ui.None()).
		WithInitializers(
			func(state *clio.State) error {
				bus.Set(state.Bus)
				log.Set(state.Logger)

				//stereoscope.SetBus(state.Bus)
				//stereoscope.SetLogger(state.Logger)
				return nil
			},
		)
	//WithPostRuns(func(_ *clio.State, _ error) {
	//	stereoscope.Cleanup()
	//})

	app := clio.New(*clioCfg)

	rootCmd := command.Root(app)

	rootCmd.AddCommand(
		clio.VersionCommand(id),
		clio.ConfigCommand(app, nil),
		command.Build(app),
	)

	return app, rootCmd
}
