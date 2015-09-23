package cli

import (
	"io"
	"os"
)

type App struct {
	EnableShellCompletion bool
	HideHelp              bool
	Name                  string
	Usage                 string
	Main                  Command
	Out                   io.Writer
}

func NewApp() *App {
	return &App{
		Name: os.Args[0],
		Out:  os.Stdout,
	}
}

func (a *App) Run(arguments []string) error {
	a.Main.appendHelp()
	ctx := &Context{
		app:  a,
		args: arguments,
	}
	a.Main.FindCommand(ctx)
	return ctx.run()
}
