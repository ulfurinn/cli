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
			intValue := 0
			strValue := ""
			app.Main = Command{
				Commands: []Command{{
					Name: "cmd",
					Options: []Option{
						IntOption{Name: "int"},
						StringOption{Name: "str"},
					},
					Action: func(c *Context) {
						intValue = c.Int("int")
						strValue = c.String("str")
					},
				}},
			}
			app.Run([]string{"cmd", "--int", "42", "--str", "42"})
			So(intValue, ShouldEqual, 42)
			So(strValue, ShouldEqual, "42")
		})
	})
}
