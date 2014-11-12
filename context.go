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
	if strOpt, ok := opt.Value.(*options.BoolValue); ok && strOpt != nil {
		v = bool(*strOpt)
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
