package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-outscale/version"
)

func TestAccFwOthers_DataSourceQuota(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			// new() is an example function that returns a provider.Provider
			"outscale": providerserver.NewProtocol5WithError(New(version.GetVersion())),
		},
		PreCheck: func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFwDataSourceOutscaleOAPIQuotaConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

const testAccFwDataSourceOutscaleOAPIQuotaConfig = `
	data "outscale_quota" "lbuQuota1" { 
	   filter {
	      name     = "quota_names"
	      values   = ["lb_listeners_limit"]
	   }
	}`
