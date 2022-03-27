package partitioner

import (
	"reflect"
	"testing"
)

func TestKarmarkarKarp(t *testing.T) {
	expected := [][]int64{{3, 9}, {5, 7}, {10, 1}}
	actual := KarmarkarKarp([]int64{1, 7, 5, 10, 9, 3}, 3)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %+v, Actual: %+v", expected, actual)
	}
}
