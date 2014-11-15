package cli

import (
	"fmt"
	"strings"

	"text/template"
)

var tplSource = `
Usage: {{.AppName}}{{.CommandList}}{{if .Usage}}

{{.Usage}}{{end}}{{if .Subcommands}}

Subcommands:{{range .Subcommands}}
  {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}{{end}}{{end}}
`

type helpContext struct {
	AppName     string
	CommandList string
	Usage       string
	Subcommands []struct {
		Name  string
		Usage string
	}
}

// This flag prints the help for all commands and subcommands
var HelpOption = BoolOption{
	Name:  "help, h",
	Usage: "show help",
}

var HelpCommand = Command{
	Name:   "help",
	Action: helpAction,
}

func helpAction(ctx *Context) {
	tpl, _ := template.New("help").Parse(tplSource)
	helpCtx := helpContext{}
	helpCtx.setup(ctx)
	tpl.Execute(ctx.app.Out, helpCtx)
}

func (h *helpContext) setup(ctx *Context) {
	h.AppName = ctx.app.Name
	usedCommands := []Command{}
	for _, cmd := range ctx.commands {
		if cmd.Name != "help" {
			usedCommands = append(usedCommands, cmd)
		}
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
	maxLength := 0
	for _, cmd := range activeCommand.Commands {
		h.Subcommands = append(h.Subcommands, struct {
			Name  string
			Usage string
		}{cmd.Name, cmd.Usage})
		if len(cmd.Name) > maxLength {
			maxLength = len(cmd.Name)
		}
	}
	for _, cmd := range h.Subcommands {
		cmd.Name = fmt.Sprintf(fmt.Sprintf("%%%d%%s", maxLength), cmd.Name)
	}
}
