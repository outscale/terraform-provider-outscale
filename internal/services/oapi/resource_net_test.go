package oapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccNet_Bacic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

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

func TestAccNet_Migration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: testacc.FrameworkMigrationTestSteps("1.0.1", configNetBasic),
	})
}

func TestAccNet_UpdateTags(t *testing.T) {
	netName := "outscale_net.basic_net"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),

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

func TestAccNet_CreateFailureKeepsState(t *testing.T) {
	resourceName := "outscale_net.basic_net"
	invalidTagKey := strings.Repeat("a", 256)
	tagValue := "testacc-net"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: testacc.CreateFailureReplacementSteps(
			resourceName,
			configNetUpdateTagsWithKey(invalidTagKey, tagValue),
			configNetUpdateTags(tagValue),
			resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "net_id"),
				resource.TestCheckResourceAttr(resourceName, "tags.0.value", tagValue),
			),
		),
	})
}

const configNetBasic = `
	resource "outscale_net" "basic_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "Name"
			value = "testacc-net-rs"
		}
	}
`

func configNetUpdateTags(tagValue string) string {
	return configNetUpdateTagsWithKey("name", tagValue)
}

func configNetUpdateTagsWithKey(tagKey, tagValue string) string {
	return fmt.Sprintf(`
	resource "outscale_net" "basic_net" {
		ip_range = "10.0.0.0/16"
		tags {
			key = "%s"
			value = "%s"
		}
	   }
`, tagKey, tagValue)
}
