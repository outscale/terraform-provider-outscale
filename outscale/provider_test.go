package outscale

import (
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"outscale": testAccProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("OUTSCALE_ACCESSKEYID") == "" ||
		os.Getenv("OUTSCALE_REGION") == "" ||
		os.Getenv("OUTSCALE_SECRETKEYID") == "" ||
		os.Getenv("OUTSCALE_IMAGEID") == "" {
		t.Fatal("`OUTSCALE_ACCESSKEYID`, `OUTSCALE_SECRETKEYID`, `OUTSCALE_REGION` and `OUTSCALE_IMAGEID` must be set for acceptance testing")
	}
}

func testAccWait(n time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(n)
		return nil
	}
}
