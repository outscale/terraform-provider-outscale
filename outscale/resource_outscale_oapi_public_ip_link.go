package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
)

func resourceOutscaleOAPIPublicIPLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIPublicIPLinkCreate,
		Read:   resourceOutscaleOAPIPublicIPLinkRead,
		Delete: resourceOutscaleOAPIPublicIPLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"reservation_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"allow_relink": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"vm_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"nic_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"link_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"reques_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIPublicIPLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	request := oapi.LinkPublicIpRequest{}

	if v, ok := d.GetOk("reservation_id"); ok {
		fmt.Println(v.(string))
		request.ReservationId = v.(string)
	}
	if v, ok := d.GetOk("allow_relink"); ok {
		request.AllowRelink = v.(bool)
	}
	if v, ok := d.GetOk("vm_id"); ok {
		request.VmId = v.(string)
	}
	if v, ok := d.GetOk("nic_id"); ok {
		request.NicId = v.(string)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		request.PrivateIp = v.(string)
	}
	if v, ok := d.GetOk("public_ip"); ok {
		request.PublicIp = v.(string)
	}

	log.Printf("[DEBUG] EIP association configuration: %#v", request)

	var resp *oapi.LinkPublicIpResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {

		response, err := conn.POST_LinkPublicIp(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		resp = response.OK
		return nil
	})

	if err != nil {
		fmt.Printf("[WARN] ERROR resourceOutscaleOAPIPublicIPLinkCreate (%s)", err)
		return err
	}

	//Missing on swagger spec
	// if resp != nil && resp.ReservationId != "" && len(*resp.ReservationId) > 0 {
	// 	d.SetId(*resp.AssociationId)
	// } else {
	// 	d.SetId(*request.PublicIp)
	// }

	//Using validation with request.
	if resp != nil && request.ReservationId != "" && len(request.ReservationId) > 0 {
		d.SetId(resp.LinkId)
	} else {
		d.SetId(request.PublicIp)
	}
	return resourceOutscaleOAPIPublicIPLinkRead(d, meta)
}

func resourceOutscaleOAPIPublicIPLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	id := d.Id()
	var request oapi.ReadPublicIpsRequest

	if strings.Contains(id, "eipassoc") {
		request = oapi.ReadPublicIpsRequest{
			Filters: oapi.ReadPublicIpsFilters{
				ReservationIds: []string{id},
			},
		}
	} else {
		request = oapi.ReadPublicIpsRequest{
			Filters: oapi.ReadPublicIpsFilters{
				PublicIps: []string{id},
			},
		}
	}

	var response *oapi.ReadPublicIpsResponse
	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {

		res, err := conn.POST_ReadPublicIps(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = res.OK
		return nil
	})

	fmt.Printf("[WARN] ERROR resourceOutscaleOAPIPublicIPLinkRead (%s)", err)

	if err != nil {
		return fmt.Errorf("Error reading Outscale VM Public IP %s: %#v", d.Get("reservation_id").(string), err)
	}

	if response.PublicIps == nil || len(response.PublicIps) == 0 {
		log.Printf("[INFO] EIP Association ID Not Found. Refreshing from state")
		d.SetId("")
		return nil
	}

	d.Set("request_id", response.ResponseContext.RequestId)
	return readOutscaleOAPIPublicIPLink(d, response.PublicIps[0])
}

func resourceOutscaleOAPIPublicIPLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	assocID := d.Get("link_id")

	opts := oapi.UnlinkPublicIpRequest{
		LinkId: assocID.(string),
	}

	var err error

	err = resource.Retry(60*time.Second, func() *resource.RetryError {

		_, err = conn.POST_UnlinkPublicIp(opts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	fmt.Printf("[WARN] ERROR resourceOutscaleOAPIPublicIPLinkDelete (%s)", err)

	if err != nil {
		return fmt.Errorf("Error deleting Elastic IP association: %s", err)
	}

	return nil
}

func readOutscaleOAPIPublicIPLink(d *schema.ResourceData, address oapi.PublicIps) error {
	if err := d.Set("reservation_id", address.ReservationId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink1 (%s)", err)

		return err
	}
	if err := d.Set("vm_id", address.VmId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink2 (%s)", err)

		return err
	}
	if err := d.Set("nic_id", address.NicId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink3 (%s)", err)

		return err
	}
	if err := d.Set("private_ip", address.PrivateIp); err != nil {
		fmt.Printf("[WARN] ERROR readOutscalePublicIPLink4 (%s)", err)

		return err
	}
	if err := d.Set("public_ip", address.PublicIp); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleOAPIPublicIPLink (%s)", err)

		return err
	}

	if err := d.Set("link_id", address.LinkId); err != nil {
		fmt.Printf("[WARN] ERROR readOutscaleOAPIPublicIPLink (%s)", err)

		return err
	}

	return nil
}
