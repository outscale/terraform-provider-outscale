package outscale

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func flattenedInstanceSet(instances []*fcu.Instance) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		flattened[i] = map[string]interface{}{
			"ami_launch_index":   *instance.AmiLaunchIndex,
			"ebs_optimized":      *instance.EbsOptimized,
			"architecture":       *instance.Architecture,
			"client_token":       *instance.ClientToken,
			"hypervisor":         *instance.Hypervisor,
			"image_id":           *instance.ImageId,
			"instance_id":        *instance.InstanceId,
			"instance_type":      *instance.InstanceType,
			"kernel_id":          *instance.KernelId,
			"key_name":           *instance.KeyName,
			"private_ip_address": *instance.PrivateDnsName,
			"private_dns_name":   *instance.PrivateDnsName,
			"root_device_name":   *instance.RootDeviceName,
		}

		if instance.InstanceLifecycle != nil {
			flattened[i]["instance_lifecycle"] = *instance.InstanceLifecycle
		}
		if instance.RootDeviceType != nil {
			flattened[i]["root_device_type"] = *instance.RootDeviceType
		}

		if instance.DnsName != nil {
			flattened[i]["dns_name"] = *instance.DnsName
		}

		if instance.IpAddress != nil {
			flattened[i]["ip_address"] = *instance.IpAddress
		}
		if instance.Platform != nil {
			flattened[i]["platform"] = *instance.Platform
		}
		if instance.RamdiskId != nil {
			flattened[i]["ramdisk_id"] = *instance.RamdiskId
		}
		if instance.Reason != nil {
			flattened[i]["reason"] = *instance.Reason
		}
		if instance.SourceDestCheck != nil {
			flattened[i]["source_dest_check"] = *instance.SourceDestCheck
		}
		if instance.SpotInstanceRequestId != nil {
			flattened[i]["spot_instance_request_id"] = *instance.SpotInstanceRequestId
		}
		if instance.SriovNetSupport != nil {
			flattened[i]["sriov_net_support"] = *instance.SriovNetSupport
		}
		if instance.SubnetId != nil {
			flattened[i]["subnet_id"] = *instance.SubnetId
		}
		if instance.VirtualizationType != nil {
			flattened[i]["virtualization_type"] = *instance.VirtualizationType
		}
		if instance.VpcId != nil {
			flattened[i]["vpc_id"] = *instance.VpcId
		}

		flattened[i]["block_device_mapping"] = flattenedBlockDeviceMapping(instance.BlockDeviceMappings)

		flattened[i]["group_set"] = getGroupSet(instance.GroupSet)
		flattened[i]["iam_instance_profile"] = getIAMInstanceProfile(instance.IamInstanceProfile)
		flattened[i]["instance_state"] = getInstanceState(instance.State)
		flattened[i]["monitoring"] = getMonitoring(instance.Monitoring)
		flattened[i]["network_interface_set"] = getNetworkInterfaceSet(instance.NetworkInterfaces)
		flattened[i]["placement"] = getPlacement(instance.Placement)
		flattened[i]["state_reason"] = getStateReason(instance.StateReason)
		flattened[i]["product_codes"] = getProductCodes(instance.ProductCodes)
		flattened[i]["tag_set"] = getTagSet(instance.Tags)

	}

	return flattened
}

