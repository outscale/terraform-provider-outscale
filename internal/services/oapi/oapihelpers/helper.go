package oapihelpers

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func GetError(err error) osc.Errors {
	e := osc.AsErrorResponse(err)
	if e != nil && len(e.Errors) > 0 {
		return e.Errors[0]
	}
	return osc.Errors{}
}

func GetBsuId(vmResp osc.Vm, deviceName string) string {
	diskID := ""
	blocks := vmResp.BlockDeviceMappings

	for _, v := range blocks {
		if v.DeviceName == deviceName {
			diskID = v.Bsu.VolumeId
			break
		}
	}
	return diskID
}

func getBsuTags(ctx context.Context, client *osc.Client, timeout time.Duration, volumeId string) ([]osc.ResourceTag, error) {
	request := osc.ReadVolumesRequest{
		Filters: &osc.FiltersVolume{VolumeIds: &[]string{volumeId}},
	}
	resp, err := client.ReadVolumes(ctx, request, options.WithRetryTimeout(timeout))
	if err != nil {
		return nil, err
	}
	if resp.Volumes == nil || len(*resp.Volumes) == 0 {
		return nil, fmt.Errorf("volume %s not found", volumeId)
	}

	return (*resp.Volumes)[0].Tags, nil
}

func GetBsuTagsMaps(ctx context.Context, client *osc.Client, timeout time.Duration, vmResp osc.Vm) (map[string]interface{}, error) {
	blocks := vmResp.BlockDeviceMappings
	bsuTagsMaps := make(map[string]interface{})
	for _, v := range blocks {
		volumeId := v.Bsu.VolumeId
		bsuTags, err := getBsuTags(ctx, client, timeout, volumeId)
		if err != nil {
			return nil, err
		}
		if bsuTags != nil {
			bsuTagsMaps[v.DeviceName] = bsuTags
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
		"delete_on_vm_deletion": l.DeleteOnVmDeletion,
		"device_number":         strconv.Itoa(int(l.DeviceNumber)),
		"link_nic_id":           l.LinkNicId,
		"state":                 l.State,
	}}
}

func GetOAPILinkNic(l osc.LinkNic) []map[string]interface{} {
	return []map[string]interface{}{{
		"delete_on_vm_deletion": l.DeleteOnVmDeletion,
		"device_number":         l.DeviceNumber,
		"link_nic_id":           l.LinkNicId,
		"state":                 l.State,
		"vm_account_id":         l.VmAccountId,
		"vm_id":                 l.VmId,
	}}
}

func GetOAPILinkPublicIPsForNic(l osc.LinkPublicIp) []map[string]interface{} {
	return []map[string]interface{}{{
		"link_public_ip_id":    l.LinkPublicIpId,
		"public_dns_name":      l.PublicDnsName,
		"public_ip":            l.PublicIp,
		"public_ip_account_id": l.PublicIpAccountId,
		"public_ip_id":         l.PublicIpId,
	}}
}

func getOAPILinkPublicIpsForVm(l osc.LinkPublicIpLightForVm) []map[string]interface{} {
	return []map[string]interface{}{{
		"public_dns_name":      l.PublicDnsName,
		"public_ip":            l.PublicIp,
		"public_ip_account_id": l.PublicIpAccountId,
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
		"public_dns_name":      l.PublicDnsName,
		"public_ip":            l.PublicIp,
		"public_ip_account_id": l.PublicIpAccountId,
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
			"is_primary":       p.IsPrimary,
			"private_dns_name": p.PrivateDnsName,
			"private_ip":       p.PrivateIp,
		}

		if p.LinkPublicIp != nil {
			r["link_public_ip"] = getOAPILinkPublicIPLight(*p.LinkPublicIp)
		}

		res.Add(r)
	}
	return res
}

func GetPrivateIPsForNic(privateIPs []osc.PrivateIp) (res []map[string]interface{}) {
	for _, p := range privateIPs {
		r := map[string]interface{}{
			"is_primary":       p.IsPrimary,
			"private_dns_name": p.PrivateDnsName,
			"private_ip":       p.PrivateIp,
		}
		if p.LinkPublicIp != nil {
			r["link_public_ip"] = GetOAPILinkPublicIPsForNic(*p.LinkPublicIp)
		}
		res = append(res, r)
	}
	return
}

func GetOAPIVMNetworkInterfaceLightSet(respNics []osc.NicLight) ([]map[string]interface{}, []map[string]interface{}) {
	primaryNic := make([]map[string]interface{}, 0, 1)
	nics := make([]map[string]interface{}, 0, len(respNics))

	for _, nic := range respNics {
		securityGroups, securityGroupIds := GetSecurityGroups(nic.SecurityGroups)

		nicMap := map[string]interface{}{
			"delete_on_vm_deletion":  ptr.From(nic.LinkNic).DeleteOnVmDeletion, // Workaround.
			"device_number":          ptr.From(nic.LinkNic).DeviceNumber,
			"account_id":             nic.AccountId,
			"is_source_dest_checked": nic.IsSourceDestChecked,
			"mac_address":            nic.MacAddress,
			"net_id":                 nic.NetId,
			"nic_id":                 nic.NicId,
			"private_dns_name":       nic.PrivateDnsName,
			"security_groups":        securityGroups,
			"security_group_ids":     securityGroupIds,
			"state":                  nic.State,
			"subnet_id":              nic.SubnetId,
		}
		nicMap["description"] = nic.Description
		nicMap["private_ips"] = getOAPIPrivateIPsLight(nic.PrivateIps)

		if nic.LinkPublicIp != nil {
			nicMap["link_public_ip"] = getOAPILinkPublicIpsForVm(*nic.LinkPublicIp)
		}
		if nic.LinkNic != nil {
			nicMap["link_nic"] = getOAPILinkNicLight(*nic.LinkNic)
		}
		if nic.LinkNic.DeviceNumber == 0 {
			primaryNic = append(primaryNic, nicMap)
		}

		nics = append(nics, nicMap)
	}
	return primaryNic, nics
}

func GetSecurityGroups(groups []osc.SecurityGroupLight) (SecurityGroup []map[string]interface{}, SecurityGroupIds []string) {
	for _, g := range groups {
		SecurityGroup = append(SecurityGroup, map[string]interface{}{
			"security_group_id":   g.SecurityGroupId,
			"security_group_name": g.SecurityGroupName,
		})
		SecurityGroupIds = append(SecurityGroupIds, g.SecurityGroupId)
	}
	return
}

func ImageHasLaunchPermission(ctx context.Context, client *osc.Client, timeout time.Duration, imageID string) (bool, error) {
	resp, err := client.ReadImages(ctx, osc.ReadImagesRequest{
		Filters: &osc.FiltersImage{
			ImageIds: &[]string{imageID},
		},
	}, options.WithRetryTimeout(timeout))

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

	if resp.Images == nil || len(*resp.Images) == 0 {
		log.Printf("[DEBUG] %s no longer exists, so we'll drop launch permission for the state", imageID)
		return false, nil
	}

	result := (*resp.Images)[0]

	if len(ptr.From(result.PermissionsToLaunch.AccountIds)) > 0 {
		return true, nil
	}
	return false, nil
}

func ParseVPNConnectionRouteID(ID string) (destinationIPRange, vpnconnectionID string) {
	parts := strings.SplitN(ID, ":", 2)
	return parts[0], parts[1]
}

func RetryOnCodes(ctx context.Context, codes []string, fun func() (resp any, err error), timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err := fun()
		if err != nil {
			oscErr := GetError(err)
			if slices.Contains(codes, oscErr.Code) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
}
