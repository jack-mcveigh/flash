package main

import (
	"log"
	"os"

	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/core/updating"
	"github.com/jmcveigh55/flash/pkg/interface/cli"
	"github.com/jmcveigh55/flash/pkg/storage/memory"
)

func main() {
	r := memory.New()
	a := adding.New(r)
	d := deleting.New(r)
	g := getting.New(r)
	u := updating.New(r)

	app := cli.New(a, d, g, u)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
