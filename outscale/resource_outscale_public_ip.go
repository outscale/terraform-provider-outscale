package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
)

func resourceOutscaleOAPIPublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPublicIPCreate,
		Read:   resourceOutscaleOAPIPublicIPRead,
		Delete: resourceOutscaleOAPIPublicIPDelete,
		Update: resourceOutscaleOAPIPublicIPUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: getOAPIPublicIPSchema(),
	}
}

func resourceOutscaleOAPIPublicIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI
	allocOpts := oapi.CreatePublicIpRequest{}

	log.Printf("[DEBUG] EIP create configuration: %#v", allocOpts)
	resp, err := conn.POST_CreatePublicIp(allocOpts)
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	if resp.OK == nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}

	allocResp := resp.OK

	log.Printf("[DEBUG] EIP Allocate: %#v", allocResp)

	d.SetId(allocResp.PublicIp.PublicIpId)

	log.Printf("[DEBUG], allocResp: %+v", allocResp)

	//SetTags
	if tags, ok := d.GetOk("tags"); ok {
		err := assignOapiTags(tags.([]interface{}), allocResp.PublicIp.PublicIpId, conn)
		if err != nil {
			return err
		}
	}

	return resourceOutscaleOAPIPublicIPRead(d, meta)
}

func resourceOutscaleOAPIPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()

	req := oapi.ReadPublicIpsRequest{}

	req.Filters.PublicIpIds = []string{id}

	var describeAddresses *oapi.ReadPublicIpsResponse
	resp, err := conn.POST_ReadPublicIps(req)
	describeAddresses = resp.OK

	if err != nil {
		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving EIP: %s", err)
	}

	if len(describeAddresses.PublicIps) != 1 || describeAddresses.PublicIps[0].PublicIpId != id {
		return fmt.Errorf("Unable to find Public IP: %#v", describeAddresses.PublicIps)

	}

	publicIP := describeAddresses.PublicIps[0]

	log.Printf("[DEBUG] EIP read configuration: %+v", publicIP)

	if publicIP.LinkPublicIpId != "" {
		d.Set("link_public_ip_id", publicIP.LinkPublicIpId)
	} else {
		d.Set("link_public_ip_id", "")
	}
	if publicIP.VmId != "" {
		d.Set("vm_id", publicIP.VmId)
	} else {
		d.Set("vm_id", "")
	}
	if publicIP.NicId != "" {
		d.Set("nic_id", publicIP.NicId)
	} else {
		d.Set("nic_id", "")
	}
	if publicIP.NicAccountId != "" {
		d.Set("nic_account_id", publicIP.NicAccountId)
	} else {
		d.Set("nic_account_id", "")
	}
	d.Set("private_ip", publicIP.PrivateIp)
	d.Set("public_ip", publicIP.PublicIp)

	d.Set("public_ip_id", publicIP.PublicIpId)

	if err := d.Set("tags", tagsOAPIToMap(publicIP.Tags)); err != nil {
		log.Printf("[WARN] error setting tags for PublicIp(%s): %s", publicIP.PublicIp, err)
	}

	return d.Set("request_id", describeAddresses.ResponseContext.RequestId)
}

func resourceOutscaleOAPIPublicIPUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	d.Partial(true)

	if err := setOAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)

	return resourceOutscaleOAPIPublicIPRead(d, meta)
}

func resourceOutscaleOAPIPublicIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	if err := resourceOutscaleOAPIPublicIPRead(d, meta); err != nil {
		return err
	}
	if d.Id() == "" {
		return nil
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		var err error

		log.Printf("[DEBUG] EIP release (destroy) IP address: %v", d.Id())
		_, err = conn.POST_DeletePublicIp(oapi.DeletePublicIpRequest{
			PublicIpId: d.Id(),
		})

		if e := fmt.Sprint(err); strings.Contains(e, "InvalidAllocationID.NotFound") || strings.Contains(e, "InvalidAddress.NotFound") {
			return nil
		}

		if err == nil {
			return nil
		}
		if _, ok := err.(awserr.Error); !ok {
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(err)
	})
}

func getOAPIPublicIPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"public_ip_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"link_public_ip_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nic_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"nic_account_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tags": tagsListOAPISchema(),
	}
}
