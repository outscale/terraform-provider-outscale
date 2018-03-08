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

func dataSourceOutscalePublicIP() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscalePublicIPRead,
		Schema: getPublicIPDataSourceSchema(),
	}
}

func getPublicIPDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"filter": dataSourceFiltersSchema(),
		"allocation_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"association_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"domain": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
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
			Optional: true,
		},
	}
}

func dataSourceOutscalePublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeAddressesInput{}

	if id := d.Get("allocation_id"); id != "" {
		req.AllocationIds = []*string{aws.String(id.(string))}
	}
	if id := d.Get("public_ip"); id != "" {
		req.PublicIps = []*string{aws.String(id.(string))}
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var describeAddresses *fcu.DescribeAddressesOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		describeAddresses, err = conn.VM.DescribeAddressesRequest(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}

		return resource.NonRetryableError(err)
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

	if len(describeAddresses.Addresses) > 1 {
		return fmt.Errorf("multiple External IPs matched; use additional constraints to reduce matches to a single External IP")
	}

	address := describeAddresses.Addresses[0]

	fmt.Printf("[DEBUG] EIP read configuration: %+v", *address)

	if address.AssociationId != nil {
		d.Set("association_id", *address.AssociationId)
	} else {
		d.Set("association_id", "")
	}
	if address.InstanceId != nil {
		d.Set("instance_id", *address.InstanceId)
	} else {
		d.Set("instance_id", "")
	}
	if address.NetworkInterfaceId != nil {
		d.Set("network_interface_id", *address.NetworkInterfaceId)
	} else {
		d.Set("network_interface_id", "")
	}
	if address.NetworkInterfaceOwnerId != nil {
		d.Set("network_interface_owner_id", *address.NetworkInterfaceOwnerId)
	} else {
		d.Set("network_interface_owner_id", "")
	}
	if address.PrivateIpAddress != nil {
		d.Set("private_ip", *address.PrivateIpAddress)
	} else {
		d.Set("private_ip", "")
	}

	d.Set("allocation_id", *address.AllocationId)
	d.Set("public_ip", *address.PublicIp)
	d.Set("domain", *address.Domain)
	d.SetId(*address.AllocationId)

	return nil
}
