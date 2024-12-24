package cut

import (
	"bytes"
	"strings"
	"testing"
)

func TestCut(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  Options
		expected string
	}{
		{
			name:  "basic functionality",
			input: "name\tage\tcountry\nAlice\t30\tUSA\nBob\t25\tUK",
			options: Options{
				fields:    []int{0, 2},
				delimiter: "\t",
				separated: false,
			},
			expected: "name\tcountry\nAlice\tUSA\nBob\tUK\n",
		},
		{
			name:  "custom delimiter",
			input: "name|age|country\nAlice|30|USA\nBob|25|UK",
			options: Options{
				fields:    []int{1},
				delimiter: "|",
				separated: false,
			},
			expected: "age\n30\n25\n",
		},
		{
			name:  "separated only",
			input: "name\nage\t30\tUSA",
			options: Options{
				fields:    []int{0},
				delimiter: "\t",
				separated: true,
			},
			expected: "age\n",
		},
		{
			name:  "out of range fields",
			input: "name\tage\tcountry",
			options: Options{
				fields:    []int{5},
				delimiter: "\t",
				separated: false,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			var output bytes.Buffer

			err := cut(input, &output, tt.options)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output.String())
			}
		})
	}
}
