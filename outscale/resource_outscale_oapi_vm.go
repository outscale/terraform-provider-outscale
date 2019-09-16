package outscale

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func resourceOutscaleOApiVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVMCreate,
		Read:   resourceOAPIVMRead,
		Update: resourceOAPIVMAttributesUpdate,
		Delete: resourceOAPIVMDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getOApiVMSchema(),
	}
}

func resourceOAPIVMCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	instanceOpts, err := buildCreateVmsRequest(d, meta)
	if err != nil {
		return err
	}

	// Create the instance
	var runResp *oapi.CreateVmsResponse
	var resp *oapi.POST_CreateVmsResponses
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		var err error
		resp, err = conn.POST_CreateVms(*instanceOpts)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error launching source instance: %s", err)
	}

	runResp = resp.OK

	if runResp == nil || len(runResp.Vms) == 0 {
		return errors.New("Error launching source instance: no instances returned in response")
	}

	vm := runResp.Vms[0]
	fmt.Printf("[INFO] Instance ID: %s", vm.VmId)

	d.SetId(vm.VmId)

	if d.IsNewResource() {
		if err := setOAPITags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tag")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"running"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, vm.VmId, "terminated"),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to stop: %s", d.Id(), err)
	}

	// Initialize the connection info
	if vm.PublicIp != "" {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": vm.PublicIp,
		})
	} else if vm.PrivateIp != "" {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": vm.PrivateIp,
		})
	}

	//Check if source dest check is enabled.
	if v, ok := d.GetOk("is_source_dest_checked"); ok {
		opts := &oapi.UpdateVmRequest{
			VmId:                vm.VmId,
			IsSourceDestChecked: v.(bool),
		}
		log.Printf("[DEBGUG] is_source_dest_checked argument is not in CreateVms, we have to update the vm (%s)", vm.VmId)
		if err := oapiModifyInstanceAttr(conn, opts); err != nil {
			return err
		}
	}

	return resourceOAPIVMRead(d, meta)
}

func resourceOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	filters := oapi.FiltersVm{
		VmIds: []string{d.Id()},
	}

	input := &oapi.ReadVmsRequest{
		Filters: filters,
	}

	var resp *oapi.ReadVmsResponse
	var rs *oapi.POST_ReadVmsResponses
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		rs, err = conn.POST_ReadVms(*input)

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error reading the VM %s", err)
	}

	resp = rs.OK

	if err != nil {
		// If the instance was not found, return nil so that we can show
		// that the instance is gone.
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidInstanceID.NotFound" {
			d.SetId("")
			return nil
		}

		// Some other error, report it
		return err
	}

	// If nothing was found, then return no state
	if len(resp.Vms) == 0 {
		d.SetId("")
		return nil
	}

	instance := resp.Vms[0]

	d.Set("request_id", resp.ResponseContext.RequestId)
	return resourceDataAttrSetter(d, &instance)
}

func resourceOAPIVMDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	fmt.Printf("[INFO] Terminating instance: %s", id)
	req := &oapi.DeleteVmsRequest{
		VmIds: []string{id},
	}

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err = conn.POST_DeleteVms(*req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				fmt.Printf("[INFO] Request limit exceeded")
				return resource.RetryableError(err)
			}
		}

		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error deleting the instance")
	}

	fmt.Printf("[DEBUG] Waiting for instance (%s) to become terminated", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"terminated"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, id, ""),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to terminate: %s", id, err)
	}

	return nil
}

