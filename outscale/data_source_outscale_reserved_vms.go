package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIReservedVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIReservedVMSRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"reserved_vms_id": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sub_region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"offering_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reserved_vm": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sub_region_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"currency_code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_count": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenancy": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"offering_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"recurring_charge": &schema.Schema{
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
						"reserved_vms_id": &schema.Schema{
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

func dataSourceOutscaleOAPIReservedVMSRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	az, azok := d.GetOk("sub_region_name")
	ot, otok := d.GetOk("offering_type")
	ri, riok := d.GetOk("reserved_vms_id")
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
		r["sub_region_name"] = *v.AvailabilityZone
		r["currency_code"] = *v.CurrencyCode
		r["vm_count"] = *v.InstanceCount
		r["tenancy"] = *v.InstanceTenancy
		r["type"] = *v.InstanceType
		r["offering_type"] = *v.OfferingType
		r["product_type"] = *v.ProductDescription

		rcs := make([]map[string]interface{}, len(v.RecurringCharges))
		for k1, v1 := range v.RecurringCharges {
			rc := make(map[string]interface{})
			rc["frequency"] = v1.Frequency
			rcs[k1] = rc
		}

		r["recurring_charge"] = rcs
		r["reserved_vms_id"] = *v.ReservedInstancesId
		r["state"] = *v.State
		rsi[k] = r
	}

	d.Set("reserved_vm", rsi)
	d.Set("request_id", resp.RequestId)

	return nil
}
