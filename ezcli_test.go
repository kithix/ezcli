package ezcli

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func subject() *App {
	return &App{
		Cmd: &cobra.Command{},
	}
}

func doGVarTest[T any](t *testing.T, flagValue string, val T) {
	var testType T
	t.Run(reflect.TypeOf(testType).String(), func(t *testing.T) {
		name := t.Name()
		app := subject()
		var testVar T
		app.genericVar(&testVar, VarName(name))

		// Mimic setting our flag from command line
		app.Cmd.ParseFlags([]string{
			fmt.Sprintf("--%s=%s", name, flagValue),
		})

		// Ensure the flag got created
		flag := app.Cmd.PersistentFlags().Lookup(name)
		if flag == nil {
			t.Error("did not set flag")
			return
		}

		// Initialise the application
		app.initConfig(t.TempDir(), "doesntexist")()

		// Ensure the values match as their type
		if !reflect.DeepEqual(val, testVar) {
			t.Errorf("Expected '%v' got '%v'", val, testVar)
		}
	})
}

func TestApp_Vars(t *testing.T) {
	// Single values
	doGVarTest[bool](t, "true", true)
	doGVarTest[string](t, "testString", "testString")
	doGVarTest[int](t, "1337", 1337)
	doGVarTest[time.Duration](t, "5s", 5*time.Second)

	// Slices
	doGVarTest[[]string](t, "s1,s2", []string{"s1", "s2"})
	doGVarTest[[]time.Duration](t, "5s,2h", []time.Duration{5 * time.Second, 2 * time.Hour})
}
