package oapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

// Creates a network interface in the specified subnet
func ResourceOutscaleNic() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleNicCreate,
		Read:   ResourceOutscaleNicRead,
		Update: ResourceOutscaleNicUpdate,
		Delete: ResourceOutscaleNicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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
			Type:     schema.TypeSet,
			Computed: true,
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
			Type:     schema.TypeSet,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
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
		"tags": TagsSchemaSDK(),
		"net_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// Create OAPINic
func ResourceOutscaleNicCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutCreate)

	request := oscgo.CreateNicRequest{
		SubnetId: d.Get("subnet_id").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		request.SetDescription(v.(string))
	}
	if sgIDs := utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set)); len(sgIDs) > 0 {
		request.SetSecurityGroupIds(sgIDs)
	}
	if v, ok := d.GetOk("private_ips"); ok {
		request.SetPrivateIps(expandPrivateIPLight(v.(*schema.Set).List()))
	}

	var resp oscgo.CreateNicResponse
	err := retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.NicApi.CreateNic(context.Background()).CreateNicRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error creating nic: %s", err)
	}

	d.SetId(resp.Nic.GetNicId())

	if d.IsNewResource() {
		if err := updateOAPITagsSDK(conn, d); err != nil {
			return err
		}
	}

	if err := d.Set("tags", make([]map[string]interface{}, 0)); err != nil {
		return err
	}
	if err := d.Set("private_ip", ""); err != nil {
		return err
	}
	return ResourceOutscaleNicRead(d, meta)
}

// Read OAPINic
func ResourceOutscaleNicRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutRead)
	dnir := oscgo.ReadNicsRequest{
		Filters: &oscgo.FiltersNic{
			NicIds: &[]string{d.Id()},
		},
	}

	var resp oscgo.ReadNicsResponse
	err := retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
		rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(dnir).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error describing network interfaces : %s", err)
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
	if err := d.Set("security_groups", getSecurityGroups(eni.GetSecurityGroups())); err != nil {
		return err
	}
	if err := d.Set("security_group_ids", getSecurityGroupIds(eni.GetSecurityGroups())); err != nil {
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
	if privIps, ok := eni.GetPrivateIpsOk(); ok {
		if err := d.Set("private_ips", oapihelpers.GetPrivateIPsForNic(*privIps)); err != nil {
			return err
		}
	}
	if err := d.Set("is_source_dest_checked", eni.GetIsSourceDestChecked()); err != nil {
		return err
	}
	if err := d.Set("state", eni.GetState()); err != nil {
		return err
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(eni.GetTags())); err != nil {
		return err
	}
	if err := d.Set("net_id", eni.GetNetId()); err != nil {
		return err
	}

	return nil
}

// Delete OAPINic
func ResourceOutscaleNicDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutDelete)

	err := ResourceOutscaleNicDetach(meta, d.Id(), timeout)
	if err != nil {
		return err
	}

	deleteEniOpts := oscgo.DeleteNicRequest{
		NicId: d.Id(),
	}

	err = retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
		_, httpResp, err := conn.NicApi.DeleteNic(context.Background()).DeleteNicRequest(deleteEniOpts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting eni: %s", err)
	}

	return nil
}

func ResourceOutscaleNicDetach(meta interface{}, nicID string, timeout time.Duration) error {
	// if there was an old nic_link, remove it
	conn := meta.(*client.OutscaleClient).OSCAPI

	stateConf := &retry.StateChangeConf{
		Pending: []string{"attaching", "detaching"},
		Target:  []string{"attached", "detached", "failed"},
		Refresh: nicLinkRefreshFunc(conn, nicID, timeout),
		Timeout: timeout,
		Delay:   1 * time.Second,
	}
	resp, err := stateConf.WaitForStateContext(context.Background())
	if err != nil {
		return fmt.Errorf(
			"error waiting for eni (%s) to become dettached: %s", nicID, err)
	}

	r := resp.(oscgo.ReadNicsResponse)
	linkNic := r.GetNics()[0].GetLinkNic()

	if !reflect.DeepEqual(linkNic, oscgo.LinkNic{}) {
		log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", nicID)

		req := oscgo.UnlinkNicRequest{
			LinkNicId: linkNic.GetLinkNicId(),
		}
		var statusCode int
		err := retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
			_, httpResp, err := conn.NicApi.UnlinkNic(context.Background()).UnlinkNicRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			statusCode = httpResp.StatusCode
			return nil
		})
		if err != nil {
			if statusCode == http.StatusNotFound {
				return fmt.Errorf("error detaching eni: %s", err)
			}
		}
	}
	return nil
}

// Update OAPINic
func ResourceOutscaleNicUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutUpdate)
	var err error

	if d.HasChange("private_ips") {
		oldIps, newIps := d.GetChange("private_ips")
		removed, created := oapihelpers.GetTypeSetDifferencesForUpdating(oldIps.(*schema.Set), newIps.(*schema.Set))

		// Unassign old IP addresses
		if listOldIps := removed.List(); len(listOldIps) != 0 {
			input := oscgo.UnlinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: flattenPrivateIPLightToStringSlice(listOldIps),
			}

			err := retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
				_, httpResp, err := conn.NicApi.UnlinkPrivateIps(context.Background()).UnlinkPrivateIpsRequest(input).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failure to unassign private ips: %s", err)
			}
		}

		// Assign new IP addresses
		if listNewIps := created.List(); len(listNewIps) != 0 {
			stringSlice := flattenPrivateIPLightToStringSlice(listNewIps)
			input := oscgo.LinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: &stringSlice,
			}

			err = retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
				_, httpResp, err := conn.NicApi.LinkPrivateIps(context.Background()).LinkPrivateIpsRequest(input).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failure to assign private ips: %s", err)
			}
		}
	}

	if d.HasChange("security_group_ids") {
		request := oscgo.UpdateNicRequest{
			NicId: d.Id(),
		}
		request.SetSecurityGroupIds(utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set)))

		err = retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
			_, httpResp, err := conn.NicApi.UpdateNic(context.Background()).UpdateNicRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failure updating eni: %s", err)
		}
	}

	if d.HasChange("description") {
		request := oscgo.UpdateNicRequest{
			NicId: d.Id(),
		}
		request.SetDescription(d.Get("description").(string))
		err := retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
			_, httpResp, err := conn.NicApi.UpdateNic(context.Background()).UpdateNicRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failure updating eni: %s", err)
		}
	}

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}
	return ResourceOutscaleNicRead(d, meta)
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
