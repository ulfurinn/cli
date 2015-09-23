package cli_test

import (
	"io/ioutil"
	"testing"
)
import . "bitbucket.org/ulfurinn/cli/flags"
import . "github.com/smartystreets/goconvey/convey"

func TestOptionSet(t *testing.T) {
	Convey("Given an option set", t, func() {
		set := NewSet()
		set.Out = ioutil.Discard

		Convey("Declaring flags", func() {
			var s string
			set.StringVar(&s, "option", "defvalue", "", false, false)
			Convey("Looking up", func() {
				So(set.Lookup("option"), ShouldNotBeNil)
			})
			Convey("Redeclaring", func() {
				So(func() { set.StringVar(&s, "option", "defvalue", "", false, false) }, ShouldPanic)
			})
		})

		Convey("Returning args", func() {
			So(set.Arg(0), ShouldEqual, "")
		})

		Convey("Parsing", func() {
			Convey("Stop at non-flags", func() {
				var s string
				set.StringVar(&s, "option", "defvalue", "", false, false)
				set.Parse([]string{"--option", "value", "extra"})
				So(set.Arg(0), ShouldEqual, "extra")
			})
			Convey("Stop at --", func() {
				var s string
				set.StringVar(&s, "option", "defvalue", "", false, false)
				set.Parse([]string{"--option", "value", "--", "--option", "othervalue"})
				So(set.Arg(0), ShouldEqual, "--option")
				So(set.Arg(1), ShouldEqual, "othervalue")
				So(s, ShouldEqual, "value")
			})
			Convey("Fail with malformed flag", func() {
				err := set.Parse([]string{"--=value"})
				So(err, ShouldNotBeNil)
			})
			Convey("Parse equal signs", func() {
				var s string
				set.StringVar(&s, "option", "defvalue", "", false, false)
				set.Parse([]string{"--option=value"})
				So(s, ShouldEqual, "value")
			})
			Convey("Undeclared flag", func() {
				var s string
				set.StringVar(&s, "option", "defvalue", "", false, false)
				err := set.Parse([]string{"--option2", "value"})
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Using a string value", func() {
			var s string
			set.StringVar(&s, "option", "defvalue", "", false, false)
			Convey("Using the default value", func() {
				So(s, ShouldEqual, "defvalue")
			})
			Convey("Parsing a value", func() {
				set.Parse([]string{"--option", "value"})
				So(s, ShouldEqual, "value")
			})
		})
		Convey("Using an automatic string value", func() {
			s := set.String("option", "defvalue", "", nil, false)
			set.Parse([]string{"--option", "value"})
			So(*s, ShouldEqual, "value")
		})

		Convey("Using an int value", func() {
			var s int
			set.IntVar(&s, "option", 42, "", false, false)
			Convey("Using the default value", func() {
				So(s, ShouldEqual, 42)
			})
			Convey("Parsing a value", func() {
				set.Parse([]string{"--option", "43"})
				So(s, ShouldEqual, 43)
			})
			Convey("Parsing a wrong value", func() {
				err := set.Parse([]string{"--option", "xxx"})
				So(s, ShouldEqual, 42)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Using an automatic int value", func() {
			s := set.Int("option", 42, "", nil, false)
			set.Parse([]string{"--option", "43"})
			So(*s, ShouldEqual, 43)
		})

		Convey("Using an float64 value", func() {
			var s float64
			set.Float64Var(&s, "option", 42.42, "", false, false)
			Convey("Using the default value", func() {
				So(s, ShouldEqual, 42.42)
			})
			Convey("Parsing a value", func() {
				set.Parse([]string{"--option", "43.43"})
				So(s, ShouldEqual, 43.43)
			})
			Convey("Parsing a wrong value", func() {
				err := set.Parse([]string{"--option", "xxx"})
				So(s, ShouldEqual, 42.42)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("Using an automatic float64 value", func() {
			s := set.Float64("option", 42.42, "", nil, false)
			set.Parse([]string{"--option", "43.43"})
			So(*s, ShouldEqual, 43.43)
		})

		Convey("Using a bool value", func() {
			Convey("Explicit true value", func() {
				s := set.Bool("option", false, "", nil, false)
				So(*s, ShouldBeFalse)
				err := set.Parse([]string{"--option=true"})
				So(err, ShouldBeNil)
				So(*s, ShouldBeTrue)
			})
			Convey("Explicit false value", func() {
				s := set.Bool("option", true, "", nil, false)
				So(*s, ShouldBeTrue)
				err := set.Parse([]string{"--option=false"})
				So(err, ShouldBeNil)
				So(*s, ShouldBeFalse)
			})
			Convey("Positive form", func() {
				s := set.Bool("option", false, "", nil, false)
				So(*s, ShouldBeFalse)
				err := set.Parse([]string{"--option"})
				So(err, ShouldBeNil)
				So(*s, ShouldBeTrue)
			})
			Convey("Negative form", func() {
				s := set.Bool("option", true, "", nil, false)
				So(*s, ShouldBeTrue)
				err := set.Parse([]string{"--no-option"})
				So(err, ShouldBeNil)
				So(*s, ShouldBeFalse)
			})
		})

		Convey("Should record the last flag without a value", func() {
			var s string
			set.StringVar(&s, "option", "defvalue", "", false, false)
			err := set.Parse([]string{"--option"})
			So(err, ShouldNotBeNil)
			So(set.MissingValue, ShouldNotBeNil)
		})

	})
}
