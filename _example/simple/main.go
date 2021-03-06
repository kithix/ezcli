package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/kithix/ezcli"
	"github.com/spf13/cobra"
)

type printMeArgs struct {
	Uppercase bool `flag:"upper" env:""`
}

var appArgs = &printMeArgs{}

func Do(cmd *cobra.Command, args []string) {
	var transform func(string) string

	if appArgs.Uppercase {
		transform = strings.ToUpper
	} else {
		transform = strings.ToLower
	}

	for _, s := range args {
		fmt.Println(transform(s))
	}
}

var app = ezcli.New(&cobra.Command{
	Use:   "printme",
	Short: "print some stuff",
	Run:   Do,
	Args:  cobra.MinimumNArgs(1),
})

func main() {
	app.StructVar(appArgs)
	app.InitNoConfig()

	err := app.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
