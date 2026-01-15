package image

import (
	"github.com/khulnasoft/drive/drive/filetree"
)

type Image struct {
	Request string
	Trees   []*filetree.FileTree
	Layers  []*Layer
}
