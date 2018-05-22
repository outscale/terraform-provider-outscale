package outscale

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

// Creates a network interface in the specified subnet
func resourceOutscaleOAPINic() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINicCreate,
		Read:   resourceOutscaleOAPINicRead,
		Delete: resourceOutscaleOAPINicDelete,
		Update: resourceOutscaleOAPINicUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: getOAPINicSchema(),
	}
}

func getOAPINicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//  This is attribute part for schema OAPINic
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"private_ip": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"firewall_rules_set_id": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		// Attributes
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

		"private_ips": {
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
					"pip": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"private_dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip": {
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
		"lin_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

//Create OAPINic
func resourceOutscaleOAPINicCreate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	request := &fcu.CreateNetworkInterfaceInput{
		SubnetId: aws.String(d.Get("subnet_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("firewall_rules_set_id"); ok {
		m := v.([]interface{})
		a := make([]*string, len(m))
		for k, v := range m {
			a[k] = aws.String(v.(string))
		}
		request.Groups = a
	}

	if v, ok := d.GetOk("private_ip"); ok {
		request.PrivateIpAddress = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating network interface")

	var resp *fcu.CreateNetworkInterfaceOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		resp, err = conn.VM.CreateNetworkInterface(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating ENI: %s", err)
	}

	d.SetId(*resp.NetworkInterface.NetworkInterfaceId)

	if d.IsNewResource() {
		if err := setTags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tag")
	}

	d.Set("tag", make([]map[string]interface{}, 0))
	d.Set("private_ip", make([]map[string]interface{}, 0))

	log.Printf("[INFO] ENI ID: %s", d.Id())

	return resourceOutscaleOAPINicRead(d, meta)

}

//Read OAPINic
func resourceOutscaleOAPINicRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU
	dnir := &fcu.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{aws.String(d.Id())},
	}

	var describeResp *fcu.DescribeNetworkInterfacesOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		describeResp, err = conn.VM.DescribeNetworkInterfaces(dnir)
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
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidNetworkInterfaceID.NotFound" {
			// The ENI is gone now, so just remove it from the state
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving ENI: %s", err)
	}
	if len(describeResp.NetworkInterfaces) != 1 {
		return fmt.Errorf("Unable to find ENI: %#v", describeResp.NetworkInterfaces)
	}

	eni := describeResp.NetworkInterfaces[0]
	if eni.Description != nil {
		d.Set("description", eni.Description)
	}
	d.Set("subnet_id", eni.SubnetId)

	b := make(map[string]interface{})
	if eni.Association != nil {
		b["reservation_id"] = aws.StringValue(eni.Association.AllocationId)
		b["link_id"] = aws.StringValue(eni.Association.AssociationId)
		b["public_ip_account_id"] = aws.StringValue(eni.Association.IpOwnerId)
		b["public_dns_name"] = aws.StringValue(eni.Association.PublicDnsName)
		b["public_ip"] = aws.StringValue(eni.Association.PublicIp)
	}
	if err := d.Set("public_ip_link", b); err != nil {
		return err
	}

	aa := make([]map[string]interface{}, 1)
	bb := make(map[string]interface{})
	if eni.Attachment != nil {
		bb["nic_link_id"] = aws.StringValue(eni.Attachment.AttachmentId)
		bb["delete_on_vm_deletion"] = aws.BoolValue(eni.Attachment.DeleteOnTermination)
		bb["nic_sort_number"] = aws.Int64Value(eni.Attachment.DeviceIndex)
		bb["vm_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
		bb["vm_account_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
		bb["state"] = aws.StringValue(eni.Attachment.Status)
	}
	aa[0] = bb
	if err := d.Set("nic_link", aa); err != nil {
		return err
	}

	d.Set("sub_region_name", aws.StringValue(eni.AvailabilityZone))

	x := make([]map[string]interface{}, len(eni.Groups))
	for k, v := range eni.Groups {
		b := make(map[string]interface{})
		b["firewall_rules_set_id"] = aws.StringValue(v.GroupId)
		b["firewall_rules_set_name"] = aws.StringValue(v.GroupName)
		x[k] = b
	}
	if err := d.Set("firewall_rules_set", x); err != nil {
		return err
	}

	d.Set("mac_address", aws.StringValue(eni.MacAddress))
	d.Set("nic_id", aws.StringValue(eni.NetworkInterfaceId))
	d.Set("account_id", aws.StringValue(eni.OwnerId))
	d.Set("private_dns_name", aws.StringValue(eni.PrivateDnsName))
	d.Set("private_ip", aws.StringValue(eni.PrivateIpAddress))

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
			b["pip"] = aws.BoolValue(v.Primary)
			b["private_dns_name"] = aws.StringValue(v.PrivateDnsName)
			b["private_ip"] = aws.StringValue(v.PrivateIpAddress)

			y[k] = b
		}
	}
	if err := d.Set("private_ip", y); err != nil {
		return err
	}

	d.Set("request_id", describeResp.RequestId)

	d.Set("requester_managed", aws.BoolValue(eni.RequesterManaged))

	d.Set("activated_check", aws.BoolValue(eni.SourceDestCheck))
	d.Set("state", aws.StringValue(eni.Status))
	// Tags
	d.Set("tags", tagsToMap(eni.TagSet))
	d.Set("lin_id", aws.StringValue(eni.VpcId))

	return nil
}

//Delete OAPINic
func resourceOutscaleOAPINicDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	log.Printf("[INFO] Deleting ENI: %s", d.Id())

	err := resourceOutscaleOAPINicDetach(d.Get("nic_link").([]interface{}), meta, d.Id())
	if err != nil {
		return err
	}

	deleteEniOpts := fcu.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(d.Id()),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DeleteNetworkInterface(&deleteEniOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {

		return fmt.Errorf("Error Deleting ENI: %s", err)
	}

	return nil
}

func resourceOutscaleOAPINicDetach(oa []interface{}, meta interface{}, eniId string) error {
	// if there was an old nic_link, remove it
	if oa != nil && len(oa) > 0 && oa[0] != nil {
		oa := oa[0].(map[string]interface{})
		dr := &fcu.DetachNetworkInterfaceInput{
			AttachmentId: aws.String(oa["nic_link_id"].(string)),
			Force:        aws.Bool(true),
		}
		conn := meta.(*OutscaleClient).FCU

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.VM.DetachNetworkInterface(dr)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
				return fmt.Errorf("Error detaching ENI: %s", err)
			}
		}

		log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", eniId)
		stateConf := &resource.StateChangeConf{
			Pending: []string{"true"},
			Target:  []string{"false"},
			Refresh: networkInterfaceAttachmentRefreshFunc(conn, eniId),
			Timeout: 10 * time.Minute,
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for ENI (%s) to become dettached: %s", eniId, err)
		}
	}

	return nil
}