func flattenedInstanceSetPassword(instances []*fcu.Instance, conn fcu.VMService) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		flattened[i] = map[string]interface{}{
			"ami_launch_index":   *instance.AmiLaunchIndex,
			"ebs_optimized":      *instance.EbsOptimized,
			"architecture":       *instance.Architecture,
			"client_token":       *instance.ClientToken,
			"hypervisor":         *instance.Hypervisor,
			"image_id":           *instance.ImageId,
			"instance_id":        *instance.InstanceId,
			"instance_type":      *instance.InstanceType,
			"kernel_id":          *instance.KernelId,
			"key_name":           *instance.KeyName,
			"private_ip_address": *instance.PrivateDnsName,
			"private_dns_name":   *instance.PrivateDnsName,
			"root_device_name":   *instance.RootDeviceName,
		}

		if instance.InstanceLifecycle != nil {
			flattened[i]["instance_lifecycle"] = *instance.InstanceLifecycle
		}
		if instance.RootDeviceType != nil {
			flattened[i]["root_device_type"] = *instance.RootDeviceType
		}

		if instance.DnsName != nil {
			flattened[i]["dns_name"] = *instance.DnsName
		}

		if instance.IpAddress != nil {
			flattened[i]["ip_address"] = *instance.IpAddress
		}
		if instance.Platform != nil {
			flattened[i]["platform"] = *instance.Platform
			if *instance.Platform == "windows" {
				pass, _ := conn.GetPasswordData(&fcu.GetPasswordDataInput{
					InstanceId: instance.InstanceId,
				})
				fmt.Println(*pass.PasswordData)
				flattened[i]["password_data"] = *pass.PasswordData
			}
		}
		if instance.RamdiskId != nil {
			flattened[i]["ramdisk_id"] = *instance.RamdiskId
		}
		if instance.Reason != nil {
			flattened[i]["reason"] = *instance.Reason
		}
		if instance.SourceDestCheck != nil {
			flattened[i]["source_dest_check"] = *instance.SourceDestCheck
		}
		if instance.SpotInstanceRequestId != nil {
			flattened[i]["spot_instance_request_id"] = *instance.SpotInstanceRequestId
		}
		if instance.SriovNetSupport != nil {
			flattened[i]["sriov_net_support"] = *instance.SriovNetSupport
		}
		if instance.SubnetId != nil {
			flattened[i]["subnet_id"] = *instance.SubnetId
		}
		if instance.VirtualizationType != nil {
			flattened[i]["virtualization_type"] = *instance.VirtualizationType
		}
		if instance.VpcId != nil {
			flattened[i]["vpc_id"] = *instance.VpcId
		}

		flattened[i]["block_device_mapping"] = flattenedBlockDeviceMapping(instance.BlockDeviceMappings)

		flattened[i]["group_set"] = getGroupSet(instance.GroupSet)
		flattened[i]["iam_instance_profile"] = getIAMInstanceProfile(instance.IamInstanceProfile)
		flattened[i]["instance_state"] = getInstanceState(instance.State)
		flattened[i]["monitoring"] = getMonitoring(instance.Monitoring)
		flattened[i]["network_interface_set"] = getNetworkInterfaceSet(instance.NetworkInterfaces)
		flattened[i]["placement"] = getPlacement(instance.Placement)
		flattened[i]["state_reason"] = getStateReason(instance.StateReason)
		flattened[i]["product_codes"] = getProductCodes(instance.ProductCodes)
		flattened[i]["tag_set"] = getTagSet(instance.Tags)

	}

	return flattened
}

func flattenedBlockDeviceMapping(mappings []*fcu.InstanceBlockDeviceMapping) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(mappings))
	for i, mapping := range mappings {
		flattened[i] = map[string]interface{}{
			"device_name": *mapping.DeviceName,
			"ebs":         flattenEBS(mapping.Ebs),
		}
	}

	return flattened
}

func resourceInstancSetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["ami_launch_index"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["ami_launch_index"].(int)))
	}

	if m["architecture"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["architecture"].(string)))
	}
	if m["ip_address"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["architecture"].(string)))
	}

	return hashcode.String(buf.String())

}

func getPrivateIPAddressSet(privateIPs []*fcu.InstancePrivateIpAddress) []map[string]interface{} {
	res := []map[string]interface{}{}
	if privateIPs != nil {
		for _, p := range privateIPs {
			inter := make(map[string]interface{})
			assoc := make(map[string]interface{})

			if p.Association != nil {
				assoc["ip_owner_id"] = *p.Association.IpOwnerId
				assoc["public_dns_name"] = *p.Association.PublicDnsName
				assoc["public_ip"] = *p.Association.PublicIp
			}

			inter["association"] = assoc
			inter["private_dns_name"] = *p.Primary
			inter["private_ip_address"] = *p.PrivateIpAddress

		}
	}
	return res
}

func getPlacement(placement *fcu.Placement) map[string]interface{} {
	res := map[string]interface{}{}

	if placement != nil {
		if placement.Affinity != nil {
			res["affinity"] = *placement.Affinity
		}
		res["availability_zone"] = *placement.AvailabilityZone
		res["group_name"] = *placement.GroupName
		if placement.HostId != nil {
			res["host_id"] = *placement.HostId
		}
		res["tenancy"] = *placement.Tenancy
	}

	return res
}

