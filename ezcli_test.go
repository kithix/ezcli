package ezcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func subject() *App {
	return New(&cobra.Command{})
}

func doGVarTest[T any](t *testing.T, flagValue, envValue string, configValue any, val T) {
	var testType T
	t.Run(reflect.TypeOf(testType).String(), gVarTest(flagValue, envValue, configValue, val))
}

func doGVarFlagTest[T any](t *testing.T, flagValue string, val T) {
	doGVarTest(t, flagValue, "", nil, val)
}

func doGVarEnvTest[T any](t *testing.T, envValue string, val T) {
	doGVarTest(t, "", envValue, nil, val)
}

func doGVarConfigTest[T any](t *testing.T, configValue any, val T) {
	doGVarTest(t, "", "", configValue, val)
}

func doGVarEnvAndFlagTest[T any](t *testing.T, flagValue, envValue string, val T) {
	doGVarTest(t, flagValue, envValue, nil, val)
}

func gVarTest[T any](flagValue, envValue string, configValue any, val T) func(t *testing.T) {
	return func(t *testing.T) {
		name := t.Name()
		opts := []varOptFn{VarName(name)}

		if envValue != "" {
			// Mimic environment setting
			err := os.Setenv(name, envValue)
			if err != nil {
				t.Error("unable to set env", err)
				return
			}
			// Clean-up after ourselves
			defer os.Setenv(name, "")
			// Validate we set it and it matches
			if os.Getenv(name) != envValue {
				t.Error("unable to retrieve value")
				return
			}
			opts = append(opts, VarEnv(name))
		}

		app := subject()

		// Prepare our configuration if we have a value for one
		if configValue != nil {
			jsonVal, err := json.Marshal(configValue)
			if err != nil {
				t.Error("unable to turn value to json", val)
				return
			}
			jsonConfBytes := []byte(fmt.Sprintf(
				"{\"%s\":%s}", name, string(jsonVal),
			))
			app.Viper.SetConfigType("json")
			err = app.Viper.ReadConfig(bytes.NewReader(jsonConfBytes))
			if err != nil {
				t.Error("unable to read config", err)
				return
			}
			t.Log("Input Config:", string(jsonConfBytes))
		}

		var testVar T
		app.genericVar(&testVar, opts...)

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
		app.InitNoConfig()

		// Ensure the values match as their type
		if !reflect.DeepEqual(val, testVar) {
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

func TestApp_FromConfig_JSON(t *testing.T) {
	// Single values
	doGVarConfigTest[bool](t, true, true)
	doGVarConfigTest[string](t, "testString", "testString")

	doGVarConfigTest[int](t, 1337, 1337)
	doGVarConfigTest[int8](t, 16, 16)
	doGVarConfigTest[int16](t, 3200, 3200)
	doGVarConfigTest[int32](t, 5678123, 5678123)
	doGVarConfigTest[int64](t, 1234567890, 1234567890)

	doGVarConfigTest[uint](t, 7331, 7331)
	doGVarConfigTest[uint8](t, 32, 32)
	doGVarConfigTest[uint16](t, 2509, 2509)
	doGVarConfigTest[uint32](t, 8123567, 8123567)
	doGVarConfigTest[uint64](t, 10987654321, 10987654321)

	doGVarConfigTest[time.Duration](t, "5s", 5*time.Second)
	doGVarConfigTest[net.IP](t, "127.0.0.1", net.IPv4(127, 0, 0, 1))
	doGVarEnvTest[net.IP](t, "ff02::1", net.IPv6linklocalallnodes)

	// Slices
	doGVarConfigTest[[]string](t, []string{"s1", "s2"}, []string{"s1", "s2"})
	doGVarConfigTest[[]time.Duration](t, []string{"5s", "2h"}, []time.Duration{5 * time.Second, 2 * time.Hour})
}

func TestApp_FromEnvAndFlag(t *testing.T) {
	// Test that flags take priority over environment variables
	// Single values
	doGVarEnvAndFlagTest[bool](t, "false", "true", false)
	doGVarEnvAndFlagTest[string](t, "flagString", "envString", "flagString")

	doGVarEnvAndFlagTest[int](t, "1337", "7331", 1337)
	doGVarEnvAndFlagTest[int8](t, "16", "61", 16)
	doGVarEnvAndFlagTest[int16](t, "3200", "1234", 3200)
	doGVarEnvAndFlagTest[int32](t, "5678123", "3218765", 5678123)
	doGVarEnvAndFlagTest[int64](t, "1234567890", "9876543210", 1234567890)

	doGVarEnvAndFlagTest[uint](t, "7331", "1337", 7331)
	doGVarEnvAndFlagTest[uint8](t, "32", "23", 32)
	doGVarEnvAndFlagTest[uint16](t, "2509", "9052", 2509)
	doGVarEnvAndFlagTest[uint32](t, "8123567", "7652812", 8123567)
	doGVarEnvAndFlagTest[uint64](t, "10987654321", "12345678901", 10987654321)

	doGVarEnvAndFlagTest[time.Duration](t, "5s", "10h", 5*time.Second)
	doGVarEnvAndFlagTest[net.IP](t, "127.0.0.1", "1.2.3.4", net.IPv4(127, 0, 0, 1))
	doGVarEnvAndFlagTest[net.IP](t, "ff02::1", "ab:cd:ef:12:34:56:78::90", net.IPv6linklocalallnodes)

	// Slices
	doGVarEnvAndFlagTest[[]string](t, "s1,s2", "first second", []string{"s1", "s2"})
	doGVarEnvAndFlagTest[[]time.Duration](t, "5s,2h", "6m 3ms", []time.Duration{5 * time.Second, 2 * time.Hour})
}

// TODO fuzz tests
func TestApp_ManualFromConfig(t *testing.T) {
	config := []byte("{\"name\":\"value\"}")
	app := subject()
	app.Viper.SetConfigType("json")
	err := app.Viper.ReadConfig(bytes.NewReader(config))
	if err != nil {
		t.Error(err)
		return
	}
	in := ""
	app.StringVar(&in, "name", "default", "usage")
	app.InitNoConfig()

	if in != "value" {
		t.Errorf("expected 'value' got '%s'", in)
	}
}

func TestApp_VarsThatPanic(t *testing.T) {
	// Type aliases
	// TODO - it is likely possible to hanlde custom types, just needs investigation
	type StringAlias string
	assertPanics[StringAlias](t, gVarTest("testStringAlias", "", "", StringAlias("testStringAlias")))
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
