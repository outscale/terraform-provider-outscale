package outscale

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func resourceOutscaleOAPIVMAttributes() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVMAttributesCreate,
		Read:   dataSourceOutscaleOAPIVMRead,
		Update: resourceOAPIVMAttributesUpdate,
		Delete: resourceOAPIVMAttributesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getDataSourceOAPIVMAttrsSchemas(),
	}
}

func getDataSourceOAPIVMAttrsSchemas() map[string]*schema.Schema {
	wholeSchema := map[string]*schema.Schema{}

	attrsSchema := getOApiVMAttributesSchema()

	for k, v := range attrsSchema {
		wholeSchema[k] = v
	}

	wholeSchema["request_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return wholeSchema
}

func resourceOAPIVMAttributesCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Get("vm_id").(string)

	if v, ok := d.GetOk("deletion_protection"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:               id,
			DeletionProtection: v.(bool),
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := oapiModifyInstanceAttr(conn, opts, "deletion_protection"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("keypair_name"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:        id,
			KeypairName: v.(string),
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := oapiModifyInstanceAttr(conn, opts, "keypair_name"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:             id,
			SecurityGroupIds: v.([]string),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "security_group_ids"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("vm_initiated_shutdown_behavior"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:                        id,
			VmInitiatedShutdownBehavior: v.(string),
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := oapiModifyInstanceAttr(conn, opts, "vm_initiated_shutdown_behavior"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("is_source_dest_checked"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			IsSourceDestChecked: v.(bool),
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := oapiModifyInstanceAttr(conn, opts, "is_source_dest_checked"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("vm_type"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:   id,
			VmType: v.(string),
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := oapiModifyInstanceAttr(conn, opts, "vm_type"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:     id,
			UserData: v.(string),
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := oapiModifyInstanceAttr(conn, opts, "user_data"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("bsu_optimized"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:         id,
			BsuOptimized: v.(bool),
		}

		fmt.Printf("\n\n[DEBUG] CHANGES %+v, \n\n", opts)

		if err := oapiModifyInstanceAttr(conn, opts, "bsu_optimized"); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("block_device_mappings"); ok {
		maps := v.(*schema.Set).List()
		mappings := []oapi.BlockDeviceMappingVmUpdate{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := oapi.BlockDeviceMappingVmUpdate{
				DeviceName:        f["device_name"].(string),
				NoDevice:          f["no_device"].(string),
				VirtualDeviceName: f["virtual_device_name"].(string),
			}

			e := f["bsu"].(map[string]interface{})

			bsu := oapi.BsuToUpdateVm{
				DeleteOnVmDeletion: e["delete_on_vm_deletion"].(bool),
				VolumeId:           e["volume_id"].(string),
			}

			mapping.Bsu = bsu

			mappings = append(mappings, mapping)
		}

		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			BlockDeviceMappings: mappings,
		}
		if err := oapiModifyInstanceAttr(conn, opts, "block_device_mappings"); err != nil {
			return err
		}
	}

	d.SetId(resource.UniqueId())

	return dataSourceOutscaleOAPIVMRead(d, meta)
}

func resourceOAPIVMAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	d.Partial(true)

	id := d.Get("vm_id").(string)

	log.Printf("[DEBUG] updating the instance %s", id)

	if d.HasChange("deletion_protection") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:               id,
			DeletionProtection: d.Get("deletion_protection").(bool),
		}

		if err := oapiModifyInstanceAttr(conn, opts, "deletion_protection"); err != nil {
			return err
		}
	}

	if d.HasChange("keypair_name") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:        id,
			KeypairName: d.Get("keypair_name").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "keypair_name"); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_ids") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:             id,
			SecurityGroupIds: d.Get("security_group_ids").([]string),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "security_group_ids"); err != nil {
			return err
		}
	}

	if d.HasChange("vm_initiated_shutdown_behavior") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:                        id,
			VmInitiatedShutdownBehavior: d.Get("vm_initiated_shutdown_behavior").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "vm_initiated_shutdown_behavior"); err != nil {
			return err
		}
	}

	if d.HasChange("is_source_dest_checked") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			IsSourceDestChecked: d.Get("is_source_dest_checked").(bool),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "is_source_dest_checked"); err != nil {
			return err
		}
	}

	if d.HasChange("vm_type") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:   id,
			VmType: d.Get("vm_type").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "vm_type"); err != nil {
			return err
		}
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:     id,
			UserData: d.Get("user_data").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "user_data"); err != nil {
			return err
		}
	}

	if d.HasChange("bsu_optimized") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:         id,
			BsuOptimized: d.Get("bsu_optimized").(bool),
		}
		if err := oapiModifyInstanceAttr(conn, opts, "bsu_optimized"); err != nil {
			return err
		}
	}

	if d.HasChange("block_device_mappings") && !d.IsNewResource() {
		maps := d.Get("block_device_mappings").(*schema.Set).List()
		mappings := []oapi.BlockDeviceMappingVmUpdate{}

		for _, m := range maps {
			f := m.(map[string]interface{})
			mapping := oapi.BlockDeviceMappingVmUpdate{
				DeviceName:        f["device_name"].(string),
				NoDevice:          f["no_device"].(string),
				VirtualDeviceName: f["virtual_device_name"].(string),
			}

			e := f["bsu"].(map[string]interface{})

			bsu := oapi.BsuToUpdateVm{
				DeleteOnVmDeletion: e["delete_on_vm_deletion"].(bool),
				VolumeId:           e["volume_id"].(string),
			}

			mapping.Bsu = bsu

			mappings = append(mappings, mapping)
		}

		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			BlockDeviceMappings: mappings,
		}

		if err := oapiModifyInstanceAttr(conn, opts, "block_device_mappings"); err != nil {
			return err
		}
	}

	d.Partial(false)

	return dataSourceOutscaleOAPIVMRead(d, meta)
}

func resourceOAPIVMAttributesDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}

func oapiStopInstance(instanceAttrOpts *oapi.UpdateVmRequest, conn *oapi.Client, attr string) (*resource.StateChangeConf, error) {
	_, err := conn.POST_StopVms(oapi.StopVmsRequest{
		VmIds: []string{instanceAttrOpts.VmId},
	})

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"stopped"},
		Refresh:    oapiInstanceStateRefreshFunc(conn, instanceAttrOpts.VmId, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf(
			"Error waiting for instance (%s) to stop: %s", instanceAttrOpts.VmId, err)
	}

	return stateConf, nil
}

func oapiStartInstance(instanceAttrOpts *oapi.UpdateVmRequest, stateConf *resource.StateChangeConf, conn *oapi.Client, attr string) error {
	if _, err := conn.POST_StartVms(oapi.StartVmsRequest{
		VmIds: []string{instanceAttrOpts.VmId},
	}); err != nil {
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"pending", "stopped"},
		Target:     []string{"running"},
		Refresh:    oapiInstanceStateRefreshFunc(conn, instanceAttrOpts.VmId, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", instanceAttrOpts.VmId, err)
	}

	return nil
}

func oapiInstanceStateRefreshFunc(conn *oapi.Client, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *oapi.POST_ReadVmsResponses
		var err error

		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			resp, err = conn.POST_ReadVms(oapi.ReadVmsRequest{
				Filters: oapi.FiltersVm{VmIds: []string{instanceID}},
			})
			return resource.RetryableError(err)
		})

		if err != nil {
			return nil, "", err
		}

		if resp == nil || len(resp.OK.Vms) == 0 {
			return nil, "", nil
		}

		i := resp.OK.Vms[0]
		state := i.State

		if state == failState {
			return i, state, fmt.Errorf("Failed to reach target state. Reason: %v",
				i.StateReason)

		}

		return i, state, nil
	}
}

func needsVmRestart(attr string) bool {
	restart := false
	switch attr {
	case "vm_type":
		fallthrough
	case "user_data":
		fallthrough
	case "ebs_optimized":
		fallthrough
	case "deletion_protection":
		restart = true
	}
	return restart
}

func oapiModifyInstanceAttr(conn *oapi.Client, instanceAttrOpts *oapi.UpdateVmRequest, attr string) error {

	var err error
	var stateConf *resource.StateChangeConf

	if needsVmRestart(attr) {
		stateConf, err = oapiStopInstance(instanceAttrOpts, conn, attr)
	}

	if err != nil {
		return err
	}

	if _, err := conn.POST_UpdateVm(*instanceAttrOpts); err != nil {
		return err
	}

	if needsVmRestart(attr) {
		err = oapiStartInstance(instanceAttrOpts, stateConf, conn, attr)
	}

	if err != nil {
		return err
	}

	return nil
}
