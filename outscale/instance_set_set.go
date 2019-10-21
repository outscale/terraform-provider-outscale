package outscale

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func flattenedInstanceSet(instances []*fcu.Instance) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {

		flattened[i] = map[string]interface{}{
			"ami_launch_index":         aws.Int64Value(instance.AmiLaunchIndex),
			"ebs_optimized":            aws.BoolValue(instance.EbsOptimized),
			"architecture":             aws.StringValue(instance.Architecture),
			"client_token":             aws.StringValue(instance.ClientToken),
			"hypervisor":               aws.StringValue(instance.Hypervisor),
			"image_id":                 aws.StringValue(instance.ImageId),
			"instance_id":              aws.StringValue(instance.InstanceId),
			"instance_type":            aws.StringValue(instance.InstanceType),
			"kernel_id":                aws.StringValue(instance.KernelId),
			"key_name":                 aws.StringValue(instance.KeyName),
			"private_ip_address":       aws.StringValue(instance.PrivateDnsName),
			"private_dns_name":         aws.StringValue(instance.PrivateDnsName),
			"root_device_name":         aws.StringValue(instance.RootDeviceName),
			"instance_lifecycle":       aws.StringValue(instance.InstanceLifecycle),
			"root_device_type":         aws.StringValue(instance.RootDeviceType),
			"dns_name":                 aws.StringValue(instance.DnsName),
			"ip_address":               aws.StringValue(instance.IpAddress),
			"platform":                 aws.StringValue(instance.Platform),
			"ramdisk_id":               aws.StringValue(instance.RamdiskId),
			"reason":                   aws.StringValue(instance.Reason),
			"source_dest_check":        aws.BoolValue(instance.SourceDestCheck),
			"spot_instance_request_id": aws.StringValue(instance.SpotInstanceRequestId),
			"sriov_net_support":        aws.StringValue(instance.SriovNetSupport),
			"subnet_id":                aws.StringValue(instance.SubnetId),
			"virtualization_type":      aws.StringValue(instance.VirtualizationType),
			"vpc_id":                   aws.StringValue(instance.VpcId),
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
		flattened[i]["tag_set"] = tagsToMap(instance.Tags)
	}

	return flattened
}

