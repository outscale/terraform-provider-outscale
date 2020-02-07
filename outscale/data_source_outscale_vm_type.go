package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func dataSourceOutscaleOAPIVMType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVMTypeRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"bsu_optimized": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_private_ip": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_size": &schema.Schema{
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
			"vcore_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIVMTypeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	filter, filterOk := d.GetOk("filter")

	req := &fcu.DescribeInstanceTypesInput{}

	if filterOk {
		req.Filters = buildOutscaleDataSourceFilters(filter.(*schema.Set))
	}

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
	if len(resp.InstanceTypeSet) > 1 {
		return fmt.Errorf("multiple vm types matched; use additional constraints to reduce matches to a single vm type")
	}

	vm := resp.InstanceTypeSet[0]

	d.SetId(*vm.Name)
	d.Set("bsu_optimized", *vm.EbsOptimizedAvailable)
	d.Set("max_private_ip", *vm.MaxIpAddresses)
	d.Set("memory_size", *vm.Memory)
	d.Set("name", *vm.Name)
	d.Set("storage_count", *vm.StorageCount)
	if vm.StorageSize != nil {
		d.Set("storage_size", *vm.StorageSize)
	} else {
		d.Set("storage_size", 0)
	}
	d.Set("vcore_count", *vm.Vcpu)
	d.Set("request_id", resp.RequestId)

	return nil
}
