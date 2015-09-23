package cli

import "bitbucket.org/ulfurinn/cli/flags"

type Context struct {
	app        *App
	args       []string
	commands   []Command
	options    *flags.OptionSet
	parseError error
}

func (c *Context) Arg(i int) string {
	if i >= len(c.args) {
		return ""
	}
	return c.args[i]
}

func (c *Context) ArgLen() int {
	return len(c.args)
}

func (c *Context) String(name string) (v string) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if strOpt, ok := opt.Value.(*flags.StringValue); ok && strOpt != nil {
		v = string(*strOpt)
	}
	return
}

func (c *Context) Bool(name string) (v bool) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if boolOpt, ok := opt.Value.(*flags.BoolValue); ok && boolOpt != nil {
		v = bool(*boolOpt)
	}
	return
}

func (c *Context) Int(name string) (v int) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if intOpt, ok := opt.Value.(*flags.IntValue); ok && intOpt != nil {
		v = int(*intOpt)
	}
	return
}

func (c *Context) Float64(name string) (v float64) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if floatOpt, ok := opt.Value.(*flags.Float64Value); ok && floatOpt != nil {
		v = float64(*floatOpt)
	}
	return
}

func (c *Context) StringSlice(name string) (v []string) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if sliceOpt, ok := opt.Value.(*StringSlice); ok && sliceOpt != nil {
		v = *sliceOpt
	}
	return
}

func (c *Context) Command() *Command { return &c.commands[len(c.commands)-1] }

func (c *Context) run() (err error) {
	c.setupOptions()
	err = c.parseOptions()
	// we can't check the result now because with shell completion errors can be a legitimate case

	completion := c.Bool("generate-shell-completion")
	help := c.Bool("help")

	if completion {
		if err == nil || c.options.MissingValue != nil {
			c.Command().showCompletion(c)
			err = nil
		}
		return
	}

	if help {
		err = helpOptionAction(c)
		return
	}

	//	now we can check the result from parseOptions
	if err != nil {
		return
	}

	err = c.validateOptions()
	if err != nil {
		return err
	}

	for _, cmd := range c.commands {
		if cmd.Before != nil {
			if err := cmd.Before(c); err != nil {
				return err
			}
		}
	}

	if err == nil && c.Command().Action != nil {
		err = c.Command().Action(c)
	}

	return

}

func (c *Context) setupOptions() {
	if c.options == nil {
		c.options = flags.NewOptionSet()
	}
	for i, com := range c.commands {
		for _, arg := range com.Args {
			//	only the direct command may take a positional
			if i == len(c.commands)-1 {
				arg.ApplyPositional(c.options)
			}
		}
		for _, opt := range com.Options {
			//	local options are not inherited by subcommands
			if i == len(c.commands)-1 || !opt.local() {
				opt.Apply(c.options)
			}
		}
	}
	HelpOption.Apply(c.options)
	if c.app.EnableShellCompletion {
		ShellCompletionOption.Apply(c.options)
	}
}

func (c *Context) parseOptions() (err error) {
	err = c.options.Parse(c.args)
	c.parseError = err
	return
}

func (c *Context) validateOptions() error {
	for _, opt := range c.Command().Options {
		if opt.validation() != nil {
			if err := opt.validation()(c, opt); err != nil {
				return nil
			}
		}
	}
	return nil
}

func (c *Context) findOption(name string) (option Option) {
	for _, cmd := range c.commands {
		for _, opt := range cmd.Args {
			if opt.getName() == name {
				option = opt
			}
		}
		for _, opt := range cmd.Options {
			if opt.getName() == name {
				option = opt
			}
		}
	}
	return
}
