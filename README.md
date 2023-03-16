# simplecli

A simple CLI library for Go.

## Installation

```bash
go get github.com/joeychilson/simplecli
```

## Example

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/joeychilson/simplecli"
)

func main() {
	var name string

	fs := flag.NewFlagSet("hellocmd", flag.ExitOnError)
	fs.StringVar(&name, "name", "world", "name to say hello to")

	helloCmd := &simplecli.Command{
		Name:    "hello",
		Aliases: []string{"hi", "hey"},
		Usage:   "say hello to someone",
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Printf("Hello %s!\n", name)
			return nil
		},
	}

	rootCmd := &simplecli.Command{
		Name:        "mycli",
		Aliases:     []string{"my"},
		SubCommands: []*simplecli.Command{helloCmd},
	}

	if err := rootCmd.ParseAndExec(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```
