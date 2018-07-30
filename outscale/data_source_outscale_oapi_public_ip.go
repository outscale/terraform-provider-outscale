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
		"reservation_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"link_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"placement": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
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
		Filters: oapi.ReadPublicIpsFilters{},
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleOAPIDataSourcePublicIpsFilters(filters.(*schema.Set))
	}

	if id := d.Get("reservation_id"); id != "" {
		req.Filters.ReservationIds = []string{id.(string)}
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

	if address.LinkId != "" {
		d.Set("link_id", address.LinkId)
	} else {
		d.Set("link_id", "")
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
	d.Set("reservation_id", address.ReservationId)
	d.Set("public_ip", address.PublicIp)
	d.Set("placement", address.Placement)

	if address.Placement == "vpc" {
		d.SetId(address.ReservationId)
	} else {
		d.SetId(address.PublicIp)
	}

	return d.Set("request_id", describeAddresses.ResponseContext.RequestId)
}

func buildOutscaleOAPIDataSourcePublicIpsFilters(set *schema.Set) oapi.ReadPublicIpsFilters {
	var filters oapi.ReadPublicIpsFilters
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "reservation-ids":
			filters.ReservationIds = filterValues
		case "link-ids":
			filters.LinkIds = filterValues
		case "placements":
			filters.Placements = filterValues
		case "vm-ids":
			filters.VmIds = filterValues
		case "nic-ids":
			filters.NicIds = filterValues
		case "nic-account-ids":
			filters.NicAccountIds = filterValues
		case "private-ips":
			filters.PrivateIps = filterValues
		case "public-ips":
			filters.PublicIps = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}
