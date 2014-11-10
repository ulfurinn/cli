package cli

import "os"

type App struct {
	EnableBashCompletion bool
	HideHelp             bool
	Name                 string
	Usage                string
	Main                 Command
}

func NewApp() *App {
	return &App{
		Name: os.Args[0],
	}
}

func (a *App) Run(arguments []string) error {
	ctx := &Context{
		app:  a,
		args: arguments,
	}
	a.Main.FindCommand(ctx)
	cmd := ctx.Command()
	return cmd.Run(ctx)
}
