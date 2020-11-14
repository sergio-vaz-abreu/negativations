package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func main() {
	var port int
	app := cli.NewApp()
	app.Name = ApplicationName
	app.Description = "A application for legacy negativations integration"
	app.Version = Version + "(" + GitCommit + ")"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			Usage:       "Port which the API will run",
			Value:       80,
			Destination: &port,
		},
	}
	app.Action = func(cliCtx *cli.Context) error {
		logrus.Info(port)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
