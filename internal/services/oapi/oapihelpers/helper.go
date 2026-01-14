package oapihelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func GetError(err error) osc.Errors {
	if e, ok := err.(osc.GenericOpenAPIError); ok {
		var errorResponse osc.ErrorResponse
		if json.Unmarshal(e.Body(), &errorResponse) == nil {
			errors := errorResponse.GetErrors()
			if len(errors) > 0 {
				return errors[0]
			}
		}
	}
	return osc.Errors{}
}

func GetBsuId(vmResp osc.Vm, deviceName string) string {
	diskID := ""
	blocks := vmResp.GetBlockDeviceMappings()

	for _, v := range blocks {
		if v.GetDeviceName() == deviceName {
			diskID = ptr.From(v.GetBsu().VolumeId)
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
		volumeId := ptr.From(v.GetBsu().VolumeId)
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

func getOAPILinkNicLight(l osc.LinkNicLight) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": l.GetDeleteOnVmDeletion(),
		"device_number":         strconv.Itoa(int(l.GetDeviceNumber())),
		"link_nic_id":           l.GetLinkNicId(),
		"state":                 l.GetState(),
	}}
}

func GetOAPILinkNic(l osc.LinkNic) []map[string]interface{} {
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
			var buf strings.Builder
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["public_ip_account_id"].(string)))
			return schema.HashString(buf.String())
		},
	}

	res.Add(map[string]interface{}{
		"public_dns_name":      l.GetPublicDnsName(),
		"public_ip":            l.GetPublicIp(),
		"public_ip_account_id": l.GetPublicIpAccountId(),
	})
	return res
}

func getOAPIPrivateIPsLight(privateIPs []osc.PrivateIpLightForVm) *schema.Set {
	res := &schema.Set{
		F: func(v interface{}) int {
			var buf bytes.Buffer
			m := v.(map[string]interface{})
			buf.WriteString(fmt.Sprintf("%s-", m["private_ip"].(string)))
			buf.WriteString(fmt.Sprintf("%s-", m["private_dns_name"].(string)))
			return schema.HashString(buf.String())
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

func GetPrivateIPsForNic(privateIPs []osc.PrivateIp) (res []map[string]interface{}) {
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

func GetOAPIVMNetworkInterfaceLightSet(respNics []osc.NicLight) ([]map[string]interface{}, []map[string]interface{}) {
	primaryNic := make([]map[string]interface{}, 0, 1)
	nics := make([]map[string]interface{}, 0, len(respNics))

	for _, nic := range respNics {
		securityGroups, securityGroupIds := GetSecurityGroups(nic.GetSecurityGroups())

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
		nics = append(nics, nicMap)
	}
	return primaryNic, nics
}

func GetSecurityGroups(groups []osc.SecurityGroupLight) (SecurityGroup []map[string]interface{}, SecurityGroupIds []string) {
	for _, g := range groups {
		SecurityGroup = append(SecurityGroup, map[string]interface{}{
			"security_group_id":   g.GetSecurityGroupId(),
			"security_group_name": g.GetSecurityGroupName(),
		})
		SecurityGroupIds = append(SecurityGroupIds, g.GetSecurityGroupId())
	}
	return
}

func ImageHasLaunchPermission(conn *osc.APIClient, imageID string) (bool, error) {
	var resp osc.ReadImagesResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.ImageApi.ReadImages(context.Background()).ReadImagesRequest(osc.ReadImagesRequest{
			Filters: &osc.FiltersImage{
				ImageIds: &[]string{imageID},
			},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string
	if err != nil {
		// When an AMI disappears out from under a launch permission resource, we will
		// see either InvalidAMIID.NotFound or InvalidAMIID.Unavailable.
		if strings.Contains(fmt.Sprint(err), "InvalidAMIID") {
			log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", imageID)
			return false, nil
		}
		errString = err.Error()

		return false, fmt.Errorf("error creating outscale vm volume: %s", errString)
	}

	if len(resp.GetImages()) == 0 {
		log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", imageID)
		return false, nil
	}

	result := resp.GetImages()[0]

	if len(result.PermissionsToLaunch.GetAccountIds()) > 0 {
		return true, nil
	}
	return false, nil
}

func ParseVPNConnectionRouteID(ID string) (destinationIPRange, vpnConnectionID string) {
	parts := strings.SplitN(ID, ":", 2)
	return parts[0], parts[1]
}
