package main

import (
	"fmt"
	"log"

	"github.com/kithix/ezcli"
	"github.com/kithix/ezcli/_example/multiple/cmds"
	"github.com/spf13/cobra"
)

// Example of how you can one line your configuration if desired
var config = &struct {
	Extra bool `flag:"extra" env:""`
}{}

var app = ezcli.New(&cobra.Command{
	Use: "multi",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("multi flags: %+v\n", config)
		fmt.Println("args:", args)
	},
})

func init() {
	app.StructVar(config)
	app.Child(cmds.FooApp)
}

func main() {
	app.InitNoConfig()

	err := app.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
