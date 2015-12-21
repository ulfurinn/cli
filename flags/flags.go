package flags

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Value interface {
	String() string
	Set(string) error
	Explicit() bool
}

type Option struct {
	Name     string
	Usage    string
	Value    Value
	Optional bool
	Default  string
}

type Set struct {
	arguments        []*Option
	declared, actual map[string]*Option
	args             []string
	MissingValue     *Option
	Out              io.Writer
}

func NewSet() *Set {
	return &Set{
		Out: os.Stderr,
	}
}

func (s *Set) Argument(v Value, name, usage string, optional bool) {
	opt := &Option{Name: name, Usage: usage, Value: v, Default: v.String(), Optional: optional}
	_, declared := s.declared[name]
	if declared {
		msg := fmt.Errorf("flag redeclared: %s", name)
		fmt.Fprintln(s.out(), msg)
		panic(msg)
	}
	if s.declared == nil {
		s.declared = make(map[string]*Option)
	}
	s.declared[name] = opt
	s.arguments = append(s.arguments, opt)
}

func (s *Set) Var(v Value, name, usage string, optional bool) {
	opt := &Option{Name: name, Usage: usage, Value: v, Default: v.String(), Optional: optional}
	_, declared := s.declared[name]
	if declared {
		msg := fmt.Errorf("flag redeclared: %s", name)
		fmt.Fprintln(s.out(), msg)
		panic(msg)
	}
	if s.declared == nil {
		s.declared = make(map[string]*Option)
	}
	s.declared[name] = opt
}

func isOption(s string) (option bool, name string) {
	if len(s) < 2 {
		return // "-x" is the shortest possible option
	}
	if len(s) == 2 {
		if s[0] != '-' {
			return
		}
		if s[1] == '-' {
			return // "--"
		}
		return true, s[1:]
	}
	if len(s) > 2 {
		if s[0] == '-' {
			s = s[1:]
		} else {
			return
		}
		if s[0] == '-' {
			s = s[1:]
		}
		return true, s
	}
	return
}

func (s *Set) Parse(args []string) (err error) {
	var next string
	s.args = args
	if s.actual == nil {
		s.actual = make(map[string]*Option)
	}
	positional := s.arguments
	for len(s.args) > 0 {
		next = s.args[0]
		if next == "--" {
			s.args = s.args[1:]
			break
		}
		if option, name := isOption(next); option {
			s.args = s.args[1:]
			if name[0] == '-' || name[0] == '=' {
				return fmt.Errorf("bad flag syntax: %s", next)
			}
			var value string
			var opt *Option
			var inverted bool
			split := strings.Split(name, "=")
			if len(split) == 1 {
				opt = s.declared[name]
				if opt == nil {
					if strings.HasPrefix(name, "no-") {
						name = name[3:]
						opt = s.declared[name]
						if opt == nil {
							return fmt.Errorf("unknown argument --no-%s", name)
						}
						inverted = true
					} else {
						return fmt.Errorf("unknown argument --%s", name)
					}
				}
				if len(s.args) > 0 {
					switch opt.Value.(type) {
					case *BoolValue:
						if s.args[0] == "true" || s.args[0] == "false" {
							value, s.args = s.args[0], s.args[1:]
						} else {
							if inverted {
								value = "false"
							} else {
								value = "true"
							}
						}
					default:
						value, s.args = s.args[0], s.args[1:]
					}
				} else {
					switch opt.Value.(type) {
					case *BoolValue:
						if inverted {
							value = "false"
						} else {
							value = "true"
						}
					default:
						s.MissingValue = opt
						return fmt.Errorf("no value provided for argument --%s", name)
					}
				}
			} else if len(split) == 2 {
				name, value = split[0], split[1]
				opt = s.declared[name]
				if opt == nil {
					return fmt.Errorf("unknown argument --%s", name)
				}
			}
			if err = opt.Value.Set(value); err != nil {
				s.MissingValue = opt
				return
			}
			s.actual[name] = opt
		} else {
			if len(positional) > 0 {
				s.args = s.args[1:]
				arg := positional[0]
				if err = arg.Value.Set(next); err != nil {
					return
				}
				positional = positional[1:]
			} else {
				break // not an option and no more positionals
			}
		}
	}
	for _, opt := range positional {
		if !opt.Optional {
			s.MissingValue = positional[0]
			return fmt.Errorf("no value provided for argument %s", s.MissingValue.Name)
		}
	}
	return
}

func (s *Set) Lookup(name string) *Option {
	return s.declared[name]
}

func (s *Set) Arg(n int) string {
	if n < len(s.args) {
		return s.args[n]
	}
	return ""
}