//Update OAPINic
func resourceOutscaleOAPINicUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU
	d.Partial(true)

	if d.HasChange("nic_link") {
		oa, na := d.GetChange("nic_link")

		err := resourceOutscaleOAPINicDetach(oa.([]interface{}), meta, d.Id())
		if err != nil {
			return err
		}

		// if there is a new nic_link, attach it
		if na != nil && len(na.([]interface{})) > 0 {
			na := na.([]interface{})[0].(map[string]interface{})
			di := na["nic_sort_number"].(int)
			ar := &fcu.AttachNetworkInterfaceInput{
				DeviceIndex:        aws.Int64(int64(di)),
				InstanceId:         aws.String(na["instance"].(string)),
				NetworkInterfaceId: aws.String(d.Id()),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.VM.AttachNetworkInterface(ar)
				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {

				return fmt.Errorf("Error Attaching Network Interface: %s", err)
			}
		}

		d.SetPartial("nic_link")
	}

	if d.HasChange("private_ip") {
		o, n := d.GetChange("private_ip")
		if o == nil {
			o = new([]interface{})
		}
		if n == nil {
			n = new([]interface{})
		}

		// Unassign old IP addresses
		if len(o.([]interface{})) != 0 {
			input := &fcu.UnassignPrivateIpAddressesInput{
				NetworkInterfaceId: aws.String(d.Id()),
				PrivateIpAddresses: expandStringList(o.([]interface{})),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {

				_, err = conn.VM.UnassignPrivateIpAddresses(input)
				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure to unassign Private IPs: %s", err)
			}
		}

		// Assign new IP addresses
		if len(n.([]interface{})) != 0 {
			input := &fcu.AssignPrivateIpAddressesInput{
				NetworkInterfaceId: aws.String(d.Id()),
				PrivateIpAddresses: expandStringList(n.([]interface{})),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {

				_, err = conn.VM.AssignPrivateIpAddresses(input)
				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("Failure to assign Private IPs: %s", err)
			}
		}

		d.SetPartial("private_ip")
	}

	request := &fcu.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(d.Id()),
		SourceDestCheck:    &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("activated_check").(bool))},
	}

	// _, err := conn.VM.ModifyNetworkInterfaceAttribute(request)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, err = conn.VM.ModifyNetworkInterfaceAttribute(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure updating ENI: %s", err)
	}

	d.SetPartial("activated_check")

	if d.HasChange("private_ips_count") {
		o, n := d.GetChange("private_ips_count")
		pips := d.Get("pips").(*schema.Set).List()
		prips := pips[:0]
		pip := d.Get("private_ip")

		for _, ip := range pips {
			if ip != pip {
				prips = append(prips, ip)
			}
		}

		if o != nil && o != 0 && n != nil && n != len(prips) {

			diff := n.(int) - o.(int)

			// Surplus of IPs, add the diff
			if diff > 0 {
				input := &fcu.AssignPrivateIpAddressesInput{
					NetworkInterfaceId:             aws.String(d.Id()),
					SecondaryPrivateIpAddressCount: aws.Int64(int64(diff)),
				}
				// _, err := conn.VM.AssignPrivateIpAddresses(input)

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					var err error
					_, err = conn.VM.AssignPrivateIpAddresses(input)
					if err != nil {
						if strings.Contains(err.Error(), "RequestLimitExceeded:") {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					return fmt.Errorf("Failure to assign Private IPs: %s", err)
				}
			}

			if diff < 0 {
				input := &fcu.UnassignPrivateIpAddressesInput{
					NetworkInterfaceId: aws.String(d.Id()),
					PrivateIpAddresses: expandStringList(prips[0:int(math.Abs(float64(diff)))]),
				}

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					var err error
					_, err = conn.VM.UnassignPrivateIpAddresses(input)
					if err != nil {
						if strings.Contains(err.Error(), "RequestLimitExceeded:") {
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})

				if err != nil {
					return fmt.Errorf("Failure to unassign Private IPs: %s", err)
				}
			}

			d.SetPartial("private_ips_count")
		}
	}

	if d.HasChange("security_groups") {
		request := &fcu.ModifyNetworkInterfaceAttributeInput{
			NetworkInterfaceId: aws.String(d.Id()),
			Groups:             expandStringList(d.Get("security_groups").(*schema.Set).List()),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.VM.ModifyNetworkInterfaceAttribute(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Failure updating ENI: %s", err)
		}

		d.SetPartial("security_groups")
	}

	if d.HasChange("description") {
		request := &fcu.ModifyNetworkInterfaceAttributeInput{
			NetworkInterfaceId: aws.String(d.Id()),
			Description:        &fcu.AttributeValue{Value: aws.String(d.Get("description").(string))},
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.VM.ModifyNetworkInterfaceAttribute(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Failure updating ENI: %s", err)
		}

		d.SetPartial("description")
	}

	if err := setTags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tag")

	d.Partial(false)

	return resourceOutscaleOAPINicRead(d, meta)
}

func networkOAPIInterfaceAttachmentRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		dnir := &fcu.DescribeNetworkInterfacesInput{
			NetworkInterfaceIds: []*string{aws.String(id)},
		}

		var describeResp *fcu.DescribeNetworkInterfacesOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.VM.DescribeNetworkInterfaces(dnir)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			log.Printf("[ERROR] Could not find network interface %s. %s", id, err)
			return nil, "", err
		}

		eni := describeResp.NetworkInterfaces[0]
		hasAttachment := strconv.FormatBool(eni.Attachment != nil)
		log.Printf("[DEBUG] ENI %s has nic_link state %s", id, hasAttachment)
		return eni, hasAttachment, nil
	}
}
