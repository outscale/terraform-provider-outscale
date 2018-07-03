package outscale

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func datasourceOutscaleEIMServerCertificate() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleEIMServerCertificateRead,
		Schema: map[string]*schema.Schema{
			"server_certificate_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateMaxLength(128),
			},
			"certificate_body": {
				Type:      schema.TypeString,
				Computed:  true,
				StateFunc: normalizeCert,
			},
			"certificate_chain": {
				Type:      schema.TypeString,
				Computed:  true,
				StateFunc: normalizeCert,
			},
			"server_certificate_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_certificate_name": {
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
		},
	}
}

func datasourceOutscaleEIMServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	resp, err := conn.API.GetServerCertificate(&eim.GetServerCertificateInput{
		ServerCertificateName: aws.String(d.Get("server_certificate_name").(string)),
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

	d.Set("certificate_body", aws.StringValue(server.ServerCertificate.CertificateBody))
	d.Set("certificate_chain", aws.StringValue(server.ServerCertificate.CertificateChain))

	serverMeta := make(map[string]interface{})
	id := d.Get("server_certificate_name").(string)
	if server.ServerCertificate.ServerCertificateMetadata != nil {
		id = aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.ServerCertificateID)
		serverMeta["arn"] = aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.Arn)
		serverMeta["path"] = aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.Path)
		serverMeta["server_certificate_id"] = aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.ServerCertificateID)
		serverMeta["server_certificate_name"] = aws.StringValue(server.ServerCertificate.ServerCertificateMetadata.ServerCertificateName)
	}

	if errMeta := d.Set("server_certificate_metadata", serverMeta); errMeta != nil {
		return errMeta
	}

	d.SetId(id)

	return d.Set("request_id", aws.StringValue(resp.ResponseMetadata.RequestID))
}
