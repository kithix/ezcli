package ezcli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestApp_StructVar(t *testing.T) {
	type TestStruct struct {
		unexported     string
		String         string `flag:"custom" env:""`
		Bool           bool   `flag:""`
		Int            int    `env:""`
		UInt           uint
		InnerNoPointer struct {
			InnerStr string
		}
		InnerPointer struct {
			InnterStr string
		}
	}
	s := &TestStruct{
		unexported: "untouched",
		Int:        1337,
		Bool:       true,
	}

	app := New(&cobra.Command{
		Use:   "cmd",
		Short: "short",
		Long:  "long",
		Run:   func(cmd *cobra.Command, args []string) {}})
	app.StructVar(s)

	// Mimic setting our flag from command line
	app.Cmd.ParseFlags([]string{
		// Flag names match the exact field name
		// Should these be automatically lower cased?
		"--Bool=false",
		"--custom=teststring",
	})

	app.initConfig(t.TempDir(), "doesntexist")()

	if s.unexported != "untouched" {
		t.Error("we touched an unexported field")
	}
	if s.Bool {
		t.Error("didn't set to false")
	}
	if s.Int != 1337 {
		t.Errorf("expected default value '%d' got '%d'\n", 1337, s.Int)
	}
	if s.String != "teststring" {
		t.Errorf("expected '%s' got '%s'\n", "teststring", s.String)
	}
}
