package ezcli

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type App struct {
	cmd           *cobra.Command
	configName    string
	postLoadFuncs []func()
}

func New(name, short, long string, run func(cmd *cobra.Command, args []string)) *App {
	return &App{
		cmd: &cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
			Run:   run,
		},
		postLoadFuncs: make([]func(), 0),
	}
}

func (app *App) Child(child *App) *App {
	app.cmd.AddCommand(child.cmd)
	return child
}

func (app *App) StringVar(variable *string, name, value, usage string) {
	app.cmd.PersistentFlags().StringVar(variable, name, value, usage)
	bindFlagAndConfig(app.cmd, name)
	app.postLoadFuncs = append(app.postLoadFuncs, func() {
		*variable = viper.GetString(name)
	})
}

func (app *App) IntVar(variable *int, name string, value int, usage string) {
	app.cmd.PersistentFlags().IntVar(variable, name, value, usage)
	bindFlagAndConfig(app.cmd, name)
	app.postLoadFuncs = append(app.postLoadFuncs, func() {
		*variable = viper.GetInt(name)
	})
}

func (app *App) DurationVar(variable *time.Duration, name string, value time.Duration, usage string) {
	app.cmd.PersistentFlags().DurationVar(variable, name, value, usage)
	bindFlagAndConfig(app.cmd, name)
	app.postLoadFuncs = append(app.postLoadFuncs, func() {
		*variable = viper.GetDuration(name)
	})
}

func (app *App) BoolVar(variable *bool, name string, value bool, usage string) {
	app.cmd.PersistentFlags().BoolVar(variable, name, value, usage)
	bindFlagAndConfig(app.cmd, name)
	app.postLoadFuncs = append(app.postLoadFuncs, func() {
		*variable = viper.GetBool(name)
	})
}

// TODO options
func (app *App) StructVar(s interface{}) {
	// Must be a pointer
	t := reflect.TypeOf(s)
	fmt.Println(t.Name())
	fmt.Println(t.Kind())

	t = t.Elem()
	fmt.Println(t.Name())
	fmt.Println(t.Kind())
	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("ezcli")

		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}

func (app *App) Init(pathToConfigFile, configName string) {
	cobra.OnInitialize(app.initConfig(pathToConfigFile, configName))
}

func (app *App) initConfig(pathToConfigFile, configName string) func() {
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
		for _, fn := range app.postLoadFuncs {
			fn()
		}
	}
}

func (app *App) Execute() error {
	return app.cmd.Execute()
}

func bindFlagAndConfig(cmd *cobra.Command, names ...string) {
	for _, s := range names {
		err := viper.BindPFlag(s, cmd.PersistentFlags().Lookup(s))
		if err != nil {
			// Must bind these flags
			panic(err)
		}
	}
}
