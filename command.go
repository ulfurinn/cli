package cli

import (
	"fmt"
	"io"
	"strings"
)

type Command struct {
	Commands   []Command
	Options    []Option
	Args       []Option
	Name       string
	ShortName  string
	Usage      string
	Before     func(*Context) error
	Action     func(*Context) error
	Completion func(*Context)
}

func (c *Command) HasName(name string) bool {
	return c.Name == name || c.ShortName == name
}

func (c *Command) FindCommandByName(name string) (cmd *Command) {
	for i := range c.Commands {
		if c.Commands[i].HasName(name) {
			cmd = &c.Commands[i]
			break
		}
	}
	return
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

func (c *Command) appendHelp() {
	for _, com := range c.Commands {
		if com.HasName("help") {
			return
		}
	}
	c.Commands = append(c.Commands, HelpCommand)
}

// Invokes the command given the context, parses ctx.Args() to generate command-specific flags
func (c *Command) Run(ctx *Context) (err error) {
	ctx.setupOptions(ctx.commands)
	err = ctx.parseOptions()
	// we can't check the result now because with shell completion errors can be a legitimate case

	completion := ctx.Bool("generate-shell-completion")
	help := ctx.Bool("help")

	if completion {
		if err == nil || ctx.options.MissingValue != nil {
			c.showCompletion(ctx)
			err = nil
		}
		return
	}

	if help {
		err = helpOptionAction(ctx)
		return
	}

	//	now we can check the result from parseOptions
	if err != nil {
		return
	}

	err = ctx.validateOptions(c.Options)
	if err != nil {
		return err
	}

	for _, cmd := range ctx.commands {
		if cmd.Before != nil {
			if err := cmd.Before(ctx); err != nil {
				return err
			}
		}
	}

	if err == nil && c.Action != nil {
		err = c.Action(ctx)
	}

	return

}

func (c *Command) showCompletion(ctx *Context) {
	if c.Completion != nil {
		c.Completion(ctx)
		return
	}
	if missing := ctx.options.MissingValue; missing != nil {
		opt := ctx.findOption(missing.Name)
		if f := opt.completion(); f != nil {
			showCompletion(ctx.app.Out, f(ctx, opt))
			return
		}
	}
	if ctx.parseError != nil {
		return
	}
	list := []string{}
	for _, cmd := range c.Commands {
		if cmd.Name != "help" {
			list = append(list, cmd.Name, cmd.ShortName)
		}
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

func StdCompletionFlags() string {
	return "$stdcomp=-fd"
}

func StdCompletion(*Context, Option) []string {
	return []string{StdCompletionFlags()}
}

func ValueListCompletion(ctx *Context, opt Option) []string {
	switch o := opt.(type) {
	case StringOption:
		return o.ValueList
	default:
		return []string{}
	}
}

func ValueListValidation(ctx *Context, opt Option) error {
	// default values always pass
	if lowOpt := ctx.options.Lookup(opt.getName()); lowOpt != nil && !lowOpt.Value.Explicit() {
		return nil
	}
	switch o := opt.(type) {
	case StringOption:
		givenValue := ctx.String(o.Name)
		for _, allowedValue := range o.ValueList {
			if givenValue == allowedValue {
				return nil
			}
		}
		return fmt.Errorf("%s accepts one of the following values: %s", o.Name, strings.Join(o.ValueList, ","))
	default:
		return nil
	}
}
