package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func help(ctx *cli.Context) error {
	fmt.Println("oh SHT!")
	fmt.Println("")
	fmt.Println("Usage: sht [command] <args>")
	return nil
}

func repl(ctx *cli.Context) error {
	fmt.Println("...should open a repl here...")
	return nil
}

func run(ctx *cli.Context) error {
	fmt.Println("...should run the scripts: ", ctx.Args())
	return nil
}

func exec(ctx *cli.Context) error {
	fmt.Println("...should run the code: ", ctx.Args())
	return nil
}

func main() {
	app := &cli.App{
		Name:            "sht",
		Usage:           "A shell language",
		Version:         "v0.0.0",
		HideVersion:     true,
		HideHelp:        true,
		HideHelpCommand: false,
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() > 0 {
				fmt.Println("command not found: ", ctx.Args().Get(0))
				return nil
			}

			return repl(ctx)
		},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run your sht script",
				Action: run,
			},
			{
				Name:   "exec",
				Usage:  "execute your sht code",
				Action: exec,
			},
			{
				Name:   "help",
				Usage:  "prints this",
				Action: help,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
