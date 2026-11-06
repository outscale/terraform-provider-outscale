package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"

	"github.com/outscale/terraform-provider-outscale/provider"
	"github.com/outscale/terraform-provider-outscale/version"
)

func main() {
	ctx := context.Background()
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	server, err := provider.NewProtocol6Server(ctx, version.Version)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/outscale/outscale",
		func() tfprotov6.ProviderServer { return server },
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
