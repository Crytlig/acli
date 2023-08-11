package main

import (
	"log"
	"os"

	"github.com/crytlig/acli/lib"
	"github.com/crytlig/acli/version"

	cli "github.com/urfave/cli/v2"
)

var DebugFlag = "debug"

// TODO Add flag for choosing model
func main() {
	app := &cli.App{
		Name:      "acli",
		Version:   version.Version,
		Usage:     "Use acli to query a tool if you have forgotten a command or simply need help.",
		UsageText: "Example usage for Azure CLI.\nacli query 'get application id of app registration myapp123'",
		Commands: []*cli.Command{
			{
				Name:   "models",
				Usage:  "Lists the available models. Defaults to using gpt-3.5-turbo-0613",
				Action: lib.AvailableModels,
			},
			{
				Name:      "query",
				Aliases:   []string{"q"},
				Usage:     "Gets a command from your desired cli",
				UsageText: "acli [q]uery 'az cli get application id of app registration myapp123'\nacli [q]uery 'kubectl command for listing pods in dev namespace'",
				Action:    processQuery,
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  DebugFlag,
				Usage: "Enable debugging mode. --debug has to be specified before command",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func processQuery(c *cli.Context) error {
	query := c.Args().First()
	debugMode := c.Bool(DebugFlag)
	return lib.HandleRequest(c, query, debugMode)
}
