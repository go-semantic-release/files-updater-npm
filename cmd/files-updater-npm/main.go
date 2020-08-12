package main

import (
	npmUpdater "github.com/go-semantic-release/files-updater-npm/pkg/updater"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		FilesUpdater: func() updater.FilesUpdater {
			return &npmUpdater.Updater{}
		},
	})
}
