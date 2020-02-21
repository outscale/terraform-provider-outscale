package outscale

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/openlyinc/pointy"
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

	conn := meta.(*OutscaleClient).OSCAPI

	request := oscgo.CreateNicRequest{
		SubnetId: d.Get("subnet_id").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		request.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		m := v.([]interface{})
		a := make([]string, len(m))
		for k, v := range m {
			a[k] = v.(string)
		}
		request.SetSecurityGroupIds(a)
	}

	if v, ok := d.GetOk("private_ips"); ok {
		request.SetPrivateIps(expandPrivateIPLight(v.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] Creating network interface")

	var resp oscgo.CreateNicResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.NicApi.CreateNic(context.Background(), &oscgo.CreateNicOpts{CreateNicRequest: optional.NewInterface(request)})
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

	d.SetId(resp.Nic.GetNicId())

	if d.IsNewResource() {
		if err := setOSCAPITags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	if err := d.Set("tags", make([]map[string]interface{}, 0)); err != nil {
		return err
	}
	if err := d.Set("private_ip", make([]map[string]interface{}, 0)); err != nil {
		return err
	}

	log.Printf("[INFO] ENI ID: %s", d.Id())

	return resourceOutscaleOAPINicRead(d, meta)

}

//Read OAPINic
func resourceOutscaleOAPINicRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OSCAPI
	dnir := oscgo.ReadNicsRequest{
		Filters: &oscgo.FiltersNic{
			NicIds: &[]string{d.Id()},
		},
	}

	var resp oscgo.ReadNicsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		resp, _, err = conn.NicApi.ReadNics(context.Background(), &oscgo.ReadNicsOpts{ReadNicsRequest: optional.NewInterface(dnir)})
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
	if len(resp.GetNics()) != 1 {
		return fmt.Errorf("Unable to find ENI: %#v", resp.GetNics())
	}

	eni := resp.GetNics()[0]
	if err := d.Set("description", eni.GetDescription()); err != nil {
		return err
	}

	if err := d.Set("subnet_id", eni.GetSubnetId()); err != nil {
		return err
	}

	b := make(map[string]interface{})
	link := eni.GetLinkPublicIp()
	b["public_ip_id"] = link.GetPublicIpId()
	b["link_public_ip_id"] = link.GetLinkPublicIpId()
	b["public_ip_account_id"] = link.GetPublicIpAccountId()
	b["public_dns_name"] = link.GetPublicDnsName()
	b["public_ip"] = link.GetPublicIp()

	if err := d.Set("link_public_ip", b); err != nil {
		return err
	}

	//aa := make([]map[string]interface{}, 1)
	bb := make(map[string]interface{})
	att := eni.GetLinkNic()
	bb["link_nic_id"] = att.GetLinkNicId()
	bb["delete_on_vm_deletion"] = strconv.FormatBool(att.GetDeleteOnVmDeletion())
	bb["device_number"] = fmt.Sprintf("%d", att.GetDeviceNumber())
	bb["vm_id"] = att.GetVmId()
	bb["vm_account_id"] = att.GetVmAccountId()
	bb["state"] = att.GetState()

	//aa[0] = bb
	// if err := d.Set("link_nic", aa); err != nil {
	// 	return err
	// }

	if err := d.Set("link_nic", bb); err != nil {
		return err
	}

	if err := d.Set("subregion_name", eni.GetSubregionName()); err != nil {
		return err
	}

	x := make([]map[string]interface{}, len(eni.GetSecurityGroups()))
	for k, v := range eni.GetSecurityGroups() {
		b := make(map[string]interface{})
		b["security_group_id"] = v.GetSecurityGroupId()
		b["security_group_name"] = v.GetSecurityGroupName()
		x[k] = b
	}
	if err := d.Set("security_groups", x); err != nil {
		return err
	}

	if err := d.Set("mac_address", eni.GetMacAddress()); err != nil {
		return err
	}
	if err := d.Set("nic_id", eni.GetNicId()); err != nil {
		return err
	}
	if err := d.Set("account_id", eni.GetAccountId()); err != nil {
		return err
	}
	if err := d.Set("private_dns_name", eni.GetPrivateDnsName()); err != nil {
		return err
	}
	//d.Set("private_ip", eni.)

	y := make([]map[string]interface{}, len(eni.GetPrivateIps()))
	if eni.PrivateIps != nil {
		for k, v := range eni.GetPrivateIps() {
			b := make(map[string]interface{})

			d := make(map[string]interface{})
			assoc := v.GetLinkPublicIp()
			d["public_ip_id"] = assoc.GetPublicIpId()
			d["link_public_ip_id"] = assoc.GetLinkPublicIpId()
			d["public_ip_account_id"] = assoc.GetPublicIpAccountId()
			d["public_dns_name"] = assoc.GetPublicDnsName()
			d["public_ip"] = assoc.GetPublicIp()

			b["link_public_ip"] = d
			b["private_dns_name"] = v.GetPrivateDnsName()
			b["private_ip"] = v.GetPrivateIp()
			b["is_primary"] = v.GetIsPrimary()

			y[k] = b
		}
	}
	if err := d.Set("private_ips", y); err != nil {
		return err
	}

	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}
	if err := d.Set("is_source_dest_checked", eni.GetIsSourceDestChecked()); err != nil {
		return err
	}
	if err := d.Set("state", eni.GetState()); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(eni.GetTags())); err != nil {
		return err
	}
	if err := d.Set("net_id", eni.GetNetId()); err != nil {
		return err
	}

	return nil
}

