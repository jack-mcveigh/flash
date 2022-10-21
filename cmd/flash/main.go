package main

import (
	"github.com/jmcveigh55/flash/pkg/core/adding"
	"github.com/jmcveigh55/flash/pkg/core/deleting"
	"github.com/jmcveigh55/flash/pkg/core/getting"
	"github.com/jmcveigh55/flash/pkg/interface/cli"
	"github.com/jmcveigh55/flash/pkg/storage/memory"
)

func main() {
	r := memory.New()
	a := adding.New(r)
	d := deleting.New(r)
	g := getting.New(r)

	cli.New(a, d, g)
}
