package outscale

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

// Creates a network interface in the specified subnet
func resourceOutscaleNic() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleNicCreate,
		Read:   resourceOutscaleNicRead,
		Delete: resourceOutscaleNicDelete,
		Update: resourceOutscaleNicUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: getNicSchema(),
	}
}

func getNicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//  This is attribute part for schema Nic
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"private_ip_adress": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"security_group_id": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		// Attributes
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
			Type:     schema.TypeList,
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

		"private_ip_address_set": {
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
		"request_id": {
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
	}
}

//Create Nic
func resourceOutscaleNicCreate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	request := &fcu.CreateNetworkInterfaceInput{
		SubnetId: aws.String(d.Get("subnet_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		m := v.([]interface{})
		a := make([]*string, len(m))
		for k, v := range m {
			a[k] = aws.String(v.(string))
		}
		request.Groups = a
	}

	if v, ok := d.GetOk("private_ip_adress"); ok {
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
		d.SetPartial("tag_set")
	}

	d.Set("tag_set", make([]map[string]interface{}, 0))
	d.Set("private_ip_address_set", make([]map[string]interface{}, 0))

	log.Printf("[INFO] ENI ID: %s", d.Id())

	return resourceOutscaleNicRead(d, meta)

}

//Read Nic
func resourceOutscaleNicRead(d *schema.ResourceData, meta interface{}) error {

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
		b["allocation_id"] = aws.StringValue(eni.Association.AllocationId)
		b["association_id"] = aws.StringValue(eni.Association.AssociationId)
		b["ip_owner_id"] = aws.StringValue(eni.Association.IpOwnerId)
		b["public_dns_name"] = aws.StringValue(eni.Association.PublicDnsName)
		b["public_ip"] = aws.StringValue(eni.Association.PublicIp)
	}
	if err := d.Set("association", b); err != nil {
		return err
	}

	aa := make([]map[string]interface{}, 1)
	bb := make(map[string]interface{})
	if eni.Attachment != nil {
		bb["attachment_id"] = aws.StringValue(eni.Attachment.AttachmentId)
		bb["delete_on_termination"] = aws.BoolValue(eni.Attachment.DeleteOnTermination)
		bb["device_index"] = aws.Int64Value(eni.Attachment.DeviceIndex)
		bb["instance_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
		bb["instance_owner_id"] = aws.StringValue(eni.Attachment.InstanceOwnerId)
		bb["status"] = aws.StringValue(eni.Attachment.Status)
	}
	aa[0] = bb
	if err := d.Set("attachment", aa); err != nil {
		return err
	}

	d.Set("availability_zone", aws.StringValue(eni.AvailabilityZone))

	x := make([]map[string]interface{}, len(eni.Groups))
	for k, v := range eni.Groups {
		b := make(map[string]interface{})
		b["group_id"] = aws.StringValue(v.GroupId)
		b["group_name"] = aws.StringValue(v.GroupName)
		x[k] = b
	}
	if err := d.Set("group_set", x); err != nil {
		return err
	}

	d.Set("mac_address", aws.StringValue(eni.MacAddress))
	d.Set("network_interface_id", aws.StringValue(eni.NetworkInterfaceId))
	d.Set("owner_id", aws.StringValue(eni.OwnerId))
	d.Set("private_dns_name", aws.StringValue(eni.PrivateDnsName))
	d.Set("private_ip_address", aws.StringValue(eni.PrivateIpAddress))

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
	if err := d.Set("private_ip_address_set", y); err != nil {
		return err
	}

	d.Set("request_id", describeResp.RequestId)

	d.Set("requester_managed", aws.BoolValue(eni.RequesterManaged))

	d.Set("source_dest_check", aws.BoolValue(eni.SourceDestCheck))
	d.Set("status", aws.StringValue(eni.Status))
	// Tags
	d.Set("tags", tagsToMap(eni.TagSet))
	d.Set("vpc_id", aws.StringValue(eni.VpcId))

	return nil
}

//Delete Nic
func resourceOutscaleNicDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	log.Printf("[INFO] Deleting ENI: %s", d.Id())

	detach_err := resourceOutscaleNicDetach(d.Get("attachment").([]interface{}), meta, d.Id())
	if detach_err != nil {
		return detach_err
	}

	deleteEniOpts := fcu.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(d.Id()),
	}

	var err error
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

func resourceOutscaleNicDetach(oa []interface{}, meta interface{}, eniId string) error {
	// if there was an old attachment, remove it
	if oa != nil && len(oa) > 0 && oa[0] != nil {
		old_attachment := oa[0].(map[string]interface{})
		detach_request := &fcu.DetachNetworkInterfaceInput{
			AttachmentId: aws.String(old_attachment["attachment_id"].(string)),
			Force:        aws.Bool(true),
		}
		conn := meta.(*OutscaleClient).FCU

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.VM.DetachNetworkInterface(detach_request)
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

//Update Nic
func resourceOutscaleNicUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU
	d.Partial(true)

	if d.HasChange("attachment") {
		oa, na := d.GetChange("attachment")

		detach_err := resourceOutscaleNicDetach(oa.([]interface{}), meta, d.Id())
		if detach_err != nil {
			return detach_err
		}

		// if there is a new attachment, attach it
		if na != nil && len(na.([]interface{})) > 0 {
			new_attachment := na.([]interface{})[0].(map[string]interface{})
			di := new_attachment["device_index"].(int)
			attach_request := &fcu.AttachNetworkInterfaceInput{
				DeviceIndex:        aws.Int64(int64(di)),
				InstanceId:         aws.String(new_attachment["instance"].(string)),
				NetworkInterfaceId: aws.String(d.Id()),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.VM.AttachNetworkInterface(attach_request)
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

		d.SetPartial("attachment")
	}

	if d.HasChange("private_ip_address_set") {
		o, n := d.GetChange("private_ip_address_set")
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

		d.SetPartial("private_ip_address_set")
	}

	request := &fcu.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(d.Id()),
		SourceDestCheck:    &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("source_dest_check").(bool))},
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

	d.SetPartial("source_dest_check")

	if d.HasChange("private_ips_count") {
		o, n := d.GetChange("private_ips_count")
		private_ips := d.Get("private_ips").(*schema.Set).List()
		private_ips_filtered := private_ips[:0]
		primary_ip := d.Get("private_ip")

		for _, ip := range private_ips {
			if ip != primary_ip {
				private_ips_filtered = append(private_ips_filtered, ip)
			}
		}

		if o != nil && o != 0 && n != nil && n != len(private_ips_filtered) {

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
					PrivateIpAddresses: expandStringList(private_ips_filtered[0:int(math.Abs(float64(diff)))]),
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
	} else {
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceOutscaleNicRead(d, meta)
}

func resourceOutscaleEniAttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["instance"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m["device_index"].(int)))
	return hashcode.String(buf.String())
}

func networkInterfaceAttachmentRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
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
		log.Printf("[DEBUG] ENI %s has attachment state %s", id, hasAttachment)
		return eni, hasAttachment, nil
	}
}
