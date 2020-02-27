package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func datasourceOutscaleOApiVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOApiVMSRead,

		Schema: datasourceOutscaleOApiVMSSchema(),
	}
}

func dataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func datasourceOutscaleOApiVMSSchema() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"vms": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: getOApiVMAttributesSchema(),
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return wholeSchema
}

func dataSourceOutscaleOApiVMSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	vmID, vmIDOk := d.GetOk("vm_id")

	if !filtersOk && !vmIDOk {
		return fmt.Errorf("One of filters, and vm ID must be assigned")
	}

	// Build up search parameters
	params := oscgo.ReadVmsRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVMFilters(filters.(*schema.Set))
	}
	if vmIDOk {
		params.Filters.VmIds = &[]string{vmID.(string)}
	}

	var resp oscgo.ReadVmsResponse
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		r, _, err := client.VmApi.ReadVms(context.Background(), &oscgo.ReadVmsOpts{
			ReadVmsRequest: optional.NewInterface(params),
		})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		resp = r
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading the VMs %s", err)
	}

	// If no instances were returned, return
	if !resp.HasVms() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var filteredVms []oscgo.Vm

	// loop through reservations, and remove terminated instances, populate vm slice
	for _, res := range resp.GetVms() {
		if res.GetState() != "terminated" {
			filteredVms = append(filteredVms, res)
		}
	}

	if err := d.Set("request_id", resp.GetResponseContext().RequestId); err != nil {
		return err
	}

	if len(filteredVms) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	d.SetId(resource.UniqueId())
	return d.Set("vms", dataSourceOAPIVMS(filteredVms))
}

func dataSourceOAPIVMS(i []oscgo.Vm) []map[string]interface{} {
	vms := make([]map[string]interface{}, len(i))
	for index, v := range i {
		vm := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			vm[key] = value
			return nil
		}

		if err := oapiVMDescriptionAttributes(setterFunc, &v); err != nil {
			log.Fatalf("[DEBUG] oapiVMDescriptionAttributes ERROR %+v", err)
		}

		vm["tags"] = getOscAPITagSet(v.GetTags())
		vms[index] = vm
	}
	return vms
}
