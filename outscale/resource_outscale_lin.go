package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPINet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINetCreate,
		Read:   resourceOutscaleOAPINetRead,
		Update: resourceOutscaleOAPINetUpdate,
		Delete: resourceOutscaleOAPINetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getOAPINetSchema(),
	}
}

func resourceOutscaleOAPINetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateNetRequest{
		IpRange: d.Get("ip_range").(string),
	}

	if c, ok := d.GetOk("tenancy"); ok {
		tenancy := c.(string)
		if tenancy == "default" || tenancy == "dedicated" {
			req.SetTenancy(tenancy)
		} else {
			return fmt.Errorf("tenancy option not supported: %s", tenancy)
		}
	}

	var resp oscgo.CreateNetResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, _, err = conn.NetApi.CreateNet(context.Background(), &oscgo.CreateNetOpts{CreateNetRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("error creating Outscale Net: %s", utils.GetErrorResponse(err))
	}

	//SetTags
	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.([]interface{}), resp.Net.GetNetId(), conn)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.Net.GetNetId())

	return resourceOutscaleOAPINetRead(d, meta)
}

func resourceOutscaleLinRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

	req := &fcu.DescribeVpcsInput{
		VpcIds: []*string{aws.String(id)},
	}

	var resp *fcu.DescribeVpcsOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpcs(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", err)
	}

	if resp == nil {
		d.SetId("")
		return fmt.Errorf("Lin not found")
	}

	if len(resp.Vpcs) == 0 {
		d.SetId("")
		return fmt.Errorf("Lin not found")
	}

	d.Set("cidr_block", resp.Vpcs[0].CidrBlock)
	d.Set("instance_tenancy", resp.Vpcs[0].InstanceTenancy)
	d.Set("dhcp_options_id", resp.Vpcs[0].DhcpOptionsId)
	d.Set("request_id", resp.RequestId)
	d.Set("state", resp.Vpcs[0].State)
	d.Set("vpc_id", resp.Vpcs[0].VpcId)

	return d.Set("tag_set", tagsToMap(resp.Vpcs[0].Tags))
}

func resourceOutscaleOAPINetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	filters := oscgo.FiltersNet{
		NetIds: &[]string{id},
	}

	req := oscgo.ReadNetsRequest{
		Filters: &filters,
	}

	var resp oscgo.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, _, err = conn.NetApi.ReadNets(context.Background(), &oscgo.ReadNetsOpts{ReadNetsRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading network (%s)", err)
	}

	if len(resp.GetNets()) == 0 {
		d.SetId("")
		return fmt.Errorf("oAPI network not found")
	}

	d.Set("ip_range", resp.GetNets()[0].GetIpRange())
	d.Set("tenancy", resp.GetNets()[0].Tenancy)
	d.Set("dhcp_options_set_id", resp.GetNets()[0].GetDhcpOptionsSetId())
	d.Set("net_id", resp.GetNets()[0].GetNetId())
	d.Set("state", resp.GetNets()[0].GetState())
	d.Set("request_id", resp.ResponseContext.GetRequestId())
	return d.Set("tags", tagsOSCAPIToMap(resp.GetNets()[0].GetTags()))
}

func resourceOutscaleOAPINetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)
	return resourceOutscaleOAPINetRead(d, meta)
}

func resourceOutscaleOAPINetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	id := d.Id()

	req := oscgo.DeleteNetRequest{
		NetId: id,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, _, err = conn.NetApi.DeleteNet(context.Background(), &oscgo.DeleteNetOpts{DeleteNetRequest: optional.NewInterface(req)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func getOAPINetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_range": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
		"tenancy": {
			Type:     schema.TypeString,
			ForceNew: true,
			Computed: true,
			Optional: true,
		},

		// Attributes
		"dhcp_options_set_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": tagsListOAPISchema(),
		"net_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
