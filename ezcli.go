package ezcli

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type App struct {
	Cmd           *cobra.Command
	postLoadFuncs []func()
}

func New(cmd *cobra.Command) *App {
	a := &App{
		Cmd:           cmd,
		postLoadFuncs: make([]func(), 0),
	}

	return a
}

func (a *App) Child(child *App) *App {
	a.Cmd.AddCommand(child.Cmd)
	return child
}

func (a *App) StringVar(variable *string, name, value, usage string) {
	a.genericVar(variable, VarName(name), VarDefaultValue(value), VarUsage(usage))
}

func (a *App) IntVar(variable *int, name string, value int, usage string) {
	a.genericVar(variable, VarName(name), VarDefaultValue(value), VarUsage(usage))
}

func (a *App) DurationVar(variable *time.Duration, name string, value time.Duration, usage string) {
	a.genericVar(variable, VarName(name), VarDefaultValue(value), VarUsage(usage))
}

func (a *App) BoolVar(variable *bool, name string, value bool, usage string) {
	a.genericVar(variable, VarName(name), VarDefaultValue(value), VarUsage(usage))
}

func (a *App) Var(variable any, name string, value any, usage string) {
	a.genericVar(variable, VarName(name), VarDefaultValue(value), VarUsage(usage))
}

func (a *App) genericVar(variable any, optFns ...varOptFn) {
	// Setup our variable options
	opts := defaultOpts()
	for _, optFn := range optFns {
		optFn(opts)
	}
	if opts.Name == "" {
		panic("no name provided for variable")
	}
	// Get the appropriate cobra flagSet for use later
	// Local flags only apply to this command
	// Persistent flags apply to all sub-commands
	flagSet := a.Cmd.LocalFlags()
	if opts.Persistent {
		flagSet = a.Cmd.PersistentFlags()
	}

	typeOf := reflect.TypeOf(variable)
	// We must have a pointer before continuing
	if typeOf.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("Type must be a pointer, got %s:", typeOf))
	}
	// Get the
	elem := typeOf.Elem()
	// Ensure we have a zero'd value for our type
	if opts.DefaultValue == nil {
		opts.DefaultValue = reflect.Zero(elem).Interface()
	}

	var postLoadFunc func()
	// Set the flag for the kind of data
	switch elem.String() {
	// TODO Could we and should we allow type aliases from users?
	case "bool":
		flagSet.BoolVar(variable.(*bool), opts.Name, opts.DefaultValue.(bool), opts.Usage)
		postLoadFunc = func() { variable = viper.GetBool(opts.Name) }

	case "int":
		flagSet.IntVar(variable.(*int), opts.Name, opts.DefaultValue.(int), opts.Usage)
		postLoadFunc = func() { variable = viper.GetInt(opts.Name) }

	case "string":
		flagSet.StringVar(variable.(*string), opts.Name, opts.DefaultValue.(string), opts.Usage)
		postLoadFunc = func() { variable = viper.GetString(opts.Name) }

	case "[]string":
		flagSet.StringSliceVar(variable.(*[]string), opts.Name, opts.DefaultValue.([]string), opts.Usage)
		postLoadFunc = func() { variable = viper.GetStringSlice(opts.Name) }

	case "time.Duration":
		flagSet.DurationVar(variable.(*time.Duration), opts.Name, opts.DefaultValue.(time.Duration), opts.Usage)
		postLoadFunc = func() { variable = viper.GetDuration(opts.Name) }

	case "[]time.Duration":
		flagSet.DurationSliceVar(variable.(*[]time.Duration), opts.Name, opts.DefaultValue.([]time.Duration), opts.Usage)
		postLoadFunc = func() {
			// Format of: [durationString,durationString]
			durationSliceAsString := viper.Get(opts.Name).(string)
			// Remove the brackets and split on the commas
			durationStrings := strings.Split(durationSliceAsString[1:len(durationSliceAsString)-1], ",")
			// Populate our slice with the duration values
			durations := make([]time.Duration, len(durationStrings))
			for i, durationString := range durationStrings {
				d, err := time.ParseDuration(durationString)
				if err != nil {
					panic(err)
				}
				durations[i] = d
			}
			variable = durations
		}
	default:
		panic(fmt.Sprintf("unable to use variable type %s", elem))
	}

	// Prepare our post load function
	a.postLoadFuncs = append(a.postLoadFuncs, postLoadFunc)

	// Bind the cobra flag to Viper for configuration file and environment mapping
	viper.BindPFlag(opts.Name, flagSet.Lookup(opts.Name))
}

func (app *App) Init(pathToConfigFile, configName string) {
	cobra.OnInitialize(app.initConfig(pathToConfigFile, configName))
}

func (a *App) initConfig(pathToConfigFile, configName string) func() {
	return func() {
		if pathToConfigFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(pathToConfigFile)
		} else {
			// Find home directory for config
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".nrf-go" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(configName)
		}

		// Allow any environment variables to match
		viper.AutomaticEnv()
		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
		// Now we have our config, override the things
		for _, fn := range a.postLoadFuncs {
			fn()
		}
	}
}

func (a *App) Execute() error {
	return a.Cmd.Execute()
}
