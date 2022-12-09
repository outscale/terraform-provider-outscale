package outscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPIVM() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOutscaleOAPIVMRead,
		Schema: GetDataSourceSchema(VMSchema()),
	}
}
func dataSourceOutscaleOAPIVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	instanceID, instanceIDOk := d.GetOk("vm_id")

	if !filtersOk && !instanceIDOk {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}
	// Build up search parameters
	params := oscgo.ReadVmsRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceVMFilters(filters.(*schema.Set))
	}
	if instanceIDOk {
		params.Filters.VmIds = &[]string{instanceID.(string)}
	}

	log.Printf("[DEBUG] ReadVmsRequest -> %+v\n", params)

	var resp oscgo.ReadVmsResponse
	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		rp, httpResp, err := client.VmApi.ReadVms(context.Background()).ReadVmsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("error reading the VM %s", err)
	}

	if !resp.HasVms() {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var filteredVms []oscgo.Vm

	// loop through reservations, and remove terminated instances, populate vm slice
	for _, res := range resp.GetVms() {
		if res.GetState() != "terminated" {
			filteredVms = append(filteredVms, res)
		}
	}

	var vm oscgo.Vm
	if len(filteredVms) < 1 {
		return errors.New("Your query returned no results. Please change your search criteria and try again")
	}

	if len(filteredVms) > 1 {
		return errors.New("Your query returned more than one result. Please try a more " +
			"specific search criteria")
	}

	vm = filteredVms[0]

	// Populate vm attribute fields with the returned vm
	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(vm.GetVmId())
		return oapiVMDescriptionAttributes(set, &vm)
	})
}

// Populate instance attribute fields with the returned instance
func oapiVMDescriptionAttributes(set AttributeSetter, vm *oscgo.Vm) error {

	if err := set("architecture", vm.GetArchitecture()); err != nil {
		return err
	}

	if err := set("block_device_mappings_created", getOscAPIVMBlockDeviceMapping(vm.GetBlockDeviceMappings())); err != nil {
		log.Printf("[DEBUG] BLOCKING DEVICE MAPPING ERR %+v", err)
		return err
	}
	if err := set("bsu_optimized", vm.GetBsuOptimized()); err != nil {
		return err
	}
	if err := set("client_token", vm.GetClientToken()); err != nil {
		return err
	}
	if err := set("creation_date", vm.GetCreationDate()); err != nil {
		return err
	}
	if err := set("deletion_protection", vm.GetDeletionProtection()); err != nil {
		return err
	}
	if err := set("hypervisor", vm.GetHypervisor()); err != nil {
		return err
	}
	if err := set("image_id", vm.GetImageId()); err != nil {
		return err
	}
	if err := set("is_source_dest_checked", vm.GetIsSourceDestChecked()); err != nil {
		return err
	}
	if err := set("keypair_name", vm.GetKeypairName()); err != nil {
		return err
	}
	if err := set("launch_number", vm.GetLaunchNumber()); err != nil {
		return err
	}
	if err := set("net_id", vm.GetNetId()); err != nil {
		return err
	}
	if err := set("nested_virtualization", vm.GetNestedVirtualization()); err != nil {
		return err
	}
	if err := set("nics", getOAPIVMNetworkInterfaceLightSet(vm.GetNics())); err != nil {
		return err
	}

	if err := set("os_family", vm.GetOsFamily()); err != nil {
		return err
	}
	if err := set("performance", vm.GetPerformance()); err != nil {
		return err
	}
	if err := set("placement_subregion_name", aws.StringValue(vm.GetPlacement().SubregionName)); err != nil {
		return err
	}
	if err := set("placement_tenancy", aws.StringValue(vm.GetPlacement().Tenancy)); err != nil {
		return err
	}
	if err := set("private_dns_name", vm.GetPrivateDnsName()); err != nil {
		return err
	}
	if err := set("private_ip", vm.GetPrivateIp()); err != nil {
		return err
	}
	if err := set("product_codes", vm.GetProductCodes()); err != nil {
		return err
	}
	if err := set("public_dns_name", vm.GetPublicDnsName()); err != nil {
		return err
	}
	if err := set("public_ip", vm.GetPublicIp()); err != nil {
		return err
	}
	if err := set("reservation_id", vm.GetReservationId()); err != nil {
		return err
	}
	if err := set("root_device_name", vm.GetRootDeviceName()); err != nil {
		return err
	}
	if err := set("root_device_type", vm.GetRootDeviceType()); err != nil {
		return err
	}
	if err := set("security_groups", getOAPIVMSecurityGroups(vm.GetSecurityGroups())); err != nil {
		log.Printf("[DEBUG] SECURITY GROUPS ERR %+v", err)
		return err
	}
	if err := set("state", vm.GetState()); err != nil {
		return err
	}
	if err := set("state_reason", vm.GetStateReason()); err != nil {
		return err
	}
	if err := set("subnet_id", vm.GetSubnetId()); err != nil {
		return err
	}
	if err := set("user_data", vm.GetUserData()); err != nil {
		return err
	}
	if err := set("vm_id", vm.GetVmId()); err != nil {
		return err
	}
	if err := set("vm_initiated_shutdown_behavior", vm.GetVmInitiatedShutdownBehavior()); err != nil {
		return err
	}
	if err := set("tags", getOscAPITagSet(vm.GetTags())); err != nil {
		return err
	}

	return set("vm_type", vm.GetVmType())
}

func getOscAPIVMBlockDeviceMapping(blockDeviceMappings []oscgo.BlockDeviceMappingCreated) []map[string]interface{} {
	blockDeviceMapping := make([]map[string]interface{}, len(blockDeviceMappings))

	for k, v := range blockDeviceMappings {
		blockDeviceMapping[k] = map[string]interface{}{
			"device_name": aws.StringValue(v.DeviceName),
			"bsu": map[string]interface{}{
				"delete_on_vm_deletion": fmt.Sprintf("%t", aws.BoolValue(v.GetBsu().DeleteOnVmDeletion)),
				"volume_id":             aws.StringValue(v.GetBsu().VolumeId),
				"state":                 aws.StringValue(v.GetBsu().State),
				"link_date":             aws.StringValue(v.GetBsu().LinkDate),
			},
		}
	}
	return blockDeviceMapping
}

func getOAPIVMSecurityGroups(groupSet []oscgo.SecurityGroupLight) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, g := range groupSet {
		r := map[string]interface{}{
			"security_group_id":   g.GetSecurityGroupId(),
			"security_group_name": g.GetSecurityGroupName(),
		}
		res = append(res, r)
	}

	return res
}

func buildOutscaleOAPIDataSourceVMFilters(set *schema.Set) *oscgo.FiltersVm {
	filters := new(oscgo.FiltersVm)

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "vm_ids":
			filters.VmIds = &filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
