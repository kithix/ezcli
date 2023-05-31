package cmds

import (
	"fmt"

	"github.com/kithix/ezcli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var fooArgs = &struct {
	Bar string `flag:"bar" env:""`
}{
	"default value",
}

var FooApp = ezcli.New(&cobra.Command{
	Use: "foo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("parent:", cmd.Parent().Name())
		fmt.Println("parent flags:")
		cmd.Parent().Flags().VisitAll(func(flag *pflag.Flag) {
			fmt.Printf("\t%s=%s\n", flag.Name, flag.Value)
		})
		fmt.Printf("foo flags: %+v\n", fooArgs)
		fmt.Println("args:", args)
	},
})

func init() {
	FooApp.StructVar(fooArgs)
	FooApp.Child(BarApp)
}
