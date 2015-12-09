package cli

import (
	"fmt"
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

func (a *App) RunMain() {
	if err := a.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if exit, ok := err.(Exit); ok {
			os.Exit(exit.StatusCode)
		} else {
			os.Exit(1)
		}
	}
	os.Exit(0)
}
