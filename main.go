package main

import (
	"github.com/clausa/packer-builder-shell/builder/shell"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(shell.Builder))
	server.Serve()
}
