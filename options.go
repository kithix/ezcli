package ezcli

type varOptFn func(*VarOpts)

type VarOpts struct {
	Name         string
	DefaultValue any
	Usage        string
	Persistent   bool   // Option will persist to sub-commands
	Env          string // If not "" - will
}

func defaultOpts() *VarOpts {
	return &VarOpts{
		Persistent: true,
	}
}

func WithOptions(newOpts *VarOpts) varOptFn {
	return func(opts *VarOpts) {
		// Overwrite the value of incoming options
		*opts = *newOpts
	}
}

func VarName(name string) varOptFn {
	return func(opts *VarOpts) {
		opts.Name = name
	}
}

// VarUsage sets the help information for a flag
func VarUsage(usage string) varOptFn {
	return func(opts *VarOpts) {
		opts.Usage = usage
	}
}

// VarDefaultValue provides a custom default value for an unset option
func VarDefaultValue(val any) varOptFn {
	return func(opts *VarOpts) {
		opts.DefaultValue = val
	}
}

// VarLocal makes the option only set for the current command
func VarLocal() varOptFn {
	return func(opts *VarOpts) {
		opts.Persistent = false
	}
}

func VarEnv(name string) varOptFn {
	return func(opts *VarOpts) {
		opts.Env = name
	}
}
