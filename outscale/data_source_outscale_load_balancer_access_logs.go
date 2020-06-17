package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOutscaleLoadBalancerAccessLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleLoadBalancerAccessLogsRead,

		Schema: map[string]*schema.Schema{
			"emit_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"s3_bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"s3_bucket_prefix": {
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
		},
	}
}

func dataSourceOutscaleLoadBalancerAccessLogsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	elbName, ok1 := d.GetOk("load_balancer_name")

	if !ok1 {
		return fmt.Errorf("please provide the load_balancer_name required attribute")
	}

	filter := &oscgo.FiltersLoadBalancer{
		LoadBalancerNames: &[]string{elbName.(string)},
	}

	req := oscgo.ReadLoadBalancersRequest{
		Filters: filter,
	}

	describeElbOpts := &oscgo.ReadLoadBalancersOpts{
		ReadLoadBalancersRequest: optional.NewInterface(req),
	}

	var resp oscgo.ReadLoadBalancersResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.LoadBalancerApi.ReadLoadBalancers(
			context.Background(),
			describeElbOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if isLoadBalancerNotFound(err) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving LBU Attr: %s", err)
	}

	lbs := *resp.LoadBalancers
	if len(lbs) != 1 {
		return fmt.Errorf("Unable to find LBU: %s", elbName.(string))
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

	d.SetId(elbName.(string))
	d.Set("request_id", resp.ResponseContext.RequestId)

	return nil
}
