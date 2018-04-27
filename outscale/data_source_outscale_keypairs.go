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

func datasourceOutscaleKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.DescribeKeyPairsInput{}

	filters, filtersOk := d.GetOk("filter")
	KeyName, KeyNameisOk := d.GetOk("key_name")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if KeyNameisOk {
		var names []*string
		for _, v := range KeyName.([]interface{}) {
			names = append(names, aws.String(v.(string)))
		}
		req.KeyNames = names
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

	d.SetId(resource.UniqueId())
	keypairs := make([]map[string]interface{}, len(resp.KeyPairs))
	for k, v := range resp.KeyPairs {
		keypair := make(map[string]interface{})
		if v.KeyName != nil {
			keypair["key_name"] = *v.KeyName
		}
		if v.KeyFingerprint != nil {
			keypair["key_fingerprint"] = *v.KeyFingerprint
		}
		keypairs[k] = keypair
	}
	d.Set("key_set", keypairs)
	d.Set("request_id", resp.RequestId)
	return nil
}

func datasourceOutscaleKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleKeyPairsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// Attributes
			"key_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"key_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
