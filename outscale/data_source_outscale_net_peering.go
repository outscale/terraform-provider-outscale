package outscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPILinPeeringConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILinPeeringConnectionRead,

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
			"tags": dataSourceTagsSchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPILinPeeringConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[DEBUG] Reading Net Peering Connections.")

	req := oscgo.ReadNetPeeringsRequest{}

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters must be assigned")
	}
	req.SetFilters(buildOutscaleOAPILinPeeringConnectionFilters(filters.(*schema.Set)))

	var resp oscgo.ReadNetPeeringsResponse
	var err error
	var statusCode int
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.NetPeeringApi.ReadNetPeerings(context.Background()).ReadNetPeeringsRequest(req).Execute()

		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})

	if err != nil {
		if statusCode == http.StatusNotFound {
			return fmt.Errorf("no matching Net Peering Connection found")
		}
		return fmt.Errorf("Error reading Net Peering Connection details: %s", err)
	}

	if len(resp.GetNetPeerings()) > 1 {
		return fmt.Errorf("multiple Net Peering connections matched; use additional constraints to reduce matches to a single Net Peering Connection")
	}
	netPeering := resp.GetNetPeerings()[0]

	// The failed status is a status that we can assume just means the
	// connection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	if !reflect.DeepEqual(netPeering.State, oscgo.NetPeeringState{}) {
		status := map[string]bool{
			"deleted":  true,
			"deleting": true,
			"expired":  true,
			"failed":   true,
			"rejected": true,
		}
		if _, ok := status[netPeering.State.GetName()]; ok {
			log.Printf("[DEBUG] Net Peering Connection (%s) in state (%s), removing.",
				d.Id(), netPeering.State.GetName())
			return nil
		}
	}
	log.Printf("[DEBUG] Net Peering Connection response: %#v", netPeering)

	log.Printf("[DEBUG] Net Peering Connection Source %s, Accepter %s", netPeering.SourceNet.GetAccountId(), netPeering.AccepterNet.GetAccountId())

	if !reflect.DeepEqual(netPeering.GetAccepterNet(), oscgo.AccepterNet{}) {
		if err := d.Set("accepter_net", getOAPINetPeeringAccepterNet(*netPeering.AccepterNet)); err != nil {
			return err
		}
	}

	if !reflect.DeepEqual(netPeering.SourceNet, oscgo.SourceNet{}) {
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
	if err := d.Set("tags", tagsOSCAPIToMap(netPeering.GetTags())); err != nil {
		return errwrap.Wrapf("Error setting Net Peering tags: {{err}}", err)
	}

	d.SetId(netPeering.GetNetPeeringId())

	return nil
}

func buildOutscaleOAPILinPeeringConnectionFilters(set *schema.Set) oscgo.FiltersNetPeering {
	var filters oscgo.FiltersNetPeering
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
		case "state_names":
			filters.SetStateNames(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
