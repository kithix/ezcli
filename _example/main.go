package main

import (
	"log"

	"github.com/kithix/ezcli"
	"github.com/spf13/cobra"
)

type config struct {
	Vbool     bool
	Vint      int
	Vstring   string
	Vduration string
}

var conf = &config{}
var root = ezcli.New(
	"root",
	"just a root command",
	"just a longer root command description",
	do,
)

func do(cmd *cobra.Command, args []string) {
	log.Println(args)
	log.Printf("%+v\n", conf)
}

func init() {
	// Register configs
	root.BoolVar(&conf.Vbool, "dabool", true, "this is a bool")
	root.IntVar(&conf.Vint, "daint", 5, "this is an int")
	root.StringVar(&conf.Vstring, "dastring", "apples", "this is an int")

	// init
	root.Init("", ".root")
}

func main() {
	log.Println("running main")
	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
