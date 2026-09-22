package oapi_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/outscale/terraform-provider-outscale/internal/testacc"
)

func TestAccOthers_Tags_DuplicateKey(t *testing.T) {
	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      testAccTagsDuplicateKey,
				PlanOnly:    true,
				ExpectError: testacc.AnyError,
			},
		},
	})
}

func TestAccOthers_Tags_UnknownKey(t *testing.T) {
	testacc.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      testAccTagsUnknownKey,
				PlanOnly:    true,
				ExpectError: testacc.AnyError,
			},
		},
	})
}

const testAccTagsUnknownKey = `
variable "tag_name" {
 default = "test"
}

resource "outscale_net" "net01" {
	ip_range = "10.10.0.0/16"
	tenancy  = "default"

  tags {
    key   = "prefix/${var.tag_name}"
    value = "owned"
  }

  tags {
    key   = "prefix/${var.tag_name}"
    value = "True"
  }
}
`

const testAccTagsDuplicateKey = `
resource "outscale_net" "net01" {
	ip_range = "10.10.0.0/16"
	tenancy  = "default"

  tags {
    key   = "name"
    value = "owned"
  }

  tags {
    key   = "name"
    value = "True"
  }
}
`
