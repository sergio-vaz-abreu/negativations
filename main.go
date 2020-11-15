package main

import (
	"context"
	"github.com/negativations/api"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/signal"
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
		config := api.ArangoConfig{
			Host:     "localhost",
			Port:     8529,
			User:     "root",
			Password: "somepassword",
		}
		return run(port, config)
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func run(port int, config api.ArangoConfig) error {
	app, err := api.LoadAPI(port, "http://localhost", config)
	if err != nil {
		return err
	}
	appErr := app.Run()
	ctx := gracefullyShutdown()
	defer app.Shutdown()
	select {
	case err := <-appErr:
		return err
	case <-ctx.Done():
		logrus.Info("gracefully shutdown")
		return nil
	}
}

func gracefullyShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		cancel()
	}()
	return ctx
}
