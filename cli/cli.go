// Copyright (c) 2022 Wireleap

package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/wireleap/common/cli/fsdir"
)

type Subcmd struct {
	*flag.FlagSet
	io.Writer

	Desc     string
	Usage    func()
	Run      func(fsdir.T)
	Sections []Section
	Hidden   bool

	args0 string
}

type Section struct {
	Title   string
	Entries []Entry
}

type Entry struct{ Key, Value string }

func (sub *Subcmd) Output() io.Writer {
	if sub.Writer != nil {
		return sub.Writer
	}
	return tabwriter.NewWriter(sub.FlagSet.Output(), 0, 8, 2, ' ', 0)
}

func (sub *Subcmd) SetDefaultUsage() {
	sub.Usage = func() {
		o := sub.Output()
		// stupid go, this number is in stdlib but not exported
		n := 0
		sub.VisitAll(func(_ *flag.Flag) { n++ })
		options := ""

		if n > 0 {
			options = " [OPTIONS]"
		}
		fmt.Fprintf(o, "Usage: %s %s%s\n\n", sub.args0, sub.Name(), options)
		fmt.Fprintf(o, "%s\n", sub.Desc)
		if n > 0 {
			fmt.Fprintln(o, "\nOptions:")
		}
		sub.VisitAll(func(f *flag.Flag) {
			name, usage := flag.UnquoteUsage(f)
			dash := "--"
			if len(f.Name) == 1 {
				dash = "-"
			}
			fmt.Fprintf(o, "  %s%s %s\t%s\n", dash, f.Name, name, usage)
		})
		for _, s := range sub.Sections {
			fmt.Fprintf(o, "\n%s:\n", s.Title)
			for _, e := range s.Entries {
				fmt.Fprintf(o, "  %s\t%s\n", e.Key, e.Value)
			}
		}
		if w, ok := o.(*tabwriter.Writer); ok {
			w.Flush()
		}
		os.Exit(2)
	}
}

func (sub *Subcmd) SetMinimalUsage(example string) {
	sub.Usage = func() {
		o := sub.Output()
		if example != "" {
			example = " " + example
		}
		fmt.Fprintf(o, "Usage: %s %s%s\n\n", sub.args0, sub.Name(), example)
		fmt.Fprintf(o, "%s\n", sub.Desc)
		for _, s := range sub.Sections {
			fmt.Fprintf(o, "\n%s:\n", s.Title)
			for _, e := range s.Entries {
				fmt.Fprintf(o, "  %s\t%s\n", e.Key, e.Value)
			}
		}
		if w, ok := o.(*tabwriter.Writer); ok {
			w.Flush()
		}
		os.Exit(2)
	}
}

type CLI struct {
	Subcmds  []*Subcmd
	Sections []Section
}

func (c CLI) Usage() {
	w := tabwriter.NewWriter(os.Stderr, 0, 8, 4, ' ', 0)

	fmt.Fprintf(w, "Usage: %s COMMAND [OPTIONS]\n\n", os.Args[0])

	fmt.Fprintln(w, "Commands:")
	fmt.Fprintf(w, "  %s\t%s\n", "help", "Display this help message or help on a command")

	for _, sub := range c.Subcmds {
		// skip dummy & hidden commands
		if sub == nil || sub.Hidden {
			continue
		}

		fmt.Fprintf(w, "  %s\t%s\n", sub.Name(), sub.Desc)
	}

	fmt.Fprintln(w)

	for _, s := range c.Sections {
		fmt.Fprintf(w, "%s:\n", s.Title)

		for _, e := range s.Entries {
			fmt.Fprintf(w, "  %s\t%s\n", e.Key, e.Value)
		}

		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "Run '%s help COMMAND' for more information on a command.\n", os.Args[0])

	w.Flush()
	os.Exit(2)
}

func (c CLI) Parse(args []string) *Subcmd {
	if len(args) < 2 {
		c.Usage()
	}

	maybesub := args[1]

	var sub *Subcmd

	// special-case help
	help := false
	if maybesub == "help" {
		if len(args[2:]) != 1 {
			c.Usage()
		}

		maybesub = args[2]
		help = true
	}

	for _, maybe := range c.Subcmds {
		// skip dummy commands
		if maybe == nil {
			continue
		}

		if maybesub == maybe.Name() {
			if maybe.Usage == nil {
				maybe.SetDefaultUsage()
			}

			if help {
				maybe.args0 = args[0]
				maybe.Usage()
			}

			sub = maybe
			break
		}
	}

	if sub == nil {
		c.Usage()
	}

	sub.args0 = args[0]
	err := sub.Parse(args[2:])

	if err != nil {
		// don't need to print error here -- flag pkg does it
		// actually it exits too, but just in case...
		os.Exit(2)
	}

	return sub
}
