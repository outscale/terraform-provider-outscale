package oks_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOKSManifest_DefaultNamespace(t *testing.T) {
	projectName := acctest.RandomWithPrefix("test-acc-project")
	clusterName := acctest.RandomWithPrefix("test-acc-cluster")
	resourceName := "outscale_oks_manifest.configmap"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: clusterConfig(projectName, clusterName) + configMapManifestConfig("one"),
				Check:  resource.TestCheckResourceAttrSet(resourceName, "object"),
			},
			{
				// a namespaced object written without namespace lives in "default" and is updated in place
				Config: clusterConfig(projectName, clusterName) + configMapManifestConfig("two"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttrSet(resourceName, "object"),
			},
		},
	})
}

func configMapManifestConfig(value string) string {
	return fmt.Sprintf(`
resource "outscale_oks_manifest" "configmap" {
  cluster_id = outscale_oks_cluster.cluster.id
  manifest   = <<-EOT
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: test-acc-default-namespace
    data:
      key: %s
  EOT
}
`, value)
}
