package ezcli

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func subject() *App {
	return New(&cobra.Command{})
}

func doGVarTest[T any](t *testing.T, flagValue, envValue string, val T) {
	var testType T
	t.Run(reflect.TypeOf(testType).String(), gVarTest(flagValue, envValue, val))
}

func doGVarFlagTest[T any](t *testing.T, flagValue string, val T) {
	doGVarTest(t, flagValue, "", val)
}

func doGVarEnvTest[T any](t *testing.T, envValue string, val T) {
	doGVarTest(t, "", envValue, val)
}

func gVarTest[T any](flagValue, envValue string, val T) func(t *testing.T) {
	return func(t *testing.T) {
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
	}
}

func TestApp_Vars(t *testing.T) {
	// Single values
	doGVarFlagTest[bool](t, "true", true)
	doGVarFlagTest[string](t, "testString", "testString")

	doGVarFlagTest[int](t, "1337", 1337)
	doGVarFlagTest[int8](t, "16", 16)
	doGVarFlagTest[int16](t, "3200", 3200)
	doGVarFlagTest[int32](t, "5678123", 5678123)
	doGVarFlagTest[int64](t, "1234567890", 1234567890)

	doGVarFlagTest[uint](t, "7331", 7331)
	doGVarFlagTest[uint8](t, "32", 32)
	doGVarFlagTest[uint16](t, "2509", 2509)
	doGVarFlagTest[uint32](t, "8123567", 8123567)
	doGVarFlagTest[uint64](t, "10987654321", 10987654321)

	doGVarFlagTest[time.Duration](t, "5s", 5*time.Second)
	doGVarFlagTest[net.IP](t, "127.0.0.1", net.IPv4(127, 0, 0, 1))
	doGVarFlagTest[net.IP](t, "ff02::1", net.IPv6linklocalallnodes)

	// Slices
	doGVarFlagTest[[]string](t, "s1,s2", []string{"s1", "s2"})
	doGVarFlagTest[[]time.Duration](t, "5s,2h", []time.Duration{5 * time.Second, 2 * time.Hour})
}

func TestApp_FromEnv(t *testing.T) {
	// Do all the above tests but use environment instead
}

func TestApp_VarsThatPanic(t *testing.T) {
	// Type aliases
	type StringAlias string
	assertPanics[StringAlias](t, gVarTest("testStringAlias", "", StringAlias("testStringAlias")))
}

func assertPanics[T any](t *testing.T, fn func(t *testing.T)) {
	var testType T
	t.Run(reflect.TypeOf(testType).String(), func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("code did not panic")
			}
		}()
		fn(t)
	})
}
