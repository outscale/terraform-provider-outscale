package outscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/utils/testutils"
)

func TestAccNet_Attr_import(t *testing.T) {
	resourceName := "outscale_net_attributes.outscale_net_attributes"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleLinAttrConfig,
			},
			testutils.ImportStep(resourceName, testutils.DefaultIgnores()...),
		},
	})
}
