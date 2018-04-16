package outscale

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func TestAccOutscaleImageCopy(t *testing.T) {
	var amiId string
	o := os.Getenv("OUTSCALE_OAPI")

	oapi, err := strconv.ParseBool(o)
	if err != nil {
		oapi = false
	}

	if oapi {
		t.Skip()
	}

	snapshots := []string{}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOutscaleImageCopyConfig,
				Check: func(state *terraform.State) error {
					rs, ok := state.RootModule().Resources["outscale_image_copy.test"]
					if !ok {
						return fmt.Errorf("Image resource not found")
					}

					amiId = rs.Primary.ID

					if amiId == "" {
						return fmt.Errorf("Image id is not set")
					}

					conn := testAccProvider.Meta().(*OutscaleClient)
					req := &fcu.DescribeImagesInput{
						ImageIds: []*string{aws.String(amiId)},
					}
					describe, err := conn.DescribeImages(req)
					if err != nil {
						return err
					}

					if len(describe.Images) != 1 ||
						*describe.Images[0].ImageId != rs.Primary.ID {
						return fmt.Errorf("Image not found")
					}

					image := describe.Images[0]
					if expected := "available"; *image.State != expected {
						return fmt.Errorf("invalid image state; expected %v, got %v", expected, image.State)
					}
					if expected := "machine"; *image.ImageType != expected {
						return fmt.Errorf("wrong image type; expected %v, got %v", expected, image.ImageType)
					}
					if expected := "terraform-acc-ami-copy"; *image.Name != expected {
						return fmt.Errorf("wrong name; expected %v, got %v", expected, image.Name)
					}

					for _, bdm := range image.BlockDeviceMappings {
						// The snapshot ID might not be set,
						// even for a block device that is an
						// EBS volume.
						if bdm.Ebs != nil && bdm.Ebs.SnapshotId != nil {
							snapshots = append(snapshots, *bdm.Ebs.SnapshotId)
						}
					}

					if expected := 1; len(snapshots) != expected {
						return fmt.Errorf("wrong number of snapshots; expected %v, got %v", expected, len(snapshots))
					}

					return nil
				},
			},
		},
		CheckDestroy: func(state *terraform.State) error {
			conn := testAccProvider.Meta().(*OutscaleClient)
			diReq := &fcu.DescribeImagesInput{
				ImageIds: []*string{aws.String(amiId)},
			}
			diRes, err := conn.DescribeImages(diReq)
			if err != nil {
				return err
			}

			if len(diRes.Images) > 0 {
				state := diRes.Images[0].State
				return fmt.Errorf("Image %v remains in state %v", amiId, state)
			}

			stillExist := make([]string, 0, len(snapshots))
			checkErrors := make(map[string]error)
			for _, snapshotId := range snapshots {
				dsReq := &fcu.DescribeSnapshotsInput{
					SnapshotIds: []*string{aws.String(snapshotId)},
				}
				_, err := conn.DescribeSnapshots(dsReq)
				if err == nil {
					stillExist = append(stillExist, snapshotId)
					continue
				}

				awsErr, ok := err.(awserr.Error)
				if !ok {
					checkErrors[snapshotId] = err
					continue
				}

				if awsErr.Code() != "InvalidSnapshot.NotFound" {
					checkErrors[snapshotId] = err
					continue
				}
			}

			if len(stillExist) > 0 || len(checkErrors) > 0 {
				errParts := []string{
					"Expected all snapshots to be gone, but:",
				}
				for _, snapshotId := range stillExist {
					errParts = append(
						errParts,
						fmt.Sprintf("- %v still exists", snapshotId),
					)
				}
				for snapshotId, err := range checkErrors {
					errParts = append(
						errParts,
						fmt.Sprintf("- checking %v gave error: %v", snapshotId, err),
					)
				}
				return errors.New(strings.Join(errParts, "\n"))
			}

			return nil
		},
	})
}

var testAccOutscaleImageCopyConfig = `
provider "outscale" {
	region = "us-east-2a"
}
// An AMI can't be directly copied from one account to another, and
// we can't rely on any particular AMI being available since anyone
// can run this test in whatever account they like.
// Therefore we jump through some hoops here:
//  - Spin up an EC2 instance based on a public AMI
//  - Create an AMI by snapshotting that EC2 instance, using
//    aws_ami_from_instance .
//  - Copy the new AMI using aws_ami_copy .
//
// Thus this test can only succeed if the aws_ami_from_instance resource
// is working. If it's misbehaving it will likely cause this test to fail too.
// Since we're booting a t2.micro HVM instance we need a VPC for it to boot
// up into.
resource "outscale_vpc" "foo" {
	cidr_block = "10.1.0.0/16"
}
resource "outscale_subnet" "foo" {
	cidr_block = "10.1.1.0/24"
	vpc_id = "${aws_vpc.foo.id}"
}
resource "outscale_instance" "test" {
    // This AMI has one block device mapping, so we expect to have
    // one snapshot in our created AMI.
    // This is an Ubuntu Linux HVM AMI. A public HVM AMI is required
    // because paravirtual images cannot be copied between accounts.
    image = "ami-0f8bce65"
    instance_type = "t2.micro"
    tags {
        Name = "terraform-acc-ami-copy-victim"
    }
    subnet_id = "${outscale_subnet.foo.id}"
}
resource "outscale_ami_from_instance" "test" {
    name = "terraform-acc-ami-copy-victim"
    description = "Testing Terraform aws_ami_from_instance resource"
    source_instance_id = "${outscale_instance.test.id}"
}
resource "outscale_image_copy" "test" {
    name = "terraform-acc-ami-copy"
    description = "Testing Terraform aws_ami_copy resource"
    source_image_id = "${aws_ami_from_instance.test.id}"
    source_region = "us-east-2a"
}
`
