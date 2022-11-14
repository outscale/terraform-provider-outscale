package outscale

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/openlyinc/pointy"
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
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"security_group_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"subnet_id": {
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
					},
					"is_primary": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"primary_private_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"secondary_private_ips": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
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

// Create OAPINic
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

	primaryPrivateIp, primaryOk := d.GetOk("primary_private_ip")
	secondaryPrivateIp, secondaryOk := d.GetOk("secondary_private_ips")
	if primaryOk || secondaryOk {
		privateIPs := make([]oscgo.PrivateIpLight, 0)

		if primaryOk {
			isPrimary := true
			private := primaryPrivateIp.(string)
			privateIP := oscgo.PrivateIpLight{
				IsPrimary: &isPrimary,
				PrivateIp: &private,
			}
			privateIPs = append(privateIPs, privateIP)
		}

		if secondaryOk {
			for _, v := range secondaryPrivateIp.(*schema.Set).List() {
				isPrimary := false
				private := v.(string)
				privateIP := oscgo.PrivateIpLight{
					IsPrimary: &isPrimary,
					PrivateIp: &private,
				}
				privateIPs = append(privateIPs, privateIP)
			}
		}

		request.SetPrivateIps(privateIPs)
	}

	log.Printf("[DEBUG] Creating network interface")

	var resp oscgo.CreateNicResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		resp, _, err = conn.NicApi.CreateNic(context.Background()).CreateNicRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
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
	}

	if err := d.Set("tags", make([]map[string]interface{}, 0)); err != nil {
		return err
	}
	if err := d.Set("private_ip", ""); err != nil {
		return err
	}

	log.Printf("[INFO] ENI ID: %s", d.Id())

	return resourceOutscaleOAPINicRead(d, meta)

}

// Read OAPINic
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

		resp, _, err = conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(dnir).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	if err != nil {

		return fmt.Errorf("Error describing Network Interfaces : %s", err)
	}

	if err != nil {
		if strings.Contains(err.Error(), "Unable to find Nic") {
			// The ENI is gone now, so just remove it from the state
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving ENI: %s", err)
	}

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetNics()), "Nic"); err != nil {
		return err
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

	privateIps := make([]map[string]interface{}, len(eni.GetPrivateIps()))
	secondaryIps := make([]string, 0)
	if eni.PrivateIps != nil {
		for k, v := range eni.GetPrivateIps() {
			privateIp := make(map[string]interface{})

			link := make(map[string]interface{})
			assoc := v.GetLinkPublicIp()
			link["public_ip_id"] = assoc.GetPublicIpId()
			link["link_public_ip_id"] = assoc.GetLinkPublicIpId()
			link["public_ip_account_id"] = assoc.GetPublicIpAccountId()
			link["public_dns_name"] = assoc.GetPublicDnsName()
			link["public_ip"] = assoc.GetPublicIp()

			privateIp["link_public_ip"] = link
			privateIp["private_dns_name"] = v.GetPrivateDnsName()
			privateIp["private_ip"] = v.GetPrivateIp()
			privateIp["is_primary"] = v.GetIsPrimary()

			privateIps[k] = privateIp
			if v.GetIsPrimary() {
				if err := d.Set("primary_private_ip", v.GetPrivateIp()); err != nil {
					return err
				}
			} else {
				secondaryIps = append(secondaryIps, v.GetPrivateIp())
			}
		}
	}
	if err := d.Set("private_ips", privateIps); err != nil {
		return err
	}
	if err := d.Set("secondary_private_ips", secondaryIps); err != nil {
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

// Delete OAPINic
func resourceOutscaleOAPINicDelete(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[INFO] Deleting ENI: %s", d.Id())

	err := resourceOutscaleOAPINicDetach(meta, d.Id())
	if err != nil {
		return err
	}

	deleteEniOpts := oscgo.DeleteNicRequest{
		NicId: d.Id(),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.NicApi.DeleteNic(context.Background()).DeleteNicRequest(deleteEniOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	if err != nil {

		return fmt.Errorf("Error Deleting ENI: %s", err)
	}

	return nil
}

func resourceOutscaleOAPINicDetach(meta interface{}, nicID string) error {
	// if there was an old nic_link, remove it
	conn := meta.(*OutscaleClient).OSCAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"attaching", "detaching"},
		Target:  []string{"attached", "detached", "failed"},
		Refresh: nicLinkRefreshFunc(conn, nicID),
		Timeout: 10 * time.Minute,
	}
	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for ENI (%s) to become dettached: %s", nicID, err)
	}

	r := resp.(oscgo.ReadNicsResponse)
	linkNic := r.GetNics()[0].GetLinkNic()

	if !reflect.DeepEqual(linkNic, oscgo.LinkNic{}) {
		log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", nicID)

		req := oscgo.UnlinkNicRequest{
			LinkNicId: linkNic.GetLinkNicId(),
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err = conn.NicApi.UnlinkNic(context.Background()).UnlinkNicRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), utils.ResourceNotFound) {
				return fmt.Errorf("Error detaching ENI: %s", err)
			}
		}
	}
	return nil
}

