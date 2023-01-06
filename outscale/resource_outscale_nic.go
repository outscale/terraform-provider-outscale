package outscale

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			Type:     schema.TypeSet,
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
			Type:     schema.TypeSet,
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
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"link_public_ip": {
						Type:     schema.TypeList,
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

	if v, ok := d.GetOk("private_ips"); ok {
		request.SetPrivateIps(expandPrivateIPLight(v.([]interface{})))
	}

	log.Printf("[DEBUG] Creating network interface")

	var resp oscgo.CreateNicResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.NicApi.CreateNic(context.Background()).CreateNicRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
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
		rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(dnir).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
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
	if utils.IsResponseEmpty(len(resp.GetNics()), "Nic", d.Id()) {
		d.SetId("")
		return nil
	}
	eni := resp.GetNics()[0]
	if err := d.Set("description", eni.GetDescription()); err != nil {
		return err
	}
	if err := d.Set("subnet_id", eni.GetSubnetId()); err != nil {
		return err
	}
	if linkIp, ok := eni.GetLinkPublicIpOk(); ok {
		if err := d.Set("link_public_ip", flattenLinkPublicIp(linkIp)); err != nil {
			return err
		}
	}

	if linkNic, ok := eni.GetLinkNicOk(); ok {
		if err := d.Set("link_nic", flattenLinkNic(linkNic)); err != nil {
			return err
		}
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

	y := make([]map[string]interface{}, len(eni.GetPrivateIps()))
	if eni.PrivateIps != nil {
		for k, v := range eni.GetPrivateIps() {
			b := make(map[string]interface{})

			d := make(map[string]interface{})
			if assoc, ok := v.GetLinkPublicIpOk(); ok {
				d["public_ip_id"] = assoc.GetPublicIpId()
				d["link_public_ip_id"] = assoc.GetLinkPublicIpId()
				d["public_ip_account_id"] = assoc.GetPublicIpAccountId()
				d["public_dns_name"] = assoc.GetPublicDnsName()
				d["public_ip"] = assoc.GetPublicIp()
				b["link_public_ip"] = d
			}
			b["private_dns_name"] = v.GetPrivateDnsName()
			b["private_ip"] = v.GetPrivateIp()
			b["is_primary"] = v.GetIsPrimary()

			y[k] = b
		}
	}
	if err := d.Set("private_ips", y); err != nil {
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
		_, httpResp, err := conn.NicApi.DeleteNic(context.Background()).DeleteNicRequest(deleteEniOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
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
		var statusCode int
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.NicApi.UnlinkNic(context.Background()).UnlinkNicRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			statusCode = httpResp.StatusCode
			return nil
		})

		if err != nil {
			if statusCode == http.StatusNotFound {
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
				_, httpResp, err := conn.NicApi.LinkNic(context.Background()).LinkNicRequest(ar).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Error Attaching Network Interface: %s", err)
			}
		}
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
				_, httpResp, err := conn.NicApi.UnlinkPrivateIps(context.Background()).UnlinkPrivateIpsRequest(input).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
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
				_, httpResp, err := conn.NicApi.LinkPrivateIps(context.Background()).LinkPrivateIpsRequest(input).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("Failure to assign Private IPs: %s", err)
			}
		}
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
					SecondaryPrivateIpCount: pointy.Int32(dif),
				}

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					var err error
					_, httpResp, err := conn.NicApi.LinkPrivateIps(context.Background()).LinkPrivateIpsRequest(input).Execute()
					if err != nil {
						return utils.CheckThrottling(httpResp, err)
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
					PrivateIps: utils.InterfaceSliceToStringSlice(prips[0:int(math.Abs(float64(diff)))]),
				}

				err := resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, httpResp, err := conn.NicApi.UnlinkPrivateIps(context.Background()).UnlinkPrivateIpsRequest(input).Execute()
					if err != nil {
						return utils.CheckThrottling(httpResp, err)
					}
					return nil
				})
				if err != nil {
					return fmt.Errorf("Failure to unassign Private IPs: %s", err)
				}
			}
		}
	}

	if d.HasChange("security_group_ids") {
		stringValueList := utils.InterfaceSliceToStringSlice(d.Get("security_group_ids").([]interface{}))
		request := oscgo.UpdateNicRequest{
			NicId:            d.Id(),
			SecurityGroupIds: &stringValueList,
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, httpResp, err := conn.NicApi.UpdateNic(context.Background()).UpdateNicRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
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
			_, httpResp, err := conn.NicApi.UpdateNic(context.Background()).UpdateNicRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
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

func flattenLinkPublicIp(linkIp *oscgo.LinkPublicIp) []map[string]interface{} {
	return []map[string]interface{}{{
		"public_ip_id":         linkIp.GetPublicIpId(),
		"link_public_ip_id":    linkIp.GetLinkPublicIpId(),
		"public_ip_account_id": linkIp.GetPublicIpAccountId(),
		"public_dns_name":      linkIp.GetPublicDnsName(),
		"public_ip":            linkIp.GetPublicIp(),
	}}
}

func flattenLinkNic(linkNic *oscgo.LinkNic) []map[string]interface{} {
	return []map[string]interface{}{{
		"link_nic_id":           linkNic.GetLinkNicId(),
		"delete_on_vm_deletion": strconv.FormatBool(linkNic.GetDeleteOnVmDeletion()),
		"device_number":         linkNic.GetDeviceNumber(),
		"vm_id":                 linkNic.GetVmId(),
		"vm_account_id":         linkNic.GetVmAccountId(),
		"state":                 linkNic.GetState(),
	}}
}
