package main

import (
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: Commands,
	}
	retry.DefaultDelay = time.Second * 5
	app.Run(os.Args)
}
