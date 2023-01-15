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
	retry.DefaultDelay = time.Second
	app.Run(os.Args)
}
