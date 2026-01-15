package command

import (
	"fmt"
	"github.com/anchore/clio"
	"github.com/khulnasoft/drive/cmd/drive/cli/internal/command/adapter"
	"github.com/khulnasoft/drive/cmd/drive/cli/internal/options"
	"github.com/khulnasoft/drive/drive"
	"github.com/spf13/cobra"
)

type buildOptions struct {
	options.Application `yaml:",inline" mapstructure:",squash"`

	// reserved for future use of build-only flags
}

func Build(app clio.Application) *cobra.Command {
	opts := &buildOptions{
		Application: options.DefaultApplication(),
	}
	return app.SetupCommand(&cobra.Command{
		Use:                "build [any valid `docker build` arguments]",
		Short:              "Builds and analyzes a docker image from a Dockerfile (this is a thin wrapper for the `docker build` command).",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := setUI(app, opts.Application); err != nil {
				return fmt.Errorf("failed to set UI: %w", err)
			}

			resolver, err := drive.GetImageResolver(opts.Analysis.Source)
			if err != nil {
				return fmt.Errorf("cannot determine image provider for build: %w", err)
			}

			ctx := cmd.Context()

			img, err := adapter.ImageResolver(resolver).Build(ctx, args)
			if err != nil {
				return fmt.Errorf("cannot build image: %w", err)
			}

			return run(cmd.Context(), opts.Application, img, resolver)
		},
	}, opts)
}
