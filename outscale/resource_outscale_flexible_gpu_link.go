package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIFlexibleGpuLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIFlexibleGpuLinkCreate,
		Read:   resourceOutscaleOAPIFlexibleGpuLinkRead,
		Update: resourceFlexibleGpuLinkUpdate,
		Delete: resourceOutscaleOAPIFlexibleGpuLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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

func resourceOutscaleOAPIFlexibleGpuLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	vmId := d.Get("vm_id").(string)
	GpuIdsList := utils.SetToStringSlice(d.Get("flexible_gpu_ids").(*schema.Set))

	for _, flexGpuID := range GpuIdsList {
		var resp oscgo.LinkFlexibleGpuResponse
		reqLink := oscgo.LinkFlexibleGpuRequest{
			FlexibleGpuId: flexGpuID,
			VmId:          vmId,
		}
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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
			return fmt.Errorf("Error Link flexibe gpu: %s", err.Error())
		}
		if !resp.HasResponseContext() {
			return fmt.Errorf("Error there is not Link flexible gpu (%s)", err)
		}
	}

	if err := changeShutdownBehavior(conn, vmId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("Unable to change ShutdownBehavior: %s\n", err)
	}

	return resourceOutscaleOAPIFlexibleGpuLinkRead(d, meta)
}

func resourceOutscaleOAPIFlexibleGpuLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	vmId := d.Get("vm_id").(string)
	req := &oscgo.ReadFlexibleGpusRequest{
		Filters: &oscgo.FiltersFlexibleGpu{
			VmIds: &[]string{vmId},
		},
	}
	var resp oscgo.ReadFlexibleGpusResponse
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
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
	d.SetId(resource.UniqueId())
	return nil
}

func resourceOutscaleOAPIFlexibleGpuLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	GpuIdsList := utils.SetToStringSlice(d.Get("flexible_gpu_ids").(*schema.Set))
	vmId := d.Get("vm_id").(string)
	var err error

	for _, flexGpuID := range GpuIdsList {
		req := &oscgo.UnlinkFlexibleGpuRequest{
			FlexibleGpuId: flexGpuID,
		}
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			rp, httpResp, err := conn.FlexibleGpuApi.ReadFlexibleGpus(context.Background()).
				ReadFlexibleGpusRequest(*reqFlex).Execute()

			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
			return fmt.Errorf("error reading the FlexibleGpu %s", err)
		}

		if len(*resp.FlexibleGpus) != 1 {
			return fmt.Errorf("Unable to find Flexible GPU")
		}
		if (*resp.FlexibleGpus)[0].GetState() != "detaching" {
			return fmt.Errorf("Unable to unlink Flexible GPU")
		}
	}

	if err := changeShutdownBehavior(conn, vmId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("Unable to change ShutdownBehavior: %s\n", err)
	}

	d.SetId("")
	return nil
}

func resourceFlexibleGpuLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
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
			err = resource.Retry(20*time.Second, func() *resource.RetryError {
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
			err = resource.Retry(20*time.Second, func() *resource.RetryError {
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
	if err := changeShutdownBehavior(conn, vmId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("Unable to change ShutdownBehavior: %s\n", err)
	}

	return resourceOutscaleOAPIFlexibleGpuLinkRead(d, meta)
}

func changeShutdownBehavior(conn *oscgo.APIClient, vmId string, timeOut time.Duration) error {
	var resp oscgo.ReadVmsResponse
	err := resource.Retry(timeOut, func() *resource.RetryError {
		rp, httpResp, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{vmId},
			}}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading the VM %s", err)
	}
	if len(resp.GetVms()) == 0 {
		return fmt.Errorf("error reading the VM %s err %s ", vmId, err)
	}
	vm := resp.GetVms()[0]

	shutdownBehOpt := vm.GetVmInitiatedShutdownBehavior()
	if shutdownBehOpt != "stop" {
		sbOpts := oscgo.UpdateVmRequest{VmId: vm.GetVmId()}
		sbOpts.SetVmInitiatedShutdownBehavior("stop")
		if err := updateVmAttr(conn, sbOpts); err != nil {
			return err
		}
	}

	if err := stopVM(vmId, conn, timeOut); err != nil {
		return err
	}

	if shutdownBehOpt != "stop" {
		sbReq := oscgo.UpdateVmRequest{VmId: vmId}
		sbReq.SetVmInitiatedShutdownBehavior(shutdownBehOpt)
		if err = updateVmAttr(conn, sbReq); err != nil {
			return err
		}
	}

	if err := startVM(vmId, conn, timeOut); err != nil {
		return err
	}
	return nil
}
