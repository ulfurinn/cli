package cli_test

import (
	"bytes"
	"os"
	"testing"

	. "bitbucket.org/ulfurinn/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApp(t *testing.T) {
	Convey("Given an app", t, func() {
		app := NewApp()
		app.Name = "testapp"
		app.EnableShellCompletion = true
		Convey("Can run a simple command", func() {
			run := false
			app.Main = Command{
				Commands: []Command{{
					Name:   "testcmd",
					Action: func(c *Context) error { run = true; return nil },
				}},
			}
			app.Run([]string{"testcmd"})
			So(run, ShouldBeTrue)
		})
		Convey("Can run a subcommand", func() {
			run := false
			app.Main = Command{
				Commands: []Command{{
					Name: "testcmd",
					Commands: []Command{{
						Name:   "sub",
						Action: func(c *Context) error { run = true; return nil },
					}},
				}},
			}
			app.Run([]string{"testcmd", "sub"})
			So(run, ShouldBeTrue)
		})
		Convey("Option types", func() {
			Convey("String slice", func() {
				var o []string
				app.Main.Name = "_main"
				app.Main.Options = []Option{
					StringSliceOption{Name: "o"},
				}
				app.Main.Action = func(ctx *Context) error {
					o = ctx.StringSlice("o")
					return nil
				}
				err := app.Run([]string{"--o", "1", "--o", "2", "--o", "3"})
				So(err, ShouldBeNil)
				So(o, ShouldResemble, []string{"1", "2", "3"})
			})

		})
		Convey("Given flags", func() {
			app.Main = Command{
				Commands: []Command{{
					Name: "cmd",
					Options: []Option{
						IntOption{Name: "int"},
						StringOption{Name: "str"},
						BoolOption{Name: "bool"},
						Float64Option{Name: "float"},
					},
					Action: func(c *Context) error {
						So(c.Int("int"), ShouldEqual, 42)
						So(c.Int("nonesuch"), ShouldEqual, 0)

						So(c.String("str"), ShouldEqual, "42")
						So(c.String("nonesuch"), ShouldEqual, "")

						So(c.Bool("bool"), ShouldEqual, true)
						So(c.Bool("nonesuch"), ShouldEqual, false)

						So(c.Float64("float"), ShouldEqual, 42.42)
						So(c.Float64("nonesuch"), ShouldEqual, 0.0)

						return nil
					},
				}},
			}
			err := app.Run([]string{"cmd", "--int", "42", "--str", "42", "--bool", "--float", "42.42"})
			So(err, ShouldBeNil)
			Convey("Shell completion", func() {
				os.Setenv("_CLI_SHELL_COMPLETION", "true")
				var b bytes.Buffer
				app.Out = &b
				Convey("Empty", func() {
					app.Main = Command{}
					app.Run([]string{})
					So(b.String(), ShouldEqual, "")
				})
				Convey("Commands and flags", func() {
					app.Main = Command{
						Commands: []Command{{
							Name: "cmd1",
							Commands: []Command{{
								Name: "sub11",
							}},
						}, {
							Name: "cmd2",
						}},
						Options: []Option{
							IntOption{Name: "int"},
							StringOption{
								Name:       "string",
								Completion: func(*Context, Option) []string { return []string{"a", "b"} },
							},
						},
					}
					Convey("Top level", func() {
						app.Run([]string{})
						So(b.String(), ShouldEqual, "cmd1\ncmd2\n--int\n--string\n")
					})
					Convey("Second level", func() {
						app.Run([]string{"cmd1"})
						So(b.String(), ShouldEqual, "sub11\n")
					})
					Convey("Flag completion", func() {
						err := app.Run([]string{"--string"})
						So(err, ShouldBeNil)
						So(b.String(), ShouldEqual, "a\nb\n")
					})
				})
				Convey("Help completion", func() {
					app.Main = Command{
						Commands: []Command{{
							Name:     "cmd1",
							Commands: []Command{{Name: "sub1"}},
						}, {
							Name: "cmd2",
						}},
					}
					Convey("Top level", func() {
						app.Run([]string{"help"})
						So(b.String(), ShouldEqual, "cmd1\ncmd2\n")
					})
					Convey("Second level", func() {
						app.Run([]string{"help", "cmd1"})
						So(b.String(), ShouldEqual, "sub1\n")
					})
				})
				os.Setenv("_CLI_SHELL_COMPLETION", "false")
			})
		})
		Convey("Help", func() {
			var b bytes.Buffer
			app.Out = &b
			app.Main = Command{
				Usage: "app usage",
				Commands: []Command{{
					Name:  "cmd1",
					Usage: "cmd1 usage",
					Commands: []Command{{
						Name: "sub1",
					}},
				}, {
					Name: "cmd2",
				}},
				Options: []Option{
					IntOption{Name: "int"},
					StringOption{
						Name:       "string",
						Completion: func(*Context, Option) []string { return []string{"a", "b"} },
					},
				},
			}
			Convey("Using a subcommand", func() {
				Convey("Root", func() {
					app.Run([]string{"help"})
					So(b.String(), ShouldEqual, "\nUsage: testapp\n\napp usage\n\nSubcommands:\n  cmd1             cmd1 usage\n  cmd2         \n  help         \n  help-commands\n\nOptions:\n  --int       default = 0\n  --string    default = \"\"\n")
				})
				Convey("Subcommand", func() {
					app.Run([]string{"help", "cmd1"})
					So(b.String(), ShouldEqual, "\nUsage: testapp cmd1\n\ncmd1 usage\n\nSubcommands:\n  sub1\n\nOptions:\n  --int       default = 0\n  --string    default = \"\"\n")
				})
			})
			Convey("Using an option", func() {
				Convey("Root", func() {
					app.Run([]string{"--help"})
					So(b.String(), ShouldEqual, "\nUsage: testapp\n\napp usage\n\nSubcommands:\n  cmd1             cmd1 usage\n  cmd2         \n  help         \n  help-commands\n\nOptions:\n  --int       default = 0\n  --string    default = \"\"\n")
				})
				Convey("Subcommand", func() {
					app.Run([]string{"cmd1", "--help"})
					So(b.String(), ShouldEqual, "\nUsage: testapp cmd1\n\ncmd1 usage\n\nSubcommands:\n  sub1\n\nOptions:\n  --int       default = 0\n  --string    default = \"\"\n")
				})
			})
		})
	})
}
