package outscale

import (
	"context"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceFlexibleGpu() *schema.Resource {
	return &schema.Resource{
		Create: resourceFlexibleGpuCreate,
		Read:   resourceFlexibleGpuRead,
		Delete: resourceFlexibleGpuDelete,
		Update: resourceFlexibleGpuUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"delete_on_vm_deletion": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"model_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"generation": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"subregion_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"flexible_gpu_id": {
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

func resourceFlexibleGpuCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	req := &oscgo.CreateFlexibleGpuRequest{}

	mn := d.Get("model_name")
	req.SetModelName(mn.(string))

	sn := d.Get("subregion_name")
	req.SetSubregionName(sn.(string))

	if v, ok := d.GetOk("delete_on_vm_deletion"); ok {
		req.SetDeleteOnVmDeletion(v.(bool))
	}

	if v, ok := d.GetOk("generation"); ok {
		req.SetGeneration(v.(string))
	}

	var resp oscgo.CreateFlexibleGpuResponse
	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.FlexibleGpuApi.CreateFlexibleGpu(
			context.Background()).
			CreateFlexibleGpuRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resp.FlexibleGpu.GetFlexibleGpuId())

	return resourceFlexibleGpuRead(d, meta)
}

func resourceFlexibleGpuRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	napid := d.Id()

	filter := &oscgo.FiltersFlexibleGpu{
		FlexibleGpuIds: &[]string{napid},
	}

	req := &oscgo.ReadFlexibleGpusRequest{
		Filters: filter,
	}

	var resp oscgo.ReadFlexibleGpusResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.FlexibleGpuApi.ReadFlexibleGpus(
			context.Background()).
			ReadFlexibleGpusRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetFlexibleGpus()), "FlexibleGpu"); err != nil {
		return err
	}

	fg := (*resp.FlexibleGpus)[0]
	if err := d.Set("delete_on_vm_deletion", fg.GetDeleteOnVmDeletion()); err != nil {
		return err
	}
	if err := d.Set("subregion_name", fg.GetSubregionName()); err != nil {
		return err
	}
	if err := d.Set("generation", fg.GetGeneration()); err != nil {
		return err
	}
	if err := d.Set("flexible_gpu_id", fg.GetFlexibleGpuId()); err != nil {
		return err
	}
	if err := d.Set("vm_id", fg.GetVmId()); err != nil {
		return err
	}
	if err := d.Set("model_name", fg.GetModelName()); err != nil {
		return err
	}
	if err := d.Set("state", fg.GetState()); err != nil {
		return err
	}
	d.SetId(fg.GetFlexibleGpuId())

	return nil
}

func resourceFlexibleGpuUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	v := d.Get("delete_on_vm_deletion")
	req := &oscgo.UpdateFlexibleGpuRequest{
		FlexibleGpuId: d.Id(),
	}
	req.SetDeleteOnVmDeletion(v.(bool))

	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.FlexibleGpuApi.UpdateFlexibleGpu(
			context.Background()).
			UpdateFlexibleGpuRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return resourceFlexibleGpuRead(d, meta)

}

func resourceFlexibleGpuDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	req := &oscgo.DeleteFlexibleGpuRequest{
		FlexibleGpuId: d.Id(),
	}

	var err error
	err = resource.Retry(60*time.Second, func() *resource.RetryError {
		_, httpResp, err := conn.FlexibleGpuApi.DeleteFlexibleGpu(
			context.Background()).
			DeleteFlexibleGpuRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil

}
