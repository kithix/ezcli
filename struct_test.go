package ezcli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestApp_StructVar(t *testing.T) {
	type TestStruct struct {
		unexported string
		String     string
		Bool       bool
		Int        int
		UInt       uint
		Inner      struct {
			InnerStr string
		}
	}
	s := &TestStruct{
		unexported: "untouched",
		Bool:       true,
	}

	app := New("cmd", "short", "long", func(cmd *cobra.Command, args []string) {})
	app.StructVar(s)
}
