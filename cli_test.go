package cli_test

import (
	"testing"
	. "bitbucket.org/ulfurinn/cli"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApp(t *testing.T) {
	Convey("Given an app", t, func() {
		app := NewApp()
		Convey("Can run a simple command", func() {
			run := false
			app.Main = Command{
				Commands: []Command{{
					Name:   "testcmd",
					Action: func(c *Context) { run = true },
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
						Action: func(c *Context) { run = true },
					}},
				}},
			}
			app.Run([]string{"testcmd", "sub"})
			So(run, ShouldBeTrue)
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
					Action: func(c *Context) {
						So(c.Int("int"), ShouldEqual, 42)
						So(c.Int("nonesuch"), ShouldEqual, 0)

						So(c.String("str"), ShouldEqual, "42")
						So(c.String("nonesuch"), ShouldEqual, "")

						So(c.Bool("bool"), ShouldEqual, true)
						So(c.Bool("nonesuch"), ShouldEqual, false)

						So(c.Float64("float"), ShouldEqual, 42.42)
						So(c.Float64("nonesuch"), ShouldEqual, 0.0)

					},
				}},
			}
			err := app.Run([]string{"cmd", "--int", "42", "--str", "42", "--bool", "--float", "42.42"})
			So(err, ShouldBeNil)
		})
	})
}
