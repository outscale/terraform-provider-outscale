package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleFlexibleGpuLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleFlexibleGpuLinkCreate,
		ReadContext:   ResourceOutscaleFlexibleGpuLinkRead,
		UpdateContext: resourceFlexibleGpuLinkUpdate,
		DeleteContext: ResourceOutscaleFlexibleGpuLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"flexible_gpu_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleFlexibleGpuLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)
	vmId := d.Get("vm_id").(string)
	GpuIdsList := utils.SetToStringSlice(d.Get("flexible_gpu_ids").(*schema.Set))

	for _, flexGpuID := range GpuIdsList {
		var resp osc.LinkFlexibleGpuResponse
		reqLink := osc.LinkFlexibleGpuRequest{
			FlexibleGpuId: flexGpuID,
			VmId:          vmId,
		}
		_, err := client.LinkFlexibleGpu(ctx, reqLink, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.Errorf("error link flexibe gpu: %s", err.Error())
		}
		if resp.ResponseContext == nil {
			return diag.Errorf("error there is not link flexible gpu (%s)", err)
		}
	}

	if err := changeShutdownBehavior(ctx, client, vmId, timeout); err != nil {
		return diag.Errorf("unable to change shutdownbehavior: %s", err)
	}

	return ResourceOutscaleFlexibleGpuLinkRead(ctx, d, meta)
}

func ResourceOutscaleFlexibleGpuLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	vmId := d.Get("vm_id").(string)

	req := &osc.ReadFlexibleGpusRequest{
		Filters: &osc.FiltersFlexibleGpu{
			VmIds: &[]string{vmId},
		},
	}
	resp, err := client.ReadFlexibleGpus(ctx, *req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}
	if utils.IsResponseEmpty(len(ptr.From(resp.FlexibleGpus)), "FlexibleGpuLink", d.Id()) {
		d.SetId("")
		return nil
	}
	flexGpus := ptr.From(resp.FlexibleGpus)[:]
	readGpuIdsLink := make([]string, len(flexGpus))
	for k, flexGpu := range flexGpus {
		readGpuIdsLink[k] = *flexGpu.FlexibleGpuId
	}
	if err := d.Set("flexible_gpu_ids", readGpuIdsLink); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_id", vmId); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id.UniqueId())
	return nil
}

func ResourceOutscaleFlexibleGpuLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)
	GpuIdsList := utils.SetToStringSlice(d.Get("flexible_gpu_ids").(*schema.Set))
	vmId := d.Get("vm_id").(string)

	for _, flexGpuID := range GpuIdsList {
		req := &osc.UnlinkFlexibleGpuRequest{
			FlexibleGpuId: flexGpuID,
		}
		_, err := client.UnlinkFlexibleGpu(ctx, *req, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.FromErr(err)
		}

		reqFlex := osc.ReadFlexibleGpusRequest{
			Filters: &osc.FiltersFlexibleGpu{
				FlexibleGpuIds: &[]string{flexGpuID},
			},
		}
		resp, err := client.ReadFlexibleGpus(ctx, reqFlex, options.WithRetryTimeout(timeout))
		if err != nil {
			return diag.Errorf("error reading the flexiblegpu %s", err)
		}

		if len(*resp.FlexibleGpus) != 1 {
			return diag.Errorf("unable to find flexible gpu")
		}
		if (*(*resp.FlexibleGpus)[0].State) != osc.FlexibleGpuStateDetaching &&
			(*(*resp.FlexibleGpus)[0].State) != osc.FlexibleGpuStateAllocated {
			return diag.Errorf("unable to unlink flexible gpu")
		}
	}

	if err := changeShutdownBehavior(ctx, client, vmId, timeout); err != nil {
		return diag.Errorf("unable to change shutdownbehavior: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceFlexibleGpuLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutUpdate)
	vmId := d.Get("vm_id").(string)
	oldIds, newIds := d.GetChange("flexible_gpu_ids")

	interIds := oldIds.(*schema.Set).Intersection(newIds.(*schema.Set))
	toCreate := newIds.(*schema.Set).Difference(interIds)
	toRemove := oldIds.(*schema.Set).Difference(interIds)

	if toRemove.Len() > 0 {
		for _, flexGpuID := range utils.SetToStringSlice(toRemove) {
			req := &osc.UnlinkFlexibleGpuRequest{
				FlexibleGpuId: flexGpuID,
			}
			_, err := client.UnlinkFlexibleGpu(ctx, *req, options.WithRetryTimeout(timeout))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if toCreate.Len() > 0 {
		for _, flexGpuID := range utils.SetToStringSlice(toCreate) {
			req := &osc.LinkFlexibleGpuRequest{
				FlexibleGpuId: flexGpuID,
				VmId:          vmId,
			}
			_, err := client.LinkFlexibleGpu(ctx, *req, options.WithRetryTimeout(timeout))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if err := changeShutdownBehavior(ctx, client, vmId, timeout); err != nil {
		return diag.Errorf("unable to change shutdownbehavior: %s", err)
	}

	return ResourceOutscaleFlexibleGpuLinkRead(ctx, d, meta)
}

func changeShutdownBehavior(ctx context.Context, client *osc.Client, vmId string, timeout time.Duration) error {
	resp, err := client.ReadVms(ctx, osc.ReadVmsRequest{
		Filters: &osc.FiltersVm{
			VmIds: &[]string{vmId},
		},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return fmt.Errorf("error reading the vm %s", err)
	}
	if len(ptr.From(resp.Vms)) == 0 {
		return fmt.Errorf("error reading the vm %s err %s ", vmId, err)
	}
	vm := ptr.From(resp.Vms)[0]

	shutdownBehOpt := vm.VmInitiatedShutdownBehavior
	if shutdownBehOpt != "stop" {
		sbOpts := osc.UpdateVmRequest{VmId: vm.VmId}
		sbOpts.VmInitiatedShutdownBehavior = new("stop")
		if err := updateVmAttr(ctx, client, timeout, sbOpts); err != nil {
			return err
		}
	}

	if err := stopVM(ctx, client, timeout, vmId); err != nil {
		return err
	}

	if shutdownBehOpt != "stop" {
		sbReq := osc.UpdateVmRequest{VmId: vmId}
		sbReq.VmInitiatedShutdownBehavior = new(shutdownBehOpt)
		if err = updateVmAttr(ctx, client, timeout, sbReq); err != nil {
			return err
		}
	}

	if err := startVM(ctx, client, timeout, vmId); err != nil {
		return err
	}
	return nil
}