func getProductCodes(codes []*fcu.ProductCode) []map[string]interface{} {
	res := []map[string]interface{}{}

	if codes != nil {
		for _, c := range codes {
			code := map[string]interface{}{}

			code["product_code"] = *c.ProductCode
			code["type"] = *c.Type

			res = append(res, code)
		}
	}

	return res
}

func getStateReason(reason *fcu.StateReason) map[string]interface{} {
	res := map[string]interface{}{}
	if reason != nil {
		res["code"] = reason.Code
		res["message"] = reason.Message
	}
	return res
}

func getTagSet(tags []*fcu.Tag) []map[string]interface{} {
	res := []map[string]interface{}{}

	if tags != nil {
		for _, t := range tags {
			tag := map[string]interface{}{}

			tag["key"] = *t.Key
			tag["value"] = *t.Value

			res = append(res, tag)
		}
	}

	return res
}

func flattenEBS(ebs *fcu.EbsInstanceBlockDevice) map[string]interface{} {

	res := map[string]interface{}{
		"delete_on_termination": fmt.Sprintf("%t", *ebs.DeleteOnTermination),
		"status":                *ebs.Status,
		"volume_id":             *ebs.VolumeId,
	}

	return res
}

func getBlockDeviceMapping(blockDeviceMappings []*fcu.InstanceBlockDeviceMapping) []map[string]interface{} {
	blockDeviceMapping := []map[string]interface{}{}

	if blockDeviceMapping != nil {
		for _, mapping := range blockDeviceMappings {
			r := map[string]interface{}{}
			r["block_device_mapping"] = *mapping.DeviceName

			e := map[string]interface{}{}

			r["ebs"] = e

			blockDeviceMapping = append(blockDeviceMapping, r)
		}
	}

	return blockDeviceMapping
}

func getGroupSet(groupSet []*fcu.GroupIdentifier) []map[string]interface{} {
	res := []map[string]interface{}{}
	for _, g := range groupSet {

		r := map[string]interface{}{
			"group_id":   *g.GroupId,
			"group_name": *g.GroupName,
		}
		res = append(res, r)
	}

	return res
}

func getIAMInstanceProfile(profile *fcu.IamInstanceProfile) map[string]interface{} {
	iam := map[string]interface{}{}

	if profile != nil {
		iam["arn"] = *profile.Arn
		if profile.Id != nil {
			iam["id"] = *profile.Id
		}
	}

	return iam
}

func getInstanceState(state *fcu.InstanceState) map[string]interface{} {
	statem := map[string]interface{}{}

	statem["code"] = fmt.Sprintf("%d", *state.Code)
	statem["name"] = *state.Name

	return statem
}

func getMonitoring(monitoring *fcu.Monitoring) map[string]interface{} {
	monitoringm := map[string]interface{}{}

	monitoringm["state"] = *monitoring.State

	return monitoringm
}

func getNetworkInterfaceSet(interfaces []*fcu.InstanceNetworkInterface) []map[string]interface{} {
	res := []map[string]interface{}{}

	if interfaces != nil {
		for _, i := range interfaces {
			inter := make(map[string]interface{})
			assoc := make(map[string]interface{})

			if i.Association != nil {
				assoc["ip_owner_id"] = *i.Association.IpOwnerId
				assoc["public_dns_name"] = *i.Association.PublicDnsName
				assoc["public_ip"] = *i.Association.PublicIp
			}

			attch := make(map[string]interface{})
			assoc["attachement_id"] = *i.Attachment.AttachmentId
			assoc["delete_on_termination"] = fmt.Sprintf("%t", *i.Attachment.DeleteOnTermination)
			assoc["device_index"] = fmt.Sprintf("%d", *i.Attachment.DeviceIndex)
			assoc["status"] = *i.Attachment.Status

			inter["association"] = assoc
			inter["attachment"] = attch

			inter["description"] = *i.Description
			inter["group_set"] = getGroupSet(i.Groups)
			inter["mac_address"] = *i.MacAddress
			inter["network_interface_id"] = *i.NetworkInterfaceId
			inter["owner_id"] = *i.OwnerId
			inter["private_dns_name"] = *i.PrivateDnsName
			inter["private_ip_address"] = *i.PrivateIpAddress
			inter["private_ip_addresses_set"] = getPrivateIPAddressSet(i.PrivateIpAddresses)
			inter["source_dest_check"] = *i.SourceDestCheck
			inter["status"] = *i.Status
			inter["vpc_id"] = *i.VpcId

			res = append(res, inter)
		}
	}

	return res
}
