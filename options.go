package ezcli

import "github.com/spf13/cobra"

type CommandOptions func(app *App) error

func CheckArgs(args cobra.PositionalArgs) func(app *App) error {
	return func(app *App) error {
		app.Cmd.Args = args
		return nil
	}
}
