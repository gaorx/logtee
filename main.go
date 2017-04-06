package logtee

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func Main() {
	rand.Seed(time.Now().UnixNano())

	app := cli.NewApp()
	app.Name = "logtee"
	app.Usage = "LogTee"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config,c",
			Value: "",
			Usage: "Configuration path",
		},
		cli.BoolFlag{
			Name:  "follow,f",
			Usage: "Follow file",
		},
	}
	app.Action = func(cc *cli.Context) {
		err := main0(cc)
		if err != nil {
			panic(err)
		}
	}
	app.Run(os.Args)
}

func main0(cc *cli.Context) error {
	var err error

	// load config
	configFilename := cc.String("config")
	conf, err := LoadConfig(configFilename)
	if err != nil {
		return errors.Wrap(err, "Load config error")
	}

	// parse config
	hh, err := ParseHandlers(conf)
	if err != nil {
		return errors.Wrap(err, "Parse config error")
	}

	// source
	var source Source
	filename := cc.Args().First()
	if filename != "" {
		source, err = FileSource(filename, cc.Bool("follow"))
		if err != nil {
			return errors.Wrap(err, "Load source error")
		}
	} else {
		source = StdinSource()
	}

	// init
	err = hh.Init()
	if err != nil {
		return errors.Wrap(err, "Init error")
	}
	defer hh.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			switch sig {
			case os.Interrupt:
				hh.Flush()
				os.Exit(0)
			}
		}
	}()

	// go
	hh.Do(source)

	return nil
}
