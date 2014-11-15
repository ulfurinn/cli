package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/ulfurinn/options"
)

// This flag enables bash-completion for all commands and subcommands
var ShellCompletionOption = BoolOption{
	Name:   "generate-shell-completion",
	EnvVar: "_CLI_SHELL_COMPLETION",
	Hidden: true,
}

// // This flag prints the version for the application
// var VersionOption = BoolOption{
// 	Name:  "version, v",
// 	Usage: "print the version",
// }

// Option is a common interface related to parsing flags in cli.
// For more advanced flag parsing techniques, it is recomended that
// this interface be implemented.
type Option interface {
	HelpString() string
	CompletionStrings() []string
	// Apply Option settings to the given flag set
	Apply(*options.OptionSet)
	getName() string
	completion() func(*Context) []string
	//visible() bool
}

func eachName(longName string, fn func(string)) {
	parts := strings.Split(longName, ",")
	for _, name := range parts {
		name = strings.Trim(name, " ")
		fn(name)
	}
}

// // Generic is a generic parseable type identified by a specific flag
// type Generic interface {
// 	Set(value string) error
// 	String() string
// }

// // GenericOption is the flag type for types implementing Generic
// type GenericOption struct {
// 	Name   string
// 	Value  Generic
// 	Usage  string
// 	EnvVar string
// 	Hidden bool
// }

// func (f GenericOption) String() string {
// 	return withEnvHint(f.EnvVar, fmt.Sprintf("%s%s %v\t`%v` %s", prefixFor(f.Name), f.Name, f.Value, "-"+f.Name+" option -"+f.Name+" option", f.Usage))
// }

// func (f GenericOption) Apply(set *options.OptionSet) {
// 	val := f.Value
// 	if f.EnvVar != "" {
// 		if envVal := os.Getenv(f.EnvVar); envVal != "" {
// 			val.Set(envVal)
// 		}
// 	}

// 	eachName(f.Name, func(name string) {
// 		set.Var(f.Value, name, f.Usage)
// 	})
// }

// func (f GenericOption) getName() string {
// 	return f.Name
// }

// func (f GenericOption) visible() bool { return !f.Hidden }

// type StringSlice []string

// func (f *StringSlice) Set(value string) error {
// 	*f = append(*f, value)
// 	return nil
// }

// func (f *StringSlice) String() string {
// 	return fmt.Sprintf("%s", *f)
// }

// func (f *StringSlice) Value() []string {
// 	return *f
// }

// type StringSliceOption struct {
// 	Name   string
// 	Value  *StringSlice
// 	Usage  string
// 	EnvVar string
// 	Hidden bool
// }

// func (f StringSliceOption) String() string {
// 	firstName := strings.Trim(strings.Split(f.Name, ",")[0], " ")
// 	pref := prefixFor(firstName)
// 	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), pref+firstName+" option "+pref+firstName+" option", f.Usage))
// }

// func (f StringSliceOption) Apply(set *options.OptionSet) {
// 	if f.EnvVar != "" {
// 		if envVal := os.Getenv(f.EnvVar); envVal != "" {
// 			newVal := &StringSlice{}
// 			for _, s := range strings.Split(envVal, ",") {
// 				newVal.Set(s)
// 			}
// 			f.Value = newVal
// 		}
// 	}

// 	eachName(f.Name, func(name string) {
// 		set.Var(f.Value, name, f.Usage)
// 	})
// }

// func (f StringSliceOption) getName() string {
// 	return f.Name
// }

// func (f StringSliceOption) visible() bool { return !f.Hidden }

// type IntSlice []int

// func (f *IntSlice) Set(value string) error {

// 	tmp, err := strconv.Atoi(value)
// 	if err != nil {
// 		return err
// 	} else {
// 		*f = append(*f, tmp)
// 	}
// 	return nil
// }

// func (f *IntSlice) String() string {
// 	return fmt.Sprintf("%d", *f)
// }

// func (f *IntSlice) Value() []int {
// 	return *f
// }

// type IntSliceOption struct {
// 	Name   string
// 	Value  *IntSlice
// 	Usage  string
// 	EnvVar string
// 	Hidden bool
// }

// func (f IntSliceOption) String() string {
// 	firstName := strings.Trim(strings.Split(f.Name, ",")[0], " ")
// 	pref := prefixFor(firstName)
// 	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), pref+firstName+" option "+pref+firstName+" option", f.Usage))
// }

// func (f IntSliceOption) Apply(set *options.OptionSet) {
// 	if f.EnvVar != "" {
// 		if envVal := os.Getenv(f.EnvVar); envVal != "" {
// 			newVal := &IntSlice{}
// 			for _, s := range strings.Split(envVal, ",") {
// 				err := newVal.Set(s)
// 				if err != nil {
// 					fmt.Fprintf(os.Stderr, err.Error())
// 				}
// 			}
// 			f.Value = newVal
// 		}
// 	}