// Update OAPINic
func resourceOutscaleOAPINicUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	var err error

	if d.HasChange("link_nic") {
		_, na := d.GetChange("link_nic")

		err := resourceOutscaleOAPINicDetach(meta, d.Id())
		if err != nil {
			return err
		}

		// if there is a new nic_link, attach it
		if na != nil && len(na.([]interface{})) > 0 {
			na := na.([]interface{})[0].(map[string]interface{})
			di := na["device_number"].(int)
			ar := oscgo.LinkNicRequest{
				DeviceNumber: int32(di),
				VmId:         na["instance"].(string),
				NicId:        d.Id(),
			}

			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.NicApi.LinkNic(context.Background()).LinkNicRequest(ar).Execute()
				if err != nil {
					return utils.CheckThrottling(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Error Attaching Network Interface: %s", err)
			}
		}
	}

	if d.HasChange("secondary_private_ips") {
		oldSecondary, newSecondary := d.GetChange("secondary_private_ips")
		inter := oldSecondary.(*schema.Set).Intersection(newSecondary.(*schema.Set))
		created := newSecondary.(*schema.Set).Difference(inter)
		removed := oldSecondary.(*schema.Set).Difference(inter)

		if removed.Len() > 0 {
			input := oscgo.UnlinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: utils.InterfaceSliceToStringSlice(removed.List()),
			}

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.NicApi.UnlinkPrivateIps(context.Background()).UnlinkPrivateIpsRequest(input).Execute()
				if err != nil {
					return utils.CheckThrottling(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Failure to unassign Private IPs: %s", err)
			}
		}

		// Assign new IP addresses
		if created.Len() > 0 {
			input := oscgo.LinkPrivateIpsRequest{
				NicId: d.Id(),
			}
			input.SetPrivateIps(utils.InterfaceSliceToStringSlice(created.List()))
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, _, err = conn.NicApi.LinkPrivateIps(context.Background()).LinkPrivateIpsRequest(input).Execute()
				if err != nil {
					return utils.CheckThrottling(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Failure to assign Private IPs: %s", err)
			}
		}
	}

	if d.HasChange("security_group_ids") {
		stringValueList := expandStringValueList(d.Get("security_group_ids").([]interface{}))
		request := oscgo.UpdateNicRequest{
			NicId:            d.Id(),
			SecurityGroupIds: &stringValueList,
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err = conn.NicApi.UpdateNic(context.Background()).UpdateNicRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Failure updating ENI: %s", err)
		}
	}

	if d.HasChange("description") {
		request := oscgo.UpdateNicRequest{
			NicId:       d.Id(),
			Description: pointy.String(d.Get("description").(string)),
		}

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, _, err = conn.NicApi.UpdateNic(context.Background()).UpdateNicRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Failure updating ENI: %s", err)
		}
	}

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	return resourceOutscaleOAPINicRead(d, meta)
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
