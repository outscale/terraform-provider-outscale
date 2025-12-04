package outscale

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOKSCluster_basic(t *testing.T) {
	id := rand.Int()
	projectName := fmt.Sprintf("%s-%d", "project", id)
	clusterName := fmt.Sprintf("%s-%d", "cluster-basic", id)
	resourceName := "outscale_oks_cluster.cluster"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: oksClusterConfig(projectName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cidr_pods", "10.91.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "cidr_service", "10.92.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.32"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"request_id"},
			},
		},
	})
}

func oksClusterConfig(projectName, clusterName string) string {
	return fmt.Sprintf(`
	resource "outscale_oks_project" "project" {
	  name = "%s"
	  cidr = "10.50.0.0/18"
	  region = "eu-west-2"
	  disable_api_termination = false

	  tags = {
		test = "TestAccClusterBasic"
	  }
	}

	resource "outscale_oks_cluster" "cluster" {
		project_id = outscale_oks_project.project.id
		admin_whitelist = ["0.0.0.0/0"]
		cidr_pods = "10.91.0.0/16"
		cidr_service = "10.92.0.0/16"
		version = "1.32"
		name = "%s"
		control_planes = "cp.mono.master"

		tags = {
			test = "TestAccClusterBasic"
		}
	}
`, projectName, clusterName)
}
