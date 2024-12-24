package unpack

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"qwe\\\\5", "qwe\\\\\\", false},
		{"a10", "aaaaaaaaaa", false},
		{"a\\2b", "a2b", false},
		{"\\3\\3", "33", false},
	}

	for _, test := range tests {
		result, err := UnpackString(test.input)
		if (err != nil) != test.err {
			t.Errorf("UnpackString(%q): unexpected error status: got %v, want error: %v", test.input, err, test.err)
		}
		if result != test.expected {
			t.Errorf("UnpackString(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
