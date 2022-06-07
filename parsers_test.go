package ezcli

import (
	"reflect"
	"testing"
	"time"
)

func TestParseDurationSlice(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  []time.Duration
		err  error
	}{
		{"flag", "[5s,2h]", []time.Duration{5 * time.Second, 2 * time.Hour}, nil},
		{"env", "5s 2h", []time.Duration{5 * time.Second, 2 * time.Hour}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseDurationSlice(test.in)
			if err != test.err {
				t.Errorf("got error '%v' expected '%v'", err, test.err)
				return
			}
			expected := test.out
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("got '%s' expected '%s'", got, expected)
			}
		})
	}
}
