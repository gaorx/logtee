package logtee

import (
	"github.com/urfave/cli"
	"math/rand"
	"os"
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
		cli.StringFlag{
			Name:  "file",
			Value: "",
			Usage: "Filename",
		},
		cli.BoolFlag{
			Name:  "follow,f",
			Usage: "Follow file",
		},
	}
	app.Action = main0
	app.Run(os.Args)
}

func main0(cc *cli.Context) error {
	var err error

	// load config
	configFilename := cc.String("config")
	conf, err := LoadConfig(configFilename)
	if err != nil {
		return err
	}

	// parse config
	hh, err := ParseHandlers(conf)
	if err != nil {
		return err
	}

	// source
	var source Source
	filename := cc.Args().First()
	if filename != "" {
		source, err = FileSource(filename, cc.Bool("follow"))
		if err != nil {
			return err
		}
	} else {
		source = StdinSource()
	}

	// init
	err = hh.Init()
	if err != nil {
		return err
	}
	defer hh.Close()

	// go
	hh.Do(source)

	return nil
}
