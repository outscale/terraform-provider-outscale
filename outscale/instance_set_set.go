package outscale

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func getOAPISecurityGroups(groups []oscgo.SecurityGroupLight) (SecurityGroup []map[string]interface{}, SecurityGroupIds []string) {
	for _, g := range groups {
		SecurityGroup = append(SecurityGroup, map[string]interface{}{
			"security_group_id":   g.GetSecurityGroupId(),
			"security_group_name": g.GetSecurityGroupName(),
		})
		SecurityGroupIds = append(SecurityGroupIds, g.GetSecurityGroupId())
	}
	return
}

func getOAPILinkNicLight(l oscgo.LinkNicLight) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": l.GetDeleteOnVmDeletion(),
		"device_number":         fmt.Sprintf("%d", l.GetDeviceNumber()),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
	}}
}

func getOAPILinkNic(l oscgo.LinkNic) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": l.GetDeleteOnVmDeletion(),
		"device_number":         l.GetDeviceNumber(),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
		"vm_account_id":         l.GetVmAccountId(),
		"vm_id":                 l.GetVmId(),
	}}
}

func getOAPIBsuSet(bsu oscgo.BsuCreated) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": bsu.GetDeleteOnVmDeletion(),
		"link_date":             bsu.GetLinkDate(),
		"state":                 bsu.GetState(),
		"volume_id":             bsu.GetVolumeId(),
	}}
}

func getOAPILinkPublicIPLight(l oscgo.LinkPublicIpLightForVm) *schema.Set {
	res := &schema.Set{
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip_account_id"].(string)))
			return utils.String(buf.String())
		},
	}

	res.Add(map[string]interface{}{
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
	})
	return res
}

func getOAPILinkPublicIP(l oscgo.LinkPublicIp) map[string]interface{} {
	return map[string]interface{}{
		"link_public_ip_id":    l.GetLinkPublicIpId(),
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
		"public_ip_id":         l.GetPublicIpId(),
	}
}

func getOAPIPrivateIPsLight(privateIPs []oscgo.PrivateIpLightForVm) *schema.Set {
	res := &schema.Set{
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["private_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["private_dns_name"].(string)))
			return utils.String(buf.String())
		},
	}

	for _, p := range privateIPs {
		r := map[string]interface{}{
			"is_primary":       p.GetIsPrimary(),
			"private_dns_name": p.GetPrivateDnsName(),
			"private_ip":       p.GetPrivateIp(),
		}

		if p.HasLinkPublicIp() {
			r["link_public_ip"] = getOAPILinkPublicIPLight(p.GetLinkPublicIp())
		}

		res.Add(r)
	}
	return res
}

func getOAPIPrivateIPs(privateIPs []oscgo.PrivateIp) (res []map[string]interface{}) {
	for _, p := range privateIPs {
		r := map[string]interface{}{
			"is_primary":       p.GetIsPrimary(),
			"private_dns_name": p.GetPrivateDnsName(),
			"private_ip":       p.GetPrivateIp(),
		}
		if _, ok := p.GetLinkPublicIpOk(); ok {
			r["link_public_ip"] = getOAPILinkPublicIP(p.GetLinkPublicIp())
		}
		res = append(res, r)
	}
	return
}

func getOAPIVMNetworkInterfaceLightSet(nics []oscgo.NicLight) (primaryNic []map[string]interface{}, secondaryNic []map[string]interface{}) {
	for _, nic := range nics {
		securityGroups, securityGroupIds := getOAPISecurityGroups(nic.GetSecurityGroups())

		nicMap := map[string]interface{}{
			"delete_on_vm_deletion":  nic.LinkNic.GetDeleteOnVmDeletion(), // Workaround.
			"device_number":          nic.LinkNic.GetDeviceNumber(),
			"account_id":             nic.GetAccountId(),
			"is_source_dest_checked": nic.GetIsSourceDestChecked(),
			"mac_address":            nic.GetMacAddress(),
			"net_id":                 nic.GetNetId(),
			"nic_id":                 nic.GetNicId(),
			"private_dns_name":       nic.GetPrivateDnsName(),
			"security_groups":        securityGroups,
			"security_group_ids":     securityGroupIds,
			"state":                  nic.GetState(),
			"subnet_id":              nic.GetSubnetId(),
		}

		if nic.HasDescription() {
			nicMap["description"] = nic.GetDescription()
		}

		if nic.HasLinkPublicIp() {
			nicMap["link_public_ip"] = getOAPILinkPublicIPLight(nic.GetLinkPublicIp())
		}

		if nic.HasPrivateIps() {
			nicMap["private_ips"] = getOAPIPrivateIPsLight(nic.GetPrivateIps())
		}

		if nic.HasLinkNic() {
			nicMap["link_nic"] = getOAPILinkNicLight(nic.GetLinkNic())
		}
		if nic.LinkNic.GetDeviceNumber() == 0 {
			primaryNic = append(primaryNic, nicMap)
		}
		secondaryNic = append(secondaryNic, nicMap)
	}
	return
}

func getOAPIVMNetworkInterfaceSet(nics []oscgo.Nic) (res []map[string]interface{}) {

	for _, nic := range nics {
		securityGroups, _ := getOAPISecurityGroups(*nic.SecurityGroups)
		r := map[string]interface{}{
			"account_id":             nic.GetAccountId(),
			"description":            nic.GetDescription(),
			"is_source_dest_checked": nic.GetIsSourceDestChecked(),
			"mac_address":            nic.GetMacAddress(),
			"net_id":                 nic.GetNetId(),
			"nic_id":                 nic.GetNicId(),
			"private_dns_name":       nic.GetPrivateDnsName(),
			"private_ips":            getOAPIPrivateIPs(nic.GetPrivateIps()),
			"security_groups":        securityGroups,
			"state":                  nic.GetState(),
			"subnet_id":              nic.GetSubnetId(),
			"subregion_name":         nic.GetSubregionName(),
			"tags":                   getOapiTagSet(nic.Tags),
		}
		if _, ok := nic.GetLinkNicOk(); ok {
			r["link_nic"] = getOAPILinkNic(nic.GetLinkNic())
		}
		if _, ok := nic.GetLinkPublicIpOk(); ok {
			r["link_public_ip"] = getOAPILinkPublicIP(nic.GetLinkPublicIp())
		}
		res = append(res, r)
	}

	return
}
