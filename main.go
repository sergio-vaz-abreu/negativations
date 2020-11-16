package main

import (
	"context"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"github.com/negativations/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

func main() {
	var filename string
	app := cli.NewApp()
	app.Name = ApplicationName
	app.Description = "A application for legacy negativations integration"
	app.Version = Version + "(" + GitCommit + ")"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Configuration file",
			Value:       "config.json",
			Destination: &filename,
		},
	}
	app.Action = func(cliCtx *cli.Context) error {
		var applicationConfig api.ApplicationConfig
		err := loadFile(filename, &applicationConfig)
		if err != nil {
			return err
		}
		return run(applicationConfig)
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

func run(applicationConfig api.ApplicationConfig) error {
	logrus.Info("starting application")
	app, err := api.LoadAPI(applicationConfig)
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

func loadFile(filename string, data interface{}) error {
	if err := config.Load(file.NewSource(
		file.WithPath(filename),
	)); err != nil {
		return errors.Wrap(err, "failed to loading file")
	}
	err := config.Scan(&data)
	return errors.Wrap(err, "failed parsing data from file")
}
