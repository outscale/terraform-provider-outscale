package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleDHCPOption() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleDHCPOptionCreate,
		Read:   resourceOutscaleDHCPOptionRead,
		Delete: resourceOutscaleDHCPOptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: getDHCPOptionSchema(),
	}
}

func getDHCPOptionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Attributes
		"dhcp_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"value": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"dhcp_configuration_set": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value_set": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"dhcp_options_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tag_set": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceOutscaleDHCPOptionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	setDHCPOption := func(key string) *fcu.NewDhcpConfiguration {
		log.Printf("[DEBUG] Setting DHCP option %s...", key)
		tfKey := strings.Replace(key, "-", "_", -1)

		value, ok := d.GetOk(tfKey)
		if !ok {
			return nil
		}

		if v, ok := value.(string); ok {
			return &fcu.NewDhcpConfiguration{
				Key: aws.String(key),
				Values: []*string{
					aws.String(v),
				},
			}
		}

		if v, ok := value.([]interface{}); ok {
			var s []*string
			for _, attr := range v {
				s = append(s, aws.String(attr.(string)))
			}

			return &fcu.NewDhcpConfiguration{
				Key:    aws.String(key),
				Values: s,
			}
		}

		return nil
	}

	var createOpts *fcu.CreateDhcpOptionsInput

	if v := setDHCPOption("dhcp-configuration"); v != nil {

		fmt.Printf("[DEBUG] INPUT %s", v)

		createOpts = &fcu.CreateDhcpOptionsInput{
			DhcpConfigurations: []*fcu.NewDhcpConfiguration{
				v,
			},
		}
	} else {
		createOpts = &fcu.CreateDhcpOptionsInput{}
		createOpts.DhcpConfigurations = []*fcu.NewDhcpConfiguration{
			&fcu.NewDhcpConfiguration{
				Key:    aws.String(""),
				Values: []*string{},
			},
		}
	}

	fmt.Printf("[DEBUG] VALUE => %s", createOpts)

	var resp *fcu.CreateDhcpOptionsOutput

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.CreateDhcpOptions(createOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating DHCP Options Set: %s", err)
	}

	dos := resp.DhcpOptions
	d.SetId(*dos.DhcpOptionsId)
	log.Printf("[INFO] DHCP Options Set ID: %s", d.Id())

	// Wait for the DHCP Options to become available
	log.Printf("[DEBUG] Waiting for DHCP Options (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"created"},
		Refresh: resourceDHCPOptionsStateRefreshFunc(conn, d.Id()),
		Timeout: 1 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for DHCP Options (%s) to become available: %s",
			d.Id(), err)
	}

	dhcp := make([]map[string]interface{}, 0)
	d.Set("dhcp_configuration_set", dhcp)

	return resourceOutscaleDHCPOptionRead(d, meta)
}

func resourceOutscaleDHCPOptionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	req := &fcu.DescribeDhcpOptionsInput{
		DhcpOptionsIds: []*string{
			aws.String(d.Id()),
		},
	}

	var resp *fcu.DescribeDhcpOptionsOutput

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeDhcpOptions(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error describing DHCP Options: %s", err)
	}

	// resp, err := conn.VM.DescribeDhcpOptions(req)
	// if err != nil {
	// 	ec2err, ok := err.(awserr.Error)
	// 	if !ok {
	// 		return fmt.Errorf("Error retrieving DHCP Options: %s", err.Error())
	// 	}

	// 	if ec2err.Code() == "InvalidDhcpOptionID.NotFound" {
	// 		log.Printf("[WARN] DHCP Options (%s) not found, removing from state", d.Id())
	// 		d.SetId("")
	// 		return nil
	// 	}

	// 	return fmt.Errorf("Error retrieving DHCP Options: %s", err.Error())
	// }

	// if len(resp.DhcpOptions) == 0 {
	// 	return nil
	// }

	opts := resp.DhcpOptions[0]
	d.Set("tag_set", tagsToMap(opts.Tags))

	dhcpConfiguration := make([]map[string]interface{}, len(resp.DhcpOptions))

	for k, cfg := range opts.DhcpConfigurations {

		dhcp := make(map[string]interface{})
		var values []string
		for _, v := range cfg.Values {
			values = append(values, *v.Value)
		}
		dhcp[*cfg.Key] = values
		dhcpConfiguration[k] = dhcp

	}
	d.Set("dhcp_options_id", d.Id())
	d.Set("request_id", resp.RequestId)

	return d.Set("dhcp_configuration_set", dhcpConfiguration)
}

func resourceOutscaleDHCPOptionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		log.Printf("[INFO] Deleting DHCP Options ID %s...", d.Id())

		// _, err := conn.VM.DeleteDhcpOptions(&fcu.DeleteDhcpOptionsInput{
		// 	DhcpOptionsId: aws.String(d.Id()),
		// })

		//	var resp *fcu.DeleteDhcpOptionsOutput

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			_, err = conn.VM.DeleteDhcpOptions(&fcu.DeleteDhcpOptionsInput{
				DhcpOptionsId: aws.String(d.Id()),
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		// if err != nil {
		// 	return fmt.Errorf("Error creating DHCP Options Set: %s, err", err)
		// }

		if err == nil {
			return nil
		}

		log.Printf("[WARN] %s", err)

		ec2err, ok := err.(awserr.Error)
		if !ok {
			return resource.RetryableError(err)
		}

		switch ec2err.Code() {
		case "InvalidDhcpOptionsID.NotFound":
			return nil
		case "DependencyViolation":
			// If it is a dependency violation, we want to disassociate
			// all VPCs using the given DHCP Options ID, and retry deleting.
			vpcs, err2 := findVPCsByDHCPOptionsID(conn, d.Id())
			if err2 != nil {
				log.Printf("[ERROR] %s", err2)
				return resource.RetryableError(err2)
			}

			for _, vpc := range vpcs {
				log.Printf("[INFO] Disassociating DHCP Options Set %s from VPC %s...", d.Id(), *vpc.VpcId)
				if _, err := conn.VM.AssociateDhcpOptions(&fcu.AssociateDhcpOptionsInput{
					DhcpOptionsId: aws.String("default"),
					VpcId:         vpc.VpcId,
				}); err != nil {
					return resource.RetryableError(err)
				}
			}
			return resource.RetryableError(err)
		default:
			return resource.NonRetryableError(err)
		}
	})
}

func resourceDHCPOptionsStateRefreshFunc(conn *fcu.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		DescribeDhcpOpts := &fcu.DescribeDhcpOptionsInput{
			DhcpOptionsIds: []*string{
				aws.String(id),
			},
		}

		//resp, err := conn.VM.DescribeDhcpOptions(DescribeDhcpOpts)

		var resp *fcu.DescribeDhcpOptionsOutput

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.DescribeDhcpOptions(DescribeDhcpOpts)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		// if err != nil {

		//     return fmt.Errorf("Error creating NAT Gateway: %s", err)
		// }

		if err != nil {

			if strings.Contains(fmt.Sprint(err), "InvalidDhcpOptionsID.NotFound") {
				resp = nil
			} else {
				log.Printf("Error on DHCPOptionsStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		dos := resp.DhcpOptions[0]
		return dos, "created", nil
	}
}

func findVPCsByDHCPOptionsID(conn *fcu.Client, id string) ([]*fcu.Vpc, error) {
	req := &fcu.DescribeVpcsInput{
		Filters: []*fcu.Filter{
			&fcu.Filter{
				Name: aws.String("dhcp-options-id"),
				Values: []*string{
					aws.String(id),
				},
			},
		},
	}

	var resp *fcu.DescribeVpcsOutput

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = conn.VM.DescribeVpcs(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("InvalidVpcID.NotFound: %s", err)
	}

	// resp, err := conn.VM.DescribeVpcs(req)
	// if err != nil {
	// 	if strings.Contains(fmt.Sprint(err), "InvalidVpcID.NotFound") {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	return resp.Vpcs, nil
}
