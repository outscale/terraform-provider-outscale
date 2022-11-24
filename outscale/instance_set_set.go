package outscale

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func getSecurityGroups(groups []oscgo.SecurityGroupLight) (SecurityGroup []map[string]interface{}, SecurityGroupIds []string) {
	for _, g := range groups {
		SecurityGroup = append(SecurityGroup, map[string]interface{}{
			"security_group_id":   g.GetSecurityGroupId(),
			"security_group_name": g.GetSecurityGroupName(),
		})
		SecurityGroupIds = append(SecurityGroupIds, g.GetSecurityGroupId())
	}
	return
}

func getLinkNicLight(l oscgo.LinkNicLight) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": l.GetDeleteOnVmDeletion(),
		"device_number":         fmt.Sprintf("%d", l.GetDeviceNumber()),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
	}}
}

func getLinkNic(l oscgo.LinkNic) map[string]interface{} {
	return map[string]interface{}{
		"delete_on_vm_deletion": strconv.FormatBool(l.GetDeleteOnVmDeletion()),
		"device_number":         fmt.Sprintf("%d", l.GetDeviceNumber()),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
		"vm_account_id":         l.GetVmAccountId(),
		"vm_id":                 l.GetVmId(),
	}
}

func getLinkPublicIPLight(l oscgo.LinkPublicIpLightForVm) *schema.Set {
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
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
	})
	return res
}

func getLinkPublicIP(l oscgo.LinkPublicIp) map[string]interface{} {
	return map[string]interface{}{
		"link_public_ip_id":    l.GetLinkPublicIpId(),
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
		"public_ip_id":         l.GetPublicIpId(),
	}
}

func getPrivateIPsLight(privateIPs []oscgo.PrivateIpLightForVm) *schema.Set {
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
			"is_primary":       p.GetIsPrimary(),
			"private_dns_name": p.GetPrivateDnsName(),
			"private_ip":       p.GetPrivateIp(),
		}

		if p.HasLinkPublicIp() {
			r["link_public_ip"] = getLinkPublicIPLight(p.GetLinkPublicIp())
		}

		res.Add(r)
	}
	return res
}

func getPrivateIPs(privateIPs []oscgo.PrivateIp) (res []map[string]interface{}) {
	for _, p := range privateIPs {
		res = append(res, map[string]interface{}{
			"is_primary":       p.GetIsPrimary(),
			"link_public_ip":   getLinkPublicIP(p.GetLinkPublicIp()),
			"private_dns_name": p.GetPrivateDnsName(),
			"private_ip":       p.GetPrivateIp(),
		})
	}
	return
}

func getVMNetworkInterfaceLightSet(nics []oscgo.NicLight) (res []map[string]interface{}) {
	for _, nic := range nics {
		securityGroups, securityGroupIds := getSecurityGroups(nic.GetSecurityGroups())

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
			nicMap["link_public_ip"] = getLinkPublicIPLight(nic.GetLinkPublicIp())
		}

		if nic.HasPrivateIps() {
			nicMap["private_ips"] = getPrivateIPsLight(nic.GetPrivateIps())
		}

		if nic.HasLinkNic() {
			nicMap["link_nic"] = getLinkNicLight(nic.GetLinkNic())
		}

		res = append(res, nicMap)
	}
	return
}

func getVMNetworkInterfaceSet(nics []oscgo.Nic) (res []map[string]interface{}) {
	for _, nic := range nics {
		securityGroups, _ := getSecurityGroups(*nic.SecurityGroups)
		res = append(res, map[string]interface{}{
			"account_id":             nic.GetAccountId(),
			"description":            nic.GetDescription(),
			"is_source_dest_checked": nic.GetIsSourceDestChecked(),
			"link_nic":               getLinkNic(nic.GetLinkNic()),
			"link_public_ip":         getLinkPublicIP(nic.GetLinkPublicIp()),
			"mac_address":            nic.GetMacAddress(),
			"net_id":                 nic.GetNetId(),
			"nic_id":                 nic.GetNicId(),
			"private_dns_name":       nic.GetPrivateDnsName(),
			"private_ips":            getPrivateIPs(nic.GetPrivateIps()),
			"security_groups":        securityGroups,
			"state":                  nic.GetState(),
			"subnet_id":              nic.GetSubnetId(),
			"subregion_name":         nic.GetSubregionName(),
			"tags":                   getTagSet(nic.GetTags()),
		})
	}
	return
}
