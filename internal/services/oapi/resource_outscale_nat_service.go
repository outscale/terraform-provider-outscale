package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleNatService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOAPINatServiceCreate,
		ReadContext:   resourceOAPINatServiceRead,
		DeleteContext: resourceOAPINatServiceDelete,
		UpdateContext: ResourceOutscaleNatServiceUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"public_ip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nat_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip_id": {
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaSDK(),
		},
	}
}

func resourceOAPINatServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	req := osc.CreateNatServiceRequest{
		PublicIpId: d.Get("public_ip_id").(string),
		SubnetId:   d.Get("subnet_id").(string),
	}

	resp, err := client.CreateNatService(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error creating nat service: %s", err.Error())
	}

	if resp.NatService == nil {
		return diag.Errorf("error there is not nat service (%s)", err)
	}

	natService := resp.NatService

	// Get the ID and store it
	log.Printf("\n\n[INFO] NAT Service ID: %s", natService.NatServiceId)

	// Wait for the NAT Service to become available
	log.Printf("\n\n[DEBUG] Waiting for NAT Service (%s) to become available", natService.NatServiceId)

	filterReq := osc.ReadNatServicesRequest{
		Filters: &osc.FiltersNatService{NatServiceIds: &[]string{natService.NatServiceId}},
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Timeout: timeout,
		Refresh: NGOAPIStateRefreshFunc(ctx, client, filterReq, timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for nat service (%s) to become available: %s", natService.NatServiceId, err)
	}
	d.SetId(natService.NatServiceId)

	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOAPINatServiceRead(ctx, d, meta)
}

func resourceOAPINatServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	filterReq := osc.ReadNatServicesRequest{
		Filters: &osc.FiltersNatService{NatServiceIds: &[]string{d.Id()}},
	}

	resp, err := client.ReadNatServices(ctx, filterReq, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error waiting for nat service (%s) to become available: %s", d.Id(), err)
	}
	if resp.NatServices == nil || utils.IsResponseEmpty(len(*resp.NatServices), "NatService", d.Id()) {
		d.SetId("")
		return nil
	}
	natService := (*resp.NatServices)[0]

	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(natService.NatServiceId)

		if err := set("nat_service_id", natService.NatServiceId); err != nil {
			return err
		}
		if err := set("net_id", natService.NetId); err != nil {
			return err
		}
		if err := set("state", natService.State); err != nil {
			return err
		}
		if err := set("subnet_id", natService.SubnetId); err != nil {
			return err
		}

		public_ips := natService.PublicIps
		if err := set("public_ips", getOSCPublicIPs(public_ips)); err != nil {
			return err
		}

		if len(public_ips) > 0 {
			if err := set("public_ip_id", public_ips[0].PublicIpId); err != nil {
				return err
			}
		} else {
			if err := set("public_ip_id", ""); err != nil {
				return err
			}
		}

		if err := d.Set("tags", FlattenOAPITagsSDK(natService.Tags)); err != nil {
			fmt.Printf("[WARN] ERROR TAGS PROBLEME (%s)", err)
		}

		return nil
	}))
}

func ResourceOutscaleNatServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return resourceOAPINatServiceRead(ctx, d, meta)
}

func resourceOAPINatServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[INFO] Deleting NAT Service: %s\n", d.Id())
	_, err := client.DeleteNatService(ctx, osc.DeleteNatServiceRequest{
		NatServiceId: d.Id(),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error deleting nat service: %s", err)
	}

	filterReq := osc.ReadNatServicesRequest{
		Filters: &osc.FiltersNatService{NatServiceIds: &[]string{d.Id()}},
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted", "available"},
		Timeout: timeout,
		Refresh: NGOAPIStateRefreshFunc(ctx, client, filterReq, timeout),
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("error waiting for nat service (%s) to delete: %s", d.Id(), stateErr)
	}
	return nil
}

// NGOAPIStateRefreshFunc returns a retry.StateRefreshFunc that is used to watch
// a NAT Service.
func NGOAPIStateRefreshFunc(ctx context.Context, client *osc.Client, req osc.ReadNatServicesRequest, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ReadNatServices(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "", err
		}
		if resp.NatServices == nil {
			return nil, "", fmt.Errorf("nat service not found")
		}

		return resp, string((*resp.NatServices)[0].State), nil
	}
}

func getOSCPublicIPs(publicIps []osc.PublicIpLight) (res []map[string]interface{}) {
	for _, p := range publicIps {
		res = append(res, map[string]interface{}{
			"public_ip_id": p.PublicIpId,
			"public_ip":    p.PublicIp,
		})
	}
	return
}
