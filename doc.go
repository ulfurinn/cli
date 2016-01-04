/*

DOC DRAFT

Package cli provides a toolset for writing command line interfaces.

This package started off as a fork of package github.com/codegangsta/cli with an aim to enhance shell completion but has diverged since then. Application definition structure is similar but not identical.

The following sections briefly describe the main components; see the API index for information on further customization.

Application basics

A minimally viable application looks like this:

	func main() {
		a := cli.NewApp()
		a.Main.Action = func(*cli.Context) error {
			fmt.Println("hello, world")
			return nil
		}
		a.RunMain()
	}

Named options

Options can be used as follows:

	...
	a.Main.Options = []cli.Option{
		cli.StringOption{
			Name:  "flag",
			Value: "default value",
		},
	}
	a.Main.Action = func(ctx *cli.Context) error {
		fmt.Println(ctx.String("flag"))
		return nil
	}
	...

	$ app --flag value
	> value

	$ app
	> default value

Positional arguments

Any command line argument that cannot be identified and parsed as a named option will be available in Args(), but a formal declaration provides type-specific parsing and better help messages.

Named options and positional arguments are declared and accessed through the same interface.

	...
	a.Main.Args = []cli.Option{
		cli.StringOption{
			Name:  "flag",
			Value: "default value",
		},
	}
	a.Main.Action = func(ctx *cli.Context) error {
		fmt.Println(ctx.String("flag"))
		return nil
	}
	...

	$ app value
	> value

Subcommands

Subcommands are created as follows:

	...
	a.Main.Commands = []cli.Command{
		Name:   "cmd",
		Action: func(*ctx.Context) error {
			fmt.Println("subcommand")
			return nil
		},
	}
	...

	$ app cmd
	> subcommand

Like the root command Main, subcommands can have their own options and subcommands.

Help

The root command has an implicit "help" subcommand, showing usage instructions. For help on subcommands, it is invoked as "app help subcmd1 subcmd2 ...".

Alternatively, every command has an implicit "--help" option that has the same effect.

The implicit "help-commands" subcommand prints a recursive list of all declared subcommands.

Shell completion

All subcommand and options are available for shell completion. Additionally, they can declare custom completion functions, returning a list of accepted values.

The bash completion function is available at https://bitbucket.org/ulfurinn/cli/raw/default/bash_completion; replace $PROG with the name of your executable.

*/
package cli
