package ezcli

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"

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

func createFriendlyName(s string) string {
	s = strings.Replace(s, "[]", "SLICE_", -1)
	s = strings.Replace(s, "/", "_", -1)
	s = strings.Replace(s, ".", "_", -1)
	s = strings.Replace(s, "#", "_", -1)
	return s
}

func viperEnvTest(t *testing.T) {
	v := viper.New()
	// string
	varname := "stringvar"
	envname := "string_TEST"
	expected := "teststring"
	err := os.Setenv(envname, expected)
	if err != nil {
		t.Error(err)
		return
	}
	err = v.BindEnv(varname, envname)
	if err != nil {
		t.Error(err)
		return
	}
	got := v.GetString(varname)
	if got != expected {
		t.Errorf("got '%s' expected '%s'", got, expected)
	}

	// []string
}

func gVarTest[T any](flagValue, envValue string, val T) func(t *testing.T) {
	return func(t *testing.T) {
		name := t.Name()
		envname := createFriendlyName(name)
		// Mimic environment setting
		err := os.Setenv(envname, envValue)
		if err != nil {
			t.Error("unable to set env", err)
			return
		}
		// Clean-up after ourselves
		defer os.Setenv(envname, "")
		// Validate we set it and it matches
		if os.Getenv(envname) != envValue {
			t.Error("unable to retrieve value")
			return
		}

		app := subject()
		var testVar T

		app.genericVar(&testVar, VarName(name), VarEnv(envname))

		// Mimic setting our flag from command line if we have a flag value
		if flagValue != "" {
			app.Cmd.ParseFlags([]string{
				fmt.Sprintf("--%s=%s", name, flagValue),
			})
		}

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
			fmt.Printf("env:\"%s=%s\"\n", envname, envValue)
			t.Errorf("Expected '%v' got '%v'", val, testVar)
		}
	}
}

func TestApp_FromFlag(t *testing.T) {
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
	// Single values
	doGVarEnvTest[bool](t, "true", true)
	doGVarEnvTest[string](t, "testString", "testString")

	doGVarEnvTest[int](t, "1337", 1337)
	doGVarEnvTest[int8](t, "16", 16)
	doGVarEnvTest[int16](t, "3200", 3200)
	doGVarEnvTest[int32](t, "5678123", 5678123)
	doGVarEnvTest[int64](t, "1234567890", 1234567890)

	doGVarEnvTest[uint](t, "7331", 7331)
	doGVarEnvTest[uint8](t, "32", 32)
	doGVarEnvTest[uint16](t, "2509", 2509)
	doGVarEnvTest[uint32](t, "8123567", 8123567)
	doGVarEnvTest[uint64](t, "10987654321", 10987654321)

	doGVarEnvTest[time.Duration](t, "5s", 5*time.Second)
	doGVarEnvTest[net.IP](t, "127.0.0.1", net.IPv4(127, 0, 0, 1))
	doGVarEnvTest[net.IP](t, "ff02::1", net.IPv6linklocalallnodes)

	// Slices - not comma seperated
	doGVarEnvTest[[]string](t, "s1 s2", []string{"s1", "s2"})
	doGVarEnvTest[[]time.Duration](t, "5s 2h", []time.Duration{5 * time.Second, 2 * time.Hour})
}

func TestApp_FromEnvAndFlag(t *testing.T) {
	// Single values
	doGVarTest[bool](t, "false", "true", false)
	doGVarTest[string](t, "flagString", "envString", "flagString")

	doGVarTest[int](t, "1337", "7331", 1337)
	doGVarTest[int8](t, "16", "61", 16)
	doGVarTest[int16](t, "3200", "1234", 3200)
	doGVarTest[int32](t, "5678123", "3218765", 5678123)
	doGVarTest[int64](t, "1234567890", "9876543210", 1234567890)

	doGVarTest[uint](t, "7331", "1337", 7331)
	doGVarTest[uint8](t, "32", "23", 32)
	doGVarTest[uint16](t, "2509", "9052", 2509)
	doGVarTest[uint32](t, "8123567", "7652812", 8123567)
	doGVarTest[uint64](t, "10987654321", "12345678901", 10987654321)

	doGVarTest[time.Duration](t, "5s", "10h", 5*time.Second)
	doGVarTest[net.IP](t, "127.0.0.1", "1.2.3.4", net.IPv4(127, 0, 0, 1))
	doGVarTest[net.IP](t, "ff02::1", "ab:cd:ef:12:34:56:78::90", net.IPv6linklocalallnodes)

	// Slices
	doGVarTest[[]string](t, "s1,s2", "first second", []string{"s1", "s2"})
	doGVarTest[[]time.Duration](t, "5s,2h", "6m 3ms", []time.Duration{5 * time.Second, 2 * time.Hour})
}

func TestApp_FromConfig(t *testing.T) {
	// TODO
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
