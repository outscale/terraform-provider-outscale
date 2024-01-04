package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/terraform-providers/terraform-provider-outscale/outscale"
	vers "github.com/terraform-providers/terraform-provider-outscale/version"
)

var (
	version string = vers.GetVersion()
)

func main() {

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(outscale.New(version)), // Example terraform-plugin-framework provider
		outscale.Provider().GRPCProvider,                   // Example terraform-plugin-sdk provider
	}

	//using muxer
	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/providers/outscale/outscale/",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}

//	plugin.Serve(&plugin.ServeOpts{
//		ProviderFunc: outscale.Provider,
//	})
//}
