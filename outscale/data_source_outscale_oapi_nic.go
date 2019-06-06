package outscale

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

// Creates a network interface in the specified subnet
func dataSourceOutscaleOAPINic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPINicRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// This is attribute part for schema Nic
			// Argument
			"nic_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Attributes
			"description": &schema.Schema{
				Type: schema.TypeString,

				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type: schema.TypeString,

				Computed: true,
			},
			"security_group_id": &schema.Schema{
				Type: schema.TypeList,

				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// Attributes
			"link_public_ip": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link_public_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"link_nic": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"link_nic_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_on_vm_deletion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"subregion_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"link_public_ip": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_ip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"link_public_ip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip_account_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"private_dns_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"requester_managed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_source_dest_checked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_set": tagsSchemaComputed(),
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

//Read Nic
func dataSourceOutscaleOAPINicRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	nicID, okID := d.GetOk("nic_id")
	filters, okFilters := d.GetOk("filter")

	if okID && okFilters {
		return errors.New("nic_id and filter set")
	}

	dnri := oapi.ReadNicsRequest{}

	if okID {
		dnri.Filters = oapi.FiltersNic{
			NicIds: []string{nicID.(string)},
		}
	}

	if okFilters {
		dnri.Filters = buildOutscaleOAPIDataSourceNicFilters(filters.(*schema.Set))
	}

	var describeResp *oapi.POST_ReadNicsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		describeResp, err = conn.POST_ReadNics(dnri)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {

		return fmt.Errorf("Error describing Network Interfaces : %s", err)
	}

	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidNetworkInterfaceID.NotFound" {
			// The ENI is gone now, so just remove it from the state
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving ENI: %s", err)
	}
	if len(describeResp.OK.Nics) != 1 {
		return fmt.Errorf("Unable to find ENI: %#v", describeResp.OK.Nics)
	}

	eni := describeResp.OK.Nics[0]
	if eni.Description != "" {
		d.Set("description", eni.Description)
	}
	d.Set("nic_id", eni.SubnetId)

	b := make(map[string]interface{})

	link := eni.LinkPublicIp
	b["reservation_id"] = link.PublicIpId
	b["link_id"] = link.LinkPublicIpId
	b["public_ip_account_id"] = link.PublicIpAccountId
	b["public_dns_name"] = link.PublicDnsName
	b["public_ip"] = link.PublicIp

	if err := d.Set("public_ip_link", b); err != nil {
		return err
	}

	aa := make([]map[string]interface{}, 1)
	bb := make(map[string]interface{})

	linkNic := eni.LinkNic

	bb["nic_link_id"] = linkNic.LinkNicId
	bb["delete_on_vm_deletion"] = linkNic.DeleteOnVmDeletion
	bb["nic_sort_number"] = linkNic.DeviceNumber
	bb["vm_id"] = linkNic.VmId
	bb["vm_account_id"] = linkNic.VmAccountId
	bb["state"] = linkNic.State

	aa[0] = bb
	if err := d.Set("nic_link", aa); err != nil {
		return err
	}

	d.Set("sub_region_name", eni.SubregionName)

	x := make([]map[string]interface{}, len(eni.SecurityGroups))
	for k, v := range eni.SecurityGroups {
		b := make(map[string]interface{})
		b["firewall_rules_set_id"] = v.SecurityGroupId
		b["firewall_rules_set_name"] = v.SecurityGroupName
		x[k] = b
	}
	if err := d.Set("firewall_rules_set", x); err != nil {
		return err
	}

	d.Set("mac_address", eni.MacAddress)
	d.Set("nic_id", eni.NicId)
	d.Set("account_id", eni.AccountId)
	d.Set("private_dns_name", eni.PrivateDnsName)
	// Check this one later
	d.Set("private_ip_address", eni.NetId)

	y := make([]map[string]interface{}, len(eni.PrivateIps))
	if eni.PrivateIps != nil {
		for k, v := range eni.PrivateIps {
			b := make(map[string]interface{})

			d := make(map[string]interface{})
			linkPrivateIP := v.LinkPublicIp
			d["reservation_id"] = linkPrivateIP.PublicIpId
			d["link_id"] = linkPrivateIP.LinkPublicIpId
			d["public_ip_account_id"] = linkPrivateIP.PublicIpAccountId
			d["public_dns_name"] = linkPrivateIP.PublicDnsName
			d["public_ip"] = linkPrivateIP.PublicIpId
			b["association"] = d
			b["primary_ip"] = v.IsPrimary
			b["private_dns_name"] = v.PrivateDnsName
			b["private_ip"] = v.PrivateIp
			y[k] = b
		}
	}
	if err := d.Set("private_ip_set", y); err != nil {
		return err
	}

	d.Set("request_id", describeResp.OK.ResponseContext.RequestId)

	// Missing
	// d.Set("requester_managed", aws.BoolValue(eni.))

	// Missing
	// d.Set("activated_check", aws.BoolValue(eni.SourceDestCheck))
	d.Set("state", eni.State)
	d.Set("subnet_id", eni.SubnetId)
	// Tags
	d.Set("tags", tagsOAPIToMap(eni.Tags))
	d.Set("net_id", eni.NetId)

	return nil
}

