package shared

import (
	"sync"
	"testing"
)

func TestNewMapConv(t *testing.T) {
	expected := &MapConv{
		m: map[string]interface{}{
			"a": "a",
			"b": 1,
			"c": true,
			"d": 1.0,
		},
		mu: sync.RWMutex{},
	}

	actual := NewMapConv(map[string]interface{}{
		"a": "a",
		"b": 1,
		"c": true,
		"d": 1.0,
	})

	if actual.String() != expected.String() {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func TestMapConv_Get(t *testing.T) {
	m := NewMapConv(map[string]interface{}{
		"a": "a",
		"b": 1,
		"c": true,
		"d": 1.0,
	})

	expected := "a"
	actual, _ := m.Get("a")
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

	expected2 := 1
	actual, _ = m.Get("b")
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected2, actual)
	}

}
