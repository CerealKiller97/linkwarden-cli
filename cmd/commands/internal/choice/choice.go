// Package choice provides flags that take one of a fixed set of names (such as
// --format pdf) and map them to the integer codes the Linkwarden API expects, so
// users never have to know that 2 means PDF. Each flag also completes its names
// in the shell, with a short description of each.
package choice

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Option is one accepted name for a flag, the API value it stands for, and a
// short description shown in shell completion.
type Option[T ~int] struct {
	Name        string
	Value       T
	Description string
}

// Value is a pflag.Value that accepts only the names of its options.
type Value[T ~int] struct {
	current *Option[T]
	options []Option[T]
}

func (v *Value[T]) String() string {
	if v.current == nil {
		return ""
	}

	return v.current.Name
}

func (v *Value[T]) Set(s string) error {
	name := strings.ToLower(strings.TrimSpace(s))
	for i := range v.options {
		if v.options[i].Name == name {
			v.current = &v.options[i]

			return nil
		}
	}

	return fmt.Errorf("must be one of %s", v.names())
}

func (v *Value[T]) Type() string {
	return "string"
}

func (v *Value[T]) names() string {
	names := make([]string, len(v.options))
	for i, o := range v.options {
		names[i] = o.Name
	}

	return strings.Join(names, "|")
}

// Flag adds --<name> to cmd. It accepts only the names in options, lists them
// after usage in --help, and completes them in the shell. def is the name of
// the default option, or "" for no default.
func Flag[T ~int](cmd *cobra.Command, name, def, usage string, options []Option[T]) {
	v := &Value[T]{options: options}
	if def != "" {
		if err := v.Set(def); err != nil {
			panic(fmt.Sprintf("choice: default %q for --%s: %v", def, name, err))
		}
	}

	cmd.Flags().Var(v, name, fmt.Sprintf("%s One of: %s.", usage, v.names()))

	completions := make([]cobra.Completion, len(options))
	for i, o := range options {
		completions[i] = cobra.CompletionWithDesc(o.Name, o.Description)
	}
	if err := cmd.RegisterFlagCompletionFunc(name, cobra.FixedCompletions(completions, cobra.ShellCompDirectiveNoFileComp)); err != nil {
		panic(fmt.Sprintf("choice: register completion for --%s: %v", name, err))
	}
}

// Get returns the API value of the option chosen for --<name>, and whether one
// was chosen (or defaulted). It errors if the flag wasn't added with Flag.
func Get[T ~int](cmd *cobra.Command, name string) (T, bool, error) {
	f := cmd.Flags().Lookup(name)
	if f == nil {
		return 0, false, fmt.Errorf("flag --%s not defined", name)
	}
	v, ok := f.Value.(*Value[T])
	if !ok {
		return 0, false, fmt.Errorf("flag --%s is not a choice of %T", name, *new(T))
	}
	if v.current == nil {
		return 0, false, nil
	}

	return v.current.Value, true, nil
}
