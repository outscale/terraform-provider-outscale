package oapi

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

// Creates a network interface in the specified subnet
func ResourceOutscaleNic() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleNicCreate,
		ReadContext:   ResourceOutscaleNicRead,
		UpdateContext: ResourceOutscaleNicUpdate,
		DeleteContext: ResourceOutscaleNicDelete,
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
func ResourceOutscaleNicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	request := osc.CreateNicRequest{
		SubnetId: d.Get("subnet_id").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = new(v.(string))
	}
	if sgIDs := utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set)); len(sgIDs) > 0 {
		request.SecurityGroupIds = &sgIDs
	}
	if v, ok := d.GetOk("private_ips"); ok {
		request.PrivateIps = new(expandPrivateIPLight(v.(*schema.Set).List()))
	}

	resp, err := client.CreateNic(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error creating nic: %s", err)
	}

	d.SetId(resp.Nic.NicId)

	if d.IsNewResource() {
		if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("tags", make([]map[string]interface{}, 0)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_ip", ""); err != nil {
		return diag.FromErr(err)
	}
	return ResourceOutscaleNicRead(ctx, d, meta)
}

// Read OAPINic
func ResourceOutscaleNicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)
	dnir := osc.ReadNicsRequest{
		Filters: &osc.FiltersNic{
			NicIds: &[]string{d.Id()},
		},
	}

	resp, err := client.ReadNics(ctx, dnir, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error describing network interfaces : %s", err)
	}

	if resp.Nics == nil || utils.IsResponseEmpty(len(*resp.Nics), "Nic", d.Id()) {
		d.SetId("")
		return nil
	}
	eni := (*resp.Nics)[0]
	if err := d.Set("description", eni.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subnet_id", eni.SubnetId); err != nil {
		return diag.FromErr(err)
	}
	if eni.LinkPublicIp != nil {
		if err := d.Set("link_public_ip", flattenLinkPublicIp(eni.LinkPublicIp)); err != nil {
			return diag.FromErr(err)
		}
	}

	if eni.LinkNic != nil {
		if err := d.Set("link_nic", flattenLinkNic(eni.LinkNic)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("subregion_name", eni.SubregionName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("security_groups", getSecurityGroups(eni.SecurityGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("security_group_ids", getSecurityGroupIds(eni.SecurityGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mac_address", eni.MacAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_id", eni.NicId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account_id", eni.AccountId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_dns_name", eni.PrivateDnsName); err != nil {
		return diag.FromErr(err)
	}
	if eni.PrivateIps != nil {
		if err := d.Set("private_ips", oapihelpers.GetPrivateIPsForNic(eni.PrivateIps)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("is_source_dest_checked", eni.IsSourceDestChecked); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", eni.State); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(eni.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_id", eni.NetId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Delete OAPINic
func ResourceOutscaleNicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	err := ResourceOutscaleNicDetach(ctx, meta, d.Id(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteEniOpts := osc.DeleteNicRequest{
		NicId: d.Id(),
	}

	_, err = client.DeleteNic(ctx, deleteEniOpts, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error deleting eni: %s", err)
	}

	return nil
}

func ResourceOutscaleNicDetach(ctx context.Context, meta interface{}, nicID string, timeout time.Duration) error {
	// if there was an old nic_link, remove it
	client := meta.(*client.OutscaleClient).OSC

	stateConf := &retry.StateChangeConf{
		Pending: []string{"attaching", "detaching"},
		Target:  []string{"attached", "detached", "failed"},
		Timeout: timeout,
		Refresh: nicLinkRefreshFunc(ctx, client, nicID, timeout),
	}
	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf(
			"error waiting for eni (%s) to become dettached: %s", nicID, err)
	}
	r := resp.(*osc.ReadNicsResponse)
	if r == nil || r.Nics == nil {
		return fmt.Errorf("nic (%s) not found", nicID)
	}

	linkNic := ptr.From((*r.Nics)[0].LinkNic)

	if !reflect.DeepEqual(linkNic, osc.LinkNic{}) {
		log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", nicID)

		req := osc.UnlinkNicRequest{
			LinkNicId: linkNic.LinkNicId,
		}
		_, err := client.UnlinkNic(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return err
		}
	}
	return nil
}

// Update OAPINic
func ResourceOutscaleNicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if d.HasChange("private_ips") {
		oldIps, newIps := d.GetChange("private_ips")
		removed, created := oapihelpers.GetTypeSetDifferencesForUpdating(oldIps.(*schema.Set), newIps.(*schema.Set))

		// Unassign old IP addresses
		if listOldIps := removed.List(); len(listOldIps) != 0 {
			input := osc.UnlinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: flattenPrivateIPLightToStringSlice(listOldIps),
			}

			_, err := client.UnlinkPrivateIps(ctx, input, options.WithRetryTimeout(timeout))
			if err != nil {
				return diag.Errorf("failure to unassign private ips: %s", err)
			}
		}

		// Assign new IP addresses
		if listNewIps := created.List(); len(listNewIps) != 0 {
			stringSlice := flattenPrivateIPLightToStringSlice(listNewIps)
			input := osc.LinkPrivateIpsRequest{
				NicId:      d.Id(),
				PrivateIps: &stringSlice,
			}

			_, err := client.LinkPrivateIps(ctx, input, options.WithRetryTimeout(timeout))
			if err != nil {
				return diag.Errorf("failure to assign private ips: %s", err)
			}
		}
	}

	if d.HasChange("security_group_ids") {
		request := osc.UpdateNicRequest{
			NicId: d.Id(),
		}
		request.SecurityGroupIds = new(utils.SetToStringSlice(d.Get("security_group_ids").(*schema.Set)))
		_, err := client.UpdateNic(ctx, request, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.Errorf("failure updating eni: %s", err)
		}
	}

	if d.HasChange("description") {
		request := osc.UpdateNicRequest{
			NicId: d.Id(),
		}
		request.Description = new(d.Get("description").(string))
		_, err := client.UpdateNic(ctx, request, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.Errorf("failure updating eni: %s", err)
		}
	}

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return ResourceOutscaleNicRead(ctx, d, meta)
}

func expandPrivateIPLight(pIPs []interface{}) []osc.PrivateIpLight {
	privateIPs := make([]osc.PrivateIpLight, 0)
	for _, v := range pIPs {
		privateIPMap := v.(map[string]interface{})
		isPrimary := privateIPMap["is_primary"].(bool)
		private := privateIPMap["private_ip"].(string)
		privateIP := osc.PrivateIpLight{
			IsPrimary: isPrimary,
			PrivateIp: private,
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

func flattenLinkPublicIp(linkIp *osc.LinkPublicIp) []map[string]interface{} {
	return []map[string]interface{}{{
		"public_ip_id":         linkIp.PublicIpId,
		"link_public_ip_id":    linkIp.LinkPublicIpId,
		"public_ip_account_id": linkIp.PublicIpAccountId,
		"public_dns_name":      linkIp.PublicDnsName,
		"public_ip":            linkIp.PublicIp,
	}}
}

func flattenLinkNic(linkNic *osc.LinkNic) []map[string]interface{} {
	return []map[string]interface{}{{
		"link_nic_id":           linkNic.LinkNicId,
		"delete_on_vm_deletion": strconv.FormatBool(linkNic.DeleteOnVmDeletion),
		"device_number":         linkNic.DeviceNumber,
		"vm_id":                 linkNic.VmId,
		"vm_account_id":         linkNic.VmAccountId,
		"state":                 linkNic.State,
	}}
}
