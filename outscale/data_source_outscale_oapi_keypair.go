package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func datasourceOutscaleOApiKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.DescribeKeyPairsInput{}

	filters, filtersOk := d.GetOk("filter")
	KeyName, KeyNameisOk := d.GetOk("keypair_name")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if KeyNameisOk {
		req.KeyNames = []*string{aws.String(KeyName.(string))}
	}

	var resp *fcu.DescribeKeyPairsOutput
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeKeyPairs(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	if err != nil {
		return err
	}

	if len(resp.KeyPairs) < 1 {
		return fmt.Errorf("Unable to find key pair, please provide a better query criteria ")
	}
	if len(resp.KeyPairs) > 1 {

		return fmt.Errorf("Found to many key pairs, please provide a better query criteria ")
	}

	keypair := resp.KeyPairs[0]
	d.Set("keypair_name", keypair.KeyName)
	d.Set("keypair_fingerprint", keypair.KeyFingerprint)
	d.SetId(resource.UniqueId())
	return nil
}

func datasourceOutscaleOAPIKeyPair() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleKeyPairRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// Attributes
			"keypair_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keypair_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
