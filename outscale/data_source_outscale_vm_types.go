package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleVMTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleVMTypesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"instance_type_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ebs_optimized_available": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_ip_addresses": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_count": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_size": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vcpu": &schema.Schema{
							Type:     schema.TypeInt,
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

func dataSourceOutscaleVMTypesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filter, filterOk := d.GetOk("filter")

	req := &fcu.DescribeInstanceTypesInput{}

	if filterOk {
		req.Filters = buildOutscaleDataSourceFilters(filter.(*schema.Set))
	}

	log.Printf("[DEBUG] DescribeVMTypes %+v\n", req)

	var resp *fcu.DescribeInstanceTypesOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeInstanceTypes(req)
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

	if resp == nil || len(resp.InstanceTypeSet) == 0 {
		return fmt.Errorf("no matching regions found")
	}

	vms := make([]map[string]interface{}, len(resp.InstanceTypeSet))

	for k, v := range resp.InstanceTypeSet {
		vm := make(map[string]interface{})
		vm["ebs_optimized_available"] = *v.EbsOptimizedAvailable
		vm["max_ip_addresses"] = *v.MaxIpAddresses
		vm["memory"] = *v.Memory
		vm["name"] = *v.Name
		vm["storage_count"] = *v.StorageCount
		if v.StorageSize != nil {
			vm["storage_size"] = *v.StorageSize
		} else {
			vm["storage_size"] = 0
		}
		vm["vcpu"] = *v.Vcpu
		vms[k] = vm
	}

	if err := d.Set("instance_type_set", vms); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())
	d.Set("request_id", resp.RequestId)

	return nil
}
