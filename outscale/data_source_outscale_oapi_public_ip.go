package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func dataSourceOutscaleOAPIPublicIP() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIPublicIPRead,
		Schema: getOAPIPublicIPDataSourceSchema(),
	}
}

func getOAPIPublicIPDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"filter": dataSourceFiltersSchema(),
		"public_ip_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"link_public_ip_id": {
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
		"request_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func dataSourceOutscaleOAPIPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := oapi.ReadPublicIpsRequest{
		Filters: oapi.FiltersPublicIp{},
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIDataSourcePublicIpsFilters(filters.(*schema.Set))
	}

	if id := d.Get("public_ip_id"); id != "" {
		req.Filters.PublicIpIds = []string{id.(string)}
	}
	if id := d.Get("public_ip"); id != "" {
		req.Filters.PublicIps = []string{id.(string)}
	}

	var describeAddresses *oapi.ReadPublicIpsResponse
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err := conn.POST_ReadPublicIps(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		describeAddresses = resp.OK
		return nil
	})

	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	// Verify Outscale returned our EIP
	if describeAddresses == nil || len(describeAddresses.PublicIps) == 0 {
		return fmt.Errorf("Unable to find EIP: %#v", describeAddresses.PublicIps)
	}

	if len(describeAddresses.PublicIps) > 1 {
		return fmt.Errorf("multiple External IPs matched; use additional constraints to reduce matches to a single External IP")
	}

	address := describeAddresses.PublicIps[0]

	fmt.Printf("[DEBUG] EIP read configuration: %+v", address)

	if address.LinkPublicIpId != "" {
		d.Set("link_public_ip_id", address.LinkPublicIpId)
	} else {
		d.Set("link_public_ip_id", "")
	}
	if address.VmId != "" {
		d.Set("vm_id", address.VmId)
	} else {
		d.Set("vm_id", "")
	}
	if address.NicId != "" {
		d.Set("nic_id", address.NicId)
	} else {
		d.Set("nic_id", "")
	}
	if address.NicAccountId != "" {
		d.Set("nic_account_id", address.NicAccountId)
	} else {
		d.Set("nic_account_id", "")
	}
	if address.PrivateIp != "" {
		d.Set("private_ip", address.PrivateIp)
	} else {
		d.Set("private_ip", "")
	}

	d.Set("request_id", describeAddresses.ResponseContext.RequestId)
	d.Set("public_ip_id", address.PublicIpId)

	d.Set("public_ip", address.PublicIp)
	//missing
	// if address.Placement == "vpc" {
	// 	d.SetId(address.ReservationId)
	// } else {
	// 	d.SetId(address.PublicIp)
	// }

	d.SetId(address.PublicIp)

	return d.Set("request_id", describeAddresses.ResponseContext.RequestId)
}

func buildOutscaleOAPIDataSourcePublicIpsFilters(set *schema.Set) oapi.FiltersPublicIp {
	var filters oapi.FiltersPublicIp
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "public_ip_ids":
			filters.PublicIpIds = filterValues
		case "link_ids":
			filters.LinkPublicIpIds = filterValues
		case "placements":
			filters.Placements = filterValues
		case "vm_ids":
			filters.VmIds = filterValues
		case "nic_ids":
			filters.NicIds = filterValues
		case "nic_account_ids":
			filters.NicAccountIds = filterValues
		case "private_ips":
			filters.PrivateIps = filterValues
		case "public_ips":
			filters.PublicIps = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
