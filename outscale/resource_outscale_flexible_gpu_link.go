package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIFlexibleGpuLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIFlexibleGpuLinkCreate,
		Read:   resourceOutscaleOAPIFlexibleGpuLinkRead,
		Delete: resourceOutscaleOAPIFlexibleGpuLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"flexible_gpu_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	flexGpuID := d.Get("flexible_gpu_id").(string)
	vmId := d.Get("vm_id").(string)

	filter := &oscgo.FiltersFlexibleGpu{
		FlexibleGpuIds: &[]string{flexGpuID},
	}
	reqFlex := &oscgo.ReadFlexibleGpusRequest{
		Filters: filter,
	}
	reqLink := oscgo.LinkFlexibleGpuRequest{
		FlexibleGpuId: flexGpuID,
		VmId:          vmId,
	}

	resp, _, err := conn.FlexibleGpuApi.LinkFlexibleGpu(
		context.Background()).LinkFlexibleGpuRequest(reqLink).Execute()
	if err != nil {
		return fmt.Errorf("Error Link flexibe gpu: %s", err.Error())
	}

	if !resp.HasResponseContext() {
		return fmt.Errorf("Error there is not Link flexible gpu (%s)", err)
	}

	respGpu, _, err := conn.FlexibleGpuApi.ReadFlexibleGpus(context.Background()).
		ReadFlexibleGpusRequest(*reqFlex).Execute()
	if err != nil {
		return fmt.Errorf("error reading the FlexibleGpu %s", err)
	}

	if err := utils.IsResponseEmptyOrMutiple(len(respGpu.GetFlexibleGpus()), "FlexibleGpu"); err != nil {
		return err
	}

	if (*respGpu.FlexibleGpus)[0].GetState() != "attaching" {
		return fmt.Errorf("Unable to link Flexible GPU")
	}

	if err := changeShutdownBehavior(conn, vmId); err != nil {
		return fmt.Errorf("Unable to change ShutdownBehavior: %s\n", err)
	}

	return resourceOutscaleOAPIFlexibleGpuLinkRead(d, meta)
}

func resourceOutscaleOAPIFlexibleGpuLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	flexGpuID := d.Get("flexible_gpu_id").(string)
	if flexGpuID == "" {
		flexGpuID = d.Id()
	}
	req := &oscgo.ReadFlexibleGpusRequest{
		Filters: &oscgo.FiltersFlexibleGpu{
			FlexibleGpuIds: &[]string{flexGpuID},
		},
	}

	var resp oscgo.ReadFlexibleGpusResponse
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		resp, _, err = conn.FlexibleGpuApi.ReadFlexibleGpus(
			context.Background()).
			ReadFlexibleGpusRequest(*req).Execute()
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetFlexibleGpus()), "FlexibleGpu"); err != nil {
		return err
	}

	fg := (*resp.FlexibleGpus)[0]
	if err := d.Set("flexible_gpu_id", fg.GetFlexibleGpuId()); err != nil {
		return err
	}
	if err := d.Set("vm_id", fg.GetVmId()); err != nil {
		return err
	}
	d.SetId(fg.GetFlexibleGpuId())

	return nil
}

func resourceOutscaleOAPIFlexibleGpuLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	flexGpuID := d.Get("flexible_gpu_id").(string)
	vmId := d.Get("vm_id").(string)

	filter := &oscgo.FiltersFlexibleGpu{
		FlexibleGpuIds: &[]string{flexGpuID},
	}
	reqFlex := &oscgo.ReadFlexibleGpusRequest{
		Filters: filter,
	}

	req := &oscgo.UnlinkFlexibleGpuRequest{
		FlexibleGpuId: flexGpuID,
	}

	var err error

	err = resource.Retry(20*time.Second, func() *resource.RetryError {
		_, _, err = conn.FlexibleGpuApi.UnlinkFlexibleGpu(
			context.Background()).UnlinkFlexibleGpuRequest(*req).Execute()
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	respGpu, _, err := conn.FlexibleGpuApi.ReadFlexibleGpus(context.Background()).
		ReadFlexibleGpusRequest(*reqFlex).Execute()
	if err != nil {
		return fmt.Errorf("error reading the FlexibleGpu %s", err)
	}

	if len(*respGpu.FlexibleGpus) != 1 {
		return fmt.Errorf("Unable to find Flexible GPU")
	}
	if (*respGpu.FlexibleGpus)[0].GetState() != "detaching" {
		return fmt.Errorf("Unable to unlink Flexible GPU")
	}

	if err := changeShutdownBehavior(conn, vmId); err != nil {
		return fmt.Errorf("Unable to change ShutdownBehavior: %s\n", err)
	}

	d.SetId("")
	return nil

}

func changeShutdownBehavior(conn *oscgo.APIClient, vmId string) error {

	var respV oscgo.ReadVmsResponse

	err := resource.Retry(20*time.Second, func() *resource.RetryError {
		r, _, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{vmId},
			}}).Execute()

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		respV = r
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading the VM %s", err)
	}
	if len(respV.GetVms()) == 0 {
		return fmt.Errorf("error reading the VM %s err %s ", vmId, err)
	}
	vm := respV.GetVms()[0]

	shutdownBehOpt := vm.GetVmInitiatedShutdownBehavior()
	if shutdownBehOpt != "stop" {
		sbOpts := oscgo.UpdateVmRequest{VmId: vm.GetVmId()}
		sbOpts.SetVmInitiatedShutdownBehavior("stop")
		if err := updateVmAttr(conn, sbOpts); err != nil {
			return err
		}
	}

	if err := stopVM(vmId, conn); err != nil {
		return err
	}

	if shutdownBehOpt != "stop" {
		sbReq := oscgo.UpdateVmRequest{VmId: vmId}
		sbReq.SetVmInitiatedShutdownBehavior(shutdownBehOpt)
		if err = updateVmAttr(conn, sbReq); err != nil {
			return err
		}
	}

	if err := startVM(vmId, conn); err != nil {
		return err
	}
	return nil
}
