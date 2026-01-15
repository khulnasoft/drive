package cli

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Config(t *testing.T) {
	t.Setenv("DRIVE_CONFIG", "./testdata/image-multi-layer-dockerfile/drive-pass.yaml")

	rootCmd := getTestCommand(t, "config --load")
	all := Capture().All().Run(t, func() {
		require.NoError(t, rootCmd.Execute())
	})

	snaps.MatchSnapshot(t, all)
}
