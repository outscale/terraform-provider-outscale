package outscale

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func datasourceOutscaleOAPIEIMServerCertificate() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleOAPIEIMServerCertificateRead,
		Schema: map[string]*schema.Schema{
			"path": { //it's correct
				Type:     schema.TypeString,
				Required: true,
			},
			"server_certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_certificate_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceOutscaleOAPIEIMServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	resp, err := conn.API.GetServerCertificate(&eim.GetServerCertificateInput{
		ServerCertificateName: aws.String(d.Get("path").(string)),
	})

	if resp.GetServerCertificateResult == nil {
		return fmt.Errorf("Could not get Server Certificate information")
	}

	if err != nil {
		//it is OK?
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NoSuchEntity" {
				log.Printf("[WARN] EIM Server Cert (%s) not found, removing from state", d.Id())
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[WARN] Error reading EIM Server Certificate: %s: %s", awsErr.Code(), awsErr.Message())
		}
		return fmt.Errorf("[WARN] Error reading EIM Server Certificate: %s", err)
	}

	server := resp.GetServerCertificateResult

	id := d.Get("server_certificate_name").(string)
	if server.ServerCertificate.ServerCertificateMetadata != nil {
		id = aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.ServerCertificateID)
		d.Set("server_certificate_id", aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.ServerCertificateID))
		d.Set("server_certificate_name", aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.ServerCertificateName))
	}

	d.SetId(id)

	return d.Set("request_id", aws.StringValue(resp.ResponseMetadata.RequestID))
}
