package main

import (
	"os"
	"strings"
	"time"

	"github.com/nothinux/octo-proxy/pkg/config"
	"github.com/nothinux/octo-proxy/pkg/runner"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"
)

// var banner = `         _
//  ___ ___| |_ ___ ___ ___ ___ _ _ _ _
// | . |  _|  _| . | . |  _| . |_'_| | |
// |___|___|_| |___|  _|_| |___|_,_|_  |
//                 |_|             |___| v%s

// `
// var usage = `Usage of octo:
// octo [flag] arguments...

// Flags:
//   -config
//     Specify config location path (default: ./config.yaml)
//   -listener
//     Specify listener for running octo-proxy (default: 1.0.0.0:5000)
//   -target
//     Specify target backend which traffic will be forwarded
//   -metrics
//     Specify address and port to run the metrics server
//   -debug
//     Enable debug log messages
//   -version
//     Print octo-proxy version

// `

// var (
// 	Version = "x.X"
// 	showBanner = fmt.Sprintf(banner, Version)
// )

func main() {
	if err := runMain(); err != nil {
		log.Fatal().Err(err).Msg("failed to start")
	}
}

func setupLogger(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro

	consoleWriter := zerolog.ConsoleWriter{
		TimeFormat: time.StampMicro,
		Out:        os.Stdout,
	}

	log.Logger = log.Output(consoleWriter)

	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func runMain() error {
	app := &cli.App{
		Name:                 "tcp-proxy",
		Usage:                "for tcp port proxy",
		EnableBashCompletion: true,
		Version:              "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listener",
				Value:   "127.0.0.1:8050",
				Aliases: []string{"l"},
				Usage:   "local ip:port to use for listen",
			},
			&cli.StringFlag{
				Name:    "target",
				Value:   "",
				Aliases: []string{"t"},
				Usage:   `remote target ip:port to proxy, mutil use,"192.168.1.15:9090,192.168.1.15:9091"`,
			},
			&cli.StringFlag{
				Name:    "metrics",
				Value:   "0.0.0.0:9123",
				Aliases: []string{"m"},
				Usage:   "run metrics on this ip:port",
			},
			&cli.StringFlag{
				Name:    "config",
				Value:   "config.yaml",
				Aliases: []string{"c"},
				Usage:   "config file path, if -target set, this will not used",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "enable debug log",
			},
		},
		Action: func(c *cli.Context) error {

			listenerC := c.String("listener")
			targetC := c.String("target")
			configPathC := c.String("config")
			metricsC := c.String("metrics")
			debugC := c.Bool("debug")

			setupLogger(debugC)

			// use target first
			if targetC != "" {
				targets := strings.Split(targetC, ",")
				c, err := config.GenerateConfig(listenerC, targets, metricsC)
				if err != nil {
					return err
				}

				if err := runner.Run(c, ""); err != nil {
					return err
				}

				return nil
			} else {
				// use config file secondly
				cfg, err := config.New(configPathC)
				if err != nil {
					return err
				}

				err = runner.Run(cfg, configPathC)
				if err != nil {
					return err
				}

			}

			return nil
		},
	}

	return app.Run(os.Args)
}
