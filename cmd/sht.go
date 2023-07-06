package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	repl "sht/cmd/sht"
	"sht/lang"
	"sht/lang/runtime"

	"github.com/c-bata/go-prompt"
	"github.com/urfave/cli/v2"
)

func help(ctx *cli.Context) error {
	fmt.Println("oh SHT!")
	fmt.Println("")
	fmt.Println("Usage: sht [command] <args>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("	 run  <file>  run your sht script")
	fmt.Println("	 exec <code>  execute your sht code")
	fmt.Println("	 help         prints this")
	fmt.Println("")
	return nil
}

func replCompleter(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "exit", Description: "exit the repl"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func cmdRepl(ctx *cli.Context) error {
	runtime := runtime.CreateRuntime()
	repl.Start(runtime)
	return nil
}

func cmdRun(ctx *cli.Context) error {
	file, err := os.Open(ctx.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return nil
	}

	v, err := lang.Eval(bs)
	if err != nil {
		fmt.Println(err)
	} else if v != "" {
		// fmt.Println(v)
	}

	return nil
}

func cmdExec(ctx *cli.Context) error {
	s := ""
	args := ctx.Args()
	for i := 0; i < ctx.NArg(); i++ {
		s += args.Get(i) + " "
	}

	v, err := lang.Eval([]byte(s))
	if err != nil {
		fmt.Println(err)
	} else if v != "" {
		fmt.Println(v)
	}

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

			return cmdRepl(ctx)
		},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run your sht script",
				Action: cmdRun,
			},
			{
				Name:   "exec",
				Usage:  "execute your sht code",
				Action: cmdExec,
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
