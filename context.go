package cli

import "bitbucket.org/ulfurinn/options"

type Context struct {
	app      *App
	args     []string
	commands []Command
	options  *options.OptionSet
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

func (c *Context) Command() Command { return c.commands[len(c.commands)-1] }

func (c *Context) setupOptions(opts []Option) {
	if c.options == nil {
		c.options = options.NewOptionSet()
	}
	for _, opt := range opts {
		opt.Apply(c.options)
	}
}

func (c *Context) parseOptions() error {
	return c.options.Parse(c.args)
}
