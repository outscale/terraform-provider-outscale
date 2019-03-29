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

func datasourceOutscaleKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.DescribeKeyPairsInput{}

	filters, filtersOk := d.GetOk("filter")
	KeyName, KeyNameisOk := d.GetOk("key_name")

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
	d.Set("key_name", keypair.KeyName)
	d.Set("key_fingerprint", keypair.KeyFingerprint)
	d.Set("request_id", resp.RequestId)
	d.SetId(resource.UniqueId())
	return nil
}

func datasourceOutscaleKeyPair() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleKeyPairRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// Attributes
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
