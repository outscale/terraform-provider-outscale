package outscale

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
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
				Type:     schema.TypeMap,
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
			"tags": tagsOAPIListSchemaComputed(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPILinPeeringConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	log.Printf("[DEBUG] Reading Net Peering Connections.")

	req := oapi.ReadNetPeeringsRequest{}

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("filters must be assigned")
	}
	req.Filters = buildOutscaleOAPILinPeeringConnectionFilters(filters.(*schema.Set))

	var resp *oapi.POST_ReadNetPeeringsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadNetPeerings(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil || resp.OK == nil {
		if strings.Contains(fmt.Sprint(err), "InvalidNetPeeringConnectionID.NotFound") {
			return fmt.Errorf("no matching Net Peering Connection found")
		}
		return fmt.Errorf("Error reading Net Peering Connection details: %s", err)
	}

	if len(resp.OK.NetPeerings) > 1 {
		return fmt.Errorf("multiple Net Peering connections matched; use additional constraints to reduce matches to a single Net Peering Connection")
	}
	netPeering := resp.OK.NetPeerings[0]

	// The failed status is a status that we can assume just means the
	// connection is gone. Destruction isn't allowed, and it eventually
	// just "falls off" the console. See GH-2322
	if !reflect.DeepEqual(netPeering.State, oapi.NetPeeringState{}) {
		status := map[string]bool{
			"deleted":  true,
			"deleting": true,
			"expired":  true,
			"failed":   true,
			"rejected": true,
		}
		if _, ok := status[netPeering.State.Name]; ok {
			log.Printf("[DEBUG] Net Peering Connection (%s) in state (%s), removing.",
				d.Id(), netPeering.State.Name)
			return nil
		}
	}
	log.Printf("[DEBUG] Net Peering Connection response: %#v", netPeering)

	log.Printf("[DEBUG] Net Peering Connection Source %s, Accepter %s", netPeering.SourceNet.AccountId, netPeering.AccepterNet.AccountId)

	accepter := make(map[string]interface{})
	requester := make(map[string]interface{})
	stat := make(map[string]interface{})

	if !reflect.DeepEqual(netPeering.AccepterNet, oapi.AccepterNet{}) {
		accepter["ip_range"] = netPeering.AccepterNet.IpRange
		accepter["account_id"] = netPeering.AccepterNet.AccountId
		accepter["net_id"] = netPeering.AccepterNet.NetId
	}
	if !reflect.DeepEqual(netPeering.SourceNet, oapi.SourceNet{}) {
		requester["ip_range"] = netPeering.SourceNet.IpRange
		requester["account_id"] = netPeering.SourceNet.AccountId
		requester["net_id"] = netPeering.SourceNet.NetId
	}
	if netPeering.State.Name != "" {
		stat["name"] = netPeering.State.Name
		stat["message"] = netPeering.State.Message
	}

	if err := d.Set("accepter_net", accepter); err != nil {
		return err
	}
	if err := d.Set("source_net", requester); err != nil {
		return err
	}
	if err := d.Set("state", stat); err != nil {
		return err
	}
	if err := d.Set("net_peering_id", netPeering.NetPeeringId); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOAPIToMap(netPeering.Tags)); err != nil {
		return errwrap.Wrapf("Error setting Net Peering tags: {{err}}", err)
	}
	if err := d.Set("request_id", resp.OK.ResponseContext.RequestId); err != nil {
		return err
	}

	d.SetId(netPeering.NetPeeringId)

	return nil
}

func buildOutscaleOAPILinPeeringConnectionFilters(set *schema.Set) oapi.FiltersNetPeering {
	var filters oapi.FiltersNetPeering
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "accepter_net_account_ids":
			filters.AccepterNetAccountIds = filterValues
		case "accepter_net_ip_ranges":
			filters.AccepterNetIpRanges = filterValues
		case "accepter_net_net_ids":
			filters.AccepterNetNetIds = filterValues
		case "net_peering_ids":
			filters.NetPeeringIds = filterValues
		case "source_net_account_ids":
			filters.SourceNetAccountIds = filterValues
		case "source_net_ip_ranges":
			filters.SourceNetIpRanges = filterValues
		case "source_net_net_ids":
			filters.SourceNetNetIds = filterValues
		case "state_messages":
			filters.StateMessages = filterValues
		case "state_names":
			filters.StateNames = filterValues
		case "tag_keys":
			filters.TagKeys = filterValues
		case "tag_values":
			filters.TagValues = filterValues
		case "tags":
			filters.Tags = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