func flattenedInstanceSetPassword(instances []*fcu.Instance, conn fcu.VMService) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		flattened[i] = map[string]interface{}{
			"ami_launch_index":         aws.Int64Value(instance.AmiLaunchIndex),
			"ebs_optimized":            aws.BoolValue(instance.EbsOptimized),
			"architecture":             aws.StringValue(instance.Architecture),
			"client_token":             aws.StringValue(instance.ClientToken),
			"hypervisor":               aws.StringValue(instance.Hypervisor),
			"image_id":                 aws.StringValue(instance.ImageId),
			"instance_id":              aws.StringValue(instance.InstanceId),
			"instance_type":            aws.StringValue(instance.InstanceType),
			"kernel_id":                aws.StringValue(instance.KernelId),
			"key_name":                 aws.StringValue(instance.KeyName),
			"private_ip_address":       aws.StringValue(instance.PrivateDnsName),
			"private_dns_name":         aws.StringValue(instance.PrivateDnsName),
			"root_device_name":         aws.StringValue(instance.RootDeviceName),
			"instance_lifecycle":       aws.StringValue(instance.InstanceLifecycle),
			"root_device_type":         aws.StringValue(instance.RootDeviceType),
			"dns_name":                 aws.StringValue(instance.DnsName),
			"ip_address":               aws.StringValue(instance.IpAddress),
			"ramdisk_id":               aws.StringValue(instance.RamdiskId),
			"reason":                   aws.StringValue(instance.Reason),
			"source_dest_check":        aws.BoolValue(instance.SourceDestCheck),
			"spot_instance_request_id": aws.StringValue(instance.SpotInstanceRequestId),
			"sriov_net_support":        aws.StringValue(instance.SriovNetSupport),
			"subnet_id":                aws.StringValue(instance.SubnetId),
			"virtualization_type":      aws.StringValue(instance.VirtualizationType),
			"vpc_id":                   aws.StringValue(instance.VpcId),
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
	var res []map[string]interface{}

	if len(codes) > 0 {
		res = make([]map[string]interface{}, len(codes))
		for _, c := range codes {
			code := map[string]interface{}{}

			code["product_code"] = *c.ProductCode
			code["type"] = *c.Type

			res = append(res, code)
		}
	} else {
		res = make([]map[string]interface{}, 0)
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

func getTagDescriptionSet(tags []*fcu.TagDescription) []map[string]interface{} {
	res := []map[string]interface{}{}

	if tags != nil {
		for _, t := range tags {
			tag := map[string]interface{}{}

			tag["key"] = *t.Key
			tag["value"] = *t.Value
			tag["resourceId"] = *t.ResourceId
			tag["resourceType"] = *t.ResourceType

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
	var blockDeviceMapping []map[string]interface{}

	if len(blockDeviceMappings) > 0 {
		blockDeviceMapping = make([]map[string]interface{}, len(blockDeviceMappings))
		for _, mapping := range blockDeviceMappings {
			r := map[string]interface{}{}
			r["device_name"] = *mapping.DeviceName

			e := map[string]interface{}{}
			e["delete_on_termination"] = *mapping.Ebs.DeleteOnTermination
			e["status"] = *mapping.Ebs.Status
			e["volume_id"] = *mapping.Ebs.VolumeId
			r["ebs"] = e

			blockDeviceMapping = append(blockDeviceMapping, r)
		}
	} else {
		blockDeviceMapping = make([]map[string]interface{}, 0)
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

func getOAPISecurityGroups(groupSet []oapi.SecurityGroupLight) (SecurityGroup []map[string]interface{}, SecurityGroupIds []string) {
	for _, g := range groupSet {
		SecurityGroup = append(SecurityGroup, map[string]interface{}{
			"security_group_id":   g.SecurityGroupId,
			"security_group_name": g.SecurityGroupName,
		})
		SecurityGroupIds = append(SecurityGroupIds, g.SecurityGroupId)
	}
	return
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

func getOAPILinkNicLight(l oapi.LinkNicLight) map[string]interface{} {
	return map[string]interface{}{
		"delete_on_vm_deletion": strconv.FormatBool(*l.DeleteOnVmDeletion),
		"device_number":         strconv.FormatInt(l.DeviceNumber, 10),
		"link_nic_id":           l.LinkNicId,
		"state":                 l.State,
	}
}

func getOAPILinkNic(l oapi.LinkNic) map[string]interface{} {
	return map[string]interface{}{
		"delete_on_vm_deletion": strconv.FormatBool(*l.DeleteOnVmDeletion),
		"device_number":         strconv.FormatInt(l.DeviceNumber, 10),
		"link_nic_id":           l.LinkNicId,
		"state":                 l.State,
		"vm_account_id":         l.VmAccountId,
		"vm_id":                 l.VmId,
	}
}

func getOAPILinkPublicIPLight(l oapi.LinkPublicIpLightForVm) *schema.Set {
	res := &schema.Set{
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip_account_id"].(string)))
			return hashcode.String(buf.String())
		},
	}

	res.Add(map[string]interface{}{
		"public_dns_name":      l.PublicDnsName,
		"public_ip":            l.PublicIp,
		"public_ip_account_id": l.PublicIpAccountId,
	})
	return res
}

func getOAPILinkPublicIP(l oapi.LinkPublicIp) map[string]interface{} {
	return map[string]interface{}{
		"link_public_ip_id":    l.LinkPublicIpId,
		"public_dns_name":      l.PublicDnsName,
		"public_ip":            l.PublicIp,
		"public_ip_account_id": l.PublicIpAccountId,
		"public_ip_id":         l.PublicIpId,
	}
}

func getOAPIPrivateIPsLight(privateIPs []oapi.PrivateIpLightForVm) *schema.Set {
	res := &schema.Set{
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["private_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["private_dns_name"].(string)))
			return hashcode.String(buf.String())
		},
	}

	for _, p := range privateIPs {
		r := map[string]interface{}{
			"is_primary":       p.IsPrimary,
			"link_public_ip":   getOAPILinkPublicIPLight(p.LinkPublicIp),
			"private_dns_name": p.PrivateDnsName,
			"private_ip":       p.PrivateIp,
		}
		res.Add(r)
	}
	return res
}

func getOAPIPrivateIPs(privateIPs []oapi.PrivateIp) (res []map[string]interface{}) {
	for _, p := range privateIPs {
		res = append(res, map[string]interface{}{
			"is_primary":       p.IsPrimary,
			"link_public_ip":   getOAPILinkPublicIP(p.LinkPublicIp),
			"private_dns_name": p.PrivateDnsName,
			"private_ip":       p.PrivateIp,
		})
	}
	return
}

func getOAPIVMNetworkInterfaceLightSet(nics []oapi.NicLight) (res []map[string]interface{}) {
	if nics != nil {
		for _, nic := range nics {
			securityGroups, securityGroupIds := getOAPISecurityGroups(nic.SecurityGroups)

			res = append(res, map[string]interface{}{
				"account_id":             nic.AccountId,
				"description":            nic.Description,
				"is_source_dest_checked": nic.IsSourceDestChecked,
				"link_nic":               getOAPILinkNicLight(nic.LinkNic),
				"link_public_ip":         getOAPILinkPublicIPLight(nic.LinkPublicIp),
				"mac_address":            nic.MacAddress,
				"net_id":                 nic.NetId,
				"nic_id":                 nic.NicId,
				"private_dns_name":       nic.PrivateDnsName,
				"private_ips":            getOAPIPrivateIPsLight(nic.PrivateIps),
				"security_groups":        securityGroups,
				"security_group_ids":     securityGroupIds,
				"state":                  nic.State,
				"subnet_id":              nic.SubnetId,
			})
		}
	}
	return
}

func getOAPIVMNetworkInterfaceSet(nics []oapi.Nic) (res []map[string]interface{}) {
	if nics != nil {
		for _, nic := range nics {
			securityGroups, _ := getOAPISecurityGroups(nic.SecurityGroups)

			res = append(res, map[string]interface{}{
				"account_id":             nic.AccountId,
				"description":            nic.Description,
				"is_source_dest_checked": nic.IsSourceDestChecked,
				"link_nic":               getOAPILinkNic(nic.LinkNic),
				"link_public_ip":         getOAPILinkPublicIP(nic.LinkPublicIp),
				"mac_address":            nic.MacAddress,
				"net_id":                 nic.NetId,
				"nic_id":                 nic.NicId,
				"private_dns_name":       nic.PrivateDnsName,
				"private_ips":            getOAPIPrivateIPs(nic.PrivateIps),
				"security_groups":        securityGroups,
				"state":                  nic.State,
				"subnet_id":              nic.SubnetId,
				"subregion_name":         nic.SubregionName,
				"tags":                   getOapiTagSet(nic.Tags),
			})
		}
	}
	return
}
