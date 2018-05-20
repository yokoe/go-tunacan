package tunacan

import (
	"flag"
	"log"
	"os"
	"strings"

	tunacan "github.com/yokoe/tunacan"
)

type ConcatCommand struct {
}

func (c *ConcatCommand) Synopsis() string {
	return "Concatenates image files."
}

func (c *ConcatCommand) Help() string {
	return "Usage: tunacan concat file1.png file2.png -o output.png"
}

func (c *ConcatCommand) Run(args []string) int {
	var outputFilename string

	var inputFilenames = []string{}
	prevArg := ""
	for _, arg := range args {
		if !strings.HasPrefix(prevArg, "-") && !strings.HasPrefix(arg, "-") {
			inputFilenames = append(inputFilenames, arg)
		}
		prevArg = arg
	}

	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.StringVar(&outputFilename, "o", "", "Output filepath.")
	flags.Parse(os.Args[1:])
	for 0 < flags.NArg() {
		flags.Parse(flags.Args()[1:])
	}

	err := tunacan.Concat(inputFilenames, outputFilename)
	if err != nil {
		log.Fatalln(err)
		return 1
	}
	return 0
}
