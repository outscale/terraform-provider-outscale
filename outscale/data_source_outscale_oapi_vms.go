package outscale

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func datasourceOutscaleOApiVMS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOApiVMSRead,

		Schema: datasourceOutscaleOApiVMSSchema(),
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
	client := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	vmID, vmIDOk := d.GetOk("vm_id")

	if filtersOk == false && vmIDOk == false {
		return fmt.Errorf("One of filters, and vm ID must be assigned")
	}

	// Build up search parameters
	params := oapi.ReadVmsRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVmFilters(filters.(*schema.Set))
	}
	if vmIDOk {
		params.Filters.VmIds = []string{vmID.(string)}
	}

	var resp *oapi.POST_ReadVmsResponses
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, err = client.POST_ReadVms(params)
		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error reading the VM %s", err)
	}

	if resp.OK.Vms == nil {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	// If no instances were returned, return
	if len(resp.OK.Vms) == 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var filteredInstances []oapi.Vm

	// loop through reservations, and remove terminated instances, populate instance slice
	for _, res := range resp.OK.Vms {
		if res.State != "terminated" {
			filteredInstances = append(filteredInstances, res)
		}
	}

	d.Set("request_id", resp.OK.ResponseContext.RequestId)

	if len(filteredInstances) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	return vmsOAPIDescriptionAttributes(d, filteredInstances, client)
}

// Populate instance attribute fields with the returned instance
func vmsOAPIDescriptionAttributes(d *schema.ResourceData, instances []oapi.Vm, conn *oapi.Client) error {
	d.Set("vms", dataSourceOAPIVMS(instances))
	return nil
}

func dataSourceOAPIVMS(i []oapi.Vm) *schema.Set {
	s := &schema.Set{}
	for _, v := range i {
		instance := make(map[string]interface{})

		setterFunc := func(key string, value interface{}) error {
			instance[key] = value
			return nil
		}

		oapiVMDescriptionAttributes(setterFunc, &v)

		fmt.Println("schema set -> ", s)
		fmt.Println("instance -> ", s)
		fmt.Printf("m -> %+v\n", v)

		s.Add(instance)
	}
	return s
}

func dataSourceFiltersOApiSchema() *schema.Schema {
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
