package main

import (
	"log"

	"github.com/kithix/ezcli"
	"github.com/spf13/cobra"
)

type config struct {
	// Expect env to automatically capitalise
	Vbool     bool `ezcli:"name" env:"not_bOOL"`
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

var subconf = &config{}
var subcmd = root.Child(ezcli.New(
	"subcmd",
	"look a subcmd",
	"just a longer sub command description",
	do,
))

var subsubconf = &config{}
var subsubcmd = subcmd.Child(ezcli.New(
	"subsubcmd",
	"look a subsubcmd",
	"just a longer subsub command description",
	do,
))

func do(cmd *cobra.Command, args []string) {
	log.Println(args)
	log.Printf("root: %+v\n", conf)
	log.Printf("sub: %+v\n", subconf)
	log.Printf("subsub: %+v\n", subsubconf)
}

func init() {
	root.Cmd.Version = "1.0.2"
	// Global configs
	root.BoolVar(&conf.Vbool, "dabool", true, "this is a bool")
	root.IntVar(&conf.Vint, "daint", 5, "this is an int")
	root.StringVar(&conf.Vstring, "dastring", "apples", "this is an int")

	root.StructVar(conf)

	// Inherits from root
	subcmd.BoolVar(&subconf.Vbool, "subbool", false, "sub bool")

	// Inherits from subcmd and root
	subsubcmd.BoolVar(&subsubconf.Vbool, "subsubbool", false, "sub sub bool")

	// init
	root.Init("", ".root")
}

func main() {
	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
