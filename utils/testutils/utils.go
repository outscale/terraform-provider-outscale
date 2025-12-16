package testutils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func ImportStepWithStateIdFunc(resourceName string, importStateIdFunc resource.ImportStateIdFunc, ignore ...string) resource.TestStep {
	return importStep(resourceName, importStateIdFunc, ignore...)
}

func ImportStep(resourceName string, ignore ...string) resource.TestStep {
	idFunc := func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
	return importStep(resourceName, idFunc, ignore...)
}

func importStep(resourceName string, importStateIdFunc resource.ImportStateIdFunc, ignore ...string) resource.TestStep {
	step := resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: importStateIdFunc,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}

func DefaultIgnores() []string {
	return []string{
		"request_id",
	}
}
