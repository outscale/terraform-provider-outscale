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

// Creates a network interface in the specified subnet
func dataSourceOutscaleOAPINics() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPINicRead,
		Schema: getDSOAPINicsSchema(),
	}
}

func getDSOAPINicsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//  This is attribute part for schema Nic
		"filter": dataSourceFiltersSchema(),
		"nic_id": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"nic": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip_address": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"subnet_id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"public_ip_link": {
						Type:     schema.TypeMap,
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
								"public_ip_account_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"public_dns_name": {
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
					"nic_link": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"nic_link_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"delete_on_vm_deletion": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"nic_sort_number": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"vm_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"vm_account_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"state": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"sub_region_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"firewall_rules_set": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"firewall_rules_set_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"firewall_rules_set_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"mac_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"nic_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"account_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"private_ip": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"public_ip_link": {
									Type:     schema.TypeMap,
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
											"public_ip_account_id": {
												Type:     schema.TypeString,
												Computed: true,
											},
											"public_dns_name": {
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
								"primary_ip": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"private_dns_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"private_ip_address": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"requester_managed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"activated_check": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tag": tagsSchemaComputed(),
					"lin": {
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

//Read Nic
func dataSourceOutscaleOAPINicRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	n, nok := d.GetOk("nic_id")

	if filtersOk == false && nok == false {
		return fmt.Errorf("filters, or owner must be assigned, or nat_gateway_id must be provided")
	}

	params := &fcu.DescribeNetworkInterfacesInput{}
	if filtersOk {
		params.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}
	if nok {
		params.NetworkInterfaceIds = expandStringList(n.([]interface{}))
	}

	var describeResp *fcu.DescribeNetworkInterfacesOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		describeResp, err = conn.VM.DescribeNetworkInterfaces(params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {

		return fmt.Errorf("Error describing Network Interfaces : %s", err)
	}

	if err != nil {
		return fmt.Errorf("Error retrieving ENI: %s", err)
	}
	if len(describeResp.NetworkInterfaces) < 1 {
		return fmt.Errorf("Unable to find ENI: %#v", describeResp.NetworkInterfaces)
	}

	nics := make([]map[string]interface{}, len(describeResp.NetworkInterfaces))

	for k, eni := range describeResp.NetworkInterfaces {
		nic := make(map[string]interface{})

		nic["description"] = aws.StringValue(eni.Description)
		nic["subnet_id"] = aws.StringValue(eni.SubnetId)

		b := make(map[string]interface{})
		if eni.Association != nil {
			b["reservation_id"] = aws.StringValue(eni.Association.AllocationId)
			b["link_id"] = aws.StringValue(eni.Association.AssociationId)
			b["public_ip_account_id"] = aws.StringValue(eni.Association.IpOwnerId)
			b["public_dns_name"] = aws.StringValue(eni.Association.PublicDnsName)
			b["public_ip"] = aws.StringValue(eni.Association.PublicIp)
		}
		nic["public_ip_link"] = b

		aa := make([]map[string]interface{}, 1)
		bb := make(map[string]interface{})
		if eni.Attachment != nil {
			bb["nic_link_id"] = aws.StringValue(eni.Attachment.AttachmentId)
			bb["delete_on_vm_deletion"] = aws.BoolValue(eni.Attachment.DeleteOnTermination)
			bb["vm_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
			bb["nic_sort_number"] = aws.Int64Value(eni.Attachment.DeviceIndex)
			bb["vm_account_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
			bb["state"] = aws.StringValue(eni.Attachment.Status)
		}
		aa[0] = bb
		nic["nic_link"] = aa
		nic["sub_region_name"] = aws.StringValue(eni.AvailabilityZone)

		x := make([]map[string]interface{}, len(eni.Groups))
		for k, v := range eni.Groups {
			b := make(map[string]interface{})
			b["firewall_rules_set_id"] = aws.StringValue(v.GroupId)
			b["firewall_rules_set_name"] = aws.StringValue(v.GroupName)
			x[k] = b
		}
		nic["firewall_rules_set"] = x
		nic["mac_address"] = aws.StringValue(eni.MacAddress)
		nic["nic_id"] = aws.StringValue(eni.NetworkInterfaceId)
		nic["account_id"] = aws.StringValue(eni.OwnerId)
		nic["private_dns_name"] = aws.StringValue(eni.PrivateDnsName)
		nic["private_ip_address"] = aws.StringValue(eni.PrivateIpAddress)

		y := make([]map[string]interface{}, len(eni.PrivateIpAddresses))
		if eni.PrivateIpAddresses != nil {
			for k, v := range eni.PrivateIpAddresses {
				b := make(map[string]interface{})

				d := make(map[string]interface{})
				if v.Association != nil {
					d["reservation_id"] = aws.StringValue(v.Association.AllocationId)
					d["link_id"] = aws.StringValue(v.Association.AssociationId)
					d["public_ip_account_id"] = aws.StringValue(v.Association.IpOwnerId)
					d["public_dns_name"] = aws.StringValue(v.Association.PublicDnsName)
					d["public_ip"] = aws.StringValue(v.Association.PublicIp)
				}
				b["public_ip_link"] = d
				b["primary_ip"] = aws.BoolValue(v.Primary)
				b["private_dns_name"] = aws.StringValue(v.PrivateDnsName)
				b["private_ip_address"] = aws.StringValue(v.PrivateIpAddress)

				y[k] = b
			}
		}
		nic["private_ip"] = y
		nic["requester_managed"] = aws.BoolValue(eni.RequesterManaged)
		nic["activated_check"] = aws.BoolValue(eni.SourceDestCheck)
		nic["state"] = aws.StringValue(eni.Status)
		nic["tag"] = tagsToMap(eni.TagSet)
		nic["lin"] = aws.StringValue(eni.VpcId)

		nics[k] = nic
	}

	d.Set("request_id", describeResp.RequestId)

	d.SetId(resource.UniqueId())

	return d.Set("nic", nics)
}
