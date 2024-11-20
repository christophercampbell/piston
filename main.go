package main

import (
	"fmt"
	"os"

	"github.com/0xPolygon/maera/app"
	"github.com/urfave/cli/v2"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = app.AppName

	cliApp.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "Run the program",
			Action: app.Run,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     app.JWTKey,
					Usage:    "JWT file for engine API auth",
					Required: true,
				},
				&cli.StringFlag{
					Name:        app.EthUrlKey,
					Usage:       "URL to engine's public API",
					Required:    false,
					DefaultText: app.EthUrlDefault,
				},
				&cli.StringFlag{
					Name:        app.EngineUrlKey,
					Usage:       "URL to engine's secure API",
					Required:    false,
					DefaultText: app.EngineUrlDefault,
				},
				&cli.DurationFlag{
					Name:        app.PeriodKey,
					Usage:       "block production rate in seconds",
					Required:    false,
					DefaultText: app.PeriodDefault.String(),
				},
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		fmt.Println("fatal err:", err)
		os.Exit(1)
	}
}
