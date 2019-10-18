package outscale

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func buildOutscaleDataSourceFilters(set *schema.Set) []*fcu.Filter {
	var filters []*fcu.Filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []*string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, aws.String(e.(string)))
		}
		filters = append(filters, &fcu.Filter{
			Name:   aws.String(m["name"].(string)),
			Values: filterValues,
		})
	}
	return filters
}