// 	eachName(f.Name, func(name string) {
// 		set.Var(f.Value, name, f.Usage)
// 	})
// }

// func (f IntSliceOption) getName() string {
// 	return f.Name
// }

// func (f IntSliceOption) visible() bool { return !f.Hidden }

type BoolOption struct {
	Name   string
	Value  bool
	Usage  string
	EnvVar string
	Hidden bool
}

func (f BoolOption) HelpString() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s\t%v", prefixedNames(f.Name), f.Usage))
}

func (f BoolOption) CompletionStrings() []string {
	return []string{prefixedNames(f.Name)}
}

func (f BoolOption) Apply(set *options.OptionSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValBool, err := strconv.ParseBool(envVal)
			if err == nil {
				f.Value = envValBool
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Bool(name, f.Value, f.Usage)
	})
}

func (f BoolOption) getName() string {
	return f.Name
}

func (f BoolOption) visible() bool                       { return !f.Hidden }
func (f BoolOption) completion() func(*Context) []string { return nil }

type StringOption struct {
	Name       string
	Value      string
	Usage      string
	EnvVar     string
	Hidden     bool
	Completion func(*Context) []string
}

func (f StringOption) HelpString() string {
	var fmtString string
	fmtString = "%s %v\t%v"

	if len(f.Value) > 0 {
		fmtString = "%s '%v'\t%v"
	} else {
		fmtString = "%s %v\t%v"
	}

	return withEnvHint(f.EnvVar, fmt.Sprintf(fmtString, prefixedNames(f.Name), f.Value, f.Usage))
}

func (f StringOption) CompletionStrings() []string {
	return []string{prefixedNames(f.Name)}
}

func (f StringOption) Apply(set *options.OptionSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			f.Value = envVal
		}
	}

	eachName(f.Name, func(name string) {
		set.String(name, f.Value, f.Usage)
	})
}

func (f StringOption) getName() string {
	return f.Name
}

func (f StringOption) visible() bool                       { return !f.Hidden }
func (f StringOption) completion() func(*Context) []string { return f.Completion }

type IntOption struct {
	Name       string
	Value      int
	Usage      string
	EnvVar     string
	Completion func(*Context) []string
}

func (f IntOption) HelpString() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage))
}

func (f IntOption) CompletionStrings() []string {
	return []string{prefixedNames(f.Name)}
}

func (f IntOption) Apply(set *options.OptionSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValInt, err := strconv.ParseUint(envVal, 10, 64)
			if err == nil {
				f.Value = int(envValInt)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Int(name, f.Value, f.Usage)
	})
}

func (f IntOption) getName() string {
	return f.Name
}
func (f IntOption) completion() func(*Context) []string { return f.Completion }

// type DurationOption struct {
// 	Name   string
// 	Value  time.Duration
// 	Usage  string
// 	EnvVar string
// }

// func (f DurationOption) String() string {
// 	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage))
// }

// func (f DurationOption) Apply(set *options.OptionSet) {
// 	if f.EnvVar != "" {
// 		if envVal := os.Getenv(f.EnvVar); envVal != "" {
// 			envValDuration, err := time.ParseDuration(envVal)
// 			if err == nil {
// 				f.Value = envValDuration
// 			}
// 		}
// 	}

// 	eachName(f.Name, func(name string) {
// 		set.Duration(name, f.Value, f.Usage)
// 	})
// }

// func (f DurationOption) getName() string {
// 	return f.Name
// }

type Float64Option struct {
	Name       string
	Value      float64
	Usage      string
	EnvVar     string
	Completion func(*Context) []string
}

func (f Float64Option) HelpString() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage))
}

func (f Float64Option) CompletionStrings() []string {
	return []string{prefixedNames(f.Name)}
}

func (f Float64Option) Apply(set *options.OptionSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValFloat, err := strconv.ParseFloat(envVal, 10)
			if err == nil {
				f.Value = float64(envValFloat)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Float64(name, f.Value, f.Usage)
	})
}

func (f Float64Option) getName() string {
	return f.Name
}
func (f Float64Option) completion() func(*Context) []string { return f.Completion }

func prefixFor(name string) (prefix string) {
	if len(name) == 1 {
		prefix = "-"
	} else {
		prefix = "--"
	}

	return
}

func prefixedNames(fullName string) (prefixed string) {
	parts := strings.Split(fullName, ",")
	for i, name := range parts {
		name = strings.Trim(name, " ")
		prefixed += prefixFor(name) + name
		if i < len(parts)-1 {
			prefixed += ", "
		}
	}
	return
}

func withEnvHint(envVar, str string) string {
	envText := ""
	if envVar != "" {
		envText = fmt.Sprintf(" [$%s]", envVar)
	}
	return str + envText
}
