package tunacan

import "flag"

type ServerCommand struct {
}

func (c *ServerCommand) Synopsis() string {
	return "Launch HTTP server."
}

func (c *ServerCommand) Help() string {
	return "Usage: tunacan server"
}

func (c *ServerCommand) Run(args []string) int {
	port := "8080"
	bucket := ""
	flags := flag.NewFlagSet("server", flag.ContinueOnError)
	flags.StringVar(&port, "p", "8080", "Port number to listen.")
	flags.StringVar(&bucket, "b", "", "Cloud storage bucket name.")
	flags.Parse(args)

	launchServer(port, bucket)
	return 0
}
