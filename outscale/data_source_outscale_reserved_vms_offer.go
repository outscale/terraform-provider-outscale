package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

<<<<<<< HEAD
func dataSourceOutscaleSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleSubnetRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cidr_block": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
=======
func dataSourceOutscaleReservedVMOffer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleReservedVMOfferRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"reserved_instances_offering_id": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"pricing_details_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_tenancy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"offering_type": &schema.Schema{
>>>>>>> TPD-451
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
<<<<<<< HEAD

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tag_set": dataSourceTagsSchema(),

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"available_ip_address_count": {
				Type:     schema.TypeInt,
=======
			"product_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"currency_code": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"marketplace": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"recurring_charges": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
>>>>>>> TPD-451
				Computed: true,
			},
		},
	}
}

<<<<<<< HEAD
func dataSourceOutscaleSubnetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	req := &fcu.DescribeSubnetsInput{}

	if id := d.Get("subnet_id"); id != "" {
		req.SubnetIds = []*string{aws.String(id.(string))}
	}

	filters, filtersOk := d.GetOk("filter")

	if filtersOk {
		req.Filters = buildOutscaleDataSourceFilters(filters.(*schema.Set))
	}

	var resp *fcu.DescribeSubnetsOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeSubNet(req)
=======
func dataSourceOutscaleReservedVMOfferRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	az, azok := d.GetOk("availability_zone")
	it, itok := d.GetOk("instance_tenancy")
	ity, ityok := d.GetOk("instance_type")
	pd, pdok := d.GetOk("product_description")
	ot, otok := d.GetOk("offering_type")
	ri, riok := d.GetOk("reserved_instances_offering_id")
	filter, filterOk := d.GetOk("filter")

	req := &fcu.DescribeReservedInstancesOfferingsInput{}

	if azok {
		req.AvailabilityZone = aws.String(az.(string))
	}
	if otok {
		req.OfferingType = aws.String(ot.(string))
	}
	if itok {
		req.InstanceTenancy = aws.String(it.(string))
	}
	if ityok {
		req.InstanceTenancy = aws.String(ity.(string))
	}
	if pdok {
		req.InstanceTenancy = aws.String(pd.(string))
	}
	if riok {
		var ids []*string
		for _, v := range ri.([]interface{}) {
			ids = append(ids, aws.String(v.(string)))
		}
		req.ReservedInstancesOfferingId = ids
	}
	if filterOk {
		req.Filters = buildOutscaleDataSourceFilters(filter.(*schema.Set))
	}

	var resp *fcu.DescribeReservedInstancesOfferingsOutput
	err := resource.Retry(60*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeReservedInstancesOfferings(req)
>>>>>>> TPD-451
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
		}
<<<<<<< HEAD

		return resource.NonRetryableError(err)
=======
		return nil
>>>>>>> TPD-451
	})

	if err != nil {
		return err
	}

<<<<<<< HEAD
	if resp == nil || len(resp.Subnets) == 0 {
		return fmt.Errorf("no matching subnet found")
	}

	if len(resp.Subnets) > 1 {
		return fmt.Errorf("multiple subnets matched; use additional constraints to reduce matches to a single subnet")
	}

	subnet := resp.Subnets[0]

	d.SetId(*subnet.SubnetId)
	d.Set("subnet_id", subnet.SubnetId)
	d.Set("vpc_id", subnet.VpcId)
	d.Set("availability_zone", subnet.AvailabilityZone)
	d.Set("cidr_block", subnet.CidrBlock)
	d.Set("state", subnet.State)
	d.Set("tag_set", tagsToMap(subnet.Tags))
	d.Set("available_ip_address_count", subnet.AvailableIpAddressCount)
=======
	if resp == nil || len(resp.ReservedInstancesOfferingsSet) == 0 {
		return fmt.Errorf("no matching reserved VMS Offer found")
	}

	if len(resp.ReservedInstancesOfferingsSet) > 1 {
		return fmt.Errorf("multiple VM Offer matched; use additional constraints to reduce matches to a single VM Offer")
	}

	d.SetId(resource.UniqueId())

	v := resp.ReservedInstancesOfferingsSet[0]

	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("currency_code", v.CurrencyCode)
	d.Set("instance_tenancy", v.InstanceTenancy)
	d.Set("instance_type", v.InstanceType)
	d.Set("marketplace", v.Martketplace)
	d.Set("offering_type", v.OfferingType)
	d.Set("product_description", v.ProductDescription)
	d.Set("reserved_instances_offering_id", v.ReservedInstancesOfferingId)

	rcs := make([]map[string]interface{}, len(v.RecurringCharges))
	for k1, v1 := range v.RecurringCharges {
		rc := make(map[string]interface{})
		rc["frequency"] = v1.Frequency
		rcs[k1] = rc
	}

	d.Set("recurring_charges", rcs)

	pds := make([]map[string]interface{}, len(v.PricingDetailsSet))
	for k1, v1 := range v.PricingDetailsSet {
		rc := make(map[string]interface{})
		rc["count"] = v1.Count
		rcs[k1] = rc
	}

	d.Set("pricing_details_set", pds)

>>>>>>> TPD-451
	d.Set("request_id", resp.RequestId)

	return nil
}
