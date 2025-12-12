package testutils

import (
	"fmt"

	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	fwtf "github.com/hashicorp/terraform-plugin-testing/terraform"
)

func ImportStepWithStateIdFuncSDKv2(resourceName string, importStateIdFunc sdkresource.ImportStateIdFunc, ignore ...string) sdkresource.TestStep {
	return importStepSDKv2(resourceName, importStateIdFunc, ignore...)
}

func ImportStepSDKv2(resourceName string, ignore ...string) sdkresource.TestStep {
	idFunc := func(s *sdktf.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
	return importStepSDKv2(resourceName, idFunc, ignore...)
}

func importStepSDKv2(resourceName string, importStateIdFunc sdkresource.ImportStateIdFunc, ignore ...string) sdkresource.TestStep {
	step := sdkresource.TestStep{
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

func ImportStepWithStateIdFuncFW(resourceName string, importStateIdFunc fwresource.ImportStateIdFunc, ignore ...string) fwresource.TestStep {
	return importStepFW(resourceName, importStateIdFunc, ignore...)
}

func ImportStepFW(resourceName string, ignore ...string) fwresource.TestStep {
	idFunc := func(s *fwtf.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
	return importStepFW(resourceName, idFunc, ignore...)
}

func importStepFW(resourceName string, importStateIdFunc fwresource.ImportStateIdFunc, ignore ...string) fwresource.TestStep {
	step := fwresource.TestStep{
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
