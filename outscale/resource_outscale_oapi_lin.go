package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"

	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPINet() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINetCreate,
		Read:   resourceOutscaleOAPINetRead,
		Delete: resourceOutscaleOAPINetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getOAPINetSchema(),
	}
}

func resourceOutscaleOAPINetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	req := &oapi.CreateNetRequest{
		IpRange: d.Get("ip_range").(string),
	}

	if c, ok := d.GetOk("tenancy"); ok {
		tenancy := c.(string)
		if tenancy == "default" || tenancy == "dedicated" {
			req.Tenancy = tenancy
		} else {
			return fmt.Errorf("tenancy option not supported: %s", tenancy)
		}
	}

	var resp *oapi.POST_CreateNetResponses
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_CreateNet(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("Error creating Outscale Net: %s", errString)
	}

	net := resp.OK.Net

	d.SetId(net.NetId)

	//SetTags
	if d.IsNewResource() {
		if err := setOAPITags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

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
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	filters := oapi.FiltersNet{
		NetIds: []string{id},
	}

	req := oapi.ReadNetsRequest{
		Filters: filters,
	}

	var rs *oapi.POST_ReadNetsResponses
	var resp *oapi.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rs, err = conn.POST_ReadNets(req)

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

	resp = rs.OK

	if resp == nil {
		d.SetId("")
		return fmt.Errorf("oAPI network not found")
	}

	if len(resp.Nets) == 0 {
		d.SetId("")
		return fmt.Errorf("oAPI network not found")
	}

	d.Set("ip_range", resp.Nets[0].IpRange)
	d.Set("tenancy", resp.Nets[0].Tenancy)
	d.Set("dhcp_options_set_id", resp.Nets[0].DhcpOptionsSetId)
	d.Set("net_id", resp.Nets[0].NetId)
	d.Set("state", resp.Nets[0].State)
	d.Set("request_id", resp.ResponseContext.RequestId)
	return d.Set("tags", tagsOAPIToMap(resp.Nets[0].Tags))

}

func resourceOutscaleOAPINetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	req := oapi.DeleteNetRequest{
		NetId: id,
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.POST_DeleteNet(req)

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
