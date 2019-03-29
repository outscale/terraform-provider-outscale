package outscale

import (
	"errors"
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

func resourceOutscaleLinInternetGatewayLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleLinInternetGatewayLinkCreate,
		Read:   resourceOutscaleLinInternetGatewayLinkRead,
		Delete: resourceOutscaleLinInternetGatewayLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: getLinInternetGatewayLinkSchema(),
	}
}

func resourceOutscaleLinInternetGatewayLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	vpcID := d.Get("vpc_id").(string)
	igID := d.Get("internet_gateway_id").(string)

	req := &fcu.AttachInternetGatewayInput{
		VpcId:             aws.String(vpcID),
		InternetGatewayId: aws.String(igID),
	}

	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		_, err = conn.VM.AttachInternetGateway(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		return fmt.Errorf("[DEBUG] Error linking Internet Gateway id (%s)", err)
	}

	d.SetId(igID)

	return resourceOutscaleLinInternetGatewayLinkRead(d, meta)
}

func resourceOutscaleLinInternetGatewayLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	id := d.Id()

	d.SetId(id)

	log.Printf("[DEBUG] Reading LIN Internet Gateway id (%s)", id)

	req := &fcu.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(id)},
	}

	var resp *fcu.DescribeInternetGatewaysOutput
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.VM.DescribeInternetGateways(req)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading LIN Internet Gateway id (%s)", err)
		return err
	}

	if resp == nil {
		d.SetId("")
		return errors.New("Got a nil response for Internet Gateway Link")
	}

	if resp.InternetGateways == nil {
		return errors.New("Failed to retrieve attachments Internet Gateway Link")
	}

	if len(resp.InternetGateways) > 0 {
		attchs := flattenInternetAttachements(resp.InternetGateways[0].Attachments)
		d.Set("attachment_set", attchs)
	}

	if err := d.Set("tag_set", dataSourceTags(resp.InternetGateways[0].Tags)); err != nil {
		return err
	}

	d.Set("internet_gateway_id", resp.InternetGateways[0].InternetGatewayId)
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleLinInternetGatewayLinkDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceOutscaleInternetGatewayDetach(d, meta)
}

func resourceOutscaleInternetGatewayDetach(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Get the old VPC ID to detach from
	vpcID, _ := d.GetChange("vpc_id")

	if vpcID.(string) == "" {
		log.Printf(
			"[DEBUG] Not detaching Internet Gateway '%s' as no VPC ID is set",
			d.Id())
		return nil
	}
	log.Printf(
		"[INFO] Detaching Internet Gateway '%s' from VPC '%s'",
		d.Id(),
		vpcID.(string))

	// Wait for it to be fully detached before continuing
	log.Printf("[DEBUG] Waiting for internet gateway (%s) to detach", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:        []string{"detaching"},
		Target:         []string{"detached"},
		Refresh:        detachIGStateRefreshFunc(conn, d.Id(), vpcID.(string)),
		Timeout:        15 * time.Minute,
		Delay:          10 * time.Second,
		NotFoundChecks: 30,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for internet gateway (%s) to detach: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func detachIGStateRefreshFunc(conn *fcu.Client, gatewayID, vpcID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := &fcu.DetachInternetGatewayInput{
			VpcId:             aws.String(vpcID),
			InternetGatewayId: aws.String(gatewayID),
		}

		var err error
		err = resource.Retry(120*time.Second, func() *resource.RetryError {
			_, err = conn.VM.DetachInternetGateway(req)

			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok {
				switch ec2err.Code() {
				case "InvalidInternetGatewayID.NotFound":
					log.Printf("[TRACE] Error detaching Internet Gateway '%s' from VPC '%s': %s", gatewayID, vpcID, err)
					return nil, "", nil
				case "Gateway.NotAttached":
					return 42, "detached", nil
				case "DependencyViolation":
					out, err := findPublicNetworkInterfacesForVpcID(conn, vpcID)
					if err != nil {
						return 42, "detaching", err
					}
					if len(out.NetworkInterfaces) > 0 {
						log.Printf("[DEBUG] Waiting for the following %d ENIs to be gone: %s",
							len(out.NetworkInterfaces), out.NetworkInterfaces)
					}
					return 42, "detaching", nil
				}
			}
			return 42, "", err
		}
		return 42, "detached", nil
	}
}

func findPublicNetworkInterfacesForVpcID(conn *fcu.Client, vpcID string) (*fcu.DescribeNetworkInterfacesOutput, error) {
	return conn.VM.DescribeNetworkInterfaces(&fcu.DescribeNetworkInterfacesInput{
		Filters: []*fcu.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
			{
				Name:   aws.String("association.public-ip"),
				Values: []*string{aws.String("*")},
			},
		},
	})
}

func getLinInternetGatewayLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"vpc_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"internet_gateway_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		// Attributes
		"attachment_set": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"vpc_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"tag_set": dataSourceTagsSchema(),
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
