package outscale

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	oscgo "github.com/marinsalinas/osc-sdk-go"
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

func getOAPILinkNicLight(l oscgo.LinkNicLight) map[string]interface{} {
	return map[string]interface{}{
		"delete_on_vm_deletion": strconv.FormatBool(l.GetDeleteOnVmDeletion()),
		"device_number":         fmt.Sprintf("%d", l.GetDeviceNumber()),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
	}
}

func getOAPILinkNic(l oscgo.LinkNic) map[string]interface{} {
	return map[string]interface{}{
		"delete_on_vm_deletion": strconv.FormatBool(l.GetDeleteOnVmDeletion()),
		"device_number":         fmt.Sprintf("%d", l.GetDeviceNumber()),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
		"vm_account_id":         l.GetVmAccountId(),
		"vm_id":                 l.GetVmId(),
	}
}

func getOAPILinkPublicIPLight(l oscgo.LinkPublicIpLightForVm) *schema.Set {
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
			return hashcode.String(buf.String())
		},
	}

	for _, p := range privateIPs {
		r := map[string]interface{}{
			"is_primary":       p.GetIsPrimary(),
			"link_public_ip":   getOAPILinkPublicIPLight(p.GetLinkPublicIp()),
			"private_dns_name": p.GetPrivateDnsName(),
			"private_ip":       p.GetPrivateIp(),
		}
		res.Add(r)
	}
	return res
}

func getOAPIPrivateIPs(privateIPs []oscgo.PrivateIp) (res []map[string]interface{}) {
	for _, p := range privateIPs {
		res = append(res, map[string]interface{}{
			"is_primary":       p.GetIsPrimary(),
			"link_public_ip":   getOAPILinkPublicIP(p.GetLinkPublicIp()),
			"private_dns_name": p.GetPrivateDnsName(),
			"private_ip":       p.GetPrivateIp(),
		})
	}
	return
}

func getOAPIVMNetworkInterfaceLightSet(nics []oscgo.NicLight) (res []map[string]interface{}) {
	if nics != nil {
		for _, nic := range nics {
			securityGroups, securityGroupIds := getOAPISecurityGroups(nic.GetSecurityGroups())

			nicMap := map[string]interface{}{
				"delete_on_vm_deletion":  *nic.GetLinkNic().DeleteOnVmDeletion, // Workaround.
				"account_id":             nic.GetAccountId(),
				"description":            nic.GetDescription(),
				"is_source_dest_checked": nic.GetIsSourceDestChecked(),
				"link_nic":               getOAPILinkNicLight(nic.GetLinkNic()),
				"mac_address":            nic.GetMacAddress(),
				"net_id":                 nic.GetNetId(),
				"nic_id":                 nic.GetNicId(),
				"private_dns_name":       nic.GetPrivateDnsName(),
				"private_ips":            getOAPIPrivateIPsLight(nic.GetPrivateIps()),
				"security_groups":        securityGroups,
				"security_group_ids":     securityGroupIds,
				"state":                  nic.GetState(),
				"subnet_id":              nic.GetSubnetId(),
			}

			if nic.HasLinkPublicIp() {
				nicMap["link_public_ip"] = getOAPILinkPublicIPLight(nic.GetLinkPublicIp())
			}

			res = append(res, nicMap)
		}
	}
	return
}

func getOAPIVMNetworkInterfaceSet(nics []oscgo.Nic) (res []map[string]interface{}) {
	if nics != nil {
		for _, nic := range nics {
			//securityGroups, _ := getOAPISecurityGroups(nic.SecurityGroups)
			res = append(res, map[string]interface{}{
				"account_id":             nic.GetAccountId(),
				"description":            nic.GetDescription(),
				"is_source_dest_checked": nic.GetIsSourceDestChecked(),
				"link_nic":               getOAPILinkNic(nic.GetLinkNic()),
				"link_public_ip":         getOAPILinkPublicIP(nic.GetLinkPublicIp()),
				"mac_address":            nic.GetMacAddress(),
				"net_id":                 nic.GetNetId(),
				"nic_id":                 nic.GetNicId(),
				"private_dns_name":       nic.GetPrivateDnsName(),
				"private_ips":            getOAPIPrivateIPs(nic.GetPrivateIps()),
				//"security_groups":        securityGroups,
				"state":          nic.GetState(),
				"subnet_id":      nic.GetSubnetId(),
				"subregion_name": nic.GetSubregionName(),
				// "tags":           getOapiTagSet(nic.Tags),
			})
		}
	}
	return
}
