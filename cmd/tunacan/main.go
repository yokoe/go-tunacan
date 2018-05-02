package main

import (
	"flag"
	"os"

	"github.com/yokoe/tunacan"
)

// https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang?utm_medium=organic&utm_source=google_rich_qa&utm_campaign=google_rich_qa
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "arrayFlags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	if len(os.Args) > 0 && os.Args[1] == "server" {
		launchServer()
		return
	}

	var outputFilename string
	var sourceFilenames arrayFlags

	flag.Var(&sourceFilenames, "s", "Source image filepaths.")
	flag.StringVar(&outputFilename, "o", "", "Output filepath.")
	flag.Parse()

	err := tunacan.Concat(sourceFilenames, outputFilename)
	if err != nil {
		panic(err)
	}
}