//Delete OAPINic
func resourceOutscaleOAPINicDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[INFO] Deleting ENI: %s", d.Id())

	err := resourceOutscaleOAPINicDetach(d.Get("link_nic").(interface{}), meta, d.Id())
	if err != nil {
		return err
	}

	deleteEniOpts := oscgo.DeleteNicRequest{
		NicId: d.Id(),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.NicApi.DeleteNic(context.Background(), &oscgo.DeleteNicOpts{DeleteNicRequest: optional.NewInterface(deleteEniOpts)})
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
		dr := oscgo.UnlinkNicRequest{
			LinkNicId: oa["link_nic_id"].(string),
		}

		conn := meta.(*OutscaleClient).OSCAPI

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err = conn.NicApi.UnlinkNic(context.Background(), &oscgo.UnlinkNicOpts{UnlinkNicRequest: optional.NewInterface(dr)})
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
	conn := meta.(*OutscaleClient).OSCAPI
	d.Partial(true)
	var err error

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
			ar := oscgo.LinkNicRequest{
				DeviceNumber: int64(di),
				VmId:         na["instance"].(string),
				NicId:        d.Id(),
			}

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.NicApi.LinkNic(context.Background(), &oscgo.LinkNicOpts{LinkNicRequest: optional.NewInterface(ar)})
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
			input := oscgo.UnlinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: flattenPrivateIPLightToStringSlice(o.(*schema.Set).List()),
			}

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.NicApi.UnlinkPrivateIps(context.Background(), &oscgo.UnlinkPrivateIpsOpts{UnlinkPrivateIpsRequest: optional.NewInterface(input)})
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
			stringSlice := flattenPrivateIPLightToStringSlice(n.(*schema.Set).List())
			input := oscgo.LinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: &stringSlice,
			}

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.NicApi.LinkPrivateIps(context.Background(), &oscgo.LinkPrivateIpsOpts{LinkPrivateIpsRequest: optional.NewInterface(input)})
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
				dif := int32(diff)
				input := oscgo.LinkPrivateIpsRequest{
					NicId:                   d.Id(),
					SecondaryPrivateIpCount: pointy.Int64(int64(dif)),
				}
				// _, err := conn.VM.AssignPrivateIpAddresses(input)

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					var err error
					_, _, err = conn.NicApi.LinkPrivateIps(context.Background(), &oscgo.LinkPrivateIpsOpts{LinkPrivateIpsRequest: optional.NewInterface(input)})
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
				input := oscgo.UnlinkPrivateIpsRequest{
					NicId:      d.Id(),
					PrivateIps: expandStringValueList(prips[0:int(math.Abs(float64(diff)))]),
				}

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, _, err = conn.NicApi.UnlinkPrivateIps(context.Background(), &oscgo.UnlinkPrivateIpsOpts{UnlinkPrivateIpsRequest: optional.NewInterface(input)})
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
		stringValueList := expandStringValueList(d.Get("security_group_ids").([]interface{}))
		request := oscgo.UpdateNicRequest{
			NicId:            d.Id(),
			SecurityGroupIds: &stringValueList,
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err = conn.NicApi.UpdateNic(context.Background(), &oscgo.UpdateNicOpts{UpdateNicRequest: optional.NewInterface(request)})
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
		request := oscgo.UpdateNicRequest{
			NicId:       d.Id(),
			Description: pointy.String(d.Get("description").(string)),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err = conn.NicApi.UpdateNic(context.Background(), &oscgo.UpdateNicOpts{UpdateNicRequest: optional.NewInterface(request)})
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

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)

	return resourceOutscaleOAPINicRead(d, meta)
}

func nicLinkRefreshFunc(conn *oscgo.APIClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		dnir := oscgo.ReadNicsRequest{
			Filters: &oscgo.FiltersNic{NicIds: &[]string{id}},
		}

		var resp oscgo.ReadNicsResponse
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			resp, _, err = conn.NicApi.ReadNics(context.Background(), &oscgo.ReadNicsOpts{ReadNicsRequest: optional.NewInterface(dnir)})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			errString := err.Error()
			log.Printf("[ERROR] Could not find network interface %s. %s", id, err)
			return nil, "", fmt.Errorf("Could not find network interface: %s", errString)

		}

		eni := resp.GetNics()[0]
		//hasLink := strconv.FormatBool(&eni.LinkNic != nil || !reflect.DeepEqual(eni.LinkNic, oscgo.LinkNic{}))
		hasLink := strconv.FormatBool(eni.LinkNic.GetLinkNicId() != "")
		log.Printf("[DEBUG] ENI %s has attachment state %s", id, hasLink)
		return eni, hasLink, nil
	}
}

func expandPrivateIPLight(pIPs []interface{}) []oscgo.PrivateIpLight {
	privateIPs := make([]oscgo.PrivateIpLight, 0)
	for _, v := range pIPs {
		privateIPMap := v.(map[string]interface{})
		isPrimary := privateIPMap["is_primary"].(bool)
		private := privateIPMap["private_ip"].(string)
		privateIP := oscgo.PrivateIpLight{
			IsPrimary: &isPrimary,
			PrivateIp: &private,
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
