package payload

import "github.com/khulnasoft/drive/drive/image"

type Explore struct {
	Analysis image.Analysis
	Content  image.ContentReader
}
