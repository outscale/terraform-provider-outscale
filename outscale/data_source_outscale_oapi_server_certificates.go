package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func datasourceOutscaleOAPIEIMServerCertificates() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleOAPIEIMServerCertificatesRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
			},
			"server_certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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

func datasourceOutscaleOAPIEIMServerCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	var listResp *eim.ListServerCertificatesOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		listResp, err = conn.API.ListServerCertificates(&eim.ListServerCertificatesInput{
			PathPrefix: aws.String(d.Get("path").(string)),
		})

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Throttling:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	listMeta := make([]map[string]interface{}, len(listResp.ListServerCertificatesResult.ServerCertificateMetadataList))

	for k, v := range listResp.ListServerCertificatesResult.ServerCertificateMetadataList {
		item := make(map[string]interface{})
		item["path"] = aws.StringValue(v.Path)
		item["server_certificate_id"] = aws.StringValue(v.ServerCertificateID)
		item["server_certificate_name"] = aws.StringValue(v.ServerCertificateName)
		listMeta[k] = item
	}

	d.Set("server_certificates", listMeta)

	d.SetId(resource.UniqueId())

	return d.Set("request_id", aws.StringValue(listResp.ResponseMetadata.RequestID))
}
