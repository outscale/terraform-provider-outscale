package outscale

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
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
			ForceNew: true,
		},
		"dry_run": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		// "private_ip_address": &schema.Schema{
		// 	Type:     schema.TypeString,
		// 	Optional: true,
		// 	Computed: true,
		// 	ForceNew: true,
		// },
		"security_group_id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"subnet_id": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},

		"association": {
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
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_index": {
						Type:     schema.TypeString,
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
		// "description": {
		// 	Type:     schema.TypeString,
		// 	Computed: true,
		// },

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

		"private_ip_address": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"private_ip_addresses_set": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"association": {
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
								"primary": {
									Type:     schema.TypeString,
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
						Type:     schema.TypeString,
						Computed: true,
					},
					"source_dest_check": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"subnet_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tag_set": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

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

	security_groups := d.Get("security_groups").(*schema.Set).List()
	if len(security_groups) != 0 {
		request.Groups = expandStringList(security_groups)
	}

	private_ips := d.Get("private_ips").(*schema.Set).List()
	if len(private_ips) != 0 {
		request.PrivateIpAddresses = expandPrivateIPAddresses(private_ips)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("private_ips_count"); ok {
		request.SecondaryPrivateIpAddressCount = aws.Int64(int64(v.(int)))
	}

	log.Printf("[DEBUG] Creating network interface")
	resp, err := conn.CreateNetworkInterface(request)
	if err != nil {
		return fmt.Errorf("Error creating ENI: %s", err)
	}

	d.SetId(*resp.NetworkInterface.NetworkInterfaceId)
	log.Printf("[INFO] ENI ID: %s", d.Id())
	return resourceAwsNetworkInterfaceUpdate(d, meta)
}

//Read Nic
func resourceOutscaleNicRead(d *schema.ResourceData, meta interface{}) error {
	//	conn := meta.(*OutscaleClient).FCU

	conn := meta.(*OutscaleClient).FCU
	describe_network_interfaces_request := &fcu.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{aws.String(d.Id())},
	}
	describeResp, err := conn.DescribeNetworkInterfaces(describe_network_interfaces_request)

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
	d.Set("subnet_id", eni.SubnetId)
	d.Set("private_ip", eni.PrivateIpAddress)
	d.Set("private_ips", flattenNetworkInterfacesPrivateIPAddresses(eni.PrivateIpAddresses))
	d.Set("security_groups", flattenGroupIdentifiers(eni.Groups))
	d.Set("source_dest_check", eni.SourceDestCheck)

	if eni.Description != nil {
		d.Set("description", eni.Description)
	}

	// Tags
	d.Set("tags", tagsToMap(eni.TagSet))

	if eni.Attachment != nil {
		attachment := []map[string]interface{}{flattenAttachment(eni.Attachment)}
		d.Set("attachment", attachment)
	} else {
		d.Set("attachment", nil)
	}

	return nil
}

//Delete Nic
func resourceOutscaleNicDelete(d *schema.ResourceData, meta interface{}) error {
	//	conn := meta.(*OutscaleClient).FCU
	conn := meta.(*OutscaleClient).FCU

	log.Printf("[INFO] Deleting ENI: %s", d.Id())

	detach_err := resourceAwsNetworkInterfaceDetach(d.Get("attachment").(*schema.Set), meta, d.Id())
	if detach_err != nil {
		return detach_err
	}

	deleteEniOpts := fcu.DeleteNetworkInterfaceInput{
		NetworkInterfaceId: aws.String(d.Id()),
	}
	if _, err := conn.DeleteNetworkInterface(&deleteEniOpts); err != nil {
		return fmt.Errorf("Error deleting ENI: %s", err)
	}

	return nil
}

//Update Nic
func resourceOutscaleNicUpdate(d *schema.ResourceData, meta interface{}) error {
	//	conn := meta.(*OutscaleClient).FCU
	conn := meta.(*OutscaleClient).FCU
	d.Partial(true)

	if d.HasChange("attachment") {
		oa, na := d.GetChange("attachment")

		detach_err := resourceAwsNetworkInterfaceDetach(oa.(*schema.Set), meta, d.Id())
		if detach_err != nil {
			return detach_err
		}

		// if there is a new attachment, attach it
		if na != nil && len(na.(*schema.Set).List()) > 0 {
			new_attachment := na.(*schema.Set).List()[0].(map[string]interface{})
			di := new_attachment["device_index"].(int)
			attach_request := &fcu.AttachNetworkInterfaceInput{
				DeviceIndex:        aws.Int64(int64(di)),
				InstanceId:         aws.String(new_attachment["instance"].(string)),
				NetworkInterfaceId: aws.String(d.Id()),
			}
			_, attach_err := conn.AttachNetworkInterface(attach_request)
			if attach_err != nil {
				return fmt.Errorf("Error attaching ENI: %s", attach_err)
			}
		}

		d.SetPartial("attachment")
	}

	if d.HasChange("private_ips") {
		o, n := d.GetChange("private_ips")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		// Unassign old IP addresses
		unassignIps := os.Difference(ns)
		if unassignIps.Len() != 0 {
			input := &fcu.UnassignPrivateIpAddressesInput{
				NetworkInterfaceId: aws.String(d.Id()),
				PrivateIpAddresses: expandStringList(unassignIps.List()),
			}
			_, err := conn.UnassignPrivateIpAddresses(input)
			if err != nil {
				return fmt.Errorf("Failure to unassign Private IPs: %s", err)
			}
		}

		// Assign new IP addresses
		assignIps := ns.Difference(os)
		if assignIps.Len() != 0 {
			input := &fcu.AssignPrivateIpAddressesInput{
				NetworkInterfaceId: aws.String(d.Id()),
				PrivateIpAddresses: expandStringList(assignIps.List()),
			}
			_, err := conn.AssignPrivateIpAddresses(input)
			if err != nil {
				return fmt.Errorf("Failure to assign Private IPs: %s", err)
			}
		}

		d.SetPartial("private_ips")
	}

	request := &fcu.ModifyNetworkInterfaceAttributeInput{
		NetworkInterfaceId: aws.String(d.Id()),
		SourceDestCheck:    &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("source_dest_check").(bool))},
	}

	_, err := conn.ModifyNetworkInterfaceAttribute(request)
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
				_, err := conn.AssignPrivateIpAddresses(input)
				if err != nil {
					return fmt.Errorf("Failure to assign Private IPs: %s", err)
				}
			}

			if diff < 0 {
				input := &fcu.UnassignPrivateIpAddressesInput{
					NetworkInterfaceId: aws.String(d.Id()),
					PrivateIpAddresses: expandStringList(private_ips_filtered[0:int(math.Abs(float64(diff)))]),
				}
				_, err := conn.UnassignPrivateIpAddresses(input)
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

		_, err := conn.ModifyNetworkInterfaceAttribute(request)
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

		_, err := conn.ModifyNetworkInterfaceAttribute(request)
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

	return resourceAwsNetworkInterfaceRead(d, meta)
}

func resourceAwsEniAttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["instance"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m["device_index"].(int)))
	return hashcode.String(buf.String())
}

func networkInterfaceAttachmentRefreshFunc(*fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		describe_network_interfaces_request := &fcu.DescribeNetworkInterfacesInput{
			NetworkInterfaceIds: []*string{aws.String(id)},
		}
		describeResp, err := conn.DescribeNetworkInterfaces(describe_network_interfaces_request)

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
