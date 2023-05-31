package cmds

import (
	"fmt"

	"github.com/kithix/ezcli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type barConfig struct {
	Bool          bool
	String        string
	StringDefault string
	Int           int
}

var barArgs = &barConfig{
	StringDefault: "default",
}

var BarApp = ezcli.New(&cobra.Command{
	Use: "bar",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("parent:", cmd.Parent().Name())
		fmt.Println("parent flags:")
		cmd.Parent().Flags().VisitAll(func(flag *pflag.Flag) {
			fmt.Printf("\t%s=%s\n", flag.Name, flag.Value)
		})
		fmt.Printf("bar flags: %+v\n", barArgs)
		fmt.Println("args:", args)
	},
})

func init() {
	BarApp.StructVar(barArgs)
}
