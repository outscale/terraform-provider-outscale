package oapihelpers

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

const (
	pathRegex string = "^(/[a-zA-Z0-9/_]+/)"
	pathError string = "path must begin and end with '/' and contain only alphanumeric characters and/or '/', '_' characters"
)

func GetBsuId(vmResp osc.Vm, deviceName string) string {
	diskID := ""
	blocks := vmResp.GetBlockDeviceMappings()

	for _, v := range blocks {
		if v.GetDeviceName() == deviceName {
			diskID = aws.StringValue(v.GetBsu().VolumeId)
			break
		}
	}
	return diskID
}

func getBsuTags(volumeId string, conn *osc.APIClient) ([]osc.ResourceTag, error) {
	request := osc.ReadVolumesRequest{
		Filters: &osc.FiltersVolume{VolumeIds: &[]string{volumeId}},
	}
	var resp osc.ReadVolumesResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		r, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = r
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp.GetVolumes()[0].GetTags(), nil
}

func GetBsuTagsMaps(vmResp osc.Vm, conn *osc.APIClient) (map[string]interface{}, error) {
	blocks := vmResp.GetBlockDeviceMappings()
	bsuTagsMaps := make(map[string]interface{})
	for _, v := range blocks {
		volumeId := aws.StringValue(v.GetBsu().VolumeId)
		bsuTags, err := getBsuTags(volumeId, conn)
		if err != nil {
			return nil, err
		}
		if bsuTags != nil {
			bsuTagsMaps[v.GetDeviceName()] = bsuTags
		}
	}

	return bsuTagsMaps, nil
}

func RandVpcCidr() string {
	var result string
	prefix := utils.RandIntRange(16, 29)
	switch rand.Intn(3) {
	case 0:
		// 10.0.0.0 - 10.255.255.255 (10/8 prefix)
		result = fmt.Sprintf("10.%d.0.0/%d", rand.Intn(256), prefix)
	case 1:
		// 172.16.0.0 - 172.31.255.255 (172.16/12 prefix)
		result = fmt.Sprintf("172.%d.0.0/%d", utils.RandIntRange(16, 32), prefix)
	case 2:
		// 192.168.0.0 - 192.168.255.255 (192.168/16 prefix)
		result = fmt.Sprintf("192.168.0.0/%d", prefix)
	}
	return result
}

func GetAccepterOwnerId() string {
	accountId := os.Getenv("OUTSCALE_ACCOUNT")
	if accountId == "" {
		accountId = os.Getenv("OSC_ACCOUNT")
	}
	return accountId
}

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func Strings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}

func GetTypeSetDifferencesForUpdating(oldTypeSet, newTypeSet *schema.Set) (*schema.Set, *schema.Set) {
	inter := oldTypeSet.Intersection(newTypeSet)
	toAdd := newTypeSet.Difference(inter)
	toRemove := oldTypeSet.Difference(inter)
	return toRemove, toAdd
}

func GetAttrTypes(model any) map[string]attr.Type {
	attrTypes := make(map[string]attr.Type)

	v := reflect.ValueOf(model)
	t := v.Type()

	for i := range v.NumField() {
		field := t.Field(i)
		tfsdkTag := field.Tag.Get("tfsdk")
		if tfsdkTag == "" {
			continue
		}

		switch field.Type {
		case reflect.TypeOf(types.String{}):
			attrTypes[tfsdkTag] = types.StringType
		case reflect.TypeOf(types.Bool{}):
			attrTypes[tfsdkTag] = types.BoolType
		case reflect.TypeOf(types.Int64{}):
			attrTypes[tfsdkTag] = types.Int64Type
		case reflect.TypeOf(types.Float64{}):
			attrTypes[tfsdkTag] = types.Float64Type
		case reflect.TypeOf(types.Int32{}):
			attrTypes[tfsdkTag] = types.Int32Type
		}
	}
	return attrTypes
}

func RandBgpAsn() int {
	return utils.RandIntRange(1, 50620)
}

func getOAPISecurityGroups(groups []osc.SecurityGroupLight) (SecurityGroup []map[string]interface{}, SecurityGroupIds []string) {
	for _, g := range groups {
		SecurityGroup = append(SecurityGroup, map[string]interface{}{
			"security_group_id":   g.GetSecurityGroupId(),
			"security_group_name": g.GetSecurityGroupName(),
		})
		SecurityGroupIds = append(SecurityGroupIds, g.GetSecurityGroupId())
	}
	return
}

