package outscale

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/dl"
)

func TestAccOutscaleOAPIDirectLink_basic(t *testing.T) {
	connectionName := fmt.Sprintf("tf-dx-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOutscaleOAPIDirectLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDxOAPIConnectionConfig(connectionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOutscaleOAPIDirectLinkExists("outscale_directlink.hoge"),
					resource.TestCheckResourceAttr("outscale_directlink.hoge", "direct_link_name", connectionName),
					resource.TestCheckResourceAttr("outscale_directlink.hoge", "bandwidth", "1Gbps"),
					resource.TestCheckResourceAttr("outscale_directlink.hoge", "site", "EqSe2"),
				),
			},
		},
	})
}

func testAccCheckOutscaleOAPIDirectLinkDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).DL

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_directlink" {
			continue
		}

		input := &dl.DescribeConnectionsInput{
			ConnectionID: aws.String(rs.Primary.ID),
		}

		var resp *dl.Connections
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.API.DescribeConnections(input)

			if err != nil {
				if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
		for _, v := range resp.Connections {
			if *v.ConnectionID == rs.Primary.ID && !(*v.ConnectionState == "deleted") {
				return fmt.Errorf("[DESTROY ERROR] Dx Connection (%s) not deleted", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckOutscaleOAPIDirectLinkExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}

func testAccDxOAPIConnectionConfig(n string) string {
	return fmt.Sprintf(`
resource "outscale_directlink" "hoge" {
  direct_link_name = "%s"
  bandwidth = "1Gbps"
  site = "EqSe2"
}
`, n)
}
