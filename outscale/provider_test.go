package outscale

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

var testAccProviders map[string]*schema.Provider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()

	testAccProviders = map[string]*schema.Provider{
		"outscale": testAccProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("OUTSCALE_ACCESSKEYID") == "" ||
		os.Getenv("OUTSCALE_REGION") == "" ||
		os.Getenv("OUTSCALE_SECRETKEYID") == "" ||
		os.Getenv("OUTSCALE_IMAGEID") == "" ||
		os.Getenv("OUTSCALE_ACCOUNT") == "" {
		t.Fatal("`OUTSCALE_ACCESSKEYID`, `OUTSCALE_SECRETKEYID`, `OUTSCALE_REGION`, `OUTSCALE_ACCOUNT` and `OUTSCALE_IMAGEID` must be set for acceptance testing")
	}
}

func testAccPreCheckValues(t *testing.T) {
	if utils.GetEnvVariableValue([]string{"OSC_ACCESS_KEY", "OUTSCALE_ACCESSKEYID"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_SECRET_KEY", "OUTSCALE_SECRETKEYID"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_REGION", "OUTSCALE_REGION"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_ACCOUNT_ID", "OUTSCALE_ACCOUNT"}) == "" ||
		utils.GetEnvVariableValue([]string{"OSC_IMAGE_ID", "OUTSCALE_IMAGEID"}) == "" {
		t.Fatal("`OSC_ACCESS_KEY`, `OSC_SECRET_KEY`, `OSC_REGION`, `OSC_ACCOUNT_ID` and `OUTSCALE_IMAGEID` must be set for acceptance testing")
	}
}
