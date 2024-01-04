package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOutscaleOAPINatService() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPINatServiceCreate,
		Read:   resourceOAPINatServiceRead,
		Delete: resourceOAPINatServiceDelete,
		Update: resourceOutscaleOAPINatServiceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"tags": tagsListOAPISchema(),
		},
	}
}

func resourceOAPINatServiceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateNatServiceRequest{
		PublicIpId: d.Get("public_ip_id").(string),
		SubnetId:   d.Get("subnet_id").(string),
	}

	var resp oscgo.CreateNatServiceResponse
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NatServiceApi.CreateNatService(context.Background()).CreateNatServiceRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating Nat Service: %s", err.Error())
	}

	if !resp.HasNatService() {
		return fmt.Errorf("Error there is not Nat Service (%s)", err)
	}

	natService := resp.GetNatService()

	// Get the ID and store it
	log.Printf("\n\n[INFO] NAT Service ID: %s", natService.GetNatServiceId())

	// Wait for the NAT Service to become available
	log.Printf("\n\n[DEBUG] Waiting for NAT Service (%s) to become available", natService.GetNatServiceId())

	filterReq := oscgo.ReadNatServicesRequest{
		Filters: &oscgo.FiltersNatService{NatServiceIds: &[]string{natService.GetNatServiceId()}},
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Refresh: NGOAPIStateRefreshFunc(conn, filterReq, "failed"),
		Timeout: 10 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for NAT Service (%s) to become available: %s", natService.GetNatServiceId(), err)
	}
	//SetTags
	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.(*schema.Set), natService.GetNatServiceId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(natService.GetNatServiceId())

	return resourceOAPINatServiceRead(d, meta)
}

func resourceOAPINatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filterReq := oscgo.ReadNatServicesRequest{
		Filters: &oscgo.FiltersNatService{NatServiceIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadNatServicesResponse
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NatServiceApi.ReadNatServices(context.Background()).ReadNatServicesRequest(filterReq).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error waiting for NAT Service (%s) to become available: %s", d.Id(), err)
	}
	if utils.IsResponseEmpty(len(resp.GetNatServices()), "NatService", d.Id()) {
		d.SetId("")
		return nil
	}
	natService := resp.GetNatServices()[0]

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(natService.GetNatServiceId())

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

		public_ips := natService.GetPublicIps()
		if err := set("public_ips", getOSCPublicIPs(public_ips)); err != nil {
			return err
		}

		if len(public_ips) > 0 {
			if err := set("public_ip_id", public_ips[0].GetPublicIpId()); err != nil {
				return err
			}
		} else {
			if err := set("public_ip_id", ""); err != nil {
				return err
			}
		}

		if err := d.Set("tags", tagsOSCAPIToMap(natService.GetTags())); err != nil {
			fmt.Printf("[WARN] ERROR TAGS PROBLEME (%s)", err)
		}

		return nil
	})
}

func resourceOutscaleOAPINatServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}
	return resourceOAPINatServiceRead(d, meta)
}

func resourceOAPINatServiceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[INFO] Deleting NAT Service: %s\n", d.Id())
	err := resource.Retry(120*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.NatServiceApi.DeleteNatService(context.Background()).DeleteNatServiceRequest(oscgo.DeleteNatServiceRequest{
			NatServiceId: d.Id(),
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting Nat Service: %s", err)
	}

	filterReq := oscgo.ReadNatServicesRequest{
		Filters: &oscgo.FiltersNatService{NatServiceIds: &[]string{d.Id()}},
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "available"},
		Refresh:    NGOAPIStateRefreshFunc(conn, filterReq, "failed"),
		Timeout:    30 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf("Error waiting for NAT Service (%s) to delete: %s", d.Id(), stateErr)
	}
	return nil
}

// NGOAPIStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a NAT Service.
func NGOAPIStateRefreshFunc(client *oscgo.APIClient, req oscgo.ReadNatServicesRequest, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadNatServicesResponse
		err := resource.Retry(120*time.Second, func() *resource.RetryError {
			var err error
			rp, httpResp, err := client.NatServiceApi.ReadNatServices(context.Background()).ReadNatServicesRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return nil, "failed", err
		}

		state := "deleted"
		if resp.HasNatServices() && len(resp.GetNatServices()) > 0 {
			natServices := resp.GetNatServices()
			state = natServices[0].GetState()

			if state == failState {
				return natServices[0], state, fmt.Errorf("Failed to reach target state. Reason: %v", state)
			}
		}

		return resp, state, nil
	}
}

func getOSCPublicIPs(publicIps []oscgo.PublicIpLight) (res []map[string]interface{}) {
	for _, p := range publicIps {
		res = append(res, map[string]interface{}{
			"public_ip_id": p.GetPublicIpId(),
			"public_ip":    p.GetPublicIp(),
		})
	}
	return
}
