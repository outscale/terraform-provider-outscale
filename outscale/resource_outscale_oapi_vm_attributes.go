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
	log.Printf("UPDATING ATTRS FOR: %s", id)

	vVMType, okVMType := d.GetOk("vm_type")
	vUserData, okUserData := d.GetOk("user_data")
	vBsuOptimized, okBsuOptimized := d.GetOk("bsu_optimized")

	var stateConf *resource.StateChangeConf
	var err error
	if okVMType || okUserData || okBsuOptimized {
		stateConf, err = oapiStopInstance(id, conn)
		if err != nil {
			return err
		}
	}

	if okVMType {
		opts := &oapi.UpdateVmRequest{
			VmId:   id,
			VmType: vVMType.(string),
		}
		log.Printf("UPDATE (vm_type) %+v => %+v == %+v", okVMType, d.Get("vm_type"), vVMType)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if okUserData {
		opts := &oapi.UpdateVmRequest{
			VmId:     id,
			UserData: vUserData.(string),
		}
		log.Printf("UPDATE (vm_type) %+v => %+v == %+v", okUserData, d.Get("vm_type"), vUserData)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if okBsuOptimized {
		opts := &oapi.UpdateVmRequest{
			VmId:         id,
			BsuOptimized: vBsuOptimized.(bool),
		}
		log.Printf("UPDATE (bsu_optimized) %+v => %+v == %+v", okBsuOptimized, d.Get("bsu_optimized"), vBsuOptimized)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("deletion_protection"); ok {
		deletionProtection := v.(bool)
		opts := &oapi.UpdateVmRequest{
			VmId:               id,
			DeletionProtection: &deletionProtection,
		}
		log.Printf("UPDATE (deletion_protection) %+v => %+v == %+v", ok, d.Get("deletion_protection"), v)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("keypair_name"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:        id,
			KeypairName: v.(string),
		}
		log.Printf("UPDATE (keypair_name) %+v => %+v == %+v", ok, d.Get("keypair_name"), v)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:             id,
			SecurityGroupIds: v.([]string),
		}
		log.Printf("UPDATE (security_group_ids) %+v => %+v == %+v", ok, d.Get("security_group_ids"), v)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("vm_initiated_shutdown_behavior"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:                        id,
			VmInitiatedShutdownBehavior: v.(string),
		}
		log.Printf("UPDATE (vm_initiated_shutdown_behavior) %+v => %+v == %+v", ok, d.Get("vm_initiated_shutdown_behavior"), v)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("is_source_dest_checked"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			IsSourceDestChecked: v.(bool),
		}
		log.Printf("UPDATE (is_source_dest_checked) %+v => %+v == %+v", ok, d.Get("is_source_dest_checked"), v)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
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
		log.Printf("UPDATE (block_device_mappings) %+v => %+v == %+v", ok, d.Get("block_device_mappings"), mappings)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	d.SetId(resource.UniqueId())

	if err := oapiStartInstance(id, stateConf, conn); err != nil {
		return err
	}

	return dataSourceOutscaleOAPIVMRead(d, meta)
}

func resourceOAPIVMAttributesUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	d.Partial(true)

	id := d.Get("vm_id").(string)

	var stateConf *resource.StateChangeConf
	var err error
	if d.HasChange("vm_type") && !d.IsNewResource() ||
		d.HasChange("user_data") && !d.IsNewResource() ||
		d.HasChange("bsu_optimized") && !d.IsNewResource() {
		stateConf, err = oapiStopInstance(id, conn)
		if err != nil {
			return err
		}
	}

	if d.HasChange("vm_type") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:   id,
			VmType: d.Get("vm_type").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("user_data") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:     id,
			UserData: d.Get("user_data").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("bsu_optimized") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:         id,
			BsuOptimized: d.Get("bsu_optimized").(bool),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") && !d.IsNewResource() {
		deletionProtection := d.Get("deletion_protection").(bool)
		opts := &oapi.UpdateVmRequest{
			VmId:               id,
			DeletionProtection: &deletionProtection,
		}

		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("keypair_name") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:        id,
			KeypairName: d.Get("keypair_name").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_ids") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:             id,
			SecurityGroupIds: expandStringValueList(d.Get("security_group_ids").([]interface{})),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("security_group_names") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:             id,
			SecurityGroupIds: expandStringValueList(d.Get("security_group_names").([]interface{})),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("vm_initiated_shutdown_behavior") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:                        id,
			VmInitiatedShutdownBehavior: d.Get("vm_initiated_shutdown_behavior").(string),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	if d.HasChange("is_source_dest_checked") && !d.IsNewResource() {
		opts := &oapi.UpdateVmRequest{
			VmId:                id,
			IsSourceDestChecked: d.Get("is_source_dest_checked").(bool),
		}
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
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

		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	d.Partial(false)

	if err := oapiStartInstance(id, stateConf, conn); err != nil {
		return err
	}

	return dataSourceOutscaleOAPIVMRead(d, meta)
}

func resourceOAPIVMAttributesDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	return nil
}

func oapiStopInstance(vmID string, conn *oapi.Client) (*resource.StateChangeConf, error) {
	log.Printf("STOPPING VM... %+v", vmID)
	_, err := conn.POST_StopVms(oapi.StopVmsRequest{
		VmIds: []string{vmID},
	})

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"stopped"},
		Refresh:    oapiInstanceStateRefreshFunc(conn, vmID, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf(
			"Error waiting for instance (%s) to stop: %s", vmID, err)
	}

	return stateConf, nil
}

func oapiStartInstance(vmID string, stateConf *resource.StateChangeConf, conn *oapi.Client) error {
	log.Printf("STARTING VM... %+v", vmID)
	if _, err := conn.POST_StartVms(oapi.StartVmsRequest{VmIds: []string{vmID}}); err != nil {
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"pending", "stopped"},
		Target:     []string{"running"},
		Refresh:    oapiInstanceStateRefreshFunc(conn, vmID, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", vmID, err)
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

func oapiModifyInstanceAttr(conn *oapi.Client, instanceAttrOpts *oapi.UpdateVmRequest) error {
	log.Printf("UPDATE VM Payload %+v", instanceAttrOpts)
	if _, err := conn.POST_UpdateVm(*instanceAttrOpts); err != nil {
		return err
	}
	return nil
}
