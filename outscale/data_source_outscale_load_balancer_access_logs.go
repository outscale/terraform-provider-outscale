package outscale

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleLoadBalancerAccessLogs() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleLoadBalancerAccessLogsRead,
		Schema: getDataSourceSchemas(attrLBAccessLogsSchema()),
	}
}

func attrLBAccessLogsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"publication_interval": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"is_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"osu_bucket_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"osu_bucket_prefix": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"load_balancer_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceOutscaleLoadBalancerAccessLogsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	resp, elbName, err := readLbs(conn, d)
	if err != nil {
		return err
	}
	lbs := *resp.LoadBalancers
	if len(lbs) != 1 {
		return fmt.Errorf("Unable to find LBU: %s", elbName)
	}

	lb := (lbs)[0]

	if lb.AccessLog == nil {
		return fmt.Errorf("NO Attributes FOUND")
	}

	//utils.PrintToJSON(resp, "RESPONSE =>")

	a := lb.AccessLog

	d.Set("publication_interval", a.PublicationInterval)
	d.Set("is_enabled", a.IsEnabled)
	d.Set("osu_bucket_name", a.OsuBucketName)
	d.Set("osu_bucket_prefix", a.OsuBucketPrefix)

	d.SetId(*elbName)
	d.Set("request_id", resp.ResponseContext.RequestId)

	return nil
}