func getOApiVMSchema() map[string]*schema.Schema {
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

func buildCreateVmsRequest(
	d *schema.ResourceData, meta interface{}) (*oapi.CreateVmsRequest, error) {
	request := &oapi.CreateVmsRequest{
		DeletionProtection:          d.Get("deletion_protection").(bool),
		BsuOptimized:                d.Get("bsu_optimized").(bool),
		ImageId:                     d.Get("image_id").(string),
		VmType:                      d.Get("vm_type").(string),
		VmInitiatedShutdownBehavior: d.Get("vm_initiated_shutdown_behavior").(string),
		UserData:                    d.Get("user_data").(string),
		MaxVmsCount:                 int64(1),
		MinVmsCount:                 int64(1),
	}

	request.Placement = oapi.Placement{}

	if v, ok := d.GetOk("placement_subregion_name"); ok {
		request.Placement.SubregionName = v.(string)
	}

	if v, ok := d.GetOk("placement_tenancy"); ok {
		request.Placement.Tenancy = v.(string)
	}

	sgNames := make([]string, 0)
	if v := d.Get("security_group_names"); v != nil {
		sgNames = expandStringValueList(v.([]interface{}))
	}

	sgIds := make([]string, 0)
	if v := d.Get("security_group_ids"); v != nil {
		sgIds = expandStringValueList(v.([]interface{}))
	}

	privateIPS := make([]string, 0)
	if v := d.Get("private_ips"); v != nil {
		privateIPS = expandStringValueList(v.([]interface{}))
	}

	subnetID, hasSubnet := d.GetOk("subnet_id")
	networkInterfaces, interfacesOk := d.GetOk("nics")

	if hasSubnet && interfacesOk {
		return nil, errors.New("Error you need to specify only one: subnet_id or nics")
	}

	if interfacesOk {
		request.Nics = buildNetworkOApiInterfaceOpts(d, sgNames, networkInterfaces)
	}
	if hasSubnet {
		request.SubnetId = subnetID.(string)
		request.SecurityGroupIds = sgIds
		request.SecurityGroups = sgNames
		request.PrivateIps = privateIPS
	}

	if v, ok := d.GetOk("private_ip"); ok {
		request.PrivateIps = []string{v.(string)}
	}

	if v, ok := d.GetOk("keypair_name"); ok {
		request.KeypairName = v.(string)
	}

	var err error
	if v, ok := d.GetOk("block_device_mappings"); ok {
		request.BlockDeviceMappings = expandBlockDeviceOApiMappings(v.([]interface{}))
	}
	if err != nil {
		return nil, err
	}

	return request, nil
}

func expandPrivatePublicIps(privateIPS *schema.Set) []oapi.PrivateIpLight {
	privatePublicIPS := make([]oapi.PrivateIpLight, len(privateIPS.List()))

	for i, v := range privateIPS.List() {
		value := v.(map[string]interface{})
		privatePublicIPS[i].IsPrimary = value["is_primary"].(bool)
		privatePublicIPS[i].PrivateIp = value["private_ip"].(string)
	}
	return privatePublicIPS
}

func buildNetworkOApiInterfaceOpts(d *schema.ResourceData, groups []string, nInterfaces interface{}) []oapi.NicForVmCreation {
	networkInterfaces := []oapi.NicForVmCreation{}
	vL := nInterfaces.([]interface{})

	for _, v := range vL {
		ini := v.(map[string]interface{})

		ni := oapi.NicForVmCreation{
			DeleteOnVmDeletion: ini["delete_on_vm_deletion"].(bool),
			Description:        ini["description"].(string),
			DeviceNumber:       int64(ini["device_number"].(int)),
		}

		ni.PrivateIps = expandPrivatePublicIps(ini["private_ips"].(*schema.Set))
		ni.SubnetId = ini["subnet_id"].(string)
		ni.SecurityGroupIds = expandStringValueList(ini["security_group_ids"].([]interface{}))
		ni.SecondaryPrivateIpCount = int64(ini["secondary_private_ip_count"].(int))
		ni.NicId = ini["nic_id"].(string)

		if v, ok := d.GetOk("private_ip"); ok {
			ni.PrivateIps = []oapi.PrivateIpLight{oapi.PrivateIpLight{
				PrivateIp: v.(string),
			}}
		}

		networkInterfaces = append(networkInterfaces, ni)
	}

	return networkInterfaces
}

func expandBlockDeviceOApiMappings(block []interface{}) []oapi.BlockDeviceMappingVmCreation {
	blockDevices := make([]oapi.BlockDeviceMappingVmCreation, len(block))

	for i, v := range block {
		value := v.(map[string]interface{})
		bsu := value["bsu"].(map[string]interface{})

		if deleteOnVM, ok := bsu["delete_on_vm_deletion"]; ok {
			blockDevices[i].Bsu.DeleteOnVmDeletion = deleteOnVM.(bool)
		}
		if iops, ok := bsu["iops"]; ok {
			blockDevices[i].Bsu.Iops = int64(iops.(int))
		}
		if snapshotID, ok := bsu["snapshot_id"]; ok {
			blockDevices[i].Bsu.SnapshotId = snapshotID.(string)
		}
		if v, ok := bsu["volume_size"]; ok {
			log.Printf("LOG__VALUE \n\n %+v \n\n", v)
			n, _ := strconv.Atoi(v.(string))
			blockDevices[i].Bsu.VolumeSize = int64(n)
		}
		if volumeType, ok := bsu["volume_type"]; ok {
			blockDevices[i].Bsu.VolumeType = volumeType.(string)
		}
		if deviceName, ok := value["device_name"]; ok {
			blockDevices[i].DeviceName = deviceName.(string)
		}
		if noDevice, ok := value["no_device"]; ok {
			blockDevices[i].NoDevice = noDevice.(string)
		}
		if virtualDeviceName, ok := value["virtual_device_name"]; ok {
			blockDevices[i].VirtualDeviceName = virtualDeviceName.(string)
		}
	}
	return blockDevices
}

// InstanceStateOApiRefreshFunc ...
func InstanceStateOApiRefreshFunc(conn *oapi.Client, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *oapi.ReadVmsResponse
		var rs *oapi.POST_ReadVmsResponses
		var err error

		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			rs, err = conn.POST_ReadVms(oapi.ReadVmsRequest{
				Filters: getVMsFilterByVMID(instanceID),
			})
			return resource.RetryableError(err)
		})

		if err != nil {
			fmt.Printf("Error on InstanceStateRefresh: %s", err)

			return nil, "", err
		}

		resp = rs.OK

		if resp == nil || len(resp.Vms) == 0 {
			return nil, "", nil
		}

		i := resp.Vms[0]
		state := i.State

		if state == failState {
			return i, state, fmt.Errorf("Failed to reach target state. Reason: %v",
				i.State)

		}

		return i, state, nil
	}
}

func stopVM(vmID string, conn *oapi.Client, attr string) (*resource.StateChangeConf, error) {
	_, err := conn.POST_StopVms(oapi.StopVmsRequest{
		VmIds: []string{vmID},
	})

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending", "running", "shutting-down", "stopped", "stopping"},
		Target:     []string{"stopped"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, vmID, ""),
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

func startVM(vmID string, stateConf *resource.StateChangeConf, conn *oapi.Client, attr string) error {
	if _, err := conn.POST_StartVms(oapi.StartVmsRequest{
		VmIds: []string{vmID},
	}); err != nil {
		return err
	}

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"pending", "stopped"},
		Target:     []string{"running"},
		Refresh:    InstanceStateOApiRefreshFunc(conn, vmID, ""),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", vmID, err)
	}

	return nil
}

func getVMsFilterByVMID(vmID string) oapi.FiltersVm {
	return oapi.FiltersVm{
		VmIds: []string{vmID},
	}
}
