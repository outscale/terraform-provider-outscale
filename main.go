package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/outscale/terraform-provider-outscale/outscale"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: outscale.Provider,
	})
}
