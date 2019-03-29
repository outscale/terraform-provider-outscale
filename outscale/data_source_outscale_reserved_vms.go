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

func dataSourceOutscaleReservedVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleReservedVMSRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"reserved_instances_id": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"offering_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reserved_instances_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"currency_code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_count": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_tenancy": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"offering_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"recurring_charges": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"reserved_instances_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleReservedVMSRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	az, azok := d.GetOk("availability_zone")
	ot, otok := d.GetOk("offering_type")
	ri, riok := d.GetOk("reserved_instances_id")
	filter, filterOk := d.GetOk("filter")

	req := &fcu.DescribeReservedInstancesInput{}

	if azok {
		req.AvailabilityZone = aws.String(az.(string))
	}
	if otok {
		req.OfferingType = aws.String(ot.(string))
	}
	if riok {
		var ids []*string
		for _, v := range ri.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		req.ReservedInstancesIds = ids
	}
	if filterOk {
		req.Filters = buildOutscaleDataSourceFilters(filter.(*schema.Set))
	}

	var resp *fcu.DescribeReservedInstancesOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeReservedInstances(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	if resp == nil || len(resp.ReservedInstances) == 0 {
		return fmt.Errorf("no matching reserved VMS found")
	}

	d.SetId(resource.UniqueId())

	rsi := make([]map[string]interface{}, len(resp.ReservedInstances))

	for k, v := range resp.ReservedInstances {
		r := make(map[string]interface{})
		r["availability_zone"] = *v.AvailabilityZone
		r["currency_code"] = *v.CurrencyCode
		r["instance_count"] = *v.InstanceCount
		r["instance_tenancy"] = *v.InstanceTenancy
		r["instance_type"] = *v.InstanceType
		r["offering_type"] = *v.OfferingType
		r["product_description"] = *v.ProductDescription

		rcs := make([]map[string]interface{}, len(v.RecurringCharges))
		for k1, v1 := range v.RecurringCharges {
			rc := make(map[string]interface{})
			rc["frequency"] = v1.Frequency
			rcs[k1] = rc
		}

		r["recurring_charges"] = rcs
		r["reserved_instances_id"] = *v.ReservedInstancesId
		r["state"] = *v.State
		rsi[k] = r
	}

	d.Set("reserved_instances_set", rsi)
	d.Set("request_id", resp.RequestId)

	return nil
}
