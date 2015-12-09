package cli

import (
	"fmt"
	"sort"
	"strings"

	"text/template"
)

var tplSource = `
Usage: {{.AppName}}{{.CommandList}}{{range .Args}} <{{.Name}}>{{end}}{{if .Usage}}

{{.Usage}}{{end}}{{if .Subcommands}}

Subcommands:{{range .Subcommands}}
  {{.Name}}{{if .Usage}}    {{.Usage}}{{end}}{{end}}{{end}}{{if .Args}}

Arguments:{{range .Args}}
  {{.Name}}{{if .Usage}}    {{.Usage}}{{end}}{{end}}{{end}}{{if .Options}}

Options:{{range .Options}}
  {{.Name}}{{if .Usage}}    {{.Usage}}{{end}}{{end}}{{end}}
`

type helpOption struct {
	Name  string
	Usage string
}

type helpContext struct {
	AppName     string
	CommandList string
	Usage       string
	Subcommands []struct {
		Name  string
		Usage string
	}
	Args    []helpOption
	Options []helpOption
}

// This flag prints the help for all commands and subcommands
var HelpOption = BoolOption{
	Name:  "help, h",
	Usage: "show help",
}

var HelpCommand Command

func init() {
	HelpCommand = Command{
		Name:       "help",
		Action:     helpCommandAction,
		Completion: helpCompletion,
	}
}

func helpCommandAction(ctx *Context) error {
	tpl, _ := template.New("help").Parse(tplSource)
	helpCtx := helpContext{}
	helpCtx.setupCommand(ctx)
	return tpl.Execute(ctx.app.Out, helpCtx)
}

func helpOptionAction(ctx *Context) error {
	tpl, _ := template.New("help").Parse(tplSource)
	helpCtx := helpContext{}
	helpCtx.setup(ctx)
	return tpl.Execute(ctx.app.Out, helpCtx)
}

func (h *helpContext) setupCommand(ctx *Context) {
	subctx := &Context{
		app:  ctx.app,
		args: ctx.args,
	}
	ctx.app.Main.FindCommand(subctx)
	h.setup(subctx)
}

func (h *helpContext) setup(ctx *Context) {
	h.AppName = ctx.app.Name
	usedCommands := []Command{}
	for _, cmd := range ctx.commands {
		usedCommands = append(usedCommands, cmd)
	}
	cmdPath := []string{}
	for i, cmd := range usedCommands {
		if i == 0 {
			continue
		}
		cmdPath = append(cmdPath, cmd.Name)
	}
	if len(cmdPath) > 0 {
		h.CommandList = " " + strings.Join(cmdPath, " ")
	}
	activeCommand := usedCommands[len(usedCommands)-1]
	h.Usage = activeCommand.Usage
	maxSubLength := 0
	maxOptLength := 0
	for _, cmd := range activeCommand.Commands {
		h.Subcommands = append(h.Subcommands, struct {
			Name  string
			Usage string
		}{cmd.Name, cmd.Usage})
		if len(cmd.Name) > maxSubLength {
			maxSubLength = len(cmd.Name)
		}
	}
	opts := map[string]helpOption{}
	for i, cmd := range usedCommands {
		for _, opt := range cmd.Options {
			if !opt.local() || i == len(usedCommands)-1 {
				opts[opt.name()] = helpOption{"--" + opt.name(), opt.usage()}
			}
		}
	}
	optKeys := []string{}
	for k := range opts {
		optKeys = append(optKeys, k)
	}
	sort.Strings(optKeys)
	for _, key := range optKeys {
		opt := opts[key]
		h.Options = append(h.Options, opt)
		if len(opt.Name) > maxOptLength {
			maxOptLength = len(opt.Name)
		}
	}
	for _, arg := range activeCommand.Args {
		h.Args = append(h.Args, helpOption{Name: arg.name(), Usage: arg.usage()})
	}
	for k, cmd := range h.Subcommands {
		cmd.Name = fmt.Sprintf(fmt.Sprintf("%%-%ds", maxSubLength), cmd.Name)
		h.Subcommands[k] = cmd
	}
	for k, opt := range h.Options {
		opt.Name = fmt.Sprintf(fmt.Sprintf("%%-%ds", maxOptLength), opt.Name)
		h.Options[k] = opt
	}
}

func helpCompletion(ctx *Context) {
	subctx := &Context{
		app:  ctx.app,
		args: ctx.args,
	}
	subctx.app.Main.FindCommand(subctx)
	cmd := subctx.Command()
	subctx.setupOptions()
	cmd.showCompletion(subctx)
}
