package main

import (
	"flag"
	"log"

	tunacan "github.com/yokoe/go-tunacan"
)

type ConcatCommand struct {
}

func (c *ConcatCommand) Synopsis() string {
	return "Concatenates image files."
}

func (c *ConcatCommand) Help() string {
	return "Usage: tunacan concat -s file1.png -s file2.png -o output.png"
}

func (c *ConcatCommand) Run(args []string) int {
	var outputFilename string
	var sourceFilenames arrayFlags

	flags := flag.NewFlagSet("concat", flag.ContinueOnError)
	flags.Var(&sourceFilenames, "s", "Source image filepaths.")
	flags.StringVar(&outputFilename, "o", "", "Output filepath.")
	flags.Parse(args)

	err := tunacan.Concat(sourceFilenames, outputFilename)
	if err != nil {
		log.Fatalln(err)
		return 1
	}
	return 0
}

// https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "arrayFlags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
