package utilities

import (
	"testing"
	"reflect"
)

func TestArgsort(t *testing.T) {
	my_slice := []float64 {5,6,3,2,9,1}
	result := Argsort(my_slice)
	expected := []int {4,1,0,2,3,5}

	if reflect.DeepEqual(result, expected) == false {
		t.Errorf("Expected %x, but got %x", expected, result)
	}
}
