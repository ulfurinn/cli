package cli

type Command struct {
	Commands  []Command
	Options   []Option
	Name      string
	ShortName string
	Usage     string
	Before    func(*Context) error
	Action    func(*Context)
}

func (c *Command) HasName(name string) bool {
	return c.Name == name || c.ShortName == name
}

func (c *Command) FindCommand(ctx *Context) {
	ctx.commands = append(ctx.commands, *c)
	if len(ctx.args) == 0 {
		return
	}
	a := ctx.args[0]
	for i := range c.Commands {
		if c.Commands[i].HasName(a) {
			ctx.args = ctx.args[1:]
			c.Commands[i].FindCommand(ctx)
		}
	}
}

// Invokes the command given the context, parses ctx.Args() to generate command-specific flags
func (c *Command) Run(ctx *Context) (err error) {
	ctx.setupOptions(c.Options)
	if err = ctx.parseOptions(); err != nil {
		return
	}
	if c.Action != nil {
		c.Action(ctx)
	}
	return
}