func resourceOutscaleOAPIDataSourceNicDetach(oa []interface{}, meta interface{}, eniID string) error {
	// if there was an old attachment, remove it
	if oa != nil && len(oa) > 0 && oa[0] != nil {
		oa := oa[0].(map[string]interface{})
		dr := oapi.UnlinkNicRequest{
			LinkNicId: oa["attachment_id"].(string),
		}

		conn := meta.(*OutscaleClient).OAPI

		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			_, err = conn.POST_UnlinkNic(dr)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
				return fmt.Errorf("Error detaching oAPI ENI: %s", err)
			}
		}

		log.Printf("[DEBUG] Waiting for oAPI ENI (%s) to become dettached", eniID)
		stateConf := &resource.StateChangeConf{
			Pending: []string{"true"},
			Target:  []string{"false"},
			Refresh: networkInterfaceDataSourceOAPIAttachmentRefreshFunc(conn, eniID),
			Timeout: 10 * time.Minute,
		}
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf(
				"Error waiting for oAPI ENI (%s) to become dettached: %s", eniID, err)
		}
	}

	return nil
}

func networkInterfaceDataSourceOAPIAttachmentRefreshFunc(conn *oapi.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		dnri := oapi.ReadNicsRequest{
			Filters: oapi.FiltersNic{
				NicIds: []string{id},
			},
		}

		var describeResp *oapi.POST_ReadNicsResponses
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			describeResp, err = conn.POST_ReadNics(dnri)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			log.Printf("[ERROR] Could not find oAPI network interface %s. %s", id, err)
			return nil, "", err
		}

		eni := describeResp.OK.Nics[0]
		hasAttachment := strconv.FormatBool(eni.LinkNic.LinkNicId != "")
		log.Printf("[DEBUG] oAPI ENI %s has attachment state %s", id, hasAttachment)
		return eni, hasAttachment, nil
	}
}

func buildOutscaleOAPIDataSourceNicFilters(set *schema.Set) oapi.FiltersNic {
	var filters oapi.FiltersNic
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "net_ids":
			filters.NetIds = filterValues
		case "nic_ids":
			filters.NicIds = filterValues
		case "private_dns_names":
			filters.PrivateDnsNames = filterValues
		case "private_ips_link_public_ip_account_ids":
			filters.PrivateIpsLinkPublicIpAccountIds = filterValues
		case "private_ips_link_public_ip_public_ips":
			filters.PrivateIpsLinkPublicIpPublicIps = filterValues
		case "private_ips_private_ips":
			filters.PrivateIpsPrivateIps = filterValues
		case "security_group_ids":
			filters.SecurityGroupIds = filterValues
		case "security_group_names":
			filters.SecurityGroupNames = filterValues
		case "states":
			filters.States = filterValues
		case "subnet_ids":
			filters.SubnetIds = filterValues
		case "subregion_names":
			filters.SubregionNames = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
