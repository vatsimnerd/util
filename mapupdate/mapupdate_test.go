package mapupdate

import (
	"reflect"
	"sync"
	"testing"
)

type DataType string

func (d DataType) NE(o DataType) bool {
	return d != o
}

var (
	origMap = map[string]DataType{
		"a": "a",
		"b": "b",
		"d": "d",
	}

	newMap = map[string]DataType{
		"a": "A",
		"b": "b",
		"c": "C",
	}

	expectedSet = map[string]DataType{
		"a": "A",
		"c": "C",
	}

	expectedDel = map[string]DataType{
		"d": "d",
	}
)

func TestMapDiff(t *testing.T) {

	set, del := Update[DataType, Comparable[DataType]](origMap, newMap, &sync.RWMutex{})

	if !reflect.DeepEqual(set, expectedSet) {
		t.Errorf("set map %v differs from what's expected %v", set, expectedSet)
	}

	if !reflect.DeepEqual(del, expectedDel) {
		t.Errorf("set map %v differs from what's expected %v", set, expectedSet)
	}

	if !reflect.DeepEqual(origMap, newMap) {
		t.Errorf("new map %v differs from what's expected %v", origMap, newMap)
	}

}
