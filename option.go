package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"bitbucket.org/ulfurinn/cli/flags"
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

type completionFunc func(*Context, Option) []string
type validationFunc func(*Context, Option) error

// Option is a common interface related to parsing flags in cli.
// For more advanced flag parsing techniques, it is recomended that
// this interface be implemented.
type Option interface {
	HelpString() string
	CompletionStrings() []string
	// Apply Option settings to the given flag set
	ApplyNamed(*flags.Set)
	ApplyPositional(*flags.Set)
	local() bool
	name() string
	usage() string
	completion() completionFunc
	validation() validationFunc
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

type StringSlice []string

func (f *StringSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *StringSlice) String() string {
	if f == nil {
		return ""
	}
	return fmt.Sprintf("%s", *f)
}

func (f *StringSlice) Value() []string {
	if f == nil {
		return []string{}
	}
	return *f
}

func (f *StringSlice) Explicit() bool { return true }

type StringSliceOption struct {
	Name       string
	Value      *StringSlice
	Usage      string
	EnvVar     string
	Hidden     bool
	Optional   bool
	Local      bool
	Completion completionFunc
	Validation validationFunc
}

func (f StringSliceOption) String() string {
	firstName := strings.Trim(strings.Split(f.Name, ",")[0], " ")
	pref := prefixFor(firstName)
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), pref+firstName+" option "+pref+firstName+" option", f.Usage))
}

func (f StringSliceOption) HelpString() string {
	var fmtString string
	fmtString = "%s %v\t%v"

	//if len(f.Value) > 0 {
	//	fmtString = "%s '%v'\t%v"
	//} else {
	fmtString = "%s %v\t%v"
	//}

	return withEnvHint(f.EnvVar, fmt.Sprintf(fmtString, prefixedNames(f.Name), f.Value, f.Usage))
}

func (f StringSliceOption) CompletionStrings() []string {
	return []string{prefixedNames(f.Name)}
}

func (f StringSliceOption) ApplyNamed(set *flags.Set) {
	f.Value = new(StringSlice)
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			f.Value.Set(envVal)
		}
	}

	eachName(f.Name, func(name string) {
		set.Var(f.Value, name, f.Usage, f.Optional)
	})
}

func (f StringSliceOption) ApplyPositional(set *flags.Set) {
	f.Value = new(StringSlice)
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			f.Value.Set(envVal)
		}
	}

	eachName(f.Name, func(name string) {
		set.Argument(f.Value, name, f.Usage, f.Optional)
	})
}

func (f StringSliceOption) name() string {
	return f.Name
}

func (f StringSliceOption) usage() string {
	if f.Usage == "" {
		return fmt.Sprintf("default = %q", f.Value)
	} else {
		return fmt.Sprintf("%s; default = %q", f.Usage, f.Value)
	}
}

func (f StringSliceOption) visible() bool              { return !f.Hidden }
func (f StringSliceOption) local() bool                { return f.Local }
func (f StringSliceOption) completion() completionFunc { return f.Completion }
func (f StringSliceOption) validation() validationFunc { return f.Validation }

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
	Name     string
	Value    bool
	Usage    string
	EnvVar   string
	Hidden   bool
	Var      *bool
	Optional bool
	Local    bool
}

func (f BoolOption) HelpString() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s\t%v", prefixedNames(f.Name), f.Usage))
}

func (f BoolOption) CompletionStrings() []string {
	return []string{prefixedNames(f.Name), "--no-" + f.Name}
}

func (f BoolOption) ApplyNamed(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValBool, err := strconv.ParseBool(envVal)
			if err == nil {
				f.Value = envValBool
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Bool(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f BoolOption) ApplyPositional(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValBool, err := strconv.ParseBool(envVal)
			if err == nil {
				f.Value = envValBool
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.BoolArg(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f BoolOption) name() string {
	return f.Name
}

func (f BoolOption) usage() string {
	if f.Usage == "" {
		return fmt.Sprintf("default = %v", f.Value)
	} else {
		return fmt.Sprintf("%s; default = %v", f.Usage, f.Value)
	}
}

func (f BoolOption) visible() bool              { return !f.Hidden }
func (f BoolOption) local() bool                { return f.Local }
func (f BoolOption) completion() completionFunc { return nil }
func (f BoolOption) validation() validationFunc { return nil }

type StringOption struct {
	Name       string
	Value      string
	ValueList  []string
	Usage      string
	EnvVar     string
	Hidden     bool
	Var        *string
	Optional   bool
	Local      bool
	Completion completionFunc
	Validation validationFunc
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

func (f StringOption) ApplyNamed(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			f.Value = envVal
		}
	}

	eachName(f.Name, func(name string) {
		set.String(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f StringOption) ApplyPositional(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			f.Value = envVal
		}
	}

	eachName(f.Name, func(name string) {
		set.StringArg(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f StringOption) name() string {
	return f.Name
}

func (f StringOption) usage() string {
	if f.Usage == "" {
		return fmt.Sprintf("default = %q", f.Value)
	} else {
		return fmt.Sprintf("%s; default = %q", f.Usage, f.Value)
	}
}

func (f StringOption) visible() bool              { return !f.Hidden }
func (f StringOption) local() bool                { return f.Local }
func (f StringOption) completion() completionFunc { return f.Completion }
func (f StringOption) validation() validationFunc { return f.Validation }

type IntOption struct {
	Name       string
	Value      int
	Usage      string
	EnvVar     string
	Var        *int
	Optional   bool
	Local      bool
	Completion completionFunc
}

func (f IntOption) HelpString() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage))
}

func (f IntOption) CompletionStrings() []string {
	return []string{prefixedNames(f.Name)}
}

func (f IntOption) ApplyNamed(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValInt, err := strconv.ParseUint(envVal, 10, 64)
			if err == nil {
				f.Value = int(envValInt)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Int(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f IntOption) ApplyPositional(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValInt, err := strconv.ParseUint(envVal, 10, 64)
			if err == nil {
				f.Value = int(envValInt)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.IntArg(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f IntOption) name() string {
	return f.Name
}

func (f IntOption) usage() string {
	if f.Usage == "" {
		return fmt.Sprintf("default = %v", f.Value)
	} else {
		return fmt.Sprintf("%s; default = %v", f.Usage, f.Value)
	}
}

func (f IntOption) local() bool                { return f.Local }
func (f IntOption) completion() completionFunc { return f.Completion }
func (f IntOption) validation() validationFunc { return nil }

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
	Var        *float64
	Optional   bool
	Local      bool
	Completion completionFunc
}

func (f Float64Option) HelpString() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Usage))
}

func (f Float64Option) CompletionStrings() []string {
	return []string{prefixedNames(f.Name)}
}

func (f Float64Option) ApplyNamed(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValFloat, err := strconv.ParseFloat(envVal, 10)
			if err == nil {
				f.Value = float64(envValFloat)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Float64(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f Float64Option) ApplyPositional(set *flags.Set) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValFloat, err := strconv.ParseFloat(envVal, 10)
			if err == nil {
				f.Value = float64(envValFloat)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Float64Arg(name, f.Value, f.Usage, f.Var, f.Optional)
	})
}

func (f Float64Option) name() string {
	return f.Name
}

func (f Float64Option) usage() string {
	if f.Usage == "" {
		return fmt.Sprintf("default = %v", f.Value)
	} else {
		return fmt.Sprintf("%s; default = %v", f.Usage, f.Value)
	}
}

func (f Float64Option) local() bool                { return f.Local }
func (f Float64Option) completion() completionFunc { return f.Completion }
func (f Float64Option) validation() validationFunc { return nil }

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
