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

func dataSourceOutscalePublicIPS() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscalePublicIPSRead,
		Schema: getPublicIPSDataSourceSchema(),
	}
}

func getPublicIPSDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"filter": dataSourceFiltersSchema(),
		"allocation_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"public_ips": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"addresses_set": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"allocation_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"association_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"domain": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"instance_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network_interface_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"network_interface_owner_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip_address": {
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

func dataSourceOutscalePublicIPSRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeAddressesInput{}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	if id := d.Get("allocation_id"); id != nil {
		var allocs []*string
		for _, v := range id.([]interface{}) {
			allocs = append(allocs, aws.String(v.(string)))
		}
		req.AllocationIds = allocs
	}
	if id := d.Get("public_ip"); id != nil {
		var ips []*string
		for _, v := range id.([]interface{}) {
			ips = append(ips, aws.String(v.(string)))
		}

		req.PublicIps = ips
	}

	var describeAddresses *fcu.DescribeAddressesOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		describeAddresses, err = conn.VM.DescribeAddressesRequest(req)

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
	if describeAddresses == nil || len(describeAddresses.Addresses) == 0 {
		return fmt.Errorf("Unable to find EIP: %#v", describeAddresses.Addresses)
	}

	addresses := describeAddresses.Addresses

	address := make([]map[string]interface{}, len(addresses))

	for k, v := range addresses {

		add := make(map[string]interface{})

		if v.AssociationId != nil {
			add["association_id"] = *v.AssociationId
		} else {
			add["association_id"] = ""
		}
		if v.InstanceId != nil {
			add["instance_id"] = *v.InstanceId
		} else {
			add["instance_id"] = ""
		}
		if v.NetworkInterfaceId != nil {
			add["network_interface_id"] = *v.NetworkInterfaceId
		} else {
			add["network_interface_id"] = ""
		}
		if v.NetworkInterfaceOwnerId != nil {
			add["network_interface_owner_id"] = *v.NetworkInterfaceOwnerId
		} else {
			add["network_interface_owner_id"] = ""
		}
		if v.PrivateIpAddress != nil {
			add["private_ip_address"] = *v.PrivateIpAddress
		} else {
			add["private_ip_address"] = ""
		}

		add["domain"] = *v.Domain
		add["allocation_id"] = *v.AllocationId
		add["public_ip"] = *v.PublicIp

		address[k] = add

		fmt.Printf("\n[DEBUG] ADD %s \n", v)
	}

	d.SetId(resource.UniqueId())

	d.Set("request_id", *describeAddresses.RequestId)

	err = d.Set("addresses_set", address)

	return err
}
