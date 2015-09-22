package cli

import "bitbucket.org/ulfurinn/options"

type Context struct {
	app        *App
	args       []string
	commands   []Command
	options    *options.OptionSet
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
	if strOpt, ok := opt.Value.(*options.StringValue); ok && strOpt != nil {
		v = string(*strOpt)
	}
	return
}

func (c *Context) Bool(name string) (v bool) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if boolOpt, ok := opt.Value.(*options.BoolValue); ok && boolOpt != nil {
		v = bool(*boolOpt)
	}
	return
}

func (c *Context) Int(name string) (v int) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if intOpt, ok := opt.Value.(*options.IntValue); ok && intOpt != nil {
		v = int(*intOpt)
	}
	return
}

func (c *Context) Float64(name string) (v float64) {
	opt := c.options.Lookup(name)
	if opt == nil {
		return
	}
	if floatOpt, ok := opt.Value.(*options.Float64Value); ok && floatOpt != nil {
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

func (c *Context) setupOptions(cs []Command) {
	if c.options == nil {
		c.options = options.NewOptionSet()
	}
	for i, com := range cs {
		for _, arg := range com.Args {
			//	only the direct command may take a positional
			if i == len(cs)-1 {
				arg.ApplyPositional(c.options)
			}
		}
		for _, opt := range com.Options {
			//	local options are not inherited by subcommands
			if i == len(cs)-1 || !opt.local() {
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

func (c *Context) validateOptions(opts []Option) error {
	for _, opt := range opts {
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
