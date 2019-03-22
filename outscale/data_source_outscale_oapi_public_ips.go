package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIPublicIPS() *schema.Resource {
	return &schema.Resource{
		Read:   oapiDataSourceOutscalePublicIPSRead,
		Schema: oapiGetPublicIPSDataSourceSchema(),
	}
}

func oapiGetPublicIPSDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"filter": dataSourceFiltersSchema(),
		"public_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"reservation_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"link_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"placement": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vm_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"nic_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"nic_account_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"public_ip": {
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
	}
}

func oapiDataSourceOutscalePublicIPSRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := oapi.ReadPublicIpsRequest{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIDataSourcePublicIpsFilters(filters.(*schema.Set))
	}

	//if id := d.Get("reservation_id"); id != nil {
	//	var allocs []*string
	//	for _, v := range id.([]interface{}) {
	//		allocs = append(allocs, aws.String(v.(string)))
	//	}
	//	req.Filters.AllocationIds = allocs
	//}
	//if id := d.Get("public_ip"); id != nil {
	//	var ips []string
	//	for _, v := range id.([]interface{}) {
	//		ips = append(ips, v.(string))
	//	}

	//	req.Filters.PublicIps = ips
	//}

	var describeAddresses *oapi.POST_ReadPublicIpsResponses
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		describeAddresses, err = conn.POST_ReadPublicIps(req)
		return resource.RetryableError(err)
	})

	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	// Verify Outscale returned our EIP
	if describeAddresses == nil || len(describeAddresses.OK.PublicIps) == 0 {
		return fmt.Errorf("Unable to find EIP: %#v", describeAddresses.OK.PublicIps)
	}

	addresses := describeAddresses.OK.PublicIps

	address := make([]map[string]interface{}, len(addresses))

	for k, v := range addresses {

		add := make(map[string]interface{})

		add["link_id"] = v.LinkPublicIpId
		add["vm_id"] = v.VmId
		add["nic_id"] = v.NicId
		add["nic_account_id"] = v.NicAccountId
		add["private_ip"] = v.PrivateIp
		add["placement"] = ""
		add["reservation_id"] = ""
		add["public_ip"] = v.PublicIp

		address[k] = add
	}

	d.SetId(resource.UniqueId())

	d.Set("request_id", describeAddresses.OK.ResponseContext.RequestId)

	err = d.Set("public_ips", address)

	return err
}
