package ezcli

import (
	"fmt"
	"net"
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
	Viper         *viper.Viper
	postLoadFuncs []func()
}

func New(cmd *cobra.Command) *App {
	a := &App{
		Cmd:           cmd,
		Viper:         viper.New(),
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

func (a *App) genericVar(v any, optFns ...varOptFn) {
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

	typeOf := reflect.TypeOf(v)
	// We must have a pointer before continuing
	if typeOf.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("type must be a pointer, got %s:", typeOf))
	}
	// Get the value of the pointer
	elem := typeOf.Elem()
	// TODO should we panic if it's a pointer here?
	if elem.Kind() == reflect.Pointer {
		panic("unable to service a pointer to a pointer")
	}
	// Get the value out so we can set it later
	val := reflect.ValueOf(v).Elem()

	// Ensure we have a zero'd value for our type
	if opts.DefaultValue == nil {
		opts.DefaultValue = reflect.Zero(elem).Interface()
	}

	var postLoadFunc func()

	// Set the flag for the kind of data
	switch elem.String() {
	// TODO Could we and should we allow type aliases from users?
	case "bool":
		flagSet.BoolVar(v.(*bool), opts.Name, opts.DefaultValue.(bool), opts.Usage)
		postLoadFunc = func() { val.SetBool(a.Viper.GetBool(opts.Name)) }

	case "int":
		flagSet.IntVar(v.(*int), opts.Name, opts.DefaultValue.(int), opts.Usage)
		postLoadFunc = func() { val.SetInt(a.Viper.GetInt64(opts.Name)) }

	case "int8":
		flagSet.Int8Var(v.(*int8), opts.Name, opts.DefaultValue.(int8), opts.Usage)
		postLoadFunc = func() { val.SetInt(a.Viper.GetInt64(opts.Name)) }

	case "int16":
		flagSet.Int16Var(v.(*int16), opts.Name, opts.DefaultValue.(int16), opts.Usage)
		postLoadFunc = func() { val.SetInt(a.Viper.GetInt64(opts.Name)) }

	case "int32":
		flagSet.Int32Var(v.(*int32), opts.Name, opts.DefaultValue.(int32), opts.Usage)
		postLoadFunc = func() { val.SetInt(a.Viper.GetInt64(opts.Name)) }

	case "int64":
		flagSet.Int64Var(v.(*int64), opts.Name, opts.DefaultValue.(int64), opts.Usage)
		postLoadFunc = func() { val.SetInt(a.Viper.GetInt64(opts.Name)) }

	case "uint":
		flagSet.UintVar(v.(*uint), opts.Name, opts.DefaultValue.(uint), opts.Usage)
		postLoadFunc = func() { val.SetUint(a.Viper.GetUint64(opts.Name)) }

	case "uint8":
		flagSet.Uint8Var(v.(*uint8), opts.Name, opts.DefaultValue.(uint8), opts.Usage)
		postLoadFunc = func() { val.SetUint(a.Viper.GetUint64(opts.Name)) }

	case "uint16":
		flagSet.Uint16Var(v.(*uint16), opts.Name, opts.DefaultValue.(uint16), opts.Usage)
		postLoadFunc = func() { val.SetUint(a.Viper.GetUint64(opts.Name)) }

	case "uint32":
		flagSet.Uint32Var(v.(*uint32), opts.Name, opts.DefaultValue.(uint32), opts.Usage)
		postLoadFunc = func() { val.SetUint(a.Viper.GetUint64(opts.Name)) }

	case "uint64":
		flagSet.Uint64Var(v.(*uint64), opts.Name, opts.DefaultValue.(uint64), opts.Usage)
		postLoadFunc = func() { val.SetUint(a.Viper.GetUint64(opts.Name)) }

	case "net.IP":
		flagSet.IPVar(v.(*net.IP), opts.Name, opts.DefaultValue.(net.IP), opts.Usage)
		postLoadFunc = func() {
			ipString := a.Viper.GetString(opts.Name)
			v = net.ParseIP(ipString)
		}

	case "string":
		flagSet.StringVar(v.(*string), opts.Name, opts.DefaultValue.(string), opts.Usage)
		postLoadFunc = func() { val.SetString(a.Viper.GetString(opts.Name)) }

	case "[]string":
		// TODO validate with environment var
		flagSet.StringSliceVar(v.(*[]string), opts.Name, opts.DefaultValue.([]string), opts.Usage)
		postLoadFunc = func() { v = a.Viper.GetStringSlice(opts.Name) }

	case "time.Duration":
		flagSet.DurationVar(v.(*time.Duration), opts.Name, opts.DefaultValue.(time.Duration), opts.Usage)
		postLoadFunc = func() { v = a.Viper.GetDuration(opts.Name) }

	case "[]time.Duration":
		flagSet.DurationSliceVar(v.(*[]time.Duration), opts.Name, opts.DefaultValue.([]time.Duration), opts.Usage)
		postLoadFunc = func() {
			// Format of: [durationString,durationString]
			durationSliceAsString := a.Viper.GetString(opts.Name)
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
			v = durations
		}
	default:
		// TODO how to handle aliases of base types?
		// Should we even do this?
		panic(fmt.Sprintf("unable to use variable type %s", elem))
	}

	// Prepare our post load function
	a.postLoadFuncs = append(a.postLoadFuncs, postLoadFunc)

	// Bind the cobra flag to Viper for configuration file and environment mapping
	a.Viper.BindPFlag(opts.Name, flagSet.Lookup(opts.Name))
	a.Viper.BindEnv(opts.Name)
}

func (a *App) Init(pathToConfigFile, configName string) {
	cobra.OnInitialize(a.initConfig(pathToConfigFile, configName))
}

func (a *App) InitNoConfig() {
	cobra.OnInitialize(a.init)
}

func (a *App) initConfig(pathToConfigFile, configName string) func() {
	return func() {
		if pathToConfigFile != "" {
			// Use config file from the flag.
			a.Viper.SetConfigFile(pathToConfigFile)
		} else {
			// Find home directory for config
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".nrf-go" (without extension).
			a.Viper.AddConfigPath(home)
			a.Viper.SetConfigName(configName)
		}

		// If a config file is found, read it in.
		err := a.Viper.ReadInConfig()
		if err == nil {
			fmt.Println("Using config file:", a.Viper.ConfigFileUsed())
		}

		a.init()
	}
}

func (a *App) init() {
	err := a.Viper.BindPFlags(a.Cmd.Flags())
	if err != nil {
		panic(err)
	}
	// Now we have our config, override the things
	for _, fn := range a.postLoadFuncs {
		fn()
	}
}

func (a *App) Execute() error {
	return a.Cmd.Execute()
}
