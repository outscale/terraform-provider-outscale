package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/outscale/terraform-provider-outscale/version"
)

func TestAccNet_Bacic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"outscale": providerserver.NewProtocol5WithError(New(version.GetVersion())),
		},
		PreCheck: func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configNetBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("outscale_net.basic_net", "ip_range", "10.0.0.0/16"),
				),
			},
		},
	})
}
func TestAccNet_UpdateTags(t *testing.T) {
	t.Parallel()
	netName := "outscale_net.basic_net"
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"outscale": providerserver.NewProtocol5WithError(New(version.GetVersion())),
		},
		PreCheck: func() { TestAccFwPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configNetUpdateTags("NetTags"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(netName, "tags.#"),
				),
			},
			{
				Config: configNetUpdateTags("NetTags2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(netName, "tags.#"),
				),
			},
		},
	})
}

const configNetBasic = `
	resource "outscale_net" "basic_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name"
			value = "testacc-net-rs"
		}
			timeouts {
			create = "15m"
			}

	}
`

func configNetUpdateTags(tagValue string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "basic_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "name"
			value = "%s"
		}
	   }
`, tagValue)
}