func getOAPILinkNicLight(l osc.LinkNicLight) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": l.GetDeleteOnVmDeletion(),
		"device_number":         fmt.Sprintf("%d", l.GetDeviceNumber()),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
	}}
}

func getOAPILinkNic(l osc.LinkNic) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": l.GetDeleteOnVmDeletion(),
		"device_number":         l.GetDeviceNumber(),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
		"vm_account_id":         l.GetVmAccountId(),
		"vm_id":                 l.GetVmId(),
	}}
}

func GetOAPILinkPublicIPsForNic(l osc.LinkPublicIp) []map[string]interface{} {
	return []map[string]interface{}{{
		"link_public_ip_id":    l.GetLinkPublicIpId(),
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
		"public_ip_id":         l.GetPublicIpId(),
	}}
}

func getOAPIBsuSet(bsu osc.BsuCreated) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": bsu.GetDeleteOnVmDeletion(),
		"link_date":             bsu.GetLinkDate(),
		"state":                 bsu.GetState(),
		"volume_id":             bsu.GetVolumeId(),
	}}
}

func getOAPILinkPublicIpsForVm(l osc.LinkPublicIpLightForVm) []map[string]interface{} {
	return []map[string]interface{}{{
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
	}}
}

func getOAPILinkPublicIPLight(l osc.LinkPublicIpLightForVm) *schema.Set {
	res := &schema.Set{
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip_account_id"].(string)))
			return String(buf.String())
		},
	}

	res.Add(map[string]interface{}{
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
	})
	return res
}

func getOAPILinkPublicIP(l osc.LinkPublicIp) map[string]interface{} {
	return map[string]interface{}{
		"link_public_ip_id":    l.GetLinkPublicIpId(),
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
		"public_ip_id":         l.GetPublicIpId(),
	}
}

func getOAPIPrivateIPsLight(privateIPs []osc.PrivateIpLightForVm) *schema.Set {
	res := &schema.Set{
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["private_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["private_dns_name"].(string)))
			return String(buf.String())
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

func GetOAPIPrivateIPsForNic(privateIPs []osc.PrivateIp) (res []map[string]interface{}) {
	for _, p := range privateIPs {
		r := map[string]interface{}{
			"is_primary":       p.GetIsPrimary(),
			"private_dns_name": p.GetPrivateDnsName(),
			"private_ip":       p.GetPrivateIp(),
		}
		if _, ok := p.GetLinkPublicIpOk(); ok {
			r["link_public_ip"] = GetOAPILinkPublicIPsForNic(p.GetLinkPublicIp())
		}
		res = append(res, r)
	}
	return
}

func GetOAPIVMNetworkInterfaceLightSet(nics []osc.NicLight) (primaryNic []map[string]interface{}, secondaryNic []map[string]interface{}) {
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
			nicMap["link_public_ip"] = getOAPILinkPublicIpsForVm(nic.GetLinkPublicIp())
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

func GetOAPIVMNetworkInterfaceSet(nics []osc.Nic) (res []map[string]interface{}) {
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
			"private_ips":            GetOAPIPrivateIPsForNic(nic.GetPrivateIps()),
			"security_groups":        securityGroups,
			"state":                  nic.GetState(),
			"subnet_id":              nic.GetSubnetId(),
			"subregion_name":         nic.GetSubregionName(),
			"tags":                   FlattenOAPITagsSDK(nic.GetTags()),
		}
		if _, ok := nic.GetLinkNicOk(); ok {
			r["link_nic"] = getOAPILinkNic(nic.GetLinkNic())
		}
		if _, ok := nic.GetLinkPublicIpOk(); ok {
			r["link_public_ip"] = GetOAPILinkPublicIPsForNic(nic.GetLinkPublicIp())
		}
		res = append(res, r)
	}

	return
}

func FlattenOAPITagsSDK(tags []osc.ResourceTag) []map[string]string {
	result := make([]map[string]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]string{
			"key":   tag.Key,
			"value": tag.Value,
		})
	}
	return result
}
