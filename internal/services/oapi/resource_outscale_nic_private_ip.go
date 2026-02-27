package oapi

import (
	"context"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleNetworkInterfacePrivateIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleNetworkInterfacePrivateIPCreate,
		ReadContext:   ResourceOutscaleNetworkInterfacePrivateIPRead,
		DeleteContext: ResourceOutscaleNetworkInterfacePrivateIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"allow_relink": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"secondary_private_ip_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"nic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"primary_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleNetworkInterfacePrivateIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	input := osc.LinkPrivateIpsRequest{
		NicId: d.Get("nic_id").(string),
	}

	if v, ok := d.GetOk("allow_relink"); ok {
		input.AllowRelink = new(v.(bool))
	}

	if v, ok := d.GetOk("secondary_private_ip_count"); ok {
		input.SecondaryPrivateIpCount = new(v.(int))
	}

	if v, ok := d.GetOk("private_ips"); ok {
		input.PrivateIps = new(utils.InterfaceSliceToStringSlice(v.([]interface{})))
	}

	_, err := client.LinkPrivateIps(ctx, input, options.WithRetryTimeout(timeout))
	if err != nil {
		errString := err.Error()
		return diag.Errorf("failure to assign private ips: %s", errString)
	}

	d.SetId(input.NicId)

	return ResourceOutscaleNetworkInterfacePrivateIPRead(ctx, d, meta)
}

func ResourceOutscaleNetworkInterfacePrivateIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	req := osc.ReadNicsRequest{
		Filters: &osc.FiltersNic{NicIds: &[]string{d.Id()}},
	}

	resp, err := client.ReadNics(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("could not find network interface: %s", err)
	}
	if resp.Nics == nil || utils.IsResponseEmpty(len(*resp.Nics), "NicPrivateIp", d.Id()) {
		d.SetId("")
		return nil
	}
	eni := (*resp.Nics)[0]

	if eni.NicId == "" {
		// Interface is no longer attached, remove from state
		d.SetId("")
		return nil
	}

	var ips []string

	// We need to avoid to store inside private_ips when private IP is the primary IP
	// because the primary can't remove.
	var primaryPrivateID string
	secondary_private_ip_count := 0
	for _, v := range eni.PrivateIps {
		if v.IsPrimary {
			primaryPrivateID = v.PrivateIp
		} else {
			ips = append(ips, v.PrivateIp)
			secondary_private_ip_count += 1
		}
	}

	_, ok := d.GetOk("allow_relink")

	if err := d.Set("allow_relink", ok); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_ips", ips); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secondary_private_ip_count", secondary_private_ip_count); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_id", eni.NicId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("primary_private_ip", primaryPrivateID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceOutscaleNetworkInterfacePrivateIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	input := osc.UnlinkPrivateIpsRequest{
		NicId: d.Id(),
	}

	if v, ok := d.GetOk("private_ips"); ok {
		input.PrivateIps = utils.InterfaceSliceToStringSlice(v.([]interface{}))
	}

	_, err := client.UnlinkPrivateIps(ctx, input, options.WithRetryTimeout(timeout))
	if err != nil {
		errString := err.Error()
		return diag.Errorf("failure to unassign private ips: %s", errString)
	}
	d.SetId("")
	return nil
}
