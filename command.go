package cli

import (
	"fmt"
	"io"
)

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
	err = ctx.parseOptions()
	completion := ctx.Bool("generate-shell-completion")

	if completion {
		if err == nil || ctx.options.MissingValue != nil {
			c.showCompletion(ctx)
		}
		return
	}

	if err == nil && c.Action != nil {
		c.Action(ctx)
	}

	return

}

func (c *Command) showCompletion(ctx *Context) {
	if missing := ctx.options.MissingValue; missing != nil {
		opt := ctx.findOption(missing.Name)
		if f := opt.completion(); f != nil {
			showCompletion(ctx.app.Out, f(ctx))
			return
		}
	}
	if ctx.parseError != nil {
		return
	}
	list := []string{}
	for _, cmd := range c.Commands {
		list = append(list, cmd.Name, cmd.ShortName)
	}
	for _, opt := range c.Options {
		list = append(list, opt.CompletionStrings()...)
	}
	showCompletion(ctx.app.Out, list)
}

func showCompletion(out io.Writer, strings []string) {
	for _, str := range strings {
		if str != "" {
			fmt.Fprintln(out, str)
		}
	}
}
