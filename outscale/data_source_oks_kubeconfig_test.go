package outscale

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOKSKubeconfigDataSource_basic(t *testing.T) {
	id := rand.Int()
	projectName := fmt.Sprintf("%s-%d", "project-kubeconfig", id)
	clusterName := fmt.Sprintf("%s-%d", "cluster-kubeconfig", id)
	resourceName := "data.outscale_oks_kubeconfig.config"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: DefineTestProviderFactoriesV6(),

		Steps: []resource.TestStep{
			{
				Config: oksKubeconfigConfig(projectName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig"),
				),
			},
		},
	})
}

func oksKubeconfigConfig(projectName, clusterName string) string {
	return fmt.Sprintf(`
	resource "outscale_oks_project" "project" {
	  	name = "%s"
		cidr = "10.50.0.0/18"
		region = "eu-west-2"
		disable_api_termination = false

		tags = {
			test = "TestAccKubeconfigDataSourceBasic"
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
		disable_api_termination = false

		tags = {
			test = "TestAccKubeconfigDataSourceBasic"
		}
	}

	data "outscale_oks_kubeconfig" "config" {
		cluster_id = outscale_oks_cluster.cluster.id
	}
`, projectName, clusterName)
}
