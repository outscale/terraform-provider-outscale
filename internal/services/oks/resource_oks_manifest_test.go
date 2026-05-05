package oks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOKSManifest_Basic(t *testing.T) {
	projectName := acctest.RandomWithPrefix("test-acc-project")
	clusterName := acctest.RandomWithPrefix("test-acc-cluster")
	resourceName := "outscale_oks_manifest.nodepool"
	var resourceID string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: clusterConfig(projectName, clusterName) + manifestConfig("nodepool", "pool-1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "manifest"),
					resource.TestCheckResourceAttrSet(resourceName, "object"),
					resource.TestCheckResourceAttr(resourceName, "skip_delete", "true"),
					resource.TestCheckResourceAttr(resourceName, "wait_for.timeout", "10m"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						resourceID = value
						return nil
					}),
				),
			},
			{
				// spec update is applied in place with server-side apply
				Config: clusterConfig(projectName, clusterName) + manifestConfigWithVolume("nodepool", "pool-1", 30, 1100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "manifest"),
					resource.TestCheckResourceAttrSet(resourceName, "object"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value != resourceID {
							return fmt.Errorf("resource %s was recreated during in-place update", resourceName)
						}
						return nil
					}),
				),
			},
			{
				// desired node count increased and waits for both nodes to be ready
				Config: clusterConfig(projectName, clusterName) + manifestConfigWithOptions("nodepool", "pool-1", 30, 1100, 2, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "manifest"),
					resource.TestCheckResourceAttrSet(resourceName, "object"),
					resource.TestCheckResourceAttr(resourceName, "wait_for.fields.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "wait_for.fields.status.progress.ready", "2"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value != resourceID {
							return fmt.Errorf("resource %s was recreated during desired node count update", resourceName)
						}
						return nil
					}),
				),
			},
			{
				// name update recreates the resource. create_before_destroy keeps a NodePool available before deleting the old one
				Config: clusterConfig(projectName, clusterName) + manifestConfigWithLifecycle("nodepool", "pool-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "manifest"),
					resource.TestCheckResourceAttrSet(resourceName, "object"),
					resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
						if value == resourceID {
							return fmt.Errorf("resource %s was not recreated", resourceName)
						}
						resourceID = value
						return nil
					}),
				),
			},
			{
				Config: clusterConfig(projectName, clusterName) +
					manifestConfigWithLifecycle("nodepool", "pool-2") +
					manifestConfigWithOptions("nodepool_duplicate", "pool-2", 20, 1000, 1, false),
				ExpectError: regexp.MustCompile("same name already exists"),
			},
			{
				Config:      clusterConfig(projectName, clusterName) + manifestConfigWithLifecycle("nodepool", "pool-2") + invalidManifestConfig(),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("Invalid manifest"),
			},
			{
				// skip_delete allows full config destroy with a single nodepool in the cluster
				Config:  clusterConfig(projectName, clusterName) + manifestConfigWithLifecycle("nodepool", "pool-2"),
				Destroy: true,
			},
			// testacc.ImportStep(resourceName, "kubeconfig", "request_id"),
		},
	})
}

func clusterConfig(projectName, clusterName string) string {
	return fmt.Sprintf(`
	resource "outscale_oks_project" "project" {
	  name = "%s"
	  cidr = "10.50.0.0/18"
	  region = "eu-west-2"
	  disable_api_termination = false
	}

	resource "outscale_oks_cluster" "cluster" {
		project_id = outscale_oks_project.project.id
		admin_whitelist = ["0.0.0.0/0"]
		cidr_pods = "10.91.0.0/16"
		cidr_service = "10.92.0.0/16"
		version = "1.35"
		name = "%s"
		control_planes = "cp.mono.master"
	}
`, projectName, clusterName)
}

func manifestConfig(resourceLabel, poolName string) string {
	return manifestConfigWithOptions(resourceLabel, poolName, 20, 1000, 1, false)
}

func manifestConfigWithVolume(resourceLabel, poolName string, volumeSize, volumeIOPS int) string {
	return manifestConfigWithOptions(resourceLabel, poolName, volumeSize, volumeIOPS, 1, false)
}

func manifestConfigWithLifecycle(resourceLabel, poolName string) string {
	return manifestConfigWithOptions(resourceLabel, poolName, 30, 1100, 1, true)
}

func manifestConfigWithOptions(resourceLabel, poolName string, volumeSize, volumeIOPS, desiredNodes int, createBeforeDestroy bool) string {
	return fmt.Sprintf(`
resource "outscale_oks_manifest" "%s" {
  cluster_id = outscale_oks_cluster.cluster.id
  wait = true
  skip_delete = true
  wait_for = {
    timeout = "10m"
    fields = {
      "status.progress.ready" = "%d"
    }
  }

  manifest = <<-YAML
apiVersion: oks.dev/v1beta2
kind: NodePool
metadata:
  name: %s
spec:
  autoHealing: true
  desiredNodes: %d
  nodeType: tinav7.c1r1p1
  upgradeStrategy:
    autoUpgradeEnabled: true
    autoUpgradeMaintenance:
      durationHours: 1
      startHour: 12
      weekDay: Tue
    maxSurge: 0
    maxUnavailable: 1
  volumes:
  - device: root
    dir: /
    size: %d
    type: io1
    iops: %d
  zones: [eu-west-2a]
YAML

	lifecycle {
		create_before_destroy = %t
	}
}`, resourceLabel, desiredNodes, poolName, desiredNodes, volumeSize, volumeIOPS, createBeforeDestroy)
}

func invalidManifestConfig() string {
	return `
resource "outscale_oks_manifest" "invalid" {
  cluster_id = outscale_oks_cluster.cluster.id
  manifest = <<-YAML
apiVersion: oks.dev/v1beta2
kind: NodePo
metadata:
  name: invalid-pool
YAML
}`
}
