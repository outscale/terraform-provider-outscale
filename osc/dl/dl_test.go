package dl

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-outscale/osc"
)

func TestNewDLClient(t *testing.T) {
	config := osc.Config{
		Credentials: &osc.Credentials{
			AccessKey: "AKID",
			SecretKey: "SecretKey",
			Region:    "region",
		},
	}

	c, err := NewDLClient(config)
	if err != nil {
		t.Fatalf("Got error %s", err)
	}
	if c == nil {
		t.Fatalf("Bad Client")
	}
}
