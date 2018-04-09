package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIRouteTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIRouteTableAssociationCreate,
		Read:   resourceOutscaleOAPIRouteTableAssociationRead,
		Update: resourceOutscaleOAPIRouteTableAssociationUpdate,
		Delete: resourceOutscaleOAPIRouteTableAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"link_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIRouteTableAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf(
		"[INFO] Creating route table association: %s => %s",
		d.Get("subnet_id").(string),
		d.Get("route_table_id").(string))

	associationOpts := fcu.AssociateRouteTableInput{
		RouteTableId: aws.String(d.Get("route_table_id").(string)),
		SubnetId:     aws.String(d.Get("subnet_id").(string)),
	}

	var resp *fcu.AssociateRouteTableOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.AssociateRouteTable(&associationOpts)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidRouteTableID.NotFound") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Set the ID and return
	d.SetId(*resp.AssociationId)
	d.Set("link_id", d.Id())
	log.Printf("[INFO] Association ID: %s", d.Id())

	return nil
}

func resourceOutscaleOAPIRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	rtRaw, _, err := resourceOutscaleRouteTableStateRefreshFunc(
		conn, d.Get("route_table_id").(string))()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(*fcu.RouteTable)

	found := false
	for _, a := range rt.Associations {
		if *a.RouteTableAssociationId == d.Id() {
			found = true
			d.Set("subnet_id", *a.SubnetId)
			break
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func resourceOutscaleOAPIRouteTableAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf(
		"[INFO] Creating route table association: %s => %s",
		d.Get("subnet_id").(string),
		d.Get("route_table_id").(string))

	req := &fcu.ReplaceRouteTableAssociationInput{
		AssociationId: aws.String(d.Id()),
		RouteTableId:  aws.String(d.Get("route_table_id").(string)),
	}

	var resp *fcu.ReplaceRouteTableAssociationOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.ReplaceRouteTableAssociation(req)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidAssociationID.NotFound") {
			return resourceOutscaleOAPIRouteTableAssociationCreate(d, meta)
		}
		return err
	}

	d.SetId(*resp.NewAssociationId)
	log.Printf("[INFO] Association ID: %s", d.Id())

	return nil
}

func resourceOutscaleOAPIRouteTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	log.Printf("[INFO] Deleting route table association: %s", d.Id())

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DisassociateRouteTable(&fcu.DisassociateRouteTableInput{
			AssociationId: aws.String(d.Id()),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidAssociationID.NotFound") {
			return nil
		}
		return fmt.Errorf("Error deleting route table association: %s", err)
	}

	return nil
}
