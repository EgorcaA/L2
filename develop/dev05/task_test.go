package anagrams

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	tests := []struct {
		input    []string
		expected map[string][]string
	}{
		{
			[]string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "ток", "кот", "окт"},
			map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
				"ток":    {"кот", "окт", "ток"},
			},
		},
		{
			[]string{"кот", "ток", "окт", "коТ", "тоК"},
			map[string][]string{
				"кот": {"кот", "окт", "ток"},
			},
		},
		{
			[]string{"abc", "cab", "bac", "cba", "xyz"},
			map[string][]string{
				"abc": {"abc", "bac", "cab", "cba"},
			},
		},
		{
			[]string{"один", "два", "три"},
			map[string][]string{},
		},
	}

	for _, test := range tests {
		result := FindAnagrams(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For input %v expected %v but got %v", test.input, test.expected, result)
		}
	}
}