func (s *Set) Args() []string {
	return s.args
}

func (s *Set) String(name string, value string, usage string, t *string, optional bool) *string {
	if t == nil {
		t = new(string)
	}
	s.StringVar(t, name, value, usage, false, optional)
	return t
}

func (s *Set) StringArg(name string, value string, usage string, t *string, optional bool) *string {
	if t == nil {
		t = new(string)
	}
	s.StringVar(t, name, value, usage, true, optional)
	return t
}

func (s *Set) StringVar(target *string, name string, value string, usage string, positional bool, optional bool) {
	if positional {
		s.Argument(newStringValue(target, value), name, usage, optional)
	} else {
		s.Var(newStringValue(target, value), name, usage, optional)
	}
}

func (s *Set) Int(name string, value int, usage string, t *int, optional bool) *int {
	if t == nil {
		t = new(int)
	}
	s.IntVar(t, name, value, usage, false, optional)
	return t
}

func (s *Set) IntArg(name string, value int, usage string, t *int, optional bool) *int {
	if t == nil {
		t = new(int)
	}
	s.IntVar(t, name, value, usage, true, optional)
	return t
}

func (s *Set) IntVar(target *int, name string, value int, usage string, positional bool, optional bool) {
	if positional {
		s.Argument(newIntValue(target, value), name, usage, optional)
	} else {
		s.Var(newIntValue(target, value), name, usage, optional)
	}
}

func (s *Set) Float64(name string, value float64, usage string, t *float64, optional bool) *float64 {
	if t == nil {
		t = new(float64)
	}
	s.Float64Var(t, name, value, usage, false, optional)
	return t
}

func (s *Set) Float64Arg(name string, value float64, usage string, t *float64, optional bool) *float64 {
	if t == nil {
		t = new(float64)
	}
	s.Float64Var(t, name, value, usage, true, optional)
	return t
}

func (s *Set) Float64Var(target *float64, name string, value float64, usage string, positional bool, optional bool) {
	if positional {
		s.Argument(newFloat64Value(target, value), name, usage, optional)
	} else {
		s.Var(newFloat64Value(target, value), name, usage, optional)
	}
}

func (s *Set) Bool(name string, value bool, usage string, t *bool, optional bool) *bool {
	if t == nil {
		t = new(bool)
	}
	s.BoolVar(t, name, value, usage, false, optional)
	return t
}

func (s *Set) BoolArg(name string, value bool, usage string, t *bool, optional bool) *bool {
	if t == nil {
		t = new(bool)
	}
	s.BoolVar(t, name, value, usage, true, optional)
	return t
}

func (s *Set) BoolVar(target *bool, name string, value bool, usage string, positional bool, optional bool) {
	if positional {
		s.Argument(newBoolValue(target, value), name, usage, optional)
	} else {
		s.Var(newBoolValue(target, value), name, usage, optional)
	}
}

func (s *Set) out() io.Writer {
	return s.Out
}

type StringValue string

func newStringValue(target *string, value string) Value {
	*target = value
	return (*StringValue)(target)
}

func (v *StringValue) String() string      { return (string)(*v) }
func (v *StringValue) Set(nv string) error { *v = StringValue(nv); return nil }
func (v *StringValue) Explicit() bool      { return true }

type IntValue int

func newIntValue(target *int, value int) Value {
	*target = value
	return (*IntValue)(target)
}

func (v *IntValue) String() string { return fmt.Sprintf("%v", *v) }
func (v *IntValue) Set(nv string) error {
	c, err := strconv.ParseInt(nv, 0, 64)
	if err == nil {
		*v = IntValue(c)
	}
	return err
}
func (v *IntValue) Explicit() bool { return true }

type Float64Value float64

func newFloat64Value(target *float64, value float64) Value {
	*target = value
	return (*Float64Value)(target)
}

func (v *Float64Value) String() string { return fmt.Sprintf("%v", *v) }
func (v *Float64Value) Set(nv string) error {
	c, err := strconv.ParseFloat(nv, 64)
	if err == nil {
		*v = Float64Value(c)
	}
	return err
}
func (v *Float64Value) Explicit() bool { return true }

type BoolValue bool

func newBoolValue(target *bool, value bool) Value {
	*target = value
	return (*BoolValue)(target)
}

func (v *BoolValue) String() string { return fmt.Sprintf("%v", *v) }
func (v *BoolValue) Set(nv string) (err error) {
	parsed, err := strconv.ParseBool(nv)
	if err == nil {
		*v = BoolValue(parsed)
	}
	return
}
func (v *BoolValue) Explicit() bool { return true }
