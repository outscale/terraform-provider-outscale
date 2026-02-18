package oapi

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleNetPeering() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleNetPeeringRead,

		Schema: map[string]*schema.Schema{
			"filter":       dataSourceFiltersSchema(),
			"accepter_net": vpcOAPIPeeringclientectionOptionsSchema(),
			"net_peering_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_net": vpcOAPIPeeringclientectionOptionsSchema(),
			"state": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": TagsSchemaComputedSDK(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func vpcOAPIPeeringclientectionOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_range": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"account_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"net_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func DataSourceOutscaleNetPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	log.Printf("[DEBUG] Reading Net Peering clientections.")

	var err error
	req := osc.ReadNetPeeringsRequest{}

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return ErrFilterRequired
	}
	req.Filters, err = buildOutscaleLinPeeringclientectionFilters(filters.(*schema.Set))
	if err != nil {
		return err
	}

	var resp osc.ReadNetPeeringsResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.NetPeeringApi.ReadNetPeerings(ctx).ReadNetPeeringsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading net peering clientection details: %s", err)
	}

	if len(resp.GetNetPeerings()) == 0 {
		return ErrNoResults
	}
	if len(resp.GetNetPeerings()) > 1 {
		return ErrMultipleResults
	}
	netPeering := resp.GetNetPeerings()[0]

	// The failed status is a status that we can assume just means the
	// clientection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	if !reflect.DeepEqual(netPeering.State, osc.NetPeeringState{}) {
		status := map[string]bool{
			"deleted":  true,
			"deleting": true,
			"expired":  true,
			"failed":   true,
			"rejected": true,
		}
		if _, ok := status[netPeering.State.GetName()]; ok {
			log.Printf("[DEBUG] Net Peering clientection (%s) in state (%s), removing.",
				d.Id(), netPeering.State.GetName())
			return nil
		}
	}
	log.Printf("[DEBUG] Net Peering clientection response: %#v", netPeering)

	log.Printf("[DEBUG] Net Peering clientection Source %s, Accepter %s", netPeering.SourceNet.GetAccountId(), netPeering.AccepterNet.GetAccountId())

	if !reflect.DeepEqual(netPeering.GetAccepterNet(), osc.AccepterNet{}) {
		if err := d.Set("accepter_net", getOAPINetPeeringAccepterNet(*netPeering.AccepterNet)); err != nil {
			return err
		}
	}

	if !reflect.DeepEqual(netPeering.SourceNet, osc.SourceNet{}) {
		if err := d.Set("source_net", getOAPINetPeeringSourceNet(*netPeering.SourceNet)); err != nil {
			return err
		}
	}
	if netPeering.State.GetName() != "" {
		if err := d.Set("state", getOAPINetPeeringState(netPeering.GetState())); err != nil {
			return err
		}
	}
	if err := d.Set("net_peering_id", netPeering.GetNetPeeringId()); err != nil {
		return err
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(netPeering.Tags)); err != nil {
		return fmt.Errorf("error setting net peering tags: %s", err)
	}

	d.SetId(netPeering.GetNetPeeringId())

	return nil
}

func buildOutscaleLinPeeringclientectionFilters(set *schema.Set) (*osc.FiltersNetPeering, error) {
	var filters osc.FiltersNetPeering
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "accepter_net_account_ids":
			filters.SetAccepterNetAccountIds(filterValues)
		case "accepter_net_ip_ranges":
			filters.SetAccepterNetIpRanges(filterValues)
		case "accepter_net_net_ids":
			filters.SetAccepterNetNetIds(filterValues)
		case "net_peering_ids":
			filters.SetNetPeeringIds(filterValues)
		case "source_net_account_ids":
			filters.SetSourceNetAccountIds(filterValues)
		case "source_net_ip_ranges":
			filters.SetSourceNetIpRanges(filterValues)
		case "source_net_net_ids":
			filters.SetSourceNetNetIds(filterValues)
		case "state_messages":
			filters.SetStateMessages(filterValues)
		case "expiration_dates":
			expirationDates, err := utils.StringSliceToTimeSlice(
				filterValues, "expiration_dates")
			if err != nil {
				return nil, err
			}
			filters.SetExpirationDates(expirationDates)
		case "state_names":
			filters.SetStateNames(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(ctx, name)
		}
	}
	return &filters, nil
}
