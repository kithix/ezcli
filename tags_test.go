package ezcli

import (
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	tests := []struct {
		name string
		in   reflect.StructField
		app  *App
		err  error
	}{
		{"bool", reflect.StructField{}, &App{}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts := parseTags(test.in)
			var err error
			app := &App{}
			// Apply our options to our app
			for _, opt := range opts {
				err = opt(app)
				if err != nil {
					break
				}
			}

			// Check if we wanted the error
			if err != test.err {
				if test.err == nil {
					t.Errorf("Unexpected error '%s'", err)
				} else {
					t.Errorf("Expected error '%s' got '%s'", test.err, err)
				}
				return
			}

			// Validate if the application values are the same
			// TODO I expect this will not work for everything
			if !reflect.DeepEqual(test.app, app) {
				t.Errorf("expected:\n%+v\ngot:\n%+v\n", test.app, app)
			}
		})
	}
}
