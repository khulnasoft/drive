package runtime

import (
	"github.com/spf13/viper"

	"github.com/khulnasoft/drive/drive"
)

type Options struct {
	Ci           bool
	Image        string
	Source       drive.ImageSource
	IgnoreErrors bool
	ExportFile   string
	CiConfig     *viper.Viper
	BuildArgs    []string
}
