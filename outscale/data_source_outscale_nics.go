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
func dataSourceOutscaleNics() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleNicsRead,
		Schema: getDSNicsSchema(),
	}
}

func getDSNicsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//  This is attribute part for schema Nic
		"filter": dataSourceFiltersSchema(),
		"network_interface_id": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"network_interface_set": {
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
					"association": {
						Type:     schema.TypeMap,
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
								"ip_owner_id": {
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
					"attachment": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"attachment_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"delete_on_termination": {
									Type:     schema.TypeBool,
									Computed: true,
								},
								"device_index": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"instance_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"instance_owner_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"status": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},

					"availability_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"group_set": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"group_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"group_name": {
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
					"network_interface_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"owner_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"private_ip_addresses_set": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"association": {
									Type:     schema.TypeMap,
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
											"ip_owner_id": {
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
								"primary": {
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
					"requester_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"requester_managed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"source_dest_check": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tag_set": tagsSchemaComputed(),
					"vpc_id": {
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
func dataSourceOutscaleNicsRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	filters, filtersOk := d.GetOk("filter")
	n, nok := d.GetOk("network_interface_id")

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
			b["allocation_id"] = aws.StringValue(eni.Association.AllocationId)
			b["association_id"] = aws.StringValue(eni.Association.AssociationId)
			b["ip_owner_id"] = aws.StringValue(eni.Association.IpOwnerId)
			b["public_dns_name"] = aws.StringValue(eni.Association.PublicDnsName)
			b["public_ip"] = aws.StringValue(eni.Association.PublicIp)
		}
		nic["association"] = b

		attach := make(map[string]interface{})
		if eni.Attachment != nil {
			attach["attachment_id"] = aws.StringValue(eni.Attachment.AttachmentId)
			attach["delete_on_termination"] = aws.BoolValue(eni.Attachment.DeleteOnTermination)
			attach["device_index"] = aws.Int64Value(eni.Attachment.DeviceIndex)
			attach["instance_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
			attach["instance_owner_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
			attach["status"] = aws.StringValue(eni.Attachment.Status)
		}
		nic["attachment"] = attach
		nic["availability_zone"] = aws.StringValue(eni.AvailabilityZone)

		x := make([]map[string]interface{}, len(eni.Groups))
		for k, v := range eni.Groups {
			b := make(map[string]interface{})
			b["group_id"] = aws.StringValue(v.GroupId)
			b["group_name"] = aws.StringValue(v.GroupName)
			x[k] = b
		}
		nic["group_set"] = x
		nic["mac_address"] = aws.StringValue(eni.MacAddress)
		nic["network_interface_id"] = aws.StringValue(eni.NetworkInterfaceId)
		nic["owner_id"] = aws.StringValue(eni.OwnerId)
		nic["private_dns_name"] = aws.StringValue(eni.PrivateDnsName)
		nic["private_ip_address"] = aws.StringValue(eni.PrivateIpAddress)

		y := make([]map[string]interface{}, len(eni.PrivateIpAddresses))
		if eni.PrivateIpAddresses != nil {
			for k, v := range eni.PrivateIpAddresses {
				b := make(map[string]interface{})

				d := make(map[string]interface{})
				if v.Association != nil {
					d["allocation_id"] = aws.StringValue(v.Association.AllocationId)
					d["association_id"] = aws.StringValue(v.Association.AssociationId)
					d["ip_owner_id"] = aws.StringValue(v.Association.IpOwnerId)
					d["public_dns_name"] = aws.StringValue(v.Association.PublicDnsName)
					d["public_ip"] = aws.StringValue(v.Association.PublicIp)
				}
				b["association"] = d
				b["primary"] = aws.BoolValue(v.Primary)
				b["private_dns_name"] = aws.StringValue(v.PrivateDnsName)
				b["private_ip_address"] = aws.StringValue(v.PrivateIpAddress)

				y[k] = b
			}
		}
		nic["private_ip_addresses_set"] = y
		nic["requester_id"] = aws.StringValue(eni.RequesterId)
		nic["requester_managed"] = aws.BoolValue(eni.RequesterManaged)
		nic["source_dest_check"] = aws.BoolValue(eni.SourceDestCheck)
		nic["status"] = aws.StringValue(eni.Status)
		nic["tag_set"] = tagsToMap(eni.TagSet)
		nic["vpc_id"] = aws.StringValue(eni.VpcId)

		nics[k] = nic
	}

	d.Set("request_id", describeResp.RequestId)

	d.SetId(resource.UniqueId())

	return d.Set("network_interface_set", nics)
}
