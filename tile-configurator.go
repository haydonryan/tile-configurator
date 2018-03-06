package main

import (
	"os"

	"github.com/haydonryan/tile-configurator/config"
	"github.com/haydonryan/tile-configurator/diffbase"
	"github.com/haydonryan/tile-configurator/ingest"
	flags "github.com/jessevdk/go-flags"
	"github.com/xchapter7x/lo"
)

func main() {
	var opts struct {
		Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug info"`

		Ingest ingest.Ingest     `command:"ingest" description:"Reads ops manager api output into tile configurator yaml format"`
		Config config.Config     `command:"config" description:"Configures the tile in Ops Manager"`
		Diff   diffbase.Diffbase `command:"diff" description:"Shows the structured diff of two manifests"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		lo.G.Debug(err)
		os.Exit(1)
	}
}
