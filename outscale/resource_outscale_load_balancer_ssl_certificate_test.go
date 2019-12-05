package outscale

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOutscaleLBUOAPISSLCertificate_basic(t *testing.T) {
	t.Skip()

	//WIP: Missing correct test case
	rInt := acctest.RandIntRange(0, 10)
	unixFile := "test-fixtures/eim-ssl-unix-line-endings.pem"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			skipIfNoOAPI(t)
			testAccPreCheck(t)
		},
		IDRefreshName: "outscale_load_balancer_ssl_certificate.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOutscaleOAPILBUDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOutscaleOAPISSLCertificateConfig(rInt, unixFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"outscale_load_balancer_ssl_certificate.test", "ssl_certificate.load_balancer_port", "80"),
				),
			},
		},
	})
}

func testAccOutscaleOAPISSLCertificateConfig(r int, fName string) string {
	return fmt.Sprintf(`
resource "outscale_load_balancer" "bar" {
  sub_region_name = ["eu-west-2a"]
  load_balancer_name = "foobar-terraform-lbu-%d"
  listeners {
    backend_port = 8000
    backend_protocol = "HTTP"
    load_balancer_port = 80
    load_balancer_protocol = "HTTP"
  }

	tag {
		bar = "baz"
	}

}

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

resource "outscale_load_balancer_ssl_certificate" "test" {
	load_balancer_name = "${outscale_load_balancer.bar.id}"
	load_balancer_port = "${outscale_load_balancer.bar.listeners.0.load_balancer_port}"
	server_certificate_id = "${outscale_server_certificate.test_cert.id}"
}
`, r, r, fName)

}
