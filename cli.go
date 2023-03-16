package simplecli

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// Command represents a single command in the CLI.
type Command struct {
	Name        string
	Aliases     []string
	Usage       string
	FlagSet     *flag.FlagSet
	SubCommands []*Command
	Exec        func(ctx context.Context, args []string) error
}

// ParseAndExec parses the given arguments and executes the appropriate command.
func (c *Command) ParseAndExec(ctx context.Context, args []string) error {
	var help bool

	if c.FlagSet == nil {
		c.FlagSet = flag.NewFlagSet(c.Name, flag.ExitOnError)
	}
	c.FlagSet.BoolVar(&help, "help", false, "display help information")

	if err := c.FlagSet.Parse(args); err != nil {
		return err
	}

	if help {
		c.PrintHelp()
		return nil
	}

	if len(c.FlagSet.Args()) > 0 {
		for _, cmd := range c.SubCommands {
			if cmd.Name == c.FlagSet.Args()[0] || slices.Contains(cmd.Aliases, c.FlagSet.Args()[0]) {
				return cmd.ParseAndExec(ctx, c.FlagSet.Args()[1:])
			}
		}
	}

	if c.Exec != nil {
		return c.Exec(ctx, c.FlagSet.Args())
	}

	c.PrintHelp()
	return nil
}

// PrintHelp displays the help information for the command.
func (c *Command) PrintHelp() {
	fmt.Println("Usage:")
	fmt.Printf("  %s [flags] [subcommand]\n", c.Name)

	fmt.Println("\nFlags:")
	c.FlagSet.VisitAll(func(f *flag.Flag) {
		if f.Usage == "" {
			fmt.Printf("  -%s\n", f.Name)
		} else {
			fmt.Printf("  -%s - %s\n", f.Name, f.Usage)
		}
	})

	if len(c.SubCommands) > 0 {
		fmt.Println("\nSubcommands:")
		for _, cmd := range c.SubCommands {
			aliasText := ""
			if len(cmd.Aliases) > 0 {
				aliasText = fmt.Sprintf("(%s)", strings.Join(cmd.Aliases, ", "))
			}
			if cmd.Usage == "" {
				fmt.Printf("  %s %s\n", cmd.Name, aliasText)
			} else {
				fmt.Printf("  %s %s - %s\n", cmd.Name, aliasText, cmd.Usage)
			}
		}
	}
}
