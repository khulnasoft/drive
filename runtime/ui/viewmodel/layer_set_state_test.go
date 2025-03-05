package viewmodel

import (
	"testing"
)

func TestGetCompareIndexes(t *testing.T) {
	tests := []struct {
		name              string
		layerIndex        int
		compareMode       LayerCompareMode
		compareStartIndex int
		expected          [4]int
	}{
		{
			name:              "LayerIndex equals CompareStartIndex",
			layerIndex:        2,
			compareMode:       CompareSingleLayer,
			compareStartIndex: 2,
			expected:          [4]int{2, 2, 2, 2},
		},
		{
			name:              "CompareMode is CompareSingleLayer",
			layerIndex:        3,
			compareMode:       CompareSingleLayer,
			compareStartIndex: 1,
			expected:          [4]int{1, 2, 3, 3},
		},
		{
			name:              "Default CompareMode",
			layerIndex:        4,
			compareMode:       CompareAllLayers,
			compareStartIndex: 1,
			expected:          [4]int{1, 1, 2, 4},
		},
		{
			name:              "CompareMode with zero layerIndex",
			layerIndex:        0,
			compareMode:       CompareSingleLayer,
			compareStartIndex: 1,
			expected:          [4]int{1, 1, 0, 0},
		},
		{
			name:              "CompareMode with negative compareStartIndex",
			layerIndex:        2,
			compareMode:       CompareSingleLayer,
			compareStartIndex: -1,
			expected:          [4]int{0, 1, 2, 2},
		},
		{
			name:              "CompareAllLayers with equal indexes",
			layerIndex:        3,
			compareMode:       CompareAllLayers,
			compareStartIndex: 3,
			expected:          [4]int{3, 3, 3, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &LayerSetState{
				LayerIndex:        tt.layerIndex,
				CompareMode:       tt.compareMode,
				CompareStartIndex: tt.compareStartIndex,
			}
			bottomTreeStart, bottomTreeStop, topTreeStart, topTreeStop := state.GetCompareIndexes()
			actual := [4]int{bottomTreeStart, bottomTreeStop, topTreeStart, topTreeStop}
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
