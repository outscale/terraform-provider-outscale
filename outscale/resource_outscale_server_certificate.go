package outscale

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func resourceOutscaleEIMServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleEIMServerCertificateCreate,
		Read:   resourceOutscaleEIMServerCertificateRead,
		Delete: resourceOutscaleEIMServerCertificateDelete,
		Update: resourceOutscaleEIMServerCertificateUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceOutscaleEIMServerCertificateImport,
		},
		Schema: map[string]*schema.Schema{
			"certificate_body": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: normalizeCert,
			},
			"certificate_chain": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: normalizeCert,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: normalizeCert,
				Sensitive: true,
			},
			"server_certificate_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateMaxLength(128),
			},
			"arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceOutscaleEIMServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	var sslCertName string
	if v, ok := d.GetOk("server_certificate_name"); ok {
		sslCertName = v.(string)
	} else {
		sslCertName = resource.UniqueId()
	}
	createOpts := &eim.UploadServerCertificateInput{
		CertificateBody:       aws.String(d.Get("certificate_body").(string)),
		PrivateKey:            aws.String(d.Get("private_key").(string)),
		ServerCertificateName: aws.String(sslCertName),
		Path: aws.String("/"),
	}
	if v, ok := d.GetOk("certificate_chain"); ok {
		createOpts.CertificateChain = aws.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		createOpts.Path = aws.String(v.(string))
	}
	log.Printf("[DEBUG] Creating EIM Server Certificate with opts: %+v", createOpts)
	resp, err := conn.API.UploadServerCertificate(createOpts)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return fmt.Errorf("[WARN] Error uploading server certificate, error: %s: %s", awsErr.Code(), awsErr.Message())
		}
		return fmt.Errorf("[WARN] Error uploading server certificate, error: %s", err)
	}
	d.SetId(*resp.ServerCertificateMetadata.ServerCertificateId)
	d.Set("server_certificate_name", sslCertName)
	return resourceOutscaleEIMServerCertificateRead(d, meta)
}
func resourceOutscaleEIMServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM

	resp, err := conn.API.GetServerCertificate(&eim.GetServerCertificateInput{
		ServerCertificateName: aws.String(d.Get("server_certificate_name").(string)),
	})
	if err != nil {
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
	d.SetId(*resp.ServerCertificate.ServerCertificateMetadata.ServerCertificateId)
	// these values should always be present, and have a default if not set in
	// configuration, and so safe to reference with nil checks
	d.Set("certificate_body", normalizeCert(resp.ServerCertificate.CertificateBody))
	c := normalizeCert(resp.ServerCertificate.CertificateChain)
	if c != "" {
		d.Set("certificate_chain", c)
	}
	d.Set("path", resp.ServerCertificate.ServerCertificateMetadata.Path)
	d.Set("arn", resp.ServerCertificate.ServerCertificateMetadata.Arn)

	if resp.ResponseMetadata != nil {
		d.Set("request_id", resp.ResponseMetadata.RequestId)
	}
	return nil
}

func resourceOutscaleEIMServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	if d.HasChange("server_certificate_name") {
		o, n := d.GetChange("server_certificate_name")

		updateOps := &eim.UpdateServerCertificateInput{
			ServerCertificateName:    aws.String(o.(string)),
			NewServerCertificateName: aws.String(n.(string)),
			NewPath:                  aws.String(d.Get("path").(string)),
		}
		_, err := conn.API.UpdateServerCertificate(updateOps)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				return fmt.Errorf("[WARN] Error updating server certificate, error: %s: %s", awsErr.Code(), awsErr.Message())
			}
			return fmt.Errorf("[WARN] Error updating server certificate, error: %s", err)
		}
	}

	return resourceOutscaleEIMServerCertificateRead(d, meta)
}

func resourceOutscaleEIMServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).EIM
	log.Printf("[INFO] Deleting EIM Server Certificate: %s", d.Id())
	err := resource.Retry(15*time.Minute, func() *resource.RetryError {
		_, err := conn.API.DeleteServerCertificate(&eim.DeleteServerCertificateInput{
			ServerCertificateName: aws.String(d.Get("server_certificate_name").(string)),
		})
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == "DeleteConflict" && strings.Contains(awsErr.Message(), "currently in use by arn") {
					// TODO: currentlyInUseBy(awsErr.Message(), meta.(*OutscaleClient).LBU)
					log.Printf("[WARN] Conflict deleting server certificate: %s, retrying", awsErr.Message())
					return resource.RetryableError(err)
				}
				if awsErr.Code() == "NoSuchEntity" {
					return nil
				}
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func resourceOutscaleEIMServerCertificateImport(
	d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("server_certificate_name", d.Id())
	// private_key can't be fetched from any API call
	return []*schema.ResourceData{d}, nil
}
func currentlyInUseBy(awsErr string, conn *elb.ELB) {
	r := regexp.MustCompile(`currently in use by ([a-z0-9:-]+)\/([a-z0-9-]+)\.`)
	matches := r.FindStringSubmatch(awsErr)
	if len(matches) > 0 {
		lbName := matches[2]
		describeElbOpts := &elb.DescribeLoadBalancersInput{
			LoadBalancerNames: []*string{aws.String(lbName)},
		}
		if _, err := conn.DescribeLoadBalancers(describeElbOpts); err != nil {

			if strings.Contains(fmt.Sprint(err), "LoadBalancerNotFound") {
				log.Printf("[WARN] Load Balancer (%s) causing delete conflict not found", lbName)
			}
		}
	}
}
func normalizeCert(cert interface{}) string {
	if cert == nil || cert == (*string)(nil) {
		return ""
	}
	var rawCert string
	switch cert.(type) {
	case string:
		rawCert = cert.(string)
	case *string:
		rawCert = *cert.(*string)
	default:
		return ""
	}
	cleanVal := sha1.Sum(stripCR([]byte(strings.TrimSpace(rawCert))))
	return hex.EncodeToString(cleanVal[:])
}

// strip CRs from raw literals. Lifted from go/scanner/scanner.go
// See https://github.com/golang/go/blob/release-branch.go1.6/src/go/scanner/scanner.go#L479
func stripCR(b []byte) []byte {
	c := make([]byte, len(b))
	i := 0
	for _, ch := range b {
		if ch != '\r' {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}

func validateMaxLength(length int) schema.SchemaValidateFunc {
	return validation.StringLenBetween(0, length)
}
