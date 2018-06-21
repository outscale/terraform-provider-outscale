package outscale

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/eim"
)

func TestAccOutscaleOAPIServerCertificate_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if !oapi {
		t.Skip()
	}

	var cert eim.ServerCertificate
	rInt := acctest.RandInt()
	rIntUp := acctest.RandInt()
	unixFile := "test-fixtures/eim-ssl-unix-line-endings.pem"
	winFile := "test-fixtures/eim-ssl-windows-line-endings.pem.winfile"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOAPIServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOAPIServerCertConfigFile(rInt, unixFile),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCertExists("outscale_server_certificate.test_cert", &cert),
				),
			},
			{
				Config: testAccOAPIServerCertConfigFile(rInt, winFile),
				Check: resource.ComposeTestCheckFunc(
					testAccOAPICheckCertExists("outscale_server_certificate.test_cert", &cert),
				),
			},
			{
				Config: testAccOAPIServerCertConfigFile(rIntUp, winFile),
				Check: resource.ComposeTestCheckFunc(
					testAccOAPICheckCertExists("outscale_server_certificate.test_cert", &cert),
				),
			},
		},
	})
}
func testAccOAPICheckCertExists(n string, cert *eim.ServerCertificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server Cert ID is set")
		}
		conn := testAccProvider.Meta().(*OutscaleClient).EIM
		describeOpts := &eim.GetServerCertificateInput{
			ServerCertificateName: aws.String(rs.Primary.Attributes["server_certificate_name"]),
		}
		resp, err := conn.API.GetServerCertificate(describeOpts)
		if err != nil {
			return err
		}
		*cert = *resp.ServerCertificate
		return nil
	}
}
func testAccCheckOAPIServerCertificateDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*OutscaleClient).EIM
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "outscale_eim_server_certificate" {
			continue
		}
		// Try to find the Cert
		opts := &eim.GetServerCertificateInput{
			ServerCertificateName: aws.String(rs.Primary.Attributes["server_certificate_nameonroiroiroirnroinroin"]),
		}
		resp, err := conn.API.GetServerCertificate(opts)
		if err == nil {
			if resp.ServerCertificate != nil {
				return fmt.Errorf("Error: Server Cert still exists")
			}
			return nil
		}
	}
	return nil
}

// eim-ssl-unix-line-endings
func testAccOAPIServerCertConfigFile(rInt int, fName string) string {
	return fmt.Sprintf(`
resource "outscale_server_certificate" "test_cert" {
  server_certificate_name = "terraform-test-cert-%d"
	server_certificate_body = "${file("%s")}"
	private_key =  <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDKdH6BU9Q0xBVPfeX5NjCC/B2Pm3WsFGnTtRw4abkD+r4to9wD
eYUgjH2yPCyonNOA8mNiCQgDTtaLfbA8LjBYoodt7rgaTO7C0ugRtmTNK96DmYxm
f8Gs5ZS6eC3yeaFv58d1w2mow7tv0+DRk8uXwzVfaaMxoalsCtlLznmZHwIDAQAB
AoGABZj69nBu6ZaSUERW23EYHkcCOjo+Iqfd1TCouxaROv7vyytApgfyGlhIEWmA
gpjzcBlDji5Zvl2rqOesu707MOuJavZvluo+JHy/VIuU+yGUrWuO/QVCu6Jn3yns
vS7g48ConuZ962cTzRPcpPDspONBVOAhVCF33Y8PsnxV0wECQQD5RqeoqxEUupsy
QhrDui0KkYXLdT0uhrEQ69n9rvAiQoHPsiX0MswfEKnj/g9N3VwGLdgWytT0TvcI
8fDPRB4/AkEAz+qF3taX77gB69XRPQwCGWqE1fHIFMwX7QeYdEsk3iRZ0EKVcdp6
vIPCB2Cq4a4eXcaFa/bXen4yeYgyTbeNIQJBAO92dWctdoowPRiJskZmGhC1/Q6X
gH+qenyj5VSy8hInS6anH5i4F6icDGhtzmvhgx6YeaZjkTFkjiG0sb2aVWcCQQDD
WL7UwtzX/xPXB/ril5C1Xo5WESgC2ks0ielkgmGuUYsNEDInWbXtvwGjOuDyz0x6
oRYkfTSxQzabVyqkOGvhAkBtbjUxOD8wgBIjb4T6mAMokQo6PeEAZGUTyPifjJNo
detWVr2WRvgNgQvcRnNPECwfq1RtMJJpavaI3kgeaSxg
-----END RSA PRIVATE KEY-----
EOF
}
`, rInt, fName)
}
