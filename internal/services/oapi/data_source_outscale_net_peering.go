package oapi

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"
)

func DataSourceOutscaleNetPeering() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNetPeeringRead,

		Schema: map[string]*schema.Schema{
			"filter":       dataSourceFiltersSchema(),
			"accepter_net": vpcOAPIPeeringConnectionOptionsSchema(),
			"net_peering_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_net": vpcOAPIPeeringConnectionOptionsSchema(),
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

func vpcOAPIPeeringConnectionOptionsSchema() *schema.Schema {
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

func DataSourceOutscaleNetPeeringRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	log.Printf("[DEBUG] Reading Net Peering Connections.")

	var err error
	req := osc.ReadNetPeeringsRequest{}

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.FromErr(ErrFilterRequired)
	}
	req.Filters, err = buildOutscaleLinPeeringConnectionFilters(filters.(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ReadNetPeerings(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.Errorf("error reading net peering connection details: %s", err)
	}

	if resp.NetPeerings == nil || len(*resp.NetPeerings) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.NetPeerings) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}
	netPeering := (*resp.NetPeerings)[0]

	// The failed status is a status that we can assume just means the
	// connection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	if !reflect.DeepEqual(netPeering.State, osc.NetPeeringState{}) {
		status := map[string]bool{
			"deleted":  true,
			"deleting": true,
			"expired":  true,
			"failed":   true,
			"rejected": true,
		}
		if _, ok := status[string(netPeering.State.Name)]; ok {
			log.Printf("[DEBUG] Net Peering Connection (%s) in state (%s), removing.",
				d.Id(), netPeering.State.Name)
			return nil
		}
	}
	log.Printf("[DEBUG] Net Peering Connection response: %#v", netPeering)

	log.Printf("[DEBUG] Net Peering Connection Source %v, Accepter %v", netPeering.SourceNet.AccountId, netPeering.AccepterNet.AccountId)

	if !reflect.DeepEqual(netPeering.AccepterNet, osc.AccepterNet{}) {
		if err := d.Set("accepter_net", getOAPINetPeeringAccepterNet(netPeering.AccepterNet)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(netPeering.SourceNet, osc.SourceNet{}) {
		if err := d.Set("source_net", getOAPINetPeeringSourceNet(netPeering.SourceNet)); err != nil {
			return diag.FromErr(err)
		}
	}
	if netPeering.State.Name != "" {
		if err := d.Set("state", getOAPINetPeeringState(netPeering.State)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("net_peering_id", netPeering.NetPeeringId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(netPeering.Tags)); err != nil {
		return diag.Errorf("error setting net peering tags: %s", err)
	}

	d.SetId(netPeering.NetPeeringId)

	return nil
}

func buildOutscaleLinPeeringConnectionFilters(set *schema.Set) (*osc.FiltersNetPeering, error) {
	var filters osc.FiltersNetPeering
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "accepter_net_account_ids":
			filters.AccepterNetAccountIds = &filterValues
		case "accepter_net_ip_ranges":
			filters.AccepterNetIpRanges = &filterValues
		case "accepter_net_net_ids":
			filters.AccepterNetNetIds = &filterValues
		case "net_peering_ids":
			filters.NetPeeringIds = &filterValues
		case "source_net_account_ids":
			filters.SourceNetAccountIds = &filterValues
		case "source_net_ip_ranges":
			filters.SourceNetIpRanges = &filterValues
		case "source_net_net_ids":
			filters.SourceNetNetIds = &filterValues
		case "state_messages":
			filters.StateMessages = &filterValues
		case "expiration_dates":
			expirationDates, err := utils.StringSliceToTimeSlice(
				filterValues, "expiration_dates")
			if err != nil {
				return nil, err
			}
			filters.ExpirationDates = &expirationDates
		case "state_names":
			filters.StateNames = new(lo.Map(filterValues, func(s string, _ int) osc.NetPeeringStateName {
				return osc.NetPeeringStateName(s)
			}))
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
