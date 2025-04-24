package view

import (
	"github.com/khulnasoft/drive/cmd/drive/cli/internal/ui/v1/viewmodel"
)

type LayerChangeListener func(viewmodel.LayerSelection) error
