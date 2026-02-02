package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleFlexibleGpuLink() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleFlexibleGpuLinkCreate,
		Read:   ResourceOutscaleFlexibleGpuLinkRead,
		Update: resourceFlexibleGpuLinkUpdate,
		Delete: ResourceOutscaleFlexibleGpuLinkDelete,
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

func ResourceOutscaleFlexibleGpuLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutCreate)
	vmId := d.Get("vm_id").(string)
	GpuIdsList := utils.SetToStringSlice(d.Get("flexible_gpu_ids").(*schema.Set))

	for _, flexGpuID := range GpuIdsList {
		var resp oscgo.LinkFlexibleGpuResponse
		reqLink := oscgo.LinkFlexibleGpuRequest{
			FlexibleGpuId: flexGpuID,
			VmId:          vmId,
		}
		err := retry.Retry(timeout, func() *retry.RetryError {
			var err error
			rp, httpResp, err := conn.FlexibleGpuApi.LinkFlexibleGpu(
				context.Background()).LinkFlexibleGpuRequest(reqLink).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("error link flexibe gpu: %s", err.Error())
		}
		if !resp.HasResponseContext() {
			return fmt.Errorf("error there is not link flexible gpu (%s)", err)
		}
	}

	if err := changeShutdownBehavior(conn, vmId, timeout); err != nil {
		return fmt.Errorf("unable to change shutdownbehavior: %s", err)
	}

	return ResourceOutscaleFlexibleGpuLinkRead(d, meta)
}

func ResourceOutscaleFlexibleGpuLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutRead)
	vmId := d.Get("vm_id").(string)
	req := &oscgo.ReadFlexibleGpusRequest{
		Filters: &oscgo.FiltersFlexibleGpu{
			VmIds: &[]string{vmId},
		},
	}
	var resp oscgo.ReadFlexibleGpusResponse
	err := retry.Retry(timeout, func() *retry.RetryError {
		rp, httpResp, err := conn.FlexibleGpuApi.ReadFlexibleGpus(
			context.Background()).
			ReadFlexibleGpusRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if utils.IsResponseEmpty(len(resp.GetFlexibleGpus()), "FlexibleGpuLink", d.Id()) {
		d.SetId("")
		return nil
	}
	flexGpus := resp.GetFlexibleGpus()[:]
	readGpuIdsLink := make([]string, len(flexGpus))
	for k, flexGpu := range flexGpus {
		readGpuIdsLink[k] = flexGpu.GetFlexibleGpuId()
	}
	if err := d.Set("flexible_gpu_ids", readGpuIdsLink); err != nil {
		return err
	}
	if err := d.Set("vm_id", vmId); err != nil {
		return err
	}
	d.SetId(id.UniqueId())
	return nil
}

func ResourceOutscaleFlexibleGpuLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutDelete)
	GpuIdsList := utils.SetToStringSlice(d.Get("flexible_gpu_ids").(*schema.Set))
	vmId := d.Get("vm_id").(string)
	var err error

	for _, flexGpuID := range GpuIdsList {
		req := &oscgo.UnlinkFlexibleGpuRequest{
			FlexibleGpuId: flexGpuID,
		}
		err = retry.Retry(timeout, func() *retry.RetryError {
			_, httpResp, err := conn.FlexibleGpuApi.UnlinkFlexibleGpu(
				context.Background()).UnlinkFlexibleGpuRequest(*req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		var resp oscgo.ReadFlexibleGpusResponse
		reqFlex := &oscgo.ReadFlexibleGpusRequest{
			Filters: &oscgo.FiltersFlexibleGpu{
				FlexibleGpuIds: &[]string{flexGpuID},
			},
		}
		err = retry.Retry(timeout, func() *retry.RetryError {
			rp, httpResp, err := conn.FlexibleGpuApi.ReadFlexibleGpus(context.Background()).
				ReadFlexibleGpusRequest(*reqFlex).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("error reading the flexiblegpu %s", err)
		}

		if len(*resp.FlexibleGpus) != 1 {
			return fmt.Errorf("unable to find flexible gpu")
		}
		if (*resp.FlexibleGpus)[0].GetState() != "detaching" &&
			(*resp.FlexibleGpus)[0].GetState() != "allocated" {
			return fmt.Errorf("unable to unlink flexible gpu")
		}
	}

	if err := changeShutdownBehavior(conn, vmId, timeout); err != nil {
		return fmt.Errorf("unable to change shutdownbehavior: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceFlexibleGpuLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI
	timeout := d.Timeout(schema.TimeoutUpdate)
	vmId := d.Get("vm_id").(string)
	oldIds, newIds := d.GetChange("flexible_gpu_ids")

	interIds := oldIds.(*schema.Set).Intersection(newIds.(*schema.Set))
	toCreate := newIds.(*schema.Set).Difference(interIds)
	toRemove := oldIds.(*schema.Set).Difference(interIds)
	var err error

	if toRemove.Len() > 0 {
		for _, flexGpuID := range utils.SetToStringSlice(toRemove) {
			req := &oscgo.UnlinkFlexibleGpuRequest{
				FlexibleGpuId: flexGpuID,
			}
			err = retry.Retry(timeout, func() *retry.RetryError {
				_, httpResp, err := conn.FlexibleGpuApi.UnlinkFlexibleGpu(
					context.Background()).UnlinkFlexibleGpuRequest(*req).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	if toCreate.Len() > 0 {
		for _, flexGpuID := range utils.SetToStringSlice(toCreate) {
			req := &oscgo.LinkFlexibleGpuRequest{
				FlexibleGpuId: flexGpuID,
				VmId:          vmId,
			}
			err = retry.Retry(timeout, func() *retry.RetryError {
				_, httpResp, err := conn.FlexibleGpuApi.LinkFlexibleGpu(
					context.Background()).LinkFlexibleGpuRequest(*req).Execute()
				if err != nil {
					return utils.CheckThrottling(httpResp, err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	if err := changeShutdownBehavior(conn, vmId, timeout); err != nil {
		return fmt.Errorf("unable to change shutdownbehavior: %s", err)
	}

	return ResourceOutscaleFlexibleGpuLinkRead(d, meta)
}

func changeShutdownBehavior(conn *oscgo.APIClient, vmId string, timeout time.Duration) error {
	var resp oscgo.ReadVmsResponse
	err := retry.Retry(timeout, func() *retry.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{vmId},
			},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading the vm %s", err)
	}
	if len(resp.GetVms()) == 0 {
		return fmt.Errorf("error reading the vm %s err %s ", vmId, err)
	}
	vm := resp.GetVms()[0]

	shutdownBehOpt := vm.GetVmInitiatedShutdownBehavior()
	if shutdownBehOpt != "stop" {
		sbOpts := oscgo.UpdateVmRequest{VmId: vm.GetVmId()}
		sbOpts.SetVmInitiatedShutdownBehavior("stop")
		if err := updateVmAttr(conn, timeout, sbOpts); err != nil {
			return err
		}
	}

	if err := stopVM(vmId, conn, timeout); err != nil {
		return err
	}

	if shutdownBehOpt != "stop" {
		sbReq := oscgo.UpdateVmRequest{VmId: vmId}
		sbReq.SetVmInitiatedShutdownBehavior(shutdownBehOpt)
		if err = updateVmAttr(conn, timeout, sbReq); err != nil {
			return err
		}
	}

	if err := startVM(vmId, conn, timeout); err != nil {
		return err
	}
	return nil
}
