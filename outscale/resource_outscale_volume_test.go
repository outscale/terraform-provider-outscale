package outscale

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOutscaleVolume_basic(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		//IDRefreshName: "outscale_volume.test",
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsEbsVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "attachment_set.#", "0"),
				),
			},
		},
	})
}

func TestAccOutscaleVolume_kmsKey(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	var v fcu.Volume
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAwsEbsVolumeConfigWithKmsKey, ri)
	keyRegex := regexp.MustCompile("^arn:aws:([a-zA-Z0-9\\-])+:([a-z]{2}-[a-z]+-\\d{1})?:(\\d{12})?:(.*)$")

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.test", &v),
					resource.TestCheckResourceAttr("outscale_volume.test", "encrypted", "true"),
					resource.TestMatchResourceAttr("outscale_volume.test", "kms_key_id", keyRegex),
				),
			},
		},
	})
}

func TestAccOutscaleVolume_NoIops(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsEbsVolumeConfigWithNoIops,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.iops_test", &v),
				),
			},
		},
	})
}

func TestAccOutscaleVolume_withTags(t *testing.T) {
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi != false {
		t.Skip()
	}
	var v fcu.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "outscale_volume.tags_test",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsEbsVolumeConfigWithTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("outscale_volume.tags_test", &v),
				),
			},
		},
	})
}

func testAccCheckVolumeExists(n string, v *fcu.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*OutscaleClient).FCU

		request := &fcu.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(rs.Primary.ID)},
		}

		var err error
		var response *fcu.DescribeVolumesOutput

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			response, err = conn.VM.DescribeVolumes(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})
		fmt.Printf("[DEBUG] Error Test Exists: %s", err)
		fmt.Printf("[DEBUG] Volume Exists: %v ", *response)
		if err == nil {
			if response.Volumes != nil && len(response.Volumes) > 0 {
				*v = *response.Volumes[0]
				return nil
			}
		}

		return fmt.Errorf("Error finding Outscale volume %s", rs.Primary.ID)
	}
}

const testAccAwsEbsVolumeConfig = `
resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  volume_type = "gp2"
  size = 1
  tags {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testAccAwsEbsAttachedVolumeConfig = `
data "aws_ami" "debian_jessie_latest" {
  most_recent = true

  filter {
    name   = "name"
    values = ["debian-jessie-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  owners = ["379101102735"] # Debian
}

resource "aws_instance" "test" {
  ami = "${data.aws_ami.debian_jessie_latest.id}"
  associate_public_ip_address = true
  count = 1
  instance_type = "t2.medium"

  root_block_device {
    volume_size           = "10"
    volume_type           = "standard"
    delete_on_termination = true
  }

  tags {
    Name    = "test-terraform"
  }
}

resource "outscale_volume" "test" {
  depends_on = ["aws_instance.test"]
  availability_zone = "${aws_instance.test.availability_zone}"
  volume_type = "gp2"
  size = "10"
}

resource "aws_volume_attachment" "test" {
  depends_on  = ["outscale_volume.test"]
  device_name = "/dev/xvdg"
  volume_id   = "${outscale_volume.test.id}"
  instance_id = "${aws_instance.test.id}"
}
`

const testAccAwsEbsAttachedVolumeConfigUpdateSize = `
data "aws_ami" "debian_jessie_latest" {
  most_recent = true

  filter {
    name   = "name"
    values = ["debian-jessie-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  owners = ["379101102735"] # Debian
}

resource "aws_instance" "test" {
  ami = "${data.aws_ami.debian_jessie_latest.id}"
  associate_public_ip_address = true
  count = 1
  instance_type = "t2.medium"

  root_block_device {
    volume_size           = "10"
    volume_type           = "standard"
    delete_on_termination = true
  }

  tags {
    Name    = "test-terraform"
  }
}

resource "outscale_volume" "test" {
  depends_on = ["aws_instance.test"]
  availability_zone = "${aws_instance.test.availability_zone}"
  volume_type = "gp2"
  size = "20"
}

resource "aws_volume_attachment" "test" {
  depends_on  = ["outscale_volume.test"]
  device_name = "/dev/xvdg"
  volume_id   = "${outscale_volume.test.id}"
  instance_id = "${aws_instance.test.id}"
}
`

const testAccAwsEbsVolumeConfigUpdateSize = `
resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  volume_type = "gp2"
  size = 10
  tags {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testAccAwsEbsVolumeConfigUpdateType = `
resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  volume_type = "sc1"
  size = 500
  tags {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testAccAwsEbsVolumeConfigWithIops = `
resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  volume_type = "io1"
  size = 4
  iops = 100
  tags {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testAccAwsEbsVolumeConfigWithIopsUpdated = `
resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  volume_type = "io1"
  size = 4
  iops = 200
  tags {
    Name = "tf-acc-test-ebs-volume-test"
  }
}
`

const testAccAwsEbsVolumeConfigWithKmsKey = `
resource "aws_kms_key" "foo" {
  description = "Terraform acc test %d"
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "kms-tf-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  size = 1
  encrypted = true
  kms_key_id = "${aws_kms_key.foo.arn}"
}
`

const testAccAwsEbsVolumeConfigWithTags = `
resource "outscale_volume" "tags_test" {
  availability_zone = "eu-west-2a"
  size = 1
  tags {
    Name = "TerraformTest"
  }
}
`

const testAccAwsEbsVolumeConfigWithNoIops = `
resource "outscale_volume" "iops_test" {
  availability_zone = "eu-west-2a"
  size = 10
  volume_type = "gp2"
  iops = 0
  tags {
    Name = "TerraformTest"
  }
}
`
