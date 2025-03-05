package viewmodel

import "github.com/khulnasoft/drive/drive/image"

type LayerSetState struct {
	LayerIndex        int
	Layers            []*image.Layer
	CompareMode       LayerCompareMode
	CompareStartIndex int
}

func NewLayerSetState(layers []*image.Layer, compareMode LayerCompareMode) *LayerSetState {
	return &LayerSetState{
		Layers:      layers,
		CompareMode: compareMode,
	}
}

// getCompareIndexes determines the layer boundaries to use for comparison (based on the current compare mode)
func (state *LayerSetState) GetCompareIndexes() (bottomTreeStart, bottomTreeStop, topTreeStart, topTreeStop int) {
	// Handle negative CompareStartIndex
	bottomTreeStart = max(0, state.CompareStartIndex)
	topTreeStop = state.LayerIndex

	if state.LayerIndex == state.CompareStartIndex {
		bottomTreeStop = state.LayerIndex
		topTreeStart = state.LayerIndex
	} else if state.CompareMode == CompareSingleLayer {
		if state.LayerIndex == 0 {
			bottomTreeStop = 1
			topTreeStart = 0
		} else {
			bottomTreeStop = max(0, state.LayerIndex-1)
			topTreeStart = state.LayerIndex
		}
	} else {
		bottomTreeStop = bottomTreeStart
		topTreeStart = bottomTreeStart + 1
	}

	return bottomTreeStart, bottomTreeStop, topTreeStart, topTreeStop
}

// Helper function to get max of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}