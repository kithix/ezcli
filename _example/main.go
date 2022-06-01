package main

import (
	"log"
	"time"

	"github.com/kithix/ezcli"
	"github.com/spf13/cobra"
)

type config struct {
	// Expect env to automatically capitalise
	Bool        bool `flag:"name"`
	Int         int  `env:"nOt_Int"`
	String      string
	Duration    time.Duration
	SliceString []string
}

var conf = &config{}
var root = ezcli.New(&cobra.Command{
	Use:   "root",
	Short: "just a root command",
	Long:  "just a longer root command description",
	Run:   do,
})

var subconf = &config{}
var subcmd = root.Child(ezcli.New(&cobra.Command{
	Use:   "subcmd",
	Short: "look a subcmd",
	Long:  "just a longer sub command description",
	Run:   do,
}))

var subsubconf = &config{}
var subsubcmd = subcmd.Child(ezcli.New(&cobra.Command{
	Use:   "subsubcmd",
	Short: "look a subsubcmd",
	Long:  "just a longer subsub command description",
	Run:   do,
}))

func do(cmd *cobra.Command, args []string) {
	log.Println(args)
	log.Printf("root: %+v\n", conf)
	log.Printf("sub: %+v\n", subconf)
	log.Printf("subsub: %+v\n", subsubconf)
}

func init() {
	root.Cmd.Version = "1.0.2"
	root.StructVar(conf)
	// Global configs
	root.BoolVar(&conf.Bool, "dabool", true, "this is a bool")
	root.IntVar(&conf.Int, "daint", 5, "this is an int")
	root.StringVar(&conf.String, "dastring", "apples", "this is an int")

	root.Var(&conf.SliceString, "daslicestring", []string{"apples"}, "this is an int")

	// Inherits from root
	subcmd.BoolVar(&subconf.Bool, "subbool", false, "sub bool")

	// Inherits from subcmd and root
	subsubcmd.BoolVar(&subsubconf.Bool, "subsubbool", false, "sub sub bool")

	// init
	root.Init("", ".root")
}

func main() {
	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
