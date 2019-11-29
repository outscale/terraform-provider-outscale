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
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

// Creates a network interface in the specified subnet
func resourceOutscaleOAPINic() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINicCreate,
		Read:   resourceOutscaleOAPINicRead,
		Update: resourceOutscaleOAPINicUpdate,
		Delete: resourceOutscaleOAPINicDelete,
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
		"security_group_ids": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnet_id": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		// Attributes
		"link_public_ip": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"public_ip_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"link_public_ip_id": {
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
		"link_nic": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"link_nic_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"delete_on_vm_deletion": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_number": {
						Type:     schema.TypeString,
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
		"subregion_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"security_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"security_group_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"security_group_name": {
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
			Type:     schema.TypeSet,
			Computed: true,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"link_public_ip": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"public_ip_id": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"link_public_ip_id": {
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
					"private_dns_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"private_ip": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"is_primary": {
						Type:     schema.TypeBool,
						Computed: true,
						Optional: true,
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
		"is_source_dest_checked": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": tagsListOAPISchema(),
		"net_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

//Create OAPINic
func resourceOutscaleOAPINicCreate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OAPI

	request := oapi.CreateNicRequest{
		SubnetId: d.Get("subnet_id").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		m := v.([]interface{})
		a := make([]string, len(m))
		for k, v := range m {
			a[k] = v.(string)
		}
		request.SecurityGroupIds = a
	}

	if v, ok := d.GetOk("private_ips"); ok {
		request.PrivateIps = expandPrivateIPLight(v.(*schema.Set).List())
	}

	log.Printf("[DEBUG] Creating network interface")

	var resp *oapi.POST_CreateNicResponses
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		resp, err = conn.POST_CreateNic(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating NIC: %s", err)
	}

	d.SetId(resp.OK.Nic.NicId)

	if d.IsNewResource() {
		if err := setOAPITags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Set("tags", make([]map[string]interface{}, 0))
	d.Set("private_ip", make([]map[string]interface{}, 0))

	log.Printf("[INFO] ENI ID: %s", d.Id())

	return resourceOutscaleOAPINicRead(d, meta)

}

//Read OAPINic
func resourceOutscaleOAPINicRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OAPI
	dnir := oapi.ReadNicsRequest{
		Filters: oapi.FiltersNic{
			NicIds: []string{d.Id()},
		},
	}

	var describeResp *oapi.POST_ReadNicsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		describeResp, err = conn.POST_ReadNics(dnir)
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
	if len(describeResp.OK.Nics) != 1 {
		return fmt.Errorf("Unable to find ENI: %#v", describeResp.OK.Nics)
	}

	eni := describeResp.OK.Nics[0]
	d.Set("description", eni.Description)

	d.Set("subnet_id", eni.SubnetId)

	b := make(map[string]interface{})
	link := eni.LinkPublicIp
	b["public_ip_id"] = link.PublicIpId
	b["link_public_ip_id"] = link.LinkPublicIpId
	b["public_ip_account_id"] = link.PublicIpAccountId
	b["public_dns_name"] = link.PublicDnsName
	b["public_ip"] = link.PublicIp

	if err := d.Set("link_public_ip", b); err != nil {
		return err
	}

	//aa := make([]map[string]interface{}, 1)
	bb := make(map[string]interface{})
	att := eni.LinkNic
	bb["link_nic_id"] = att.LinkNicId
	bb["delete_on_vm_deletion"] = strconv.FormatBool(aws.BoolValue(att.DeleteOnVmDeletion))
	bb["device_number"] = strconv.FormatInt(att.DeviceNumber, 10)
	bb["vm_id"] = att.VmAccountId
	bb["vm_account_id"] = att.VmAccountId
	bb["state"] = att.State

	//aa[0] = bb
	// if err := d.Set("link_nic", aa); err != nil {
	// 	return err
	// }

	if err := d.Set("link_nic", bb); err != nil {
		return err
	}

	d.Set("subregion_name", eni.SubregionName)

	x := make([]map[string]interface{}, len(eni.SecurityGroups))
	for k, v := range eni.SecurityGroups {
		b := make(map[string]interface{})
		b["security_group_id"] = v.SecurityGroupId
		b["security_group_name"] = v.SecurityGroupName
		x[k] = b
	}
	if err := d.Set("security_groups", x); err != nil {
		return err
	}

	d.Set("mac_address", eni.MacAddress)
	d.Set("nic_id", eni.NicId)
	d.Set("account_id", eni.AccountId)
	d.Set("private_dns_name", eni.PrivateDnsName)
	//d.Set("private_ip", eni.)

	y := make([]map[string]interface{}, len(eni.PrivateIps))
	if eni.PrivateIps != nil {
		for k, v := range eni.PrivateIps {
			b := make(map[string]interface{})

			d := make(map[string]interface{})
			assoc := v.LinkPublicIp
			d["public_ip_id"] = assoc.PublicIpId
			d["link_public_ip_id"] = assoc.LinkPublicIpId
			d["public_ip_account_id"] = assoc.PublicIpAccountId
			d["public_dns_name"] = assoc.PublicDnsName
			d["public_ip"] = assoc.PublicIp

			b["link_public_ip"] = d
			b["private_dns_name"] = v.PrivateDnsName
			b["private_ip"] = v.PrivateIp
			b["is_primary"] = v.IsPrimary

			y[k] = b
		}
	}
	if err := d.Set("private_ips", y); err != nil {
		return err
	}

	d.Set("request_id", describeResp.OK.ResponseContext.RequestId)
	d.Set("is_source_dest_checked", eni.IsSourceDestChecked)
	d.Set("state", eni.State)
	d.Set("tags", tagsOAPIToMap(eni.Tags))
	d.Set("net_id", eni.NetId)

	return nil
}

//Delete OAPINic
func resourceOutscaleOAPINicDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OAPI

	log.Printf("[INFO] Deleting ENI: %s", d.Id())

	err := resourceOutscaleOAPINicDetach(d.Get("link_nic").(interface{}), meta, d.Id())
	if err != nil {
		return err
	}

	deleteEniOpts := oapi.DeleteNicRequest{
		NicId: d.Id(),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.POST_DeleteNic(deleteEniOpts)
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

func resourceOutscaleOAPINicDetach(oa interface{}, meta interface{}, eniID string) error {
	// if there was an old nic_link, remove it
	if oa != nil {
		oa := oa.(map[string]interface{})
		dr := oapi.UnlinkNicRequest{
			LinkNicId: oa["link_nic_id"].(string),
		}

		conn := meta.(*OutscaleClient).OAPI

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.POST_UnlinkNic(dr)
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

		log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", eniID)
		stateConf := &resource.StateChangeConf{
			Pending: []string{"true"},
			Target:  []string{"false"},
			Refresh: nicLinkRefreshFunc(conn, eniID),
			Timeout: 10 * time.Minute,
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for ENI (%s) to become dettached: %s", eniID, err)
		}
	}

	return nil
}

//Update OAPINic
func resourceOutscaleOAPINicUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OAPI
	d.Partial(true)

	if d.HasChange("link_nic") {
		oa, na := d.GetChange("link_nic")

		err := resourceOutscaleOAPINicDetach(oa.([]interface{}), meta, d.Id())
		if err != nil {
			return err
		}

		// if there is a new nic_link, attach it
		if na != nil && len(na.([]interface{})) > 0 {
			na := na.([]interface{})[0].(map[string]interface{})
			di := na["device_number"].(int)
			ar := oapi.LinkNicRequest{
				DeviceNumber: int64(di),
				VmId:         na["instance"].(string),
				NicId:        d.Id(),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.POST_LinkNic(ar)
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

		d.SetPartial("link_nic")
	}

	if d.HasChange("private_ips") {
		o, n := d.GetChange("private_ips")

		// Unassign old IP addresses
		if len(o.(*schema.Set).List()) != 0 {
			input := oapi.UnlinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: flattenPrivateIPLightToStringSlice(o.(*schema.Set).List()),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {

				_, err = conn.POST_UnlinkPrivateIps(input)
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
		if len(n.(*schema.Set).List()) != 0 {
			input := oapi.LinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: flattenPrivateIPLightToStringSlice(n.(*schema.Set).List()),
			}

			var err error
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {

				_, err = conn.POST_LinkPrivateIps(input)
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

	// Missing Sourcedestcheck
	// request := oapi.UpdateNicRequest{
	// 	NicId:           d.Id(),
	// 	SourceDestCheck: &fcu.AttributeBooleanValue{Value: aws.Bool(d.Get("is_source_dest_checked").(bool))},
	// }

	// _, err := conn.VM.ModifyNetworkInterfaceAttribute(request)

	// err := resource.Retry(5*time.Minute, func() *resource.RetryError {
	// 	var err error
	// 	_, err = conn.POST_UpdateNic(request)
	// 	if err != nil {
	// 		if strings.Contains(err.Error(), "RequestLimitExceeded:") {
	// 			return resource.RetryableError(err)
	// 		}
	// 		return resource.NonRetryableError(err)
	// 	}
	// 	return nil
	// })

	// if err != nil {
	// 	return fmt.Errorf("Failure updating ENI: %s", err)
	// }

	// d.SetPartial("is_source_dest_checked")

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
				input := oapi.LinkPrivateIpsRequest{
					NicId:                   d.Id(),
					SecondaryPrivateIpCount: int64(diff),
				}
				// _, err := conn.VM.AssignPrivateIpAddresses(input)

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					var err error
					_, err = conn.POST_LinkPrivateIps(input)
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
				input := oapi.UnlinkPrivateIpsRequest{
					NicId:      d.Id(),
					PrivateIps: expandStringValueList(prips[0:int(math.Abs(float64(diff)))]),
				}

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					var err error
					_, err = conn.POST_UnlinkPrivateIps(input)
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

	if d.HasChange("security_group_ids") {
		request := oapi.UpdateNicRequest{
			NicId:            d.Id(),
			SecurityGroupIds: expandStringValueList(d.Get("security_group_ids").([]interface{})),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err = conn.POST_UpdateNic(request)
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
		request := oapi.UpdateNicRequest{
			NicId:       d.Id(),
			Description: d.Get("description").(string),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.POST_UpdateNic(request)
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

	if err := setOAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)

	return resourceOutscaleOAPINicRead(d, meta)
}

func nicLinkRefreshFunc(conn *oapi.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		dnir := &oapi.ReadNicsRequest{
			Filters: oapi.FiltersNic{NicIds: []string{id}},
		}

		var describeResp *oapi.POST_ReadNicsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.POST_ReadNics(*dnir)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		var errString string

		if err != nil || describeResp.OK == nil {
			if err != nil {
				errString = err.Error()
			} else if describeResp.Code401 != nil {
				errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(describeResp.Code401))
			} else if describeResp.Code400 != nil {
				errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(describeResp.Code400))
			} else if describeResp.Code500 != nil {
				errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(describeResp.Code500))
			}
			log.Printf("[ERROR] Could not find network interface %s. %s", id, err)
			return nil, "", fmt.Errorf("Could not find network interface: %s", errString)

		}

		eni := describeResp.OK.Nics[0]
		//hasLink := strconv.FormatBool(&eni.LinkNic != nil || !reflect.DeepEqual(eni.LinkNic, oapi.LinkNic{}))
		hasLink := strconv.FormatBool(eni.LinkNic.LinkNicId != "")
		log.Printf("[DEBUG] ENI %s has attachment state %s", id, hasLink)
		return eni, hasLink, nil
	}
}

func networkInterfaceOAPIAttachmentRefreshFunc(conn *oapi.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		dnir := oapi.ReadNicsRequest{
			Filters: oapi.FiltersNic{
				NicIds: []string{id},
			},
		}

		var describeResp *oapi.POST_ReadNicsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.POST_ReadNics(dnir)
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

		eni := describeResp.OK.Nics[0]
		hasAttachment := strconv.FormatBool(eni.LinkNic.LinkNicId != "")
		log.Printf("[DEBUG] ENI %s has attachment state %s", id, hasAttachment)
		return eni, hasAttachment, nil
	}
}

func expandPrivateIPLight(pIPs []interface{}) []oapi.PrivateIpLight {
	privateIPs := make([]oapi.PrivateIpLight, 0)
	for _, v := range pIPs {
		privateIPMap := v.(map[string]interface{})
		privateIP := oapi.PrivateIpLight{
			IsPrimary: privateIPMap["is_primary"].(bool),
			PrivateIp: privateIPMap["private_ip"].(string),
		}
		privateIPs = append(privateIPs, privateIP)
	}
	return privateIPs
}

func flattenPrivateIPLightToStringSlice(pIPs []interface{}) []string {
	privateIPs := make([]string, 0)
	for _, v := range pIPs {
		privateIPMap := v.(map[string]interface{})
		privateIPs = append(privateIPs, privateIPMap["private_ip"].(string))
	}
	return privateIPs
}
