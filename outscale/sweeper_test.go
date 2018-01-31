package outscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (interface{}, error) {
	if os.Getenv("OUTSCALE_ACCESSKEYID") == "" {
		return nil, fmt.Errorf("empty OUTSCALE_ACCESSKEYID")
	}
	if os.Getenv("OUTSCALE_SECRETKEYID") == "" {
		return nil, fmt.Errorf("empty OUTSCALE_SECRETKEYID")
	}

	return nil, nil
}
